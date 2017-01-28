package main

import (
	"github.com/jpatel531/stickyd/frontends"
	"log"
	"strings"
)

type udpHandler struct{}

func (u udpHandler) HandleMessage(msg []byte, rinfo *frontends.RemoteInfo) {
	log.Printf("Received %s from %+v\n", strings.TrimSpace(string(msg)), rinfo)
}

func main() {
	config := &frontends.Config{
		Port: 4000,
	}
	server := frontends.UDP{}
	server.Start(config, udpHandler{})
}
