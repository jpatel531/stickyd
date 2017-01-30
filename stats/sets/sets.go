package sets

import (
	"github.com/jpatel531/stickyd/util/collections"
)

type Sets interface {
	collections.SetMap
}

func New() Sets {
	return collections.NewSetMap()
}
