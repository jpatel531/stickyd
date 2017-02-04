package main

import (
	"github.com/jpatel531/stickyd/config"
	"github.com/jpatel531/stickyd/frontend"
	"github.com/jpatel531/stickyd/mgmt"
	"github.com/jpatel531/stickyd/stats"
	"log"
	"os"
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
		log.Panicln(err)
	}

	startupTime := time.Now().Unix()

	processStats := stats.NewProcessStats(startupTime)
	appStats := stats.NewAppStats(cfg.PrefixStats)

	for _, serverCfg := range cfg.Servers {
		server := frontend.NewUDPFrontend()
		server.Start(serverCfg, handler{
			appStats:     appStats,
			processStats: processStats,
			config:       cfg,
		})
	}

	mgmtServer := mgmt.NewMgmtServer(appStats, processStats, cfg, startupTime)
	mgmtServer.Start()

	select {}
}
