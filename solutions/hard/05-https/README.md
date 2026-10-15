# Standard-library HTTPS server — explained solution

Read the [assignment contract](../../../assignments/hard/05-https/README.md) before the code.

## Reasoning and tradeoffs

TLS encrypts bytes and authenticates the server certificate before HTTP parsing. The standard http.Server can load PEM files directly with ServeTLS. Header/body/idle timeouts limit different waiting phases; they are not interchangeable. The listener binds an address, whereas a handler defines application behavior.
Serve owns a single server goroutine and joins it. When the request to stop is canceled, derive the shutdown timeout from Background: deriving from the already canceled ctx would abort shutdown immediately. Shutdown stops new accepts and waits for active handlers; Close is the fallback. The test creates a certificate with internal/testutil, serves TLS over net.Pipe, trusts that certificate explicitly and never uses InsecureSkipVerify. Production certificates and renewal are operational responsibilities, described in lesson 12.

## Architecture

```text
Listener --> TLS handshake --> HTTP server --> handler; ctx cancellation --> bounded shutdown --> join
```

The implementation is in [solution.go](solution.go). The identical acceptance tests are
copied here, so you can verify the reference without changing your work:

```sh
go test -race -count=1 -timeout=30s ./solutions/hard/05-https
```

After comparing, close the answer and implement the algorithm again from memory.
Explain why each invalid-input or resource-lifecycle test is necessary.
