package mgmt

import (
	"encoding/json"
	"fmt"
	"github.com/jpatel531/stickyd/config"
	"github.com/jpatel531/stickyd/stats"
	"io"
	"log"
	"strings"
	"time"
)

const (
	bufLen = 1024
)

type handler struct {
	appStats     *stats.AppStats
	processStats *stats.ProcessStats
	config       *config.Config
	startupTime  int64
}

func (h *handler) handleRequest(conn io.ReadWriteCloser) {
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
		case "stats":
			now := time.Now().Unix()
			uptime := now - h.startupTime

			conn.Write([]byte(fmt.Sprintf("uptime: %d\n", uptime)))
			conn.Write([]byte(fmt.Sprintf("messages.bad_lines_seen: %d\n", h.processStats.Messages.BadLinesSeen())))
			conn.Write([]byte(fmt.Sprintf("messages.last_message_seen: %d\n", now-h.processStats.Messages.LastMessageSeen())))
			conn.Write([]byte("END\n\n"))
		// TODO add backend status
		case "counters":
			if err := writeStats(conn, h.appStats.Counters); err != nil {
				conn.Write([]byte("Error marshalling counters to json\n"))
				continue
			}
		case "gauges":
			if err := writeStats(conn, h.appStats.Gauges); err != nil {
				conn.Write([]byte("Error marshalling gauges to json\n"))
				continue
			}
		case "timers":
			if err := writeStats(conn, h.appStats.Timers); err != nil {
				conn.Write([]byte("Error marshalling timers to json\n"))
				continue
			}
		case "quit":
			return
		}
	}

}

func writeStats(conn io.ReadWriteCloser, stats interface{}) (err error) {
	body, err := json.Marshal(stats)
	if err != nil {
		return
	}
	conn.Write(append(body, []byte("\nEND\n\n")...))
	return
}
