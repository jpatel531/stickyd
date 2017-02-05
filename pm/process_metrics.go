package pm

import (
	"github.com/jpatel531/stickyd/stats"
	"time"
)

func ProcessMetrics(appStats *stats.AppStats, flushInterval int) map[string]interface{} {
	startTime := time.Now().Unix()

	counters := appStats.Counters.Map()
	counterRates := map[string]float64{}

	for key, counter := range counters {
		counterRates[key] = counter / float64(flushInterval/1000)
	}

	stickydMetrics := map[string]int64{
		"processingTime": startTime - time.Now().Unix(),
	}

	return map[string]interface{}{
		"stickyd_metrics": stickydMetrics,
		"counters":        counters,
		"counter_rates":   counterRates,
		"sets":            appStats.Sets.Map(),
		"gauges":          appStats.Gauges.Map(),
	}
}
