package stats

import (
	"github.com/jpatel531/stickyd/stats/gauges"
)

type Gauges struct {
	gauge                 gauges.Gauges
	timestampLagNamespace string
}

func newGauges(statsPrefix string) *Gauges {
	return &Gauges{
		gauge: gauges.New(),
		timestampLagNamespace: statsPrefix + ".timestamp_lag",
	}
}

func (g *Gauges) SetTimestampLag(n float64) {
	g.gauge.Set(g.timestampLagNamespace, n)
}

func (g *Gauges) Set(key string, n float64) {
	g.gauge.Set(key, n)
}

func (g *Gauges) String() string {
	return g.gauge.String()
}

func (g *Gauges) MarshalJSON() ([]byte, error) {
	return g.gauge.MarshalJSON()
}
