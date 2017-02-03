package inboundmsg

type token int

const (
	illegal token = iota
	eof

	ident
	decimal

	separator
	colon
	sample
)

var eofRune = rune(0)

func isLetter(ch rune) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z')
}

func isNumber(ch rune) bool {
	return (ch >= '0' && ch <= '9')
}

func isAlphaNumeric(ch rune) bool {
	return isNumber(ch) || isLetter(ch)
}

func isValidIdentRune(ch rune) bool {
	return isNumber(ch) ||
		isLetter(ch) ||
		ch == '-' ||
		ch == '_'
}

func isSample(ch rune) bool {
	return ch == '@'
}
