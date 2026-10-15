# CSV records and validation — explained solution

Read the [assignment contract](../../../assignments/easy/07-csv/README.md) before the code.

## Reasoning and tradeoffs

CSV is a record format, not strings.Split on a comma. csv.Reader understands escaped quotes, embedded commas and embedded newlines. Its default FieldsPerRecord learns the header width and rejects uneven rows. Validate the exact header before indexing fields, and return nil on failure to avoid callers accidentally using partial data.
Time O(b), space O(n) records. A streaming API can reduce memory by invoking a callback per row, but then rollback of previously emitted rows must be defined. Blank physical lines are skipped by encoding/csv; this exercise inherits that behavior.

## Architecture

```text
io.Reader --> encoding/csv records --> validated Score slice
```

The implementation is in [solution.go](solution.go). The identical acceptance tests are
copied here, so you can verify the reference without changing your work:

```sh
go test -race -count=1 -timeout=30s ./solutions/easy/07-csv
```

After comparing, close the answer and implement the algorithm again from memory.
Explain why each invalid-input or resource-lifecycle test is necessary.
