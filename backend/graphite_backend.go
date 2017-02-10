package backend

import (
	"fmt"
	"github.com/jpatel531/stickyd/config"
	"log"
	"net"
	"strings"
	"sync"
	"time"
)

type graphiteBackend struct {
	host             string
	port             int
	prefixStats      string
	globalPrefix     string
	globalSuffix     string
	legacyNamespace  bool
	globalNamespace  []string
	counterNamespace []string
	timerNamespace   []string
	gaugesNamespace  []string
	setsNamespace    []string

	lastFlush     int64
	lastException int64
	flushTime     int64
	flushLength   int64

	flushChannel chan *FlushBundle
	wg           *sync.WaitGroup
}

func (g *graphiteBackend) Name() string {
	return "graphite"
}

func (g *graphiteBackend) Status() map[string]int64 {
	return make(map[string]int64)
}

func (g *graphiteBackend) Start() {
	g.wg.Add(1)
	go g.start()
}

func (g *graphiteBackend) start() {
	for bundle := range g.flushChannel {
		g.flush(bundle)
	}
	g.wg.Done()
}

func (g *graphiteBackend) Flush(bundle *FlushBundle) {
	g.flushChannel <- bundle
}

func (g *graphiteBackend) flush(bundle *FlushBundle) {
	defer bundle.Wait.Done()
	startTime := time.Now().Unix()
	ts := bundle.Timestamp.Unix()
	numStats := 0

	counters := bundle.Metrics.Counters
	counterRates := bundle.Metrics.CounterRates
	timerData := bundle.Metrics.TimerData
	gauges := bundle.Metrics.Gauges

	stats := new(graphiteStats)

	for key, value := range counters {
		valuePersecond := counterRates[key]
		namespace := append(g.counterNamespace, key)

		if g.legacyNamespace {
			stats.add(
				strings.Join(namespace, ".")+g.globalSuffix,
				valuePersecond, ts)
			stats.add("stats_counts."+key+g.globalSuffix, value, ts)
		} else {
			stats.add(
				strings.Join(namespace, ". ")+".rate"+g.globalSuffix,
				valuePersecond, ts)
			stats.add(strings.Join(append(namespace, "count"), ".")+g.globalSuffix, value, ts)
		}
		numStats++
	}

	for key, value := range timerData {
		namespace := append(g.timerNamespace, key)
		graphiteKey := strings.Join(namespace, ".")

		for tdKey, _ := range value {
			stats.add(graphiteKey+"."+tdKey+g.globalSuffix, value, ts)
		}
		numStats++
	}

	for key, value := range gauges {
		namespace := append(g.gaugesNamespace, key)
		stats.add(strings.Join(namespace, ".")+g.globalSuffix, value, ts)
		numStats++
	}

	if g.legacyNamespace {
		stats.add(g.prefixStats+".numStats"+g.globalSuffix, numStats, ts)
		stats.add("stats."+g.prefixStats+".graphiteStats.calculationtime"+g.globalSuffix, time.Now().Unix()-startTime, ts)

		for key, value := range bundle.Metrics.StickyDMetrics {
			stats.add("stats."+g.prefixStats+"."+key+g.globalSuffix, value, ts)
		}
	} else {
		namespace := append(g.globalNamespace, g.prefixStats)
		stats.add(strings.Join(namespace, ".")+".numStats"+g.globalSuffix, numStats, ts)
		stats.add("stats."+g.prefixStats+".graphiteStats.calculationtime"+g.globalSuffix, time.Now().Unix()-startTime, ts)

		for key, value := range bundle.Metrics.StickyDMetrics {
			graphiteKey := append(namespace, key)
			stats.add(strings.Join(graphiteKey, ".")+g.globalSuffix, value, ts)
		}
	}

	log.Println("[graphite]", "numStats", numStats)
	if err := g.postStats(stats); err != nil {
		log.Println("[graphite] An error occured while flushing:", err.Error())
	}
}

