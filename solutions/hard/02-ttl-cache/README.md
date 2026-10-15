# Concurrent TTL cache and clock injection — explained solution

Read the [assignment contract](../../../assignments/hard/02-ttl-cache/README.md) before the code.

## Reasoning and tradeoffs

Injecting the clock makes exact expiry tests instantaneous and deterministic. Expiration means !now.Before(deadline); now.After alone misses equality. Get takes the write lock because it removes stale entries. Prune captures one time for a consistent sweep. No janitor goroutine means no hidden lifecycle or ticker leak; an application can schedule Prune with an owned context.
Expected O(1) Set/Get, O(n) Prune, O(n) space. Lazy expiry retains untouched expired keys until Prune, and TTL alone does not bound total live entries. Combine with a capacity policy for memory bounds. The serial fake-clock test changes its clock only outside concurrent use.

## Architecture

```text
Concurrent callers --> mutex --> entry map; injected clock --> expiration decisions
```

The implementation is in [solution.go](solution.go). The identical acceptance tests are
copied here, so you can verify the reference without changing your work:

```sh
go test -race -count=1 -timeout=30s ./solutions/hard/02-ttl-cache
```

After comparing, close the answer and implement the algorithm again from memory.
Explain why each invalid-input or resource-lifecycle test is necessary.
