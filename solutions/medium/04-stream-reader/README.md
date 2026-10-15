# Reader interfaces and line parsing — explained solution

Read the [assignment contract](../../../assignments/medium/04-stream-reader/README.md) before the code.

## Reasoning and tradeoffs

Accepting io.Reader makes the algorithm independent of files: tests can use strings.NewReader and injected read failures. Scanner simplifies line iteration but must be checked with Err after the loop. Configure a bound explicitly; the default Scanner token limit is smaller. ParseInt catches value overflow, while the accumulator requires separate overflow checks before addition.
Time O(b), extra space O(L), bounded by the maximum line length. Returning zero on every error enforces an all-or-error contract. The wrapper %w preserves the original parse error for errors.Is/errors.As.

## Architecture

```text
Any io.Reader --> bounded bufio.Scanner --> parsed int64 accumulator
```

The implementation is in [solution.go](solution.go). The identical acceptance tests are
copied here, so you can verify the reference without changing your work:

```sh
go test -race -count=1 -timeout=30s ./solutions/medium/04-stream-reader
```

After comparing, close the answer and implement the algorithm again from memory.
Explain why each invalid-input or resource-lifecycle test is necessary.
