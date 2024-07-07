package logs

import (
	"reflect"
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

func masking(m map[string]any) {
	for k, v := range m {
		if v == nil {
			continue
		}

		switch reflect.TypeOf(v).Kind() {
		case reflect.Map:
			masking(v.(map[string]any))
		case reflect.String:
			if isKeyInList(k) {
				m[k] = MaskChar(v.(string))
			}
		}
	}
}

func MaskChar(s string) string {
	ln := utf8.RuneCountInString(s)
	if ln < 3 {
		return strings.Repeat(MaskingChar, ln)
	}

	f, _ := utf8.DecodeRuneInString(s)
	l, _ := utf8.DecodeLastRuneInString(s)

	return string(f) + strings.Repeat(MaskingChar, ln-2) + string(l)
}

func isKeyInList(key string) bool {
	_, ok := MaskingList[key]
	return ok
}
