package stats

import (
	"github.com/jpatel531/stickyd/stats/sets"
)

type Stats struct {
	Counters *Counters
	Gauges   *Gauges
	Sets     sets.Sets
}

func New(prefix string) *Stats {
	return &Stats{
		Counters: newCounters(prefix),
		Gauges:   newGauges(prefix),
		Sets:     sets.New(),
	}
}
