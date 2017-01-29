package main

import (
	"github.com/jpatel531/stickyd/config"
	"github.com/jpatel531/stickyd/frontend"
	"github.com/jpatel531/stickyd/stats"
	"log"
	"os"
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

	sts := stats.New(cfg.PrefixStats)
	for _, serverCfg := range cfg.Servers {
		server := frontend.NewUDPFrontend()
		server.Start(serverCfg, handler{
			stats:  sts,
			config: cfg,
		})
	}
	select {}
}
