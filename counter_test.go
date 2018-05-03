package counter

import (
	"log"
	"testing"
	"time"

	"github.com/Murilovisque/counter"
	cint "github.com/Murilovisque/counter/integer"
	ctime "github.com/Murilovisque/counter/time"
	mgo "gopkg.in/mgo.v2"
)

const (
	qtdeTest = 30
	dbTest   = "counter-test-db"
)

var (
	zeroDurationTest time.Duration
	zeroIntTest      int
)

func TestValNotStopped(t *testing.T) {
	dropDataBase(dbTest)
	sumTime := sumQtdeTimeToTest()
	sumInt := sumQtdeIntToTest()
	ctime.Enable()
	cint.Enable()
	counter.Start(dbTest, 10)
	for i := 0; i < qtdeTest; i++ {
		ctime.Inc("k1", time.Duration(i))
		ctime.Inc("k2", time.Duration(i))
		ctime.Inc("k3", time.Duration(i))
		cint.Inc("k1", i)
		cint.Inc("k2", i)
		cint.Inc("k3", i)
	}
	time.Sleep(90 * time.Millisecond) // Waiting to increment all values
	if sumTime != ctime.Val("k1") || sumTime != ctime.Val("k2") || sumTime != ctime.Val("k3") ||
		sumInt != cint.Val("k1") || sumInt != cint.Val("k2") || sumInt != cint.Val("k3") {
		t.FailNow()
	}
	log.Println(t.Name(), "1 ok. Should return correct values")
	counter.Stop()
}

func TestIncAndGetAndStop(t *testing.T) {
	dropDataBase(dbTest)
	sumTime := sumQtdeTimeToTest()
	sumInt := sumQtdeIntToTest()
	ctime.Enable()
	cint.Enable()
	counter.Start(dbTest, 10)
	for i := 0; i < qtdeTest; i++ {
		ctime.Inc("k1", time.Duration(i))
		ctime.Inc("k2", time.Duration(i))
		ctime.Inc("k3", time.Duration(i))
		cint.Inc("k1", i)
		cint.Inc("k2", i)
		cint.Inc("k3", i)
	}
	counter.Stop() // persist
	if sumTime != ctime.Val("k1") || sumTime != ctime.Val("k2") || sumTime != ctime.Val("k3") ||
		sumInt != cint.Val("k1") || sumInt != cint.Val("k2") || sumInt != cint.Val("k3") {
		t.FailNow()
	}
	log.Println(t.Name(), "1 ok. Should return correct values")

	counter.Start(dbTest, 10)
	if sumTime != ctime.Val("k1") || sumTime != ctime.Val("k2") || sumTime != ctime.Val("k3") ||
		sumInt != cint.Val("k1") || sumInt != cint.Val("k2") || sumInt != cint.Val("k3") {
		t.FailNow()
	}
	log.Println(t.Name(), "2 ok. Should continue works")

	for i := 0; i < qtdeTest; i++ {
		ctime.Inc("k1", time.Duration(i))
		cint.Inc("k1", i)
	}
	ctime.Inc("k2", sumTime)
	cint.Inc("k2", sumInt)
	time.Sleep(90 * time.Millisecond) // Waiting to increment all values
	if sumTime*2 != ctime.Val("k1") || sumTime*2 != ctime.Val("k2") || sumInt*2 != cint.Val("k1") || sumInt*2 != cint.Val("k2") {
		t.FailNow()
	}
	log.Println(t.Name(), "3 ok. Should return correct values after restart and not persist")

	counter.Stop()
	if sumTime*2 != ctime.Val("k1") || sumTime*2 != ctime.Val("k2") || sumTime != ctime.Val("k3") ||
		sumInt*2 != cint.Val("k1") || sumInt*2 != cint.Val("k2") || sumInt != cint.Val("k3") {
		t.FailNow()
	}
	log.Println(t.Name(), "4 ok. Should persist and return correct values after restart and persist")
}

func TestIncAndClear(t *testing.T) {
	dropDataBase(dbTest)
	ctime.Enable()
	cint.Enable()
	counter.Start(dbTest, 10)

	ctime.Inc("k1", sumQtdeTimeToTest())
	cint.Inc("k1", sumQtdeIntToTest())
	time.Sleep(60 * time.Millisecond)
	if sumQtdeTimeToTest() != ctime.Val("k1") || sumQtdeIntToTest() != cint.Val("k1") {
		t.FailNow()
	}
	log.Println(t.Name(), "1 ok. Increment should works")

	ctime.Clear("k1")
	cint.Clear("k1")
	if cint.Val("k1") != zeroIntTest || ctime.Val("k1") != zeroDurationTest {
		t.FailNow()
	}
	log.Println(t.Name(), "2 ok. Clear and should return correct values")

	counter.Stop()
	if cint.Val("k1") != zeroIntTest || ctime.Val("k1") != zeroDurationTest {
		t.FailNow()
	}
	log.Println(t.Name(), "3 ok. Clear, persist and should return correct values")

	counter.Start(dbTest, 10)
	if cint.Val("k1") != zeroIntTest || ctime.Val("k1") != zeroDurationTest {
		t.FailNow()
	}
	log.Println(t.Name(), "4 ok. Clear, persist and restart should return correct values")
	counter.Stop()
}

func TestRestartShouldClearAll(t *testing.T) {
	dropDataBase(dbTest, "other-db-test")
	ctime.Enable()
	cint.Enable()
	counter.Start(dbTest, 10)

	ctime.Inc("k1", sumQtdeTimeToTest())
	cint.Inc("k1", sumQtdeIntToTest())
	counter.Stop()
	if sumQtdeTimeToTest() != ctime.Val("k1") || sumQtdeIntToTest() != cint.Val("k1") {
		t.FailNow()
	}
	log.Println(t.Name(), "1 ok. Increment and stop should works")

	counter.Start("other-db-test", 10)
	if ctime.Val("k1") != zeroDurationTest || cint.Val("k1") != zeroIntTest {
		t.FailNow()
	}
	log.Println(t.Name(), "2 ok. Increment and restart with other db should zero the values")

	ctime.Inc("k1", sumQtdeTimeToTest())
	cint.Inc("k1", sumQtdeIntToTest())
	counter.Stop()
	if sumQtdeTimeToTest() != ctime.Val("k1") || sumQtdeIntToTest() != cint.Val("k1") {
		t.FailNow()
	}
	log.Println(t.Name(), "3 ok. Increment and restart and stop with other db should zero the values")

}

func sumQtdeTimeToTest() time.Duration {
	return time.Duration(sumQtdeIntToTest())
}

func sumQtdeIntToTest() int {
	var sum int
	for i := 0; i < qtdeTest; i++ {
		sum += i
	}
	return sum
}

func dropDataBase(dbs ...string) {
	session, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	for _, db := range dbs {
		session.DB(db).DropDatabase()
	}
}
