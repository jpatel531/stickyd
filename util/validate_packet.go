package util

import (
	"strconv"
)

func IsValidPacket(fields []string) bool {
	if len(fields) < 2 {
		return false
	}

	if len(fields) > 2 {
		if !isValidSampleRate(fields[2]) {
			return false
		}
	}

	switch fields[1] {
	case "s":
		return true
	case "g":
		_, err := strconv.Atoi(fields[0])
		return err == nil
	case "ms":
		num, err := strconv.Atoi(fields[0])
		return err == nil && num >= 0
	default:
		_, err := strconv.Atoi(fields[0])
		return err == nil
	}
}

func isValidSampleRate(str string) bool {
	if len(str) > 1 && string(str[0]) == "@" {
		num, err := strconv.Atoi(str[1:])
		return err == nil && num >= 0
	}
	return false
}
