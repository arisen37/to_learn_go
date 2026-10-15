# A bounded local UDP datagram exchange

From the workspace root:

```sh
go run ./examples/udp
```

## Architecture

```text
Local UDP client -> one datagram -> local packet listener -> echoed datagram
```

Read `main.go`, predict the result, run it, then change one input. Identify the resource
owner and every error/cleanup path. This is a complete demonstration, not an exercise starter.
See [the roadmap](../../Roadmap.md) for its place in the learning sequence.
