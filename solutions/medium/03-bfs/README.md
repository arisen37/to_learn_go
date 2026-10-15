# Shortest unweighted path — explained solution

Read the [assignment contract](../../../assignments/medium/03-bfs/README.md) before the code.

## Reasoning and tradeoffs

BFS visits vertices in increasing edge distance because the frontier is FIFO. Mark visited when enqueuing, not when dequeuing, so cycles and multiple incoming edges do not duplicate work. Save one parent per node and reconstruct backward after reaching the goal. The nil-map read for a leaf naturally produces no neighbors.
Time O(V+E), space O(V) for the reachable subgraph. A DFS can find a path but does not guarantee the shortest path. Weighted graphs require Dijkstra (nonnegative weights) or another appropriate algorithm; see hard/07-dijkstra.

## Architecture

```text
Adjacency map --> FIFO frontier and visited parents --> reconstructed path
```

The implementation is in [solution.go](solution.go). The identical acceptance tests are
copied here, so you can verify the reference without changing your work:

```sh
go test -race -count=1 -timeout=30s ./solutions/medium/03-bfs
```

After comparing, close the answer and implement the algorithm again from memory.
Explain why each invalid-input or resource-lifecycle test is necessary.
