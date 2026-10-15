package main

import (
	"testing"
	"unicode/utf8"
)

func TestReverse(t *testing.T) {
	for _, tc := range []struct{ in, want string }{{"Go世界", "界世oG"}, {"", ""}, {"a", "a"}} {
		got, err := Reverse(tc.in)
		if err != nil || got != tc.want {
			t.Fatal(got, err)
		}
	}
	if _, err := Reverse(string([]byte{0xff})); err == nil {
		t.Fatal("invalid UTF-8")
	}
}

var result string

func BenchmarkReverse(b *testing.B) {
	for i := 0; i < b.N; i++ {
		result, _ = Reverse("Go世界hello")
	}
}
func FuzzReverse(f *testing.F) {
	f.Add("Go世界")
	f.Add("")
	f.Fuzz(func(t *testing.T, s string) {
		r, err := Reverse(s)
		if !utf8.ValidString(s) {
			if err == nil {
				t.Fatal("accepted invalid UTF-8")
			}
			return
		}
		if err != nil {
			t.Fatal(err)
		}
		back, err := Reverse(r)
		if err != nil || back != s {
			t.Fatalf("round trip %q -> %q -> %q", s, r, back)
		}
	})
}
