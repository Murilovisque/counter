package time

import (
	"testing"
	"time"

	"github.com/Murilovisque/counter"
	ctime "github.com/Murilovisque/counter/time"
)

func TestInc(*testing.T) {
	ctime.Enable()
	counter.Start("counter-test-db", 10)
	ctime.Inc("k1", time.Duration(5))
	ctime.Inc("k1", time.Duration(3))
	ctime.Inc("k1", time.Duration(13))
	counter.Stop()
	//implement get
}
