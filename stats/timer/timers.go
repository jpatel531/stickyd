package timer

import (
	"github.com/jpatel531/stickyd/stats/counter"
	"github.com/jpatel531/stickyd/util/collections"
)

type Timers interface {
	collections.FloatSliceMap
}

func NewTimers() Timers {
	return collections.NewFloatSliceMap()
}

type TimerCounters interface {
	counter.Counter
}

func NewTimerCounters() TimerCounters {
	return counter.New()
}
