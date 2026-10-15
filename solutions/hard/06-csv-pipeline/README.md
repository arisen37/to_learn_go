# Streaming CSV pipeline with workers — explained solution

Read the [assignment contract](../../../assignments/hard/06-csv-pipeline/README.md) before the code.

## Reasoning and tradeoffs

One goroutine owns csv.Reader because its parser state is mutable. Workers own conversion, and the calling goroutine alone owns the map; this avoids a map mutex. Bounded channels apply backpressure so parsing cannot run arbitrarily ahead. Every channel operation observes cancellation. The WaitGroup coordinator closes results only after all senders finish, preventing send-on-closed-channel panics. On error the collector keeps draining until the stages exit, then returns no partial output.
Time O(b), memory O(unique keys + workers * row size), excluding an individual CSV record which encoding/csv may grow without a fixed bound. Conversion results arrive unordered, so overflow is detected in arrival order; mixed-sign inputs near int64 limits may fail despite an in-range mathematical final sum. A big.Int accumulator or ordered merge is an extension if order-independent exact summation is required. Reader cancellation is a caller ownership issue, explicitly documented.

## Architecture

```text
CSV reader owner --> bounded row channel --> N converters --> one map owner; cancellation joins stages
```

The implementation is in [solution.go](solution.go). The identical acceptance tests are
copied here, so you can verify the reference without changing your work:

```sh
go test -race -count=1 -timeout=30s ./solutions/hard/06-csv-pipeline
```

After comparing, close the answer and implement the algorithm again from memory.
Explain why each invalid-input or resource-lifecycle test is necessary.
