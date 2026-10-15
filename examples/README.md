# Runnable examples

All commands run from the workspace root. TCP/UDP/HTTPS demos require permission to
open local ports. Their acceptance tests use in-memory networking where applicable.
The TLS example needs a certificate: run certgen first, then use curl --cacert.

| Example | Command |
| --- | --- |
| [Types, pointers, loops and interfaces](basics/README.md) | `go run ./examples/basics` |
| [String, JSON and CSV parsing](parsing/README.md) | `go run ./examples/parsing` |
| [Goroutines and a channel-owned result](concurrency/README.md) | `go run ./examples/concurrency` |
| [Files, metadata, Readers and a child process](system/README.md) | `go run ./examples/system` |
| [Generate a local self-signed TLS certificate](certgen/README.md) | `go run ./examples/certgen -dir .local-certs` |
| [HTTPS with timeouts and graceful shutdown](https-server/README.md) | `go run ./examples/https-server -cert .local-certs/cert.pem -key .local-certs/key.pem` |
| [A bounded local UDP datagram exchange](udp/README.md) | `go run ./examples/udp` |
| [Context-aware DNS resolution](dns/README.md) | `go run ./examples/dns localhost` |
| [Connect, deadline, write and read a newline frame](tcp-client/README.md) | `go run ./examples/tcp-client` |
| [Concurrent TCP echo with connection ownership](tcp-echo/README.md) | `go run ./examples/tcp-echo` |
| [Unit tests, a benchmark and a fuzz property](testing-tools/README.md) | `go test ./examples/testing-tools` |

Network commands can also be inspected with `go run ./examples/<name> -h` where flags are provided.
See [HTTPS lesson](../lessons/12-https.md) for the complete two-terminal workflow.
