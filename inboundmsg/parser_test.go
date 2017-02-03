package inboundmsg

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestParserWithSampleRate(t *testing.T) {
	r := strings.NewReader("yolo.gauge:14000|c|@0.1")
	p := newParser(r)
	msg, err := p.parse()
	assert.NoError(t, err)

	var sampleRate float64 = 0.1
	assert.Equal(t, &Message{
		Key:        "yolo.gauge",
		Value:      "14000",
		Type:       "c",
		SampleRate: &sampleRate,
	}, msg)
}

func TestParserNoSampleRate(t *testing.T) {
	r := strings.NewReader("yolo.gauge:14000|c")
	p := newParser(r)
	msg, err := p.parse()
	assert.NoError(t, err)

	assert.Equal(t, &Message{
		Key:   "yolo.gauge",
		Value: "14000",
		Type:  "c",
	}, msg)
}

func TestParserNonNumericalValue(t *testing.T) {
	r := strings.NewReader("bla:wootwootr3|s")
	p := newParser(r)
	msg, err := p.parse()
	assert.NoError(t, err)

	assert.Equal(t, &Message{
		Key:   "bla",
		Value: "wootwootr3",
		Type:  "s",
	}, msg)
}

func TestParserFailsWithIncorrectKey(t *testing.T) {
	r := strings.NewReader(".....bla:wootwootr3|s")
	p := newParser(r)
	msg, err := p.parse()
	assert.Nil(t, msg)
	assert.EqualError(t, err, "expected key, found \".\"")
}

func TestParserFailsWithNoValue(t *testing.T) {
	r := strings.NewReader("bla:|s")
	p := newParser(r)
	msg, err := p.parse()
	assert.Nil(t, msg)
	assert.EqualError(t, err, "expected value, found \"|\"")
}

func TestParserFailsWithNoType(t *testing.T) {
	r := strings.NewReader("bla:7|")
	p := newParser(r)
	msg, err := p.parse()
	assert.Nil(t, msg)
	assert.EqualError(t, err, "expected type, found \"\"")

	r = strings.NewReader("bla:7")
	p = newParser(r)
	msg, err = p.parse()
	assert.Nil(t, msg)
	assert.EqualError(t, err, "expected type declaration, found \"\"")
}
