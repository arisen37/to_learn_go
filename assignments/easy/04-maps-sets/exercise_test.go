package lab

import (
	"maps"
	"slices"
	"testing"
)

func TestMaps(t *testing.T) {
	if got := Frequency([]string{"go", "Go", "go", ""}); !maps.Equal(got, map[string]int{"go": 2, "Go": 1, "": 1}) {
		t.Fatal(got)
	}
	if len(Frequency(nil)) != 0 {
		t.Fatal("empty frequency")
	}
	in := []int{3, -1, 3, 2, -1}
	got := UniqueSorted(in)
	if !slices.Equal(got, []int{-1, 2, 3}) {
		t.Fatal(got)
	}
	if !slices.Equal(in, []int{3, -1, 3, 2, -1}) {
		t.Fatal("mutated")
	}
	if len(UniqueSorted(nil)) != 0 {
		t.Fatal("empty unique")
	}
}
