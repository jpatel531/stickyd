package metrics

import (
	"fmt"
	"io"
	"strconv"
)

type parser struct {
	s   *scanner
	buf struct {
		tok token
		lit string
		n   int
	}
}

func newParser(r io.Reader) *parser {
	return &parser{s: newScanner(r)}
}

func (p *parser) scan() (tok token, lit string) {
	if p.buf.n != 0 {
		p.buf.n = 0
		return p.buf.tok, p.buf.lit
	}

	tok, lit = p.s.scan()
	p.buf.tok, p.buf.lit = tok, lit
	return
}

func (p *parser) unscan() {
	p.buf.n = 1
}

func (p *parser) parse() (*Metric, error) {
	metric := &Metric{}

	var (
		tok token
		lit string
	)

	// scan key
	if tok, lit = p.scan(); tok != ident {
		return nil, fmt.Errorf("expected key, found %q", lit)
	}
	metric.Key = lit

	// scan :
	if tok, lit = p.scan(); tok != colon {
		return nil, fmt.Errorf("expected value declaration, found %q", lit)
	}

	// scan value
	if tok, lit = p.scan(); tok != ident {
		return nil, fmt.Errorf("expected value, found %q", lit)
	}
	metric.Value = lit

	// scan |
	if tok, lit = p.scan(); tok != separator {
		return nil, fmt.Errorf("expected type declaration, found %q", lit)
	}

	// scan type
	if tok, lit = p.scan(); tok != ident {
		return nil, fmt.Errorf("expected type, found %q", lit)
	}
	metric.Type = lit

	// if eof, return
	if tok, lit = p.scan(); tok == eof {
		return metric, nil
	}

	// sample rate is indented. check for |
	if tok != separator {
		return nil, fmt.Errorf("expected sample rate section, found %q", lit)

	}

	if tok, lit = p.scan(); tok != decimal {
		return nil, fmt.Errorf("expected sample rate, found %q", lit)
	}

	sampleRate, err := strconv.ParseFloat(lit, 64)
	if err != nil {
		return nil, fmt.Errorf("unable to parse sample rate string: %s", err.Error())
	}
	metric.SampleRate = &sampleRate
	return metric, nil
}
