package gauges

import (
	"encoding/json"
	"fmt"
	"github.com/jpatel531/stickyd/util/collections"
)

type Gauges interface {
	Set(key string, n float64)
	json.Marshaler
	fmt.Stringer
}

func New() Gauges {
	return collections.NewFloatMap()
}
