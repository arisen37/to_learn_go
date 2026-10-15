# HTTPS with timeouts and graceful shutdown

From the workspace root:

```sh
go run ./examples/https-server -cert .local-certs/cert.pem -key .local-certs/key.pem
```

## Architecture

```text
Listener -> TLS -> ServeMux -> JSON; signal -> Shutdown -> join
```

Read `main.go`, predict the result, run it, then change one input. Identify the resource
owner and every error/cleanup path. This is a complete demonstration, not an exercise starter.
See [the roadmap](../../Roadmap.md) for its place in the learning sequence.
