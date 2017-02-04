package keylog

import (
	"fmt"
	"github.com/jpatel531/stickyd/config"
	"github.com/jpatel531/stickyd/stats/counter"
	"io"
	"log"
	"os"
	"sort"
	"time"
)

type KeyLog struct {
	keyCounter counter.Counter
	config     *config.KeyFlush
}

func New(keyCounter counter.Counter, config config.KeyFlush) *KeyLog {
	return &KeyLog{
		keyCounter: keyCounter,
		config:     &config,
	}
}

func (k *KeyLog) Run() {
	go k.run()
}

func (k *KeyLog) run() {
	ticker := time.NewTicker(time.Millisecond * time.Duration(k.config.Interval))

	if k.config.Percent == 0 {
		k.config.Percent = 100
	}

	var writer io.Writer
	if k.config.Log != "" {
		f, err := os.OpenFile(k.config.Log, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Println("Error opening log file", err.Error())
		}
		defer f.Close()
		writer = f
	} else {
		writer = os.Stdout
	}

	for range ticker.C {
		keyMap := k.keyCounter.Map()
		sortedKeys := make(keyFrequencies, 0)

		for key, frequency := range keyMap {
			sortedKeys = append(sortedKeys, keyFrequency{key, int(frequency)})
		}
		sort.Sort(sortedKeys)

		timeString := time.Now().String()

		var logMessage string
		for i := 0; i < len(sortedKeys)*k.config.Percent/100; i++ {
			freq := sortedKeys[i]
			logMessage +=
				fmt.Sprintf("%s count=%d key=%s\n", timeString, freq.frequency, freq.key)
		}
		writer.Write([]byte(logMessage))
		k.keyCounter.Clear()
	}
}
