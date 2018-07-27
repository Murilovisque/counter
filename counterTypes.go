package counter

import "time"

var (
	zeroDurationVal time.Duration
	zeroIntVal      int
)

type DurationCounter struct {
	counter
}

func NewDurationCounter() *DurationCounter {
	d := DurationCounter{}
	d.incrementor = &d
	return &d
}

func (d *DurationCounter) Inc(key string, val time.Duration) {
	d.inc(key, val)
}

func (d *DurationCounter) Val(key string) (time.Duration, bool) {
	v, ok := d.val(key)
	if ok {
		return v.(time.Duration), ok
	}
	return zeroDurationVal, ok
}

func (sc *DurationCounter) Range(f func(k string, v time.Duration) bool) {
	sc.mapCounts.Range(func(k, v interface{}) bool {
		return f(k.(string), v.(time.Duration))
	})
}

func (d *DurationCounter) ValAndClear(key string) (time.Duration, bool) {
	v, ok := d.valAndClear(key)
	if ok {
		return v.(time.Duration), ok
	}
	return zeroDurationVal, ok
}

func (d *DurationCounter) increment(vl1, vl2 interface{}) interface{} {
	return vl1.(time.Duration) + vl2.(time.Duration)
}

type IntCounter struct {
	counter
}

func NewIntCounter() *IntCounter {
	i := IntCounter{}
	i.incrementor = &i
	return &i
}

func (i *IntCounter) Inc(key string, val int) {
	i.inc(key, val)
}

func (i *IntCounter) Val(key string) (int, bool) {
	v, ok := i.val(key)
	if ok {
		return v.(int), ok
	}
	return zeroIntVal, ok
}

func (i *IntCounter) ValAndClear(key string) (int, bool) {
	v, ok := i.valAndClear(key)
	if ok {
		return v.(int), ok
	}
	return zeroIntVal, ok
}

func (sc *IntCounter) Range(f func(k string, v int) bool) {
	sc.mapCounts.Range(func(k, v interface{}) bool {
		return f(k.(string), v.(int))
	})
}

func (d *IntCounter) increment(vl1, vl2 interface{}) interface{} {
	return vl1.(int) + vl2.(int)
}
