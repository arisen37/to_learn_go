# Bounded HTTP JSON client — explained solution

Read the [assignment contract](../../../assignments/hard/04-http-client/README.md) before the code.

## Reasoning and tradeoffs

Reuse clients and transports so connection pooling works. Context bounds request lifetime if the caller provides a deadline; client.Timeout supplies an additional overall bound. defer Close immediately after a successful Do. Read at most limit+1 bytes to distinguish a complete body of exactly the limit from an oversized body. A second decode rejects trailing values. out may be partly populated before a decode error; callers should discard it on error.
Time/space O(b), b bounded to 1 MiB. Early non-200 close may prevent connection reuse because the body is not drained; that is an intentional bounded-work tradeoff. Retrying is an extension requiring idempotency, attempt bounds, backoff, Retry-After handling, and cancellation.

## Architecture

```text
Context --> reusable http.Client --> response status/limit checks --> JSON target; body always closed
```

The implementation is in [solution.go](solution.go). The identical acceptance tests are
copied here, so you can verify the reference without changing your work:

```sh
go test -race -count=1 -timeout=30s ./solutions/hard/04-http-client
```

After comparing, close the answer and implement the algorithm again from memory.
Explain why each invalid-input or resource-lifecycle test is necessary.
