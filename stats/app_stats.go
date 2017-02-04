package stats

import (
	"github.com/jpatel531/stickyd/stats/sets"
)

type AppStats struct {
	Counters *Counters
	Gauges   *Gauges
	Sets     sets.Sets
}

func NewAppStats(prefix string) *AppStats {
	return &AppStats{
		Counters: newCounters(prefix),
		Gauges:   newGauges(prefix),
		Sets:     sets.New(),
	}
}
