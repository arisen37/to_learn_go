# Top K with a priority queue — explained solution

Read the [assignment contract](../../../assignments/medium/02-top-k/README.md) before the code.

## Reasoning and tradeoffs

Keep the k best candidates in a min-heap. Its root is the worst retained candidate, so a new larger number replaces the root and heap.Fix restores order. Do not confuse heap order with sorted order: sort only the final k candidates for the required presentation. heap.Push and heap.Pop call your methods; calling h.Push directly skips heap maintenance.
Time O(n log k + k log k), conventionally O(n log k), with O(k) extra space for 1<k<=n; k=1 is O(n). Tests verify behavior; complexity and bounded heap use should also be reviewed in code.

## Architecture

```text
Input stream --> size k min-heap --> descending result
```

The implementation is in [solution.go](solution.go). The identical acceptance tests are
copied here, so you can verify the reference without changing your work:

```sh
go test -race -count=1 -timeout=30s ./solutions/medium/02-top-k
```

After comparing, close the answer and implement the algorithm again from memory.
Explain why each invalid-input or resource-lifecycle test is necessary.
