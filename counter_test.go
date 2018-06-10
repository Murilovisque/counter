package counter_test

import (
	"log"
	"testing"
	"time"

	"github.com/Murilovisque/counter"
)

const (
	qtdeSimpleTest = 10000
)

func TestIncAndValShouldWorks(t *testing.T) {
	c := counter.Counter{}
	for i := 0; i < qtdeSimpleTest; i++ {
		c.Inc("k1d", time.Duration(i))
		c.Inc("k2d", time.Duration(i))
		c.Inc("k3d", time.Duration(i))
		c.Inc("k1i", i)
		c.Inc("k2i", i)
		c.Inc("k3i", i)
	}
	sumTime := sumQtdeTimeToTest()
	sumInt := sumQtdeIntToTest()
	c.WaitForFinalizationOfIncrements()
	passIfAreEqualsDuration(t, sumTime, &c, "k1d", "k2d", "k3d")
	passIfAreEqualsInt(t, sumInt, &c, "k1i", "k2i", "k3i")
}

func TestSimpleIncAndValAndClear(t *testing.T) {
	c := counter.Counter{}
	for i := 0; i < qtdeSimpleTest; i++ {
		c.Inc("k1d", time.Duration(i))
		c.Inc("k2d", time.Duration(i))
		c.Inc("k1i", i)
		c.Inc("k2i", i)
	}
	sumTime := sumQtdeTimeToTest()
	sumInt := sumQtdeIntToTest()
	c.WaitForFinalizationOfIncrements()
	passIfAreEqualsDuration(t, sumTime, &c, "k1d", "k2d")
	passIfAreEqualsInt(t, sumInt, &c, "k1i", "k2i")
	c.Clear("k1d")
	c.Clear("k1i")
	passIfAreEqualsDuration(t, sumTime, &c, "k2d")
	passIfAreEqualsInt(t, sumInt, &c, "k2i")
	passIfAreZero(t, &c, "k1d", "k1i")
}

func passIfAreZero(t *testing.T, c *counter.Counter, keys ...string) {
	for _, k := range keys {
		if _, ok := c.Val(k); ok {
			log.Printf("Test %s failed, value of key '%s' should be zero\n", t.Name(), k)
			t.FailNow()
			break
		}
	}
}

func passIfAreEqualsDuration(t *testing.T, assertVal time.Duration, c *counter.Counter, keys ...string) {
	comp := func(a, b interface{}) bool {
		v1, ok := a.(time.Duration)
		if !ok {
			return false
		}
		v2, ok := a.(time.Duration)
		return ok && v1 == v2
	}
	passIfAreEquals(comp, t, assertVal, c, keys...)
}

func passIfAreEqualsInt(t *testing.T, assertVal int, c *counter.Counter, keys ...string) {
	comp := func(a, b interface{}) bool {
		v1, ok := a.(int)
		if !ok {
			return false
		}
		v2, ok := a.(int)
		return ok && v1 == v2
	}
	passIfAreEquals(comp, t, assertVal, c, keys...)
}

func passIfAreEquals(comparator func(interface{}, interface{}) bool, t *testing.T, assertVal interface{}, c *counter.Counter, keys ...string) {
	for _, k := range keys {
		v, ok := c.Val(k)
		if !ok || !comparator(v, assertVal) {
			log.Printf("Test %s failed, value of key '%s' should be %v, but it is %v\n", t.Name(), k, assertVal, v)
			t.FailNow()
			break
		}
	}
}

func sumQtdeTimeToTest() time.Duration {
	return time.Duration(sumQtdeIntToTest())
}

func sumQtdeIntToTest() int {
	var sum int
	for i := 0; i < qtdeSimpleTest; i++ {
		sum += i
	}
	return sum
}
