package integer

import (
	"time"

	"github.com/Murilovisque/counter"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	counterType = "integer"
)

var (
	zero int
)

func Init() counter.Incrementable {
	return integerCounter{}
}

//Inc duration of key
func Inc(key string, val time.Duration) {
	counter.Inc(counterType, key, val)
}

type integerCounter struct {
	ID  bson.ObjectId `bson:"_id,omitempty"`
	Key string
	Val int
}

func (c integerCounter) Inc(actual interface{}, add interface{}) interface{} {
	vl1 := actual.(int)
	vl2 := add.(int)
	return vl1 + vl2
}

func (c integerCounter) GetZeroVal() interface{} {
	return zero
}

func (c integerCounter) GetVal(collection *mgo.Collection, key string) (*counter.Counter, error) {
	err := collection.Find(bson.M{counter.KeyField: key}).One(&c)
	if err != nil {
		return nil, err
	}
	return &counter.Counter{ID: c.ID, Key: c.Key, Val: c.Val}, nil
}

func (c integerCounter) GetType() string {
	return counterType
}
