package lab

import (
	"math"
	"testing"
)

func TestDijkstra(t *testing.T) {
	g := map[string][]Edge{"a": {{"b", 10}, {"c", 1}}, "c": {{"b", 2}, {"a", 0}}, "b": {{"d", 3}}}
	d, ok, err := ShortestDistance(g, "a", "d")
	if err != nil || !ok || d != 6 {
		t.Fatal(d, ok, err)
	}
	if _, ok, err = ShortestDistance(g, "d", "a"); ok || err != nil {
		t.Fatal(ok, err)
	}
	if d, ok, err = ShortestDistance(nil, "x", "x"); err != nil || !ok || d != 0 {
		t.Fatal(d, ok, err)
	}
	if _, _, err = ShortestDistance(map[string][]Edge{"unused": {{"z", -1}}}, "x", "x"); err == nil {
		t.Fatal("negative edge")
	}
	if _, _, err = ShortestDistance(map[string][]Edge{"a": {{"b", math.MaxInt64}}, "b": {{"c", 1}}}, "a", "c"); err == nil {
		t.Fatal("overflow")
	}
	if d, ok, err = ShortestDistance(map[string][]Edge{"a": {{"b", 0}}}, "a", "b"); err != nil || !ok || d != 0 {
		t.Fatal(d, ok, err)
	}
}
