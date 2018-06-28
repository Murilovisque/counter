package counter

import (
	"fmt"
	"sync"
)

type Counter struct {
	mapCounts sync.Map
	mutex     sync.Mutex
}

func (sc *Counter) Inc(key string, val interface{}) {
	counterType := RetrieveCounterType(val)
	sc.mutex.Lock()
	defer sc.mutex.Unlock()
	valLoaded, found := sc.mapCounts.LoadOrStore(key, val)
	if found {
		sc.mapCounts.Store(key, counterType.Inc(valLoaded, val))
	}
}

func (sc *Counter) Val(key string) (interface{}, bool) {
	return sc.mapCounts.Load(key)
}

func (sc *Counter) Clear(key string) {
	sc.mapCounts.Delete(key)
}

func (sc *Counter) ValAndClear(key string) (interface{}, bool) {
	sc.mutex.Lock()
	defer sc.mutex.Unlock()
	v, ok := sc.Val(key)
	if ok {
		sc.Clear(key)
	}
	return v, ok
}

func (sc *Counter) Range(f func(k, v interface{}) bool) {
	sc.mapCounts.Range(f)
}

func (sc *Counter) String() string {
	return fmt.Sprint(&sc.mapCounts)
}
