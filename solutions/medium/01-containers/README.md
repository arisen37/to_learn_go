# Generic stack and queue — explained solution

Read the [assignment contract](../../../assignments/medium/01-containers/README.md) before the code.

## Reasoning and tradeoffs

A stack removes from the end of a slice. A queue advances a head index instead of shifting every element on every Pop. Compact on Push when at least half the backing slice is consumed; the copying cost is spread over earlier removals, giving amortized O(1). Reset the empty queue to release retained capacity. Clear removed elements so a pointer-rich slice does not keep objects alive.
Stack Push and queue Push are amortized O(1), Pop is O(1); storage O(n). Generics avoid interface{} assertions while allowing both string and int instances. Race safety is deliberately outside this container contract; add synchronization at the owner boundary when needed.

## Architecture

```text
Caller --> generic container API --> private slice / queue head --> value and presence
```

The implementation is in [solution.go](solution.go). The identical acceptance tests are
copied here, so you can verify the reference without changing your work:

```sh
go test -race -count=1 -timeout=30s ./solutions/medium/01-containers
```

After comparing, close the answer and implement the algorithm again from memory.
Explain why each invalid-input or resource-lifecycle test is necessary.
