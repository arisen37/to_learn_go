# Structs, methods, and errors — explained solution

Read the [assignment contract](../../../assignments/easy/05-account/README.md) before the code.

## Reasoning and tradeoffs

Pointer receivers change the same account instead of a copy. Validate before mutation so a failed operation is atomic from the caller's perspective. Check overflow by comparing against MaxInt64-current, not by adding first: an overflowing addition has already lost information. Constructor validation establishes a nonnegative-balance invariant that later methods preserve.
All operations cost O(1) time and space. This type is intentionally sequential; a mutex or transaction boundary is needed when multiple goroutines share it. nil receivers are outside this contract.

## Architecture

```text
Validated transaction --> Account owns balance --> updated balance or unchanged on error
```

The implementation is in [solution.go](solution.go). The identical acceptance tests are
copied here, so you can verify the reference without changing your work:

```sh
go test -race -count=1 -timeout=30s ./solutions/easy/05-account
```

After comparing, close the answer and implement the algorithm again from memory.
Explain why each invalid-input or resource-lifecycle test is necessary.
