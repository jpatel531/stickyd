package counter

import (
	"encoding/json"
	"fmt"
	"github.com/jpatel531/stickyd/util/collections"
)

type Counter interface {
	Incr(key string, n float64)
	json.Marshaler
	fmt.Stringer
}

func New() Counter {
	return collections.NewFloatMap()
}
