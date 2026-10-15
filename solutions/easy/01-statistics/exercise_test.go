package lab

import (
	"slices"
	"testing"
)

func TestStats(t *testing.T) {
	for _, tc := range []struct {
		in          []int
		lo, hi, sum int
		mean        float64
	}{
		{[]int{2, -1, 5}, -1, 5, 6, 2}, {[]int{-5, -3}, -5, -3, -8, -4}, {[]int{7}, 7, 7, 7, 7}, {[]int{1, 2}, 1, 2, 3, 1.5},
	} {
		before := slices.Clone(tc.in)
		lo, hi, sum, mean, err := Stats(tc.in)
		if err != nil || lo != tc.lo || hi != tc.hi || sum != tc.sum || mean != tc.mean {
			t.Fatalf("Stats(%v)=(%d,%d,%d,%v,%v)", tc.in, lo, hi, sum, mean, err)
		}
		if !slices.Equal(before, tc.in) {
			t.Fatal("mutated input")
		}
	}
	for _, in := range [][]int{nil, {}} {
		if _, _, _, _, err := Stats(in); err == nil {
			t.Fatal("empty input must fail")
		}
	}
}
