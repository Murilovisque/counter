package counter_test

import (
	"testing"
	"time"

	"github.com/Murilovisque/counter"
	"github.com/stretchr/testify/assert"
)

const (
	qtdeSimpleTest = 100000
)

func TestIncAndValShouldWorks(t *testing.T) {
	cd := counter.NewDurationCounter()
	ci := counter.NewIntCounter()
	for i := 0; i < qtdeSimpleTest; i++ {
		cd.Inc("k1d", time.Duration(i))
		cd.Inc("k2d", time.Duration(i))
		cd.Inc("k3d", time.Duration(i))
		ci.Inc("k1i", i)
		ci.Inc("k2i", i)
		ci.Inc("k3i", i)
	}
	sumTime := sumQtdeTimeToTest()
	sumInt := sumQtdeIntToTest()

	passIfAreEqualsDurationWhenUseVal(t, sumTime, cd, "k1d", "k2d", "k3d")
	passIfAreEqualsIntWhenUseVal(t, sumInt, ci, "k1i", "k2i", "k3i")
}

func TestIncAndValAndClear(t *testing.T) {
	cd := counter.NewDurationCounter()
	ci := counter.NewIntCounter()
	for i := 0; i < qtdeSimpleTest; i++ {
		cd.Inc("k1d", time.Duration(i))
		cd.Inc("k2d", time.Duration(i))
		ci.Inc("k1i", i)
		ci.Inc("k2i", i)
	}
	sumTime := sumQtdeTimeToTest()
	sumInt := sumQtdeIntToTest()

	passIfAreEqualsDurationWhenUseVal(t, sumTime, cd, "k1d", "k2d")
	passIfAreEqualsIntWhenUseVal(t, sumInt, ci, "k1i", "k2i")
	cd.Clear("k1d")
	ci.Clear("k1i")
	passIfAreEqualsDurationWhenUseVal(t, sumTime, cd, "k2d")
	passIfAreEqualsIntWhenUseVal(t, sumInt, ci, "k2i")
	passIfDurationAreZero(t, cd, "k1d")
	passIfIntAreZero(t, ci, "k1i")
}

func TestIncAndValAndClearAndInc(t *testing.T) {
	cd := counter.NewDurationCounter()
	ci := counter.NewIntCounter()
	for i := 0; i < qtdeSimpleTest; i++ {
		cd.Inc("k1d", time.Duration(i))
		cd.Inc("k2d", time.Duration(i))
		ci.Inc("k1i", i)
		ci.Inc("k2i", i)
	}
	sumTime := sumQtdeTimeToTest()
	sumInt := sumQtdeIntToTest()

	passIfIntsAreEqualsWhenUseValAndClear(t, sumInt, ci, "k1i", "k2i")
	passIfDurationsAreEqualsWhenUseValAndClear(t, sumTime, cd, "k1d", "k2d")
	passIfIntAreZero(t, ci, "k1i", "k2i", "k1d", "k2d")
	passIfDurationAreZero(t, cd, "k1d", "k2d")
	cd.Inc("k1d", sumTime)
	ci.Inc("k1i", sumInt)

	passIfAreEqualsIntWhenUseVal(t, sumInt, ci, "k1i")
	passIfAreEqualsDurationWhenUseVal(t, sumTime, cd, "k1d")
	passIfDurationAreZero(t, cd, "k2d")
	passIfIntAreZero(t, ci, "k2i")
}

func passIfIntsAreEqualsWhenUseValAndClear(t *testing.T, assertVal int, c *counter.IntCounter, keys ...string) {
	if t.Failed() {
		return
	}
	for _, k := range keys {
		val, ok := c.ValAndClear(k)
		assert.Falsef(t, !ok || assertVal != val, "Test %s failed, value of key '%s' should be %v, but it is %v - %v\n", t.Name(), k, assertVal, val, c)
	}
}

func passIfDurationsAreEqualsWhenUseValAndClear(t *testing.T, assertVal time.Duration, c *counter.DurationCounter, keys ...string) {
	if t.Failed() {
		return
	}
	for _, k := range keys {
		val, ok := c.ValAndClear(k)
		assert.Falsef(t, !ok || assertVal != val, "Test %s failed, value of key '%s' should be %v, but it is %v - %v\n", t.Name(), k, assertVal, val, c)
	}
}

func passIfDurationAreZero(t *testing.T, c *counter.DurationCounter, keys ...string) {
	if t.Failed() {
		return
	}
	for _, k := range keys {
		_, ok := c.Val(k)
		assert.Falsef(t, ok, "Test %s failed, value of key '%s' should be zero\n", t.Name(), k)
	}
}

func passIfIntAreZero(t *testing.T, c *counter.IntCounter, keys ...string) {
	if t.Failed() {
		return
	}
	for _, k := range keys {
		_, ok := c.Val(k)
		assert.Falsef(t, ok, "Test %s failed, value of key '%s' should be zero\n", t.Name(), k)
	}
}

func passIfAreEqualsDurationWhenUseVal(t *testing.T, assertVal time.Duration, c *counter.DurationCounter, keys ...string) {
	if t.Failed() {
		return
	}
	for _, k := range keys {
		v, ok := c.Val(k)
		assert.Falsef(t, !ok || v != assertVal, "Test %s failed, value of key '%s' should be %v, but it is %v\n", t.Name(), k, assertVal, v)
	}
}

func passIfAreEqualsIntWhenUseVal(t *testing.T, assertVal int, c *counter.IntCounter, keys ...string) {
	if t.Failed() {
		return
	}
	for _, k := range keys {
		v, ok := c.Val(k)
		assert.Falsef(t, !ok || v != assertVal, "Test %s failed, value of key '%s' should be %v, but it is %v\n", t.Name(), k, assertVal, v)
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
