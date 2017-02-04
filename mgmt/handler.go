package mgmt

import (
	"encoding/json"
	"github.com/jpatel531/stickyd/config"
	"github.com/jpatel531/stickyd/stats"
	"log"
	"net"
	"strings"
)

const (
	bufLen = 1024
)

type handler struct {
	stats  *stats.Stats
	config *config.Config
}

func (h *handler) handleRequest(conn net.Conn) {
	defer conn.Close()

	for {
		buf := make([]byte, bufLen)

		n, err := conn.Read(buf)
		if err != nil {
			log.Println("Error reading:", err.Error())
			return
		}

		text := string(buf[:n])
		cmdLine := strings.Split(strings.TrimSpace(text), " ")
		var cmd string
		cmd, cmdLine = cmdLine[0], cmdLine[1:]

		// TODO add more commands
		switch cmd {
		case "help":
			conn.Write([]byte("Commands: stats, counters, timers, gauges, delcounters, deltimers, delgauges, health, config, quit\n\n"))
		case "config":
			json.NewEncoder(conn).Encode(h.config)
		case "quit":
			return
		}
	}

}
