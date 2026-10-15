package lab

import (
	"errors"
	"io"
	"strings"
	"testing"
)

var errRead = errors.New("read failed")

type brokenReader struct{}

func (brokenReader) Read([]byte) (int, error) { return 0, errRead }
func TestSumLines(t *testing.T) {
	for _, tc := range []struct {
		in   string
		want int64
	}{{"1\n -2\n\n 5 \n", 4}, {"", 0}, {"9223372036854775807", 9223372036854775807}} {
		got, err := SumLines(strings.NewReader(tc.in))
		if err != nil || got != tc.want {
			t.Fatal(got, err)
		}
	}
	for _, bad := range []string{"1\nno\n", "9223372036854775808", "9223372036854775807\n1", "-9223372036854775808\n-1", strings.Repeat("1", (1<<20)+1)} {
		got, err := SumLines(strings.NewReader(bad))
		if err == nil || got != 0 {
			t.Fatal("invalid input", got, err)
		}
	}
	if _, err := SumLines(io.MultiReader(strings.NewReader("5\n"), brokenReader{})); !errors.Is(err, errRead) {
		t.Fatal("lost read error", err)
	}
}
