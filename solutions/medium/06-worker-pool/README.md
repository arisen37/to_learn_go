# Bounded goroutines and cancellation — explained solution

Read the [assignment contract](../../../assignments/medium/06-worker-pool/README.md) before the code.

## Reasoning and tradeoffs

Use a fixed number of workers, not one goroutine per input. Jobs carry indices, so each worker writes a unique result slot and completion order does not affect output order. The producer alone closes jobs. Every blocking send/receive watches cancellation. sync.Once chooses the first callback error and cancels the derived context; WaitGroup ensures there is no work still touching result memory on return.
Time O(n * callback cost / workers) ideally, space O(n + workers). A callback that ignores ctx can delay shutdown; Go cannot forcibly kill a goroutine safely. Tests check order, bound, errors and both early/mid-flight cancellation; race detection checks memory access.

## Architecture

```text
Producer --> jobs channel --> N workers --> disjoint result slots; derived context cancels all
```

The implementation is in [solution.go](solution.go). The identical acceptance tests are
copied here, so you can verify the reference without changing your work:

```sh
go test -race -count=1 -timeout=30s ./solutions/medium/06-worker-pool
```

After comparing, close the answer and implement the algorithm again from memory.
Explain why each invalid-input or resource-lifecycle test is necessary.
