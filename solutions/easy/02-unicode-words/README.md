# Unicode text normalization — explained solution

Read the [assignment contract](../../../assignments/easy/02-unicode-words/README.md) before the code.

## Reasoning and tradeoffs

A string is UTF-8 bytes; punctuation and letters may occupy multiple bytes. The strings and unicode packages operate at the right abstraction for this task. Fields handles tabs/newlines and repeated whitespace. TrimFunc affects only the ends, preserving contractions and hyphens. Lowercasing is Unicode aware but is not locale-specific full case folding or Unicode normalization.
Time O(b) for b input bytes, space O(b) for the output. Do not use strings.Split(text, " ") or trim a fixed ASCII punctuation list: both miss tested cases.

## Architecture

```text
UTF-8 text --> whitespace tokens --> Unicode trim and lowercase --> word slice
```

The implementation is in [solution.go](solution.go). The identical acceptance tests are
copied here, so you can verify the reference without changing your work:

```sh
go test -race -count=1 -timeout=30s ./solutions/easy/02-unicode-words
```

After comparing, close the answer and implement the algorithm again from memory.
Explain why each invalid-input or resource-lifecycle test is necessary.
