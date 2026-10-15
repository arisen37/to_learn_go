# Synchronization and race detection — explained solution

Read the [assignment contract](../../../assignments/medium/05-counter/README.md) before the code.

## Reasoning and tradeoffs

An increment is a read-modify-write sequence; separate concurrent sequences lose updates. Lock both writes and reads, because reading while another goroutine writes is also a race. The zero value of a mutex is ready to use. WaitGroup joins goroutines but does not protect the counter; each tool has a different responsibility.
Each call is O(1) apart from lock contention; storage O(1). Atomic.Int64 is an alternative for this independent scalar. Run go test -race; a passing ordinary test alone is not evidence of synchronization correctness. Copying a used mutex is invalid; keep a pointer to the owner.

## Architecture

```text
Many goroutines --> synchronized Counter --> consistent numeric snapshot
```

The implementation is in [solution.go](solution.go). The identical acceptance tests are
copied here, so you can verify the reference without changing your work:

```sh
go test -race -count=1 -timeout=30s ./solutions/medium/05-counter
```

After comparing, close the answer and implement the algorithm again from memory.
Explain why each invalid-input or resource-lifecycle test is necessary.
