package lab

import (
	"slices"
	"testing"
)

func TestWords(t *testing.T) {
	for _, tc := range []struct {
		in   string
		want []string
	}{
		{"  Hello, WORLD! ... café  ", []string{"hello", "world", "café"}},
		{"“GO”\tDON'T\nनमस्ते", []string{"go", "don't", "नमस्ते"}},
		{"", nil}, {"!!! ;", nil}, {"x-y", []string{"x-y"}},
	} {
		if got := Words(tc.in); !slices.Equal(got, tc.want) {
			t.Fatalf("Words(%q)=%v want %v", tc.in, got, tc.want)
		}
	}
}
