package collections

import (
	"encoding/json"
	"fmt"
	"hash/fnv"
	"sync"
)

type SetMap interface {
	Has(key string, value interface{}) bool
	Insert(key string, value interface{})
	Map() map[string][]interface{}
	fmt.Stringer
	json.Marshaler
}

func NewSetMap() SetMap {
	m := make(concurrentSetMap, shardCount)
	for i := 0; i < shardCount; i++ {
		m[i] = &setShard{
			data: make(map[string]set),
		}
	}
	return m
}

type concurrentSetMap []*setShard

func (c concurrentSetMap) Has(key string, value interface{}) bool {
	return c.getShard(key).has(key, value)
}

func (c concurrentSetMap) Insert(key string, value interface{}) {
	c.getShard(key).insert(key, value)
}

func (c concurrentSetMap) String() string {
	return fmt.Sprintf("%+v", c.Map())
}

func (c concurrentSetMap) getShard(key string) *setShard {
	h := fnv.New32a()
	h.Write([]byte(key))
	return c[h.Sum32()%shardCount]
}

func (c concurrentSetMap) Map() map[string][]interface{} {
	merger := make(map[string][]interface{})
	for _, s := range c {
		for k, v := range s.values() {
			merger[k] = v
		}
	}
	return merger
}

type setShard struct {
	data map[string]set
	sync.Mutex
}

func (c concurrentSetMap) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.Map())
}

func (s *setShard) has(key string, value interface{}) bool {
	s.Lock()
	defer s.Unlock()
	i := s.data[key]
	if i == nil {
		return false
	}
	return i.has(value)
}

func (s *setShard) insert(key string, value interface{}) {
	s.Lock()
	defer s.Unlock()
	i := s.data[key]
	if i == nil {
		s.data[key] = set(map[interface{}]bool{
			value: true,
		})
		return
	}
	i.insert(value)
}

func (s *setShard) values() map[string][]interface{} {
	s.Lock()
	defer s.Unlock()
	values := make(map[string][]interface{})
	for key, data := range s.data {
		values[key] = data.values()
	}
	return values
}

type set map[interface{}]bool

func (s set) insert(value interface{}) {
	s[value] = true
}

func (s set) has(value interface{}) bool {
	return s[value]
}

func (s set) values() []interface{} {
	values := make([]interface{}, 0)
	for k, _ := range s {
		values = append(values, k)
	}
	return values
}
