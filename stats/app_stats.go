package stats

import (
	"github.com/jpatel531/stickyd/stats/sets"
	"github.com/jpatel531/stickyd/stats/timer"
)

type AppStats struct {
	Counters      *Counters
	Gauges        *Gauges
	Sets          sets.Sets
	Timers        timer.Timers
	TimerCounters timer.TimerCounters

	// not yet implemented
	// CounterRates
	// TimerData
	// PctThreshold
	// Histogram
}

func (a *AppStats) Clear() {
	a.Counters.Clear()
	a.Gauges.Clear()
	a.Sets.Clear()
	a.Timers.Clear()
	a.TimerCounters.Clear()
}

func NewAppStats(prefix string) *AppStats {
	return &AppStats{
		Counters:      newCounters(prefix),
		Gauges:        newGauges(prefix),
		Sets:          sets.New(),
		Timers:        timer.NewTimers(),
		TimerCounters: timer.NewTimerCounters(),
	}
}
