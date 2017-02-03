package metrics

import "bytes"

type Metric struct {
	Key        string
	Value      string
	Type       string
	SampleRate *float64
}

func Parse(raw []byte) (*Metric, error) {
	p := newParser(bytes.NewReader(raw))
	return p.parse()
}
