package backend

import (
	"time"
)

type Backend interface {
	Name() string
	Flush(bundle *FlushBundle)
	Status() map[string]int64
	Start()
	Stop()
}

type FlushBundle struct {
	Timestamp time.Time
	Metrics   map[string]interface{}
}

var Backends = map[string]func(int64) Backend{
	"console": NewConsoleBackend,
}
