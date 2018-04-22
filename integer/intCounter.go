package integer

import (
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

// Enable int counter
func Enable() counter.Incrementor {
	counter.AddIncrementor(integerCounter{})
}

//Inc duration of key
func Inc(key string, val int) {
	counter.Inc(counterType, key, val)
}

type integerCounter struct {
	ID bson.ObjectId `bson:"_id,omitempty"`
	K  string
	V  int
}

func (c integerCounter) Inc(actual interface{}, add interface{}) interface{} {
	vl1 := actual.(int)
	vl2 := add.(int)
	return vl1 + vl2
}

func (c integerCounter) ZeroVal() interface{} {
	return zero
}

func (c integerCounter) Val(collection *mgo.Collection, key string) (*counter.Counter, error) {
	err := collection.Find(bson.M{counter.KeyField: key}).One(&c)
	if err != nil {
		return nil, err
	}
	return &counter.Counter{ID: c.ID, Key: c.K, Val: c.V}, nil
}

func (c integerCounter) Type() string {
	return counterType
}
