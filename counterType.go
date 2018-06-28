package counter

import (
	"fmt"
	"time"
)

type CounterType int

const (
	IntCounterType      CounterType = iota
	DurationCounterType CounterType = iota
)

func (c CounterType) Inc(actual interface{}, add interface{}) interface{} {
	switch {
	case IntCounterType == c:
		return incInt(actual, add)
	case DurationCounterType == c:
		return incDuration(actual, add)
	default:
		panic("Not supported")
	}
}

func (c CounterType) String() string {
	switch {
	case IntCounterType == c:
		return "int"
	case DurationCounterType == c:
		return "time.Duration"
	default:
		return ""
	}
}

func RetrieveCounterType(i interface{}) CounterType {
	switch i.(type) {
	case int:
		return IntCounterType
	case time.Duration:
		return DurationCounterType
	default:
		panic(fmt.Sprintf("Type %T not supported", i))
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
