# Strict JSON parsing — explained solution

Read the [assignment contract](../../../assignments/easy/06-json/README.md) before the code.

## Reasoning and tradeoffs

A plain int cannot distinguish missing age from a valid zero. Decode into pointer fields to retain presence information, then produce the public value type after validation. DisallowUnknownFields catches misspelled input keys. A second Decode must return io.EOF: the first Decode alone accepts a valid object followed by another value.
Time and space O(b) for b bytes. This is stricter than JavaScript JSON.parse plus property access in some ways, but does not reject duplicate keys. Required-field and business validation are separate from syntax parsing.

## Architecture

```text
JSON bytes --> decoder with strict field checks --> validated Profile
```

The implementation is in [solution.go](solution.go). The identical acceptance tests are
copied here, so you can verify the reference without changing your work:

```sh
go test -race -count=1 -timeout=30s ./solutions/easy/06-json
```

After comparing, close the answer and implement the algorithm again from memory.
Explain why each invalid-input or resource-lifecycle test is necessary.
