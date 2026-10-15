# HTTP routes and JSON boundaries — explained solution

Read the [assignment contract](../../../assignments/medium/07-http-handler/README.md) before the code.

## Reasoning and tradeoffs

Keep routing, input validation, calculation and serialization as distinct steps. Body bounds apply before decoding so a valid small object followed by a huge tail cannot evade the limit. Pointer fields retain required-field presence; a second decode enforces one JSON value. Validate overflow before adding. Write headers before writing the response body, because the first write commits the status.
Time O(b), bounded input memory O(1024). httptest exercises the real handler without binding a socket. Error responses are text by contract here; a production API may use a consistent JSON error envelope. Duplicate JSON keys inherit encoding/json behavior.

## Architecture

```text
HTTP request --> route/method gate --> bounded strict decoder --> calculation --> JSON response
```

The implementation is in [solution.go](solution.go). The identical acceptance tests are
copied here, so you can verify the reference without changing your work:

```sh
go test -race -count=1 -timeout=30s ./solutions/medium/07-http-handler
```

After comparing, close the answer and implement the algorithm again from memory.
Explain why each invalid-input or resource-lifecycle test is necessary.
