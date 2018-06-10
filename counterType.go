package counter

import (
	"fmt"
	"time"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type counterType int

const (
	intType      counterType = iota
	durationType counterType = iota
)

func (c counterType) Inc(actual interface{}, add interface{}) interface{} {
	switch {
	case intType == c:
		return incInt(actual, add)
	case durationType == c:
		return incDuration(actual, add)
	default:
		panic("Not supported")
	}
}

func (c counterType) String() string {
	switch {
	case intType == c:
		return "int"
	case durationType == c:
		return "duration"
	default:
		return ""
	}
}

func (c counterType) CountMongo(collection *mgo.Collection, key string) (*countMongo, error) {
	switch {
	case intType == c:
		return integerCounter{}.CountMongo(collection, key)
	case durationType == c:
		return durationCounter{}.CountMongo(collection, key)
	default:
		return nil, nil
	}
}

func incDuration(actual interface{}, add interface{}) interface{} {
	vl1 := actual.(time.Duration)
	vl2 := add.(time.Duration)
	return vl1 + vl2
}

func incInt(actual interface{}, add interface{}) interface{} {
	vl1 := actual.(int)
	vl2 := add.(int)
	return vl1 + vl2
}

type integerCounter struct {
	ID  bson.ObjectId `bson:"_id,omitempty"`
	Key string
	Val int
}

func (c integerCounter) CountMongo(collection *mgo.Collection, key string) (*countMongo, error) {
	err := collection.Find(bson.M{keyField: key}).One(&c)
	if err != nil {
		return nil, err
	}
	return &countMongo{ID: c.ID, Key: c.Key, Val: c.Val}, nil
}

type durationCounter struct {
	ID  bson.ObjectId `bson:"_id,omitempty"`
	Key string
	Val time.Duration
}

func (c durationCounter) CountMongo(collection *mgo.Collection, key string) (*countMongo, error) {
	err := collection.Find(bson.M{keyField: key}).One(&c)
	if err != nil {
		return nil, err
	}
	return &countMongo{ID: c.ID, Key: c.Key, Val: c.Val}, nil
}

func retrieveCounterType(i interface{}) counterType {
	switch i.(type) {
	case int:
		return intType
	case time.Duration:
		return durationType
	default:
		panic(fmt.Sprintf("Type %T not supported", i))
	}
}
