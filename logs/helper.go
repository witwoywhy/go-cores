package logs

import (
	"strings"
	"unicode/utf8"
)

func infoToArgs(info map[string]any) []any {
	var args []any
	for k, v := range info {
		args = append(args, k, v)
	}
	return args
}

func Masking(m map[string]any) {
	for k, v := range m {
		if v == nil {
			continue
		}

		switch val := v.(type) {
		case map[string]any:
			Masking(val)
		case string:
			if _, sensitive := MaskingList[k]; sensitive {
				m[k] = MaskChar(val)
			}
		case []any:
			for _, item := range val {
				if mapVal, ok := item.(map[string]any); ok {
					Masking(mapVal)
				}
			}
		}
	}
}

func MaskChar(s string) string {
	ln := utf8.RuneCountInString(s)
	if ln < 3 {
		return strings.Repeat(MaskingChar, ln)
	}

	f, fSize := utf8.DecodeRuneInString(s)
	l, lSize := utf8.DecodeLastRuneInString(s)
	mask := strings.Repeat(MaskingChar, ln-2)

	var b strings.Builder
	b.Grow(fSize + len(mask) + lSize)
	b.WriteRune(f)
	b.WriteString(mask)
	b.WriteRune(l)
	return b.String()
}
