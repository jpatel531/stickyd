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
		log:           log.New(os.Stdout, "", log.LstdFlags),
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
		c.flush(bundle)
	}
	c.wg.Done()
}

func (c *consoleBackend) flush(bundle *FlushBundle) {
	defer bundle.Wait.Done()

	c.log.Println(logPrefix, "Flushing stats at ", bundle.Timestamp)

	metrics := bundle.Metrics
	outJSON, err := json.Marshal(metrics)
	if err != nil {
		c.log.Println(log.Prefix, "Error marshalling metrics to JSON", err.Error())
		return
	}

	c.log.Println(logPrefix, string(outJSON))
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
