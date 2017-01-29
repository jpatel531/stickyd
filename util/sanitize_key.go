package util

import (
	"regexp"
)

var (
	whiteSpaceRegexp      = regexp.MustCompile(`\s+`)
	slashRegexp           = regexp.MustCompile(`\/`)
	nonAlphaNumericRegexp = regexp.MustCompile(`[^a-zA-Z_\-0-9\.]`)
)

func SanitizeKey(key string) string {
	key = string(
		whiteSpaceRegexp.ReplaceAll([]byte(key), []byte("_")),
	)

	key = string(
		slashRegexp.ReplaceAll([]byte(key), []byte("-")),
	)

	key = string(
		nonAlphaNumericRegexp.ReplaceAll([]byte(key), []byte("")),
	)

	return key
}
