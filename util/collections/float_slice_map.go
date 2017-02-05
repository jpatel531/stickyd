package collections

import (
	"encoding/json"
	"fmt"
	"hash/fnv"
	"sync"
)

type FloatSliceMap interface {
	Append(key string, value float64)
	Map() map[string][]float64
	Clear()
	json.Marshaler
	fmt.Stringer
}

func NewFloatSliceMap() FloatSliceMap {
	m := make(concurrentFloatSliceMap, shardCount)
	m.newShards()
	return m
}

type concurrentFloatSliceMap []*floatSliceShard

func (c concurrentFloatSliceMap) newShards() {
	for i := 0; i < shardCount; i++ {
		c[i] = &floatSliceShard{data: make(map[string][]float64)}
	}
}

func (c concurrentFloatSliceMap) Append(key string, value float64) {
	s := c.getShard(key)
	s.append(key, value)
}

func (c concurrentFloatSliceMap) Map() map[string][]float64 {
	merger := make(map[string][]float64)
	for _, s := range c {
		for k, v := range s.data {
			merger[k] = v
		}
	}
	return merger
}

func (c concurrentFloatSliceMap) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.Map())
}

func (c concurrentFloatSliceMap) String() string {
	return fmt.Sprintf("%+v", c.Map())
}

func (c concurrentFloatSliceMap) Clear() {
	c.newShards()
}

func (c concurrentFloatSliceMap) getShard(key string) *floatSliceShard {
	h := fnv.New32a()
	h.Write([]byte(key))
	return c[h.Sum32()%shardCount]
}

type floatSliceShard struct {
	data map[string][]float64
	sync.Mutex
}

func (f *floatSliceShard) append(key string, value float64) {
	f.Lock()
	defer f.Unlock()
	f.data[key] = append(f.data[key], value)
}
