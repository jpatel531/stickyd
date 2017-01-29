package main

import (
	"github.com/jpatel531/stickyd/config"
	"github.com/jpatel531/stickyd/frontend"
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

	metrics := strings.Split(strings.TrimSpace(string(msg)), "\n")

	for _, metric := range metrics {
		h.stats.Counters.IncrMetricsReceived()

		if h.config.DumpMessages {
			log.Printf("metric received: %+v\n", metric)
		}

		var (
			key   string
			parts = strings.Split(metric, ":")
		)
		key, parts = util.SanitizeKey(parts[0]), parts[1:]

		// TODO add to key counter

		if len(parts) == 0 {
			parts = append(parts, "1")
		}

		for _, part := range parts {
			fields := strings.Split(part, "|")
			if !util.IsValidPacket(fields) {
				log.Printf("Bad line: %+v in msg in %+v\n", fields, metric)
				h.stats.Counters.IncrBadLinesSeen()
				// TODO stats.messages.bad_lines_seen
				continue
			}

			sampleRate := 1
			if len(fields) > 2 {
				sampleRate, _ = strconv.Atoi(fields[2][1:])
			}

			metricType := strings.TrimSpace(fields[1])

			value, _ := strconv.Atoi(fields[0])
			switch metricType {
			// TODO add more
			default: // c
				h.stats.Counters.Incr(key, float64(value*1/sampleRate))
			}

			log.Printf("counters: %s\n", h.stats.Counters.String())
		}
	}

}
