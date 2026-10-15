# LRU cache: map plus linked list — explained solution

Read the [assignment contract](../../../assignments/hard/01-lru/README.md) before the code.

## Reasoning and tradeoffs

A map finds the list node in expected O(1). A doubly linked list moves arbitrary known nodes to the front in O(1), unlike a slice that shifts elements. The map and list must describe exactly the same entries. On eviction remove from both structures; forgetting either leaks memory or exposes stale values. Updating a key must reuse its node, not insert a duplicate.
Expected O(1) Get/Put, O(capacity) storage. Get mutates recency, so a future concurrent implementation must lock Get as well as Put. A cache is an optimization; decide separately whether misses trigger a load and whether concurrent loads should be deduplicated.

## Architecture

```text
Key --> map lookup --> doubly linked recency list (front newest, back oldest)
```

The implementation is in [solution.go](solution.go). The identical acceptance tests are
copied here, so you can verify the reference without changing your work:

```sh
go test -race -count=1 -timeout=30s ./solutions/hard/01-lru
```

After comparing, close the answer and implement the algorithm again from memory.
Explain why each invalid-input or resource-lifecycle test is necessary.
