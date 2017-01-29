package stats

import (
	"github.com/jpatel531/stickyd/stats/counter"
)

type Counters struct {
	counter counter.Counter

	keyBadLinesSeen    string
	keyPacketsReceived string
	keyMetricsReceived string
}

func newCounters(statsPrefix string) *Counters {
	return &Counters{
		counter:            counter.New(),
		keyBadLinesSeen:    statsPrefix + ".bad_lines_seen",
		keyPacketsReceived: statsPrefix + ".packets_received",
		keyMetricsReceived: statsPrefix + ".metrics_received",
	}
}

func (c *Counters) IncrBadLinesSeen() {
	c.counter.Incr(c.keyBadLinesSeen, 1)
}

func (c *Counters) IncrPacketsReceived() {
	c.counter.Incr(c.keyPacketsReceived, 1)
}

func (c *Counters) IncrMetricsReceived() {
	c.counter.Incr(c.keyMetricsReceived, 1)
}

func (c *Counters) Incr(key string, n float64) {
	c.counter.Incr(key, n)
}

func (c *Counters) String() string {
	return c.counter.String()
}
