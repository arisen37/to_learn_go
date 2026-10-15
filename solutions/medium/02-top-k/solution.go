package lab

import (
	"container/heap"
	"slices"
)

type minHeap []int

func (h minHeap) Len() int           { return len(h) }
func (h minHeap) Less(i, j int) bool { return h[i] < h[j] }
func (h minHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *minHeap) Push(x any)        { *h = append(*h, x.(int)) }
func (h *minHeap) Pop() any          { old := *h; n := len(old); x := old[n-1]; *h = old[:n-1]; return x }
func TopK(nums []int, k int) []int {
	if k <= 0 {
		return []int{}
	}
	h := &minHeap{}
	heap.Init(h)
	for _, n := range nums {
		if h.Len() < k {
			heap.Push(h, n)
		} else if n > (*h)[0] {
			(*h)[0] = n
			heap.Fix(h, 0)
		}
	}
	out := slices.Clone([]int(*h))
	slices.Sort(out)
	slices.Reverse(out)
	return out
}
