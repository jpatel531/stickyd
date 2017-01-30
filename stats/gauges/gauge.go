package gauges

import (
	"fmt"
	"github.com/jpatel531/stickyd/util/collections"
)

type Gauges interface {
	Set(key string, n float64)
	fmt.Stringer
}

func New() Gauges {
	return collections.NewFloatMap()
}
