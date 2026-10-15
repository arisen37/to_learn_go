# Weighted shortest paths — explained solution

Read the [assignment contract](../../../assignments/hard/07-dijkstra/README.md) before the code.

## Reasoning and tradeoffs

Dijkstra repeatedly settles the cheapest frontier candidate. Nonnegative weights ensure a settled distance cannot later be improved. A simple heap implementation inserts a fresh candidate on improvement; skip stale candidates whose distances no longer match the map. This avoids implementing decrease-key with an index table. Validate all edges first because early success at start==goal must not bypass the negative-weight rule.
Time O(V+E + E log(E+1)) and auxiliary O(V+E) in this duplicate-candidate variant; do not claim O(V) heap memory. Check overflow before addition. The function returns when the goal is settled, so it does not inspect relaxations beyond that point; only encountered overflow is reported. Negative graphs require Bellman-Ford, which is outside this assignment.

## Architecture

```text
Weighted adjacency map --> min-distance heap --> relax edges --> shortest distance
```

The implementation is in [solution.go](solution.go). The identical acceptance tests are
copied here, so you can verify the reference without changing your work:

```sh
go test -race -count=1 -timeout=30s ./solutions/hard/07-dijkstra
```

After comparing, close the answer and implement the algorithm again from memory.
Explain why each invalid-input or resource-lifecycle test is necessary.
