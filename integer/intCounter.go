package integer

import (
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

//Inc int of key
func Inc(key string, inc int) {
	persistenceQueue.inc(key, inc)
}

type persistanceModel struct {
	mapDurationToPersist map[string]int
	mux                  sync.Mutex
}

func (p *persistanceModel) inc(key string, inc int) {
	p.mux.Lock()
	val, ok := p.mapDurationToPersist[key]
	if ok {
		val += inc
		p.mapDurationToPersist[key] = val
	} else {
		p.mapDurationToPersist[key] = inc
	}
	p.mux.Unlock()
}
