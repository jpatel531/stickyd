package main

import (
	"github.com/jpatel531/stickyd/config"
	"github.com/jpatel531/stickyd/metrics"
	"github.com/jpatel531/stickyd/stats"
	"github.com/jpatel531/stickyd/stats/counter"
	"github.com/jpatel531/stickyd/util"
	"log"
	"net"
	"strconv"
	"strings"
	"time"
)

type handler struct {
	appStats     *stats.AppStats
	processStats *stats.ProcessStats
	config       *config.Config
	keyCounter   counter.Counter
}

func (h handler) HandleMessage(msg []byte, addr net.Addr) {
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
		key := metric.Key

		if h.config.DumpMessages {
			log.Printf("metric received: %+v\n", metric)
		}

		if h.config.KeyFlush.Interval > 0 {
			h.keyCounter.Incr(key, 1)
		}

		var sampleRate float64
		if metric.SampleRate == nil {
			sampleRate = 1
		} else {
			sampleRate = *metric.SampleRate
			// TODO handle > 1 sample rates
		}
		switch metric.Type {
		case "s":
			h.appStats.Sets.Insert(key, metric.Value)
		case "g", "c", "ms":
			value, err := strconv.ParseFloat(metric.Value, 64)
			if err != nil {
				log.Printf("Expected float value, received %q", value)
				h.processStats.Messages.IncrBadLinesSeen()
				h.appStats.Counters.IncrBadLinesSeen()
				continue
			}

			switch metric.Type {
			case "g":
				// TODO allow +- increments
				h.appStats.Gauges.Set(key, float64(value))
			case "c":
				h.appStats.Counters.Incr(key, float64(value*1/sampleRate))
			case "ms":
				h.appStats.Timers.Append(key, value)
				h.appStats.TimerCounters.Incr(key, float64(value*1/sampleRate))
			}

		default:
			log.Printf("Unsupported type %q\n", metric.Type)
		}
	}

	h.processStats.Messages.SetLastMessageSeen(time.Now().Unix())
}
