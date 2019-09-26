package main

import (
	"strings"
	"unicode/utf8"
)

func fixString(s string) string {
	if !utf8.ValidString(s) {
		v := make([]rune, 0, len(s))
		for i, r := range s {
			if r == utf8.RuneError {
				_, size := utf8.DecodeRuneInString(s[i:])
				if size == 1 {
					continue
				}
			}
			v = append(v, r)
		}
		s = string(v)
	}
	s = strings.Trim(s, "")
	s = strings.ReplaceAll(s, "\n", "")
	return s
}
