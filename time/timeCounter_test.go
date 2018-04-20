package time

import (
	"testing"
	"time"

	"github.com/Murilovisque/counter"
	ctime "github.com/Murilovisque/counter/time"
)

func TestInc(*testing.T) {
	counter.Init("counter-test-db", 10, ctime.Init())
	ctime.Inc("k1", time.Duration(5))
	ctime.Inc("k1", time.Duration(3))
	ctime.Inc("k1", time.Duration(13))
	counter.Stop()
}
