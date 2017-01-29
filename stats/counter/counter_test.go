package counter

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSet(t *testing.T) {
	m := New()
	m.Set("hello", 2)
	assert.Equal(t, m.Get("hello"), int64(2))
}

func TestIncr(t *testing.T) {
	m := New()
	m.Incr("hello", 2)
	assert.Equal(t, m.Get("hello"), int64(2))
	m.Incr("hello", 3)
	assert.Equal(t, m.Get("hello"), int64(5))
}
