# TCP framing and partial I/O — explained solution

Read the [assignment contract](../../../assignments/hard/03-tcp-frames/README.md) before the code.

## Reasoning and tradeoffs

TCP has no message boundaries: one Write can require several Reads, and several Writes can arrive together. Prefix each frame with a fixed-width length. io.ReadFull handles fragmented reads and reports incomplete input. Bound the length before allocating, or four hostile bytes can trigger a huge allocation. Writers may return short progress; loop until complete, propagating errors and rejecting no progress.
Time O(n), receive memory O(n) with n<=1 MiB. The caller owns closing and deadlines. net.Pipe tests stream behavior without opening ports. Serialize whole WriteFrame calls when multiple goroutines share one writer, or frame bytes may interleave.

## Architecture

```text
Payload --> 4-byte length + body --> TCP byte stream --> io.ReadFull --> payload
```

The implementation is in [solution.go](solution.go). The identical acceptance tests are
copied here, so you can verify the reference without changing your work:

```sh
go test -race -count=1 -timeout=30s ./solutions/hard/03-tcp-frames
```

After comparing, close the answer and implement the algorithm again from memory.
Explain why each invalid-input or resource-lifecycle test is necessary.
