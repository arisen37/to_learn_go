# Generate a local self-signed TLS certificate

From the workspace root:

```sh
go run ./examples/certgen -dir .local-certs
```

## Architecture

```text
crypto/rand + ECDSA key -> X.509 localhost certificate -> PEM files (0600)
```

Read `main.go`, predict the result, run it, then change one input. Identify the resource
owner and every error/cleanup path. This is a complete demonstration, not an exercise starter.
See [the roadmap](../../Roadmap.md) for its place in the learning sequence.
