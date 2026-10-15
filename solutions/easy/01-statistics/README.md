# Statistics and control flow — explained solution

Read the [assignment contract](../../../assignments/easy/01-statistics/README.md) before the code.

## Reasoning and tradeoffs

Initialize both extrema from the first value: zero is wrong for an all-negative maximum or an all-positive minimum. Accumulate in a range loop, then convert both operands before division. Integer division would truncate 1.5 to 1. Empty data has no minimum or mean, so return an error rather than invented values.
Time O(n), extra space O(1). The tests cover negatives, a singleton, fractional means, empty input, and ownership. This exercise assumes no int overflow; a financial or untrusted-data program needs a defined overflow policy.

## Architecture

```text
Input slice --> one pass accumulator --> statistics or error
```

The implementation is in [solution.go](solution.go). The identical acceptance tests are
copied here, so you can verify the reference without changing your work:

```sh
go test -race -count=1 -timeout=30s ./solutions/easy/01-statistics
```

After comparing, close the answer and implement the algorithm again from memory.
Explain why each invalid-input or resource-lifecycle test is necessary.
