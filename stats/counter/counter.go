package counter

import (
	"encoding/json"
	"fmt"
	"github.com/jpatel531/stickyd/util/collections"
)

type Counter interface {
	Incr(key string, n float64)
	Map() map[string]float64
	Get(key string) float64
	Clear()
	json.Marshaler
	fmt.Stringer
}

func New() Counter {
	return collections.NewFloatMap()
}
