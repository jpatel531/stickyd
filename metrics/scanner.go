package metrics

import (
	"bufio"
	"bytes"
	"io"
)

type scanner struct {
	r *bufio.Reader
}

func newScanner(r io.Reader) *scanner {
	return &scanner{r: bufio.NewReader(r)}
}

func (s *scanner) read() rune {
	ch, _, err := s.r.ReadRune()
	if err != nil {
		return eofRune
	}
	return ch
}

func (s *scanner) unread() {
	s.r.UnreadRune()
}

func (s *scanner) scan() (tok token, lit string) {
	ch := s.read()

	if isAlphaNumeric(ch) {
		s.unread()
		return s.scanIdent()
	} else if isSample(ch) {
		return s.scanDecimal()
	}

	switch ch {
	case eofRune:
		return eof, ""
	case '|':
		return separator, string(ch)
	case ':':
		return colon, string(ch)
	default:
		return illegal, string(ch)
	}
}

func (s *scanner) scanIdent() (tok token, lit string) {
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	for {
		ch := s.read()
		if ch == eofRune {
			break
		} else if ch == '.' {
			if nextCh := s.read(); !isValidIdentRune(nextCh) {
				s.unread()
				break
			} else {
				buf.WriteRune(ch)
				buf.WriteRune(nextCh)
				continue
			}
			buf.WriteRune(ch)
		} else if !isValidIdentRune(ch) {
			s.unread()
			break
		} else {
			buf.WriteRune(ch)
		}
	}

	return ident, buf.String()
}

func (s *scanner) scanDecimal() (tok token, lit string) {
	var buf bytes.Buffer

	ch := s.read()

	if ch != '0' && ch != '1' {
		return illegal, buf.String()
	}
	buf.WriteRune(ch)
	if ch == '1' {
		return decimal, buf.String()
	}

	ch = s.read()
	if ch != '.' {
		return illegal, buf.String()
	}
	buf.WriteRune(ch)

	ch = s.read()
	if !isNumber(ch) {
		return illegal, buf.String()
	}
	buf.WriteRune(ch)

	for {
		ch = s.read()
		if !isNumber(ch) {
			s.unread()
			break
		}
		buf.WriteRune(ch)
	}
	return decimal, buf.String()
}
