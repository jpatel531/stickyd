package util

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSanitizeKeyWhitespace(t *testing.T) {
	k := SanitizeKey("hello world")
	assert.Equal(t, "hello_world", k)
}

func TestSanitizeKeySlashes(t *testing.T) {
	k := SanitizeKey("hello/world")
	assert.Equal(t, "hello-world", k)
}

func TestSanitizeKeyNonAlphaNumeric(t *testing.T) {
	k := SanitizeKey("mon£y!")
	assert.Equal(t, "mony", k)
}
