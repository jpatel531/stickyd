package backend

import (
	"sync"
	"time"

	"github.com/jpatel531/stickyd/pm"
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
	Metrics   *pm.ProcessedMetrics
	Wait      *sync.WaitGroup
}
