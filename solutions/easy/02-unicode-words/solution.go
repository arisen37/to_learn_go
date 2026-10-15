package lab

import (
	"strings"
	"unicode"
)

func Words(text string) []string {
	out := make([]string, 0)
	for _, word := range strings.Fields(text) {
		word = strings.TrimFunc(strings.ToLower(word), unicode.IsPunct)
		if word != "" {
			out = append(out, word)
		}
	}
	return out
}
