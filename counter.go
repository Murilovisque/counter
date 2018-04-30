package counter

import (
	"log"
	"time"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"sync"
	"sync/atomic"
)

const (
	counterCollection = "counters"
	KeyField          = "key"
	valField          = "val"
)

var (
	db                     string
	persistInterval        int
	persistenceCountersMap map[string]*persistanceCounter
	stop                   = false
	mutex                  sync.Mutex
	goroutineCounter       int64
)

// Start inform the dbname and internal to persist
func Start(dbParam string, persistIntervalParam int) {
	log.Println("Starting counter...")
	stop = false
	if db != "" && db != dbParam {
		for _, p := range persistenceCountersMap {
			p.clearAll()
		}
	}
	db = dbParam
	persistInterval = persistIntervalParam
	startPersistence()
	log.Println("counter started")
}

// AddIncrementor to persist counts
func AddIncrementor(i Incrementor) {
	if persistenceCountersMap == nil {
		persistenceCountersMap = make(map[string]*persistanceCounter)
	}
	persistenceCountersMap[i.Type()] = &persistanceCounter{incrementable: i,
		mapValuesToPersist:  make(map[string]interface{}),
		mapLastPersistedVal: make(map[string]interface{})}
}

//Inc duration of key, it does not lock the caller
func Inc(typeCounter, key string, val interface{}) {
	go func() {
		atomic.AddInt64(&goroutineCounter, 1)
		defer atomic.AddInt64(&goroutineCounter, -1)
		p, _ := persistenceCountersMap[typeCounter]
		p.inc(key, val)
	}()
}

//Stop persist all counters and then stop them
func Stop() {
	log.Println("Stopping counter...")
	stop = true
	for {
		if goroutineCounter < 1 {
			break
		}
	}
	doPersistence()
	log.Println("Counter stopped")
}

//Val the current counter value
func Val(typeCounter, key string) interface{} {
	p, _ := persistenceCountersMap[typeCounter]
	valCur, ok := p.mapLastPersistedVal[key]
	if !ok {
		valCur = p.incrementable.ZeroVal()
	}
	valPersist, ok := p.mapValuesToPersist[key]
	if !ok {
		valPersist = p.incrementable.ZeroVal()
	}
	return p.incrementable.Inc(valCur, valPersist)
}

//Clear counts
func Clear(typeCounter, key string) {
	p, _ := persistenceCountersMap[typeCounter]
	p.clear(key)
}

func (p *persistanceCounter) inc(key string, val interface{}) {
	p.mux.Lock()
	v, ok := p.mapValuesToPersist[key]
	if ok {
		v = p.incrementable.Inc(v, val)
		p.mapValuesToPersist[key] = v
	} else {
		p.mapValuesToPersist[key] = val
	}
	p.mux.Unlock()
}

func (p *persistanceCounter) getAndClear(key string) interface{} {
	p.mux.Lock()
	defer p.mux.Unlock()
	val, v := p.mapValuesToPersist[key]
	if !v {
		return p.incrementable.ZeroVal()
	}
	delete(p.mapValuesToPersist, key)
	return val
}

func (p *persistanceCounter) clear(key string) {
	p.mux.Lock()
	defer p.mux.Unlock()
	delete(p.mapValuesToPersist, key)
	delete(p.mapLastPersistedVal, key)
}

func (p *persistanceCounter) clearAll() {
	p.mapLastPersistedVal = make(map[string]interface{})
	p.mapValuesToPersist = make(map[string]interface{})
}

func startPersistence() {
	log.Printf("counter - Starting persistance each %d second(s)\n", persistInterval)
	go func() {
		ticker := time.NewTicker(time.Duration(persistInterval) * time.Second)
		for range ticker.C {
			if stop {
				break
			}
			doPersistence()
		}
	}()
}

func doPersistence() {
	mutex.Lock()
	defer mutex.Unlock()
	session, err := mgo.Dial("localhost")
	if err != nil {
		log.Println(err)
		return
	}
	defer session.Close()
	for _, p := range persistenceCountersMap {
		if len(p.mapValuesToPersist) == 0 {
			continue
		}
		session.SetMode(mgo.Monotonic, true)
		for k := range p.mapValuesToPersist {
			log.Printf("counter type %s - Persisting key: %s...\n", p.incrementable.Type(), k)
			persist(session, p, k)
		}
	}
}

// Persist make the persistance
func persist(session *mgo.Session, p *persistanceCounter, key string) {
	param := Counter{Key: key, Val: p.getAndClear(key)}
	collection := session.DB(db).C(counterCollection + p.incrementable.Type())
	c, err := p.incrementable.Counter(collection, param.Key)
	if err == mgo.ErrNotFound {
		c = &Counter{Key: param.Key, Val: param.Val}
		err = collection.Insert(&c)
		if err != nil {
			log.Println(err)
			return
		}
	} else {
		c.Val = p.incrementable.Inc(param.Val, c.Val)
		err = collection.Update(bson.M{"_id": c.ID}, bson.M{"$set": bson.M{valField: c.Val}})
		if err != nil {
			log.Println(err)
			return
		}
	}
	p.mapLastPersistedVal[key] = c.Val
	log.Printf("counter - Persisted: key %s, Val %v\n", c.Key, c.Val)
}

// Incrementor has a incrementation definition
type Incrementor interface {
	Inc(actual interface{}, add interface{}) interface{}
	Counter(collection *mgo.Collection, key string) (*Counter, error)
	ZeroVal() interface{}
	Type() string
}

// Counter model
type Counter struct {
	ID  bson.ObjectId `bson:"_id,omitempty"`
	Key string
	Val interface{}
}

type persistanceCounter struct {
	mapValuesToPersist  map[string]interface{}
	mux                 sync.Mutex
	incrementable       Incrementor
	mapLastPersistedVal map[string]interface{}
}
