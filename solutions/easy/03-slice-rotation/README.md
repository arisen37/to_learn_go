# Slices, indexing, and ownership — explained solution

Read the [assignment contract](../../../assignments/easy/03-slice-rotation/README.md) before the code.

## Reasoning and tradeoffs

A slice is a descriptor over an array, so returning append(nums[k:], nums[:k]...) may mutate the caller depending on capacity. Allocate the output explicitly and copy the suffix followed by the prefix. Normalize k with remainder and one addition for a negative remainder; checking emptiness first avoids division by zero.
Time O(n), extra space O(n). The ownership test changes the returned slice and checks that the original remains intact. It checks a semantic guarantee rather than allocation internals.

## Architecture

```text
Caller backing array --> normalized offset --> newly allocated backing array
```

The implementation is in [solution.go](solution.go). The identical acceptance tests are
copied here, so you can verify the reference without changing your work:

```sh
go test -race -count=1 -timeout=30s ./solutions/easy/03-slice-rotation
```

After comparing, close the answer and implement the algorithm again from memory.
Explain why each invalid-input or resource-lifecycle test is necessary.
