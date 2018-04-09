package counter

import (
	"log"
	"time"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"sync"
)

const (
	counterCollection = "counters"
	KeyField          = "key"
	valField          = "val"
)

var (
	db               string
	persistInterval  = 10
	persistenceQueue persistanceModel
	incrementable    Incrementable
)

// Init inform the dbname and internal to persist
func Init(dbParam string, persistIntervalParam int, i Incrementable) {
	db = dbParam
	incrementable = i
	persistInterval = persistIntervalParam
	persistenceQueue.mapDurationToPersist = make(map[string]interface{})
	startPersistence()
}

//Inc duration of key
func Inc(key string, val interface{}) {
	persistenceQueue.inc(key, val)
}

func (p *persistanceModel) inc(key string, val interface{}) {
	p.mux.Lock()
	v, ok := p.mapDurationToPersist[key]
	if ok {
		v = incrementable.Inc(v, val)
		p.mapDurationToPersist[key] = v
	} else {
		p.mapDurationToPersist[key] = val
	}
	p.mux.Unlock()
}

func (p *persistanceModel) getAndClear(key string) interface{} {
	p.mux.Lock()
	defer p.mux.Unlock()
	val, v := p.mapDurationToPersist[key]
	if !v {
		return incrementable.GetZeroVal()
	}
	delete(p.mapDurationToPersist, key)
	return val
}

func startPersistence() {
	log.Printf("genericCounter - Starting persistance each %d second(s)\n", persistInterval)
	go func() {
		ticker := time.NewTicker(time.Duration(persistInterval) * time.Second)
		for range ticker.C {
			session, err := mgo.Dial("localhost")
			if err != nil {

				log.Println(err)
				continue
			}
			session.SetMode(mgo.Monotonic, true)
			for k := range persistenceQueue.mapDurationToPersist {
				log.Printf("counter - Persisting key: %s", k)
				count := &Counter{Key: k, Val: persistenceQueue.getAndClear(k)}
				count.Persist(session)
			}
			session.Close()
		}
	}()
}

// Persist make the persistance
func (p *Counter) Persist(session *mgo.Session) {
	collection := session.DB(db).C(counterCollection)
	c, err := incrementable.GetVal(collection, p.Key)
	if err == mgo.ErrNotFound {
		c = &Counter{Key: p.Key, Val: p.Val}
		err = collection.Insert(&c)
		if err != nil {
			log.Println(err)
			return
		}
	} else {
		c.Val = incrementable.Inc(p.Val, c.Val)
		err = collection.Update(bson.M{"_id": c.ID}, bson.M{"$set": bson.M{valField: c.Val}})
		if err != nil {
			log.Println(err)
			return
		}
	}
	log.Printf("genericCounter - Persisted: key %s, Val %v\n", c.Key, c.Val)
}

// Incrementable has a incrementation definition
type Incrementable interface {
	Inc(actual interface{}, add interface{}) interface{}
	GetVal(collection *mgo.Collection, key string) (*Counter, error)
	GetZeroVal() interface{}
}

// Counter model
type Counter struct {
	ID  bson.ObjectId `bson:"_id,omitempty"`
	Key string
	Val interface{}
}

type persistanceModel struct {
	mapDurationToPersist map[string]interface{}
	mux                  sync.Mutex
}
