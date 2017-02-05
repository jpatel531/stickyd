package stats

import (
	"github.com/jpatel531/stickyd/stats/sets"
)

type AppStats struct {
	Counters *Counters
	Gauges   *Gauges
	Sets     sets.Sets

	// not yet implemented
	// Timers
	// TimerCounters
	// CounterRates
	// TimerData
	// PctThreshold
	// Histogram
}

func (a *AppStats) Clear() {
	a.Counters.Clear()
	a.Gauges.Clear()
	a.Sets.Clear()
}

func NewAppStats(prefix string) *AppStats {
	return &AppStats{
		Counters: newCounters(prefix),
		Gauges:   newGauges(prefix),
		Sets:     sets.New(),
	}
}
