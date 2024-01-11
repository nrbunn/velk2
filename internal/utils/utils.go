package utils

import (
	"strings"
)

func HasSpecialChar(s string) bool {
	f := func(r rune) bool {
		return r < 'A' || r > 'z'
	}
	return strings.IndexFunc(s, f) != -1
}