func (g *graphiteBackend) postStats(stats *graphiteStats) error {
	if g.host == "" {
		return nil
	}

	lastFlush := g.lastFlush
	lastException := g.lastException
	flushTime := g.flushTime
	flushLength := g.flushLength

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", g.host, g.port))
	if err != nil {
		return err
	}
	defer conn.Close()
	ts := time.Now().Unix()
	namespace := strings.Join(append(g.globalNamespace, g.prefixStats), ".")

	stats.add(namespace+".graphiteStats.last_exception"+g.globalSuffix, lastException, ts)
	stats.add(namespace+".graphiteStats.last_flush"+g.globalSuffix, lastFlush, ts)
	stats.add(namespace+".graphiteStats.flush_time"+g.globalSuffix, flushTime, ts)
	stats.add(namespace+".graphiteStats.flush_length"+g.globalSuffix, flushLength, ts)

	startTime := time.Now().Unix()
	payload := stats.text()
	conn.Write([]byte(payload))

	g.flushTime = time.Now().Unix() - startTime
	g.flushLength = int64(len(payload))
	g.lastFlush = time.Now().Unix()

	return nil
}

func (g *graphiteBackend) Stop() {
	close(g.flushChannel)
	g.wg.Wait()
}

func NewGraphiteBackend(
	startupTime int64, config config.Graphite, prefixStats, graphiteHost string, graphitePort int,
) Backend {
	var (
		globalPrefix     string
		prefixCounter    string
		prefixTimer      string
		prefixGauge      string
		prefixSet        string
		legacyNamespace  bool
		globalNamespace  []string
		counterNamespace []string
		timerNamespace   []string
		gaugesNamespace  []string
		setsNamespace    []string
	)

	if graphitePort == 0 {
		graphitePort = 2003
	}

	if prefixStats == "" {
		prefixStats = "stickyd"
	}

	if config.GlobalPrefix == nil {
		globalPrefix = "stats"
	} else {
		globalPrefix = *config.GlobalPrefix
	}

	if config.PrefixCounter == "" {
		prefixCounter = "counters"
	} else {
		prefixCounter = config.PrefixCounter
	}

	if config.PrefixTimer == "" {
		prefixTimer = "timers"
	} else {
		prefixTimer = config.PrefixTimer
	}

	if config.PrefixGauge == "" {
		prefixGauge = "gauges"
	} else {
		prefixGauge = config.PrefixGauge
	}

	if config.PrefixSet == "" {
		prefixSet = "sets"
	} else {
		prefixSet = config.PrefixSet
	}

	var globalSuffix string
	if config.GlobalSuffix == nil {
		globalSuffix = ""
	} else {
		globalSuffix = "." + *config.GlobalSuffix
	}

	if config.LegacyNamespace == nil {
		legacyNamespace = true
	} else {
		legacyNamespace = *config.LegacyNamespace
	}

	if legacyNamespace == false {
		if globalPrefix != "" {
			globalNamespace = []string{globalPrefix}
			counterNamespace = []string{globalPrefix, prefixCounter}
			timerNamespace = []string{globalPrefix, prefixTimer}
			gaugesNamespace = []string{globalPrefix, prefixGauge}
			setsNamespace = []string{globalPrefix, prefixSet}
		}
	} else {
		globalNamespace = []string{"stats"}
		counterNamespace = []string{"stats"}
		timerNamespace = []string{"stats", "timers"}
		gaugesNamespace = []string{"stats", "gauges"}
		setsNamespace = []string{"stats", "sets"}
	}

	return &graphiteBackend{
		host:             graphiteHost,
		port:             graphitePort,
		prefixStats:      prefixStats,
		globalPrefix:     globalPrefix,
		globalSuffix:     globalSuffix,
		legacyNamespace:  legacyNamespace,
		globalNamespace:  globalNamespace,
		counterNamespace: counterNamespace,
		timerNamespace:   timerNamespace,
		gaugesNamespace:  gaugesNamespace,
		setsNamespace:    setsNamespace,

		lastFlush:     startupTime,
		lastException: startupTime,
		flushTime:     0,
		flushLength:   0,

		flushChannel: make(chan *FlushBundle),
		wg:           &sync.WaitGroup{},
	}
}

type graphiteStats []graphiteMetric

func (g *graphiteStats) add(key string, value interface{}, ts int64) {
	(*g) = append(*g, graphiteMetric{key, value, ts})
}

func (g *graphiteStats) text() string {
	reprs := []string{}
	for _, v := range *g {
		reprs = append(reprs, v.text())
	}
	return strings.Join(reprs, "\n") + "\n"
}

type graphiteMetric struct {
	key   string
	value interface{}
	ts    int64
}

func (g graphiteMetric) text() string {
	return fmt.Sprintf("%s %v %d", g.key, g.value, g.ts)
}
