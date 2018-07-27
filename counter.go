package counter

import (
	"fmt"
	"sync"
)

type counter struct {
	mapCounts sync.Map
	mutex     sync.Mutex
	incrementor
}

type incrementor interface {
	increment(vl1, vl2 interface{}) interface{}
}

func (sc *counter) Clear(key string) {
	sc.mapCounts.Delete(key)
}

func (sc *counter) inc(key string, val interface{}) {
	sc.mutex.Lock()
	defer sc.mutex.Unlock()
	valLoaded, found := sc.mapCounts.LoadOrStore(key, val)
	if found {
		sc.mapCounts.Store(key, sc.increment(valLoaded, val))
	}
}

func (sc *counter) val(key string) (interface{}, bool) {
	return sc.mapCounts.Load(key)
}

func (sc *counter) valAndClear(key string) (interface{}, bool) {
	sc.mutex.Lock()
	defer sc.mutex.Unlock()
	v, ok := sc.val(key)
	if ok {
		sc.Clear(key)
	}
	return v, ok
}

func (sc *counter) rangeVals(f func(k, v interface{}) bool) {
	sc.mapCounts.Range(f)
}

func (sc *counter) String() string {
	return fmt.Sprint(&sc.mapCounts)
}
