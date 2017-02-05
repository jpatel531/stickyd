package main

import (
	"github.com/jpatel531/stickyd/backend"
	"github.com/jpatel531/stickyd/config"
	"github.com/jpatel531/stickyd/frontend"
	"github.com/jpatel531/stickyd/keylog"
	"github.com/jpatel531/stickyd/mgmt"
	"github.com/jpatel531/stickyd/pm"
	"github.com/jpatel531/stickyd/stats"
	"github.com/jpatel531/stickyd/stats/counter"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		log.Println("Insufficient number of arguments")
		return
	}

	configFile := os.Args[1]
	if configFile == "" {
		log.Println("Config path must be supplied as an argument")
		return
	}

	cfg, err := config.Load(configFile)
	if err != nil {
		log.Println(err)
		return
	}

	startupTime := time.Now().Unix()

	processStats := stats.NewProcessStats(startupTime)
	appStats := stats.NewAppStats(cfg.PrefixStats)
	keyCounter := counter.New()

	backends := make([]backend.Backend, 0)
	if len(cfg.Backends) > 0 {
		for _, bName := range cfg.Backends {
			backendConstructor, ok := backend.Backends[bName]
			if !ok {
				log.Printf("No such backend as %q\n", bName)
				return
			}
			backends = append(backends, backendConstructor(startupTime))
		}
	} else {
		b := backend.Backends["console"](startupTime)
		backends = []backend.Backend{b}
	}

	for _, b := range backends {
		b.Start()
	}

	for _, serverCfg := range cfg.Servers {
		server, ok := frontend.Frontends[serverCfg.Server]
		if !ok {
			log.Printf("No such frontend as %q\n", serverCfg.Server)
			return
		}

		server.Start(serverCfg, handler{
			appStats:     appStats,
			processStats: processStats,
			config:       cfg,
			keyCounter:   keyCounter,
		})
	}

	mgmtServer := mgmt.NewMgmtServer(appStats, processStats, cfg, startupTime)
	mgmtServer.Start()

	if cfg.KeyFlush.Interval > 0 {
		keyLog := keylog.New(keyCounter, cfg.KeyFlush)
		keyLog.Run()
	}

	percentThreshold := cfg.PercentThreshold
	if len(percentThreshold) == 0 {
		percentThreshold = append(percentThreshold, 90)
	}

	flushInterval := cfg.FlushInterval
	if flushInterval == 0 {
		flushInterval = 10000
	}
	flushTicker := time.NewTicker(time.Millisecond * time.Duration(flushInterval))

	signalChannel := make(chan os.Signal, 2)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)

	for {
		select {
		case <-flushTicker.C:
			flushMetrics(appStats, flushInterval, backends, percentThreshold)
			appStats.Clear()
		case <-signalChannel:
			log.Println("Interrupted. Flushing metrics...")
			flushMetrics(appStats, flushInterval, backends, percentThreshold)
			stopBackends(backends)
			os.Exit(1)
		}
	}
}

func flushMetrics(
	appStats *stats.AppStats,
	flushInterval int,
	backends []backend.Backend,
	percentThreshold []int,
) {

	flushWait := &sync.WaitGroup{}
	bundle := &backend.FlushBundle{
		Timestamp: time.Now(),
		Metrics:   pm.ProcessMetrics(appStats, flushInterval, percentThreshold),
		Wait:      flushWait,
	}

	for _, b := range backends {
		flushWait.Add(1)
		b.Flush(bundle)
	}
	flushWait.Wait()
}

func stopBackends(backends []backend.Backend) {
	for _, b := range backends {
		b.Stop()
	}
}
