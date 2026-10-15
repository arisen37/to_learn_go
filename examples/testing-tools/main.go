package main

import (
	"errors"
	"fmt"
	"unicode/utf8"
)

func Reverse(s string) (string, error) {
	if !utf8.ValidString(s) {
		return "", errors.New("invalid UTF-8")
	}
	r := []rune(s)
	for i, j := 0, len(r)-1; i < j; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r), nil
}
func main() { s, err := Reverse("Go世界"); fmt.Println(s, err) }
