package time

import (
	"time"

	"github.com/Murilovisque/counter"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	counterType = "time"
)

var (
	zero time.Duration
)

// Enable time counter
func Enable() {
	counter.AddIncrementor(timeCounter{})
}

//Inc duration of key
func Inc(key string, val time.Duration) {
	counter.Inc(counterType, key, val)
}

func (c timeCounter) Inc(actual interface{}, add interface{}) interface{} {
	vl1 := actual.(time.Duration)
	vl2 := add.(time.Duration)
	return vl1 + vl2
}

func (c timeCounter) ZeroVal() interface{} {
	return zero
}

func (c timeCounter) Val(collection *mgo.Collection, key string) (*counter.Counter, error) {
	err := collection.Find(bson.M{counter.KeyField: key}).One(&c)
	if err != nil {
		return nil, err
	}
	return &counter.Counter{ID: c.ID, Key: c.K, Val: c.V}, nil
}

func (c timeCounter) Type() string {
	return counterType
}

type timeCounter struct {
	ID bson.ObjectId `bson:"_id,omitempty"`
	K  string
	V  time.Duration
}
