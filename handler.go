package main

import (
	"github.com/jpatel531/stickyd/config"
	"github.com/jpatel531/stickyd/frontend"
	"github.com/jpatel531/stickyd/metrics"
	"github.com/jpatel531/stickyd/stats"
	"github.com/jpatel531/stickyd/util"
	"log"
	"strconv"
	"strings"
)

type handler struct {
	stats  *stats.Stats
	config *config.Config
}

func (h handler) HandleMessage(msg []byte, rinfo *frontend.RemoteInfo) {
	log.Printf("Received %s from %+v\n", strings.TrimSpace(string(msg)), rinfo)
	h.stats.Counters.IncrPacketsReceived()

	metricsStrings := strings.Split(strings.TrimSpace(string(msg)), "\n")

	for _, metricString := range metricsStrings {
		h.stats.Counters.IncrMetricsReceived()

		metric, err := metrics.Parse([]byte(metricString))
		if err != nil {
			log.Printf("error parsing metric %q. Error received: %v\n", metricString, err)
			// 	// TODO stats.messages.bad_lines_seen
			h.stats.Counters.IncrBadLinesSeen()
			continue
		}
		metric.Key = util.SanitizeKey(metric.Key)

		if h.config.DumpMessages {
			log.Printf("metric received: %+v\n", metric)
		}

		// TODO add to key counter

		var sampleRate float64
		if metric.SampleRate == nil {
			sampleRate = 1
		} else {
			sampleRate = *metric.SampleRate
			// TODO handle > 1 sample rates
		}
		switch metric.Type {
		// TODO add more
		case "s":
			h.stats.Sets.Insert(metric.Key, metric.Value)
		case "g", "c":
			value, err := strconv.ParseFloat(metric.Value, 64)
			if err != nil {
				log.Printf("Expected float value, received %q", value)
				// 	// TODO stats.messages.bad_lines_seen
				h.stats.Counters.IncrBadLinesSeen()
				continue
			}

			if metric.Type == "g" {
				// TODO allow +- increments
				h.stats.Gauges.Set(metric.Key, float64(value))
			} else {
				h.stats.Counters.Incr(metric.Key, float64(value*1/sampleRate))
			}
		default: // c
			log.Printf("Unsupported type %q", metric.Type)
		}

		log.Printf("counters: %s\n", h.stats.Counters.String())
		log.Printf("gauges: %s\n", h.stats.Gauges.String())
		log.Printf("sets: %s\n", h.stats.Sets.String())
	}

}
