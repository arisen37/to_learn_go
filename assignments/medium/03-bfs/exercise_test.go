package lab

import (
	"slices"
	"testing"
)

func TestBFS(t *testing.T) {
	g := map[string][]string{"a": {"b", "c"}, "b": {"a", "d"}, "c": {"d"}, "d": {"e"}}
	p, ok := ShortestPath(g, "a", "e")
	if !ok || len(p) != 4 || p[0] != "a" || p[3] != "e" {
		t.Fatal(p, ok)
	}
	for i := 1; i < len(p); i++ {
		if !slices.Contains(g[p[i-1]], p[i]) {
			t.Fatal("not a path", p)
		}
	}
	if p, ok = ShortestPath(g, "e", "a"); ok || p != nil {
		t.Fatal("directed unreachable")
	}
	if p, ok = ShortestPath(nil, "x", "x"); !ok || !slices.Equal(p, []string{"x"}) {
		t.Fatal(p, ok)
	}
	if _, ok = ShortestPath(g, "a", "z"); ok {
		t.Fatal("missing goal")
	}
}
