package collections

import (
	"encoding/json"
	"fmt"
	"hash/fnv"
	"sync"
)

type FloatMap interface {
	Get(key string) float64
	Set(key string, value float64)
	Incr(key string, n float64)
	Map() map[string]float64
	json.Marshaler
	fmt.Stringer
}

func NewFloatMap() FloatMap {
	m := make(concurrentFloatMap, shardCount)
	for i := 0; i < shardCount; i++ {
		m[i] = &floatShard{
			data: make(map[string]float64),
		}
	}
	return m
}

type concurrentFloatMap []*floatShard

func (c concurrentFloatMap) Set(key string, value float64) {
	s := c.getShard(key)
	s.set(key, value)
}

func (c concurrentFloatMap) Incr(key string, n float64) {
	s := c.getShard(key)
	s.incr(key, n)
}

func (c concurrentFloatMap) Get(key string) float64 {
	return c.getShard(key).get(key)
}

func (c concurrentFloatMap) getShard(key string) *floatShard {
	h := fnv.New32a()
	h.Write([]byte(key))
	return c[h.Sum32()%shardCount]
}

func (c concurrentFloatMap) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.Map())
}

func (c concurrentFloatMap) Map() map[string]float64 {
	merger := make(map[string]float64)
	for _, s := range c {
		for k, v := range s.data {
			merger[k] = v
		}
	}
	return merger
}

func (c concurrentFloatMap) String() string {
	return fmt.Sprintf("%+v", c.Map())
}

type floatShard struct {
	data map[string]float64
	sync.Mutex
}

func (s *floatShard) get(key string) float64 {
	return s.data[key]
}

func (s *floatShard) set(key string, value float64) {
	s.Lock()
	defer s.Unlock()
	s.data[key] = value
}

func (s *floatShard) incr(key string, n float64) {
	s.Lock()
	defer s.Unlock()
	s.data[key] += n
}
