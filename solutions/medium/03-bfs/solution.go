package lab

import "slices"

func ShortestPath(graph map[string][]string, start, goal string) ([]string, bool) {
	queue := []string{start}
	parent := map[string]string{}
	seen := map[string]bool{start: true}
	for head := 0; head < len(queue); head++ {
		v := queue[head]
		if v == goal {
			path := []string{v}
			for v != start {
				v = parent[v]
				path = append(path, v)
			}
			slices.Reverse(path)
			return path, true
		}
		for _, next := range graph[v] {
			if !seen[next] {
				seen[next] = true
				parent[next] = v
				queue = append(queue, next)
			}
		}
	}
	return nil, false
}
