package counter

import (
	"encoding/json"
	"sync"

	"sync/atomic"
)

type Counter struct {
	goroutineCounter int64
	mapCounts        sync.Map
	mutex            sync.Mutex
}

func (sc *Counter) Inc(key string, val interface{}) {
	go func() {
		atomic.AddInt64(&sc.goroutineCounter, 1)
		defer atomic.AddInt64(&sc.goroutineCounter, -1)
		counterType := retrieveCounterType(val)
		sc.mutex.Lock()
		defer sc.mutex.Unlock()
		valLoaded, found := sc.mapCounts.LoadOrStore(key, val)
		if found {
			sc.mapCounts.Store(key, counterType.Inc(valLoaded, val))
		}
	}()
}

func (sc *Counter) Val(key string) (interface{}, bool) {
	return sc.mapCounts.Load(key)
}

func (sc *Counter) Clear(key string) {
	sc.mapCounts.Delete(key)
}

func (sc *Counter) WaitForFinalizationOfIncrements() {
	for {
		if sc.goroutineCounter < 1 {
			break
		}
	}
}

func (sc *Counter) String() string {
	b, err := json.Marshal(sc.mapCounts)
	if err != nil {
		return ""
	}
	return string(b)
}
