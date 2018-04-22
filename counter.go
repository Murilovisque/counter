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
	goroutineCounter       uint64
)

// Start inform the dbname and internal to persist
func Start(dbParam string, persistIntervalParam int) {
	db = dbParam
	persistInterval = persistIntervalParam
	startPersistence()
}

// AddIncrementor to persist counts
func AddIncrementor(i Incrementor) {
	if persistenceCountersMap == nil {
		persistenceCountersMap = make(map[string]*persistanceCounter)
	}
	persistenceCountersMap[i.Type()] = &persistanceCounter{incrementable: i, mapDurationToPersist: make(map[string]interface{})}
}

//Inc duration of key
func Inc(typeCounter, key string, val interface{}) {
	atomic.AddUint64(&goroutineCounter, 1)
	defer atomic.AddUint64(&goroutineCounter, -1)
	p, _ := persistenceCountersMap[typeCounter]
	p.inc(key, val)
}

//Stop persist all counters and then stop them
func Stop() {
	stop = true
	for {
		if goroutineCounter == 0 {
			break
		}
	}
	doPersistence()
}

func (p *persistanceCounter) inc(key string, val interface{}) {
	p.mux.Lock()
	v, ok := p.mapDurationToPersist[key]
	if ok {
		v = p.incrementable.Inc(v, val)
		p.mapDurationToPersist[key] = v
	} else {
		p.mapDurationToPersist[key] = val
	}
	p.mux.Unlock()
}

func (p *persistanceCounter) getAndClear(key string) interface{} {
	p.mux.Lock()
	defer p.mux.Unlock()
	val, v := p.mapDurationToPersist[key]
	if !v {
		return p.incrementable.ZeroVal()
	}
	delete(p.mapDurationToPersist, key)
	return val
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
		if len(p.mapDurationToPersist) == 0 {
			return
		}
		session.SetMode(mgo.Monotonic, true)
		for k := range p.mapDurationToPersist {
			log.Printf("counter type %s - Persisting key: %s", p.incrementable.Type(), k)
			persist(session, &Counter{Key: k, Val: p.getAndClear(k)}, p.incrementable)
		}
	}
}

// Persist make the persistance
func persist(session *mgo.Session, param *Counter, incrementable Incrementor) {
	collection := session.DB(db).C(counterCollection + incrementable.Type())
	c, err := incrementable.Val(collection, param.Key)
	if err == mgo.ErrNotFound {
		c = &Counter{Key: param.Key, Val: param.Val}
		err = collection.Insert(&c)
		if err != nil {
			log.Println(err)
			return
		}
	} else {
		c.Val = incrementable.Inc(param.Val, c.Val)
		err = collection.Update(bson.M{"_id": c.ID}, bson.M{"$set": bson.M{valField: c.Val}})
		if err != nil {
			log.Println(err)
			return
		}
	}
	log.Printf("counter - Persisted: key %s, Val %v\n", c.Key, c.Val)
}

// Incrementor has a incrementation definition
type Incrementor interface {
	Inc(actual interface{}, add interface{}) interface{}
	Val(collection *mgo.Collection, key string) (*Counter, error)
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
	mapDurationToPersist map[string]interface{}
	mux                  sync.Mutex
	incrementable        Incrementor
}
