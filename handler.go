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
	"time"
)

type handler struct {
	appStats     *stats.AppStats
	processStats *stats.ProcessStats
	config       *config.Config
}

func (h handler) HandleMessage(msg []byte, rinfo *frontend.RemoteInfo) {
	log.Printf("Received %s from %+v\n", strings.TrimSpace(string(msg)), rinfo)
	h.appStats.Counters.IncrPacketsReceived()

	metricsStrings := strings.Split(strings.TrimSpace(string(msg)), "\n")

	for _, metricString := range metricsStrings {
		h.appStats.Counters.IncrMetricsReceived()

		metric, err := metrics.Parse([]byte(metricString))
		if err != nil {
			log.Printf("error parsing metric %q. Error received: %v\n", metricString, err)
			h.processStats.Messages.IncrBadLinesSeen()
			h.appStats.Counters.IncrBadLinesSeen()
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
			h.appStats.Sets.Insert(metric.Key, metric.Value)
		case "g", "c":
			value, err := strconv.ParseFloat(metric.Value, 64)
			if err != nil {
				log.Printf("Expected float value, received %q", value)
				h.processStats.Messages.IncrBadLinesSeen()
				h.appStats.Counters.IncrBadLinesSeen()
				continue
			}

			if metric.Type == "g" {
				// TODO allow +- increments
				h.appStats.Gauges.Set(metric.Key, float64(value))
			} else {
				h.appStats.Counters.Incr(metric.Key, float64(value*1/sampleRate))
			}
		default: // c
			log.Printf("Unsupported type %q", metric.Type)
		}
	}

	h.processStats.Messages.SetLastMessageSeen(time.Now().Unix())
}
