package timer

import (
	"github.com/alex023/clock"
)

var Timer *clock.Clock

func init() {
	Timer = clock.NewClock()
}
