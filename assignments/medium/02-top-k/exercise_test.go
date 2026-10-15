package lab

import (
	"slices"
	"testing"
)

func TestTopK(t *testing.T) {
	for _, tc := range []struct {
		in   []int
		k    int
		want []int
	}{{[]int{3, 1, 9, 9, 2}, 3, []int{9, 9, 3}}, {[]int{-4, -2, -9}, 2, []int{-2, -4}}, {[]int{2, 1}, 5, []int{2, 1}}, {nil, 3, nil}, {[]int{1}, 0, nil}, {[]int{1}, -1, nil}} {
		before := slices.Clone(tc.in)
		if got := TopK(tc.in, tc.k); !slices.Equal(got, tc.want) {
			t.Fatalf("got %v want %v", got, tc.want)
		}
		if !slices.Equal(before, tc.in) {
			t.Fatal("mutated")
		}
	}
}
