package lab

import (
	"container/heap"
	"errors"
	"math"
)

type Edge struct {
	To     string
	Weight int64
}
type candidate struct {
	node     string
	distance int64
}
type frontier []candidate

func (h frontier) Len() int           { return len(h) }
func (h frontier) Less(i, j int) bool { return h[i].distance < h[j].distance }
func (h frontier) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *frontier) Push(x any)        { *h = append(*h, x.(candidate)) }
func (h *frontier) Pop() any          { a := *h; n := len(a); v := a[n-1]; *h = a[:n-1]; return v }
func ShortestDistance(graph map[string][]Edge, start, goal string) (int64, bool, error) {
	for _, edges := range graph {
		for _, e := range edges {
			if e.Weight < 0 {
				return 0, false, errors.New("negative weight")
			}
		}
	}
	dist := map[string]int64{start: 0}
	q := &frontier{{start, 0}}
	heap.Init(q)
	for q.Len() > 0 {
		cur := heap.Pop(q).(candidate)
		if dist[cur.node] != cur.distance {
			continue
		}
		if cur.node == goal {
			return cur.distance, true, nil
		}
		for _, e := range graph[cur.node] {
			if cur.distance > math.MaxInt64-e.Weight {
				return 0, false, errors.New("distance overflow")
			}
			next := cur.distance + e.Weight
			old, ok := dist[e.To]
			if !ok || next < old {
				dist[e.To] = next
				heap.Push(q, candidate{e.To, next})
			}
		}
	}
	return 0, false, nil
}
