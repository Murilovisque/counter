package time

import (
	"log"
	"time"

	"github.com/Murilovisque/counter/repository"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"sync"
)

const (
	counterCollection = "counters"
	keyField          = "key"
	durationField     = "duration"
)

var (
	db               = ""
	persistInterval  = 10
	persistenceQueue persistanceModel
)

// Init inform the dbname and internal to persist
func Init(dbParam string, persistIntervalParam int) {
	db = dbParam
	persistInterval = persistIntervalParam
	persistenceQueue.mapDurationToPersist = make(map[string]time.Duration)
	doPersistenceImpl()
}

//Inc duration of key
func Inc(key string, duration time.Duration) {
	persistenceQueue.inc(key, duration)
}

type persistanceModel struct {
	mapDurationToPersist map[string]time.Duration
	mux                  sync.Mutex
}

func (p *persistanceModel) inc(key string, duration time.Duration) {
	p.mux.Lock()
	val, ok := p.mapDurationToPersist[key]
	if ok {
		val += duration
		p.mapDurationToPersist[key] = val
	} else {
		p.mapDurationToPersist[key] = duration
	}
	p.mux.Unlock()
}

func (p *persistanceModel) getAndClear(key string) time.Duration {
	p.mux.Lock()
	defer p.mux.Unlock()
	var duration time.Duration
	duration, v := p.mapDurationToPersist[key]
	if !v {
		return duration
	}
	delete(p.mapDurationToPersist, key)
	return duration
}

type runPersistance struct{}

func (r *runPersistance) Do() {
	for k := range persistenceQueue.mapDurationToPersist {
		log.Printf("timeCounter Persisting key: %s", k)
		repository.PersistIncrementation(&counter{Key: k, Duration: persistenceQueue.getAndClear(k)})
	}
}

func doPersistenceImpl() {
	repository.RunPersistance(persistInterval, &runPersistance{})
}

type counter struct {
	ID       bson.ObjectId `bson:"_id,omitempty"`
	Key      string
	Duration time.Duration
}

func (p *counter) Persist(session *mgo.Session) {
	collection := session.DB(db).C(counterCollection)
	c := counter{}
	err := collection.Find(bson.M{keyField: p.Key}).One(&c)
	if err != nil && err != mgo.ErrNotFound {
		log.Println(err)
		return
	}
	if err == mgo.ErrNotFound {
		c.Key = p.Key
		c.Duration = p.Duration
		err = collection.Insert(&c)
		if err != nil {
			log.Println(err)
			return
		}
	} else {
		c.Duration = p.Duration + c.Duration
		err = collection.Update(bson.M{"_id": c.ID}, bson.M{"$set": bson.M{durationField: c.Duration}})
		if err != nil {
			log.Println(err)
			return
		}
	}
	log.Printf("timeCounter - Persisted: key %s, Duration %v\n", c.Key, c.Duration)
}
