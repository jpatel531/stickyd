package backend

import (
	"encoding/json"
	"log"
	"os"
	"sync"
)

const (
	logPrefix = "[console]"
)

type consoleBackend struct {
	flushChannel  chan *FlushBundle
	lastFlush     int64
	lastException int64
	wg            *sync.WaitGroup
	log           *log.Logger
}

func NewConsoleBackend(startupTime int64) Backend {
	return &consoleBackend{
		flushChannel:  make(chan *FlushBundle),
		lastFlush:     startupTime,
		lastException: startupTime,
		wg:            &sync.WaitGroup{},
		log:           log.New(os.Stdout, logPrefix, log.LstdFlags),
	}
}

func (c *consoleBackend) Name() string {
	return "console"
}

func (c *consoleBackend) Flush(bundle *FlushBundle) {
	c.flushChannel <- bundle
}

func (c *consoleBackend) Start() {
	c.wg.Add(1)
	go c.start()
}

func (c *consoleBackend) start() {
	for bundle := range c.flushChannel {
		c.log.Println("Flushing stats at ", bundle.Timestamp)

		metrics := bundle.Metrics
		out := map[string]interface{}{
			"counters": metrics["counters"],
			"gauges":   metrics["gauges"],
			"sets":     metrics["sets"],
		}

		outJSON, err := json.Marshal(out)
		if err != nil {
			c.log.Println("Error marshalling metrics to JSON", err.Error())
			continue
		}

		c.log.Println(string(outJSON))
	}
	c.wg.Done()
}

func (c *consoleBackend) Stop() {
	close(c.flushChannel)
	c.wg.Wait()
}

func (c *consoleBackend) Status() map[string]int64 {
	return map[string]int64{
		"lastFlush":     c.lastFlush,
		"lastException": c.lastException,
	}
}