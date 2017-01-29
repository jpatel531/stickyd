package main

import (
	"github.com/jpatel531/stickyd/config"
	"github.com/jpatel531/stickyd/frontend"

	"log"
	"os"
	"strings"
)

type udpHandler struct{}

func (u udpHandler) HandleMessage(msg []byte, rinfo *frontends.RemoteInfo) {
	log.Printf("Received %s from %+v\n", strings.TrimSpace(string(msg)), rinfo)
}

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

	for _, serverCfg := range cfg.Servers {
		server := frontends.UDP{}
		server.Start(serverCfg, udpHandler{})
	}
	select {}
}
