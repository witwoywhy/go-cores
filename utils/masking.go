package utils

import (
	"net/http"
	"strings"
	"unicode/utf8"
)

func MaskString(value, char string) string {
	total := utf8.RuneCountInString(value)
	if total < 3 {
		return strings.Repeat(char, total)
	}

	f, fSize := utf8.DecodeRuneInString(value)
	l, lSize := utf8.DecodeLastRuneInString(value)
	mask := strings.Repeat(char, total-2)

	var b strings.Builder
	b.Grow(fSize + len(mask) + lSize)
	b.WriteRune(f)
	b.WriteString(mask)
	b.WriteRune(l)
	return b.String()
}

func MaskHeader(key string, header http.Header) {
	if header.Get(key) == "" {
		return
	}

	header.Set(key, "***")
}
