package lab

import "slices"

func Frequency(words []string) map[string]int {
	out := make(map[string]int)
	for _, w := range words {
		out[w]++
	}
	return out
}
func UniqueSorted(nums []int) []int {
	seen := make(map[int]struct{})
	for _, n := range nums {
		seen[n] = struct{}{}
	}
	out := make([]int, 0, len(seen))
	for n := range seen {
		out = append(out, n)
	}
	slices.Sort(out)
	return out
}
