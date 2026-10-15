# Connect, deadline, write and read a newline frame

From the workspace root:

```sh
go run ./examples/tcp-client
```

## Architecture

```text
Client -> DialContext -> newline request -> bounded response -> Close
```

Read `main.go`, predict the result, run it, then change one input. Identify the resource
owner and every error/cleanup path. This is a complete demonstration, not an exercise starter.
See [the roadmap](../../Roadmap.md) for its place in the learning sequence.
