package lab

import (
	"slices"
	"testing"
)

func TestRotateLeft(t *testing.T) {
	for _, tc := range []struct {
		k    int
		want []int
	}{{0, []int{1, 2, 3, 4}}, {1, []int{2, 3, 4, 1}}, {5, []int{2, 3, 4, 1}}, {-1, []int{4, 1, 2, 3}}, {-8, []int{1, 2, 3, 4}}} {
		in := []int{1, 2, 3, 4}
		got := RotateLeft(in, tc.k)
		if !slices.Equal(got, tc.want) {
			t.Fatalf("k=%d got %v", tc.k, got)
		}
		got[0] = 99
		if in[0] != 1 {
			t.Fatal("result aliases input")
		}
	}
	if len(RotateLeft(nil, 10)) != 0 {
		t.Fatal("empty")
	}
}
