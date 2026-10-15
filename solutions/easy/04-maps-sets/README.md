# Hash maps and sets — explained solution

Read the [assignment contract](../../../assignments/easy/04-maps-sets/README.md) before the code.

## Reasoning and tradeoffs

The zero value from a missing map lookup makes m[word]++ convenient. A set is map[T]struct{}: only key membership matters. Collect set keys into a slice and sort because map iteration is unspecified. make is required before inserting into a map; a nil map can be read but cannot be written.
Frequency uses expected O(n) time and O(u) space. UniqueSorted uses expected O(n + u log u) time and O(u) space for u unique values. Unlike C++ std::set, this Go map does not keep keys ordered.

## Architecture

```text
Input --> hash map frequency or membership --> sorted deterministic output
```

The implementation is in [solution.go](solution.go). The identical acceptance tests are
copied here, so you can verify the reference without changing your work:

```sh
go test -race -count=1 -timeout=30s ./solutions/easy/04-maps-sets
```

After comparing, close the answer and implement the algorithm again from memory.
Explain why each invalid-input or resource-lifecycle test is necessary.
