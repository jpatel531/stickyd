package collections

import (
	"fmt"
)

const shardCount = 32

type Map interface {
	Get(key string) interface{}
	Set(key string, value interface{})
	fmt.Stringer
}
