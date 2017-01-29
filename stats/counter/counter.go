// Thread-safe int map, duh
package counter

import (
	"fmt"
	"hash/fnv"
	"sync"
)

const ShardCount = 32

type Counter interface {
	Get(key string) float64
	Set(key string, value float64)
	Incr(key string, n float64)
	fmt.Stringer
}

// TODO make more efficient
func New() Counter {
	m := make(concurrentFloatMap, ShardCount)
	for i := 0; i < ShardCount; i++ {
		m[i] = &shard{
			data: make(map[string]float64),
		}
	}
	return m
}

type concurrentFloatMap []*shard

func (t concurrentFloatMap) Set(key string, value float64) {
	s := t.getShard(key)
	s.set(key, value)
}

func (t concurrentFloatMap) Incr(key string, n float64) {
	s := t.getShard(key)
	s.incr(key, n)
}

func (t concurrentFloatMap) Get(key string) float64 {
	return t.getShard(key).get(key)
}

func (t concurrentFloatMap) String() string {
	merger := make(map[string]float64)
	for _, s := range t {
		for k, v := range s.data {
			merger[k] = v
		}
	}
	return fmt.Sprintf("%+v", merger)
}

func (t concurrentFloatMap) getShard(key string) *shard {
	h := fnv.New32a()
	h.Write([]byte(key))
	return t[h.Sum32()%ShardCount]
}

type shard struct {
	data map[string]float64
	sync.Mutex
}

func (s *shard) get(key string) float64 {
	return s.data[key]
}

func (s *shard) set(key string, value float64) {
	s.Lock()
	defer s.Unlock()
	s.data[key] = value
}

func (s *shard) incr(key string, n float64) {
	s.Lock()
	defer s.Unlock()
	s.data[key] += n
}

var _ Counter = new(concurrentFloatMap)
