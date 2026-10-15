# Goroutines and a channel-owned result

From the workspace root:

```sh
go run ./examples/concurrency
```

## Architecture

```text
Producer -> buffered jobs -> 3 workers -> results -> one sum owner
```

Read `main.go`, predict the result, run it, then change one input. Identify the resource
owner and every error/cleanup path. This is a complete demonstration, not an exercise starter.
See [the roadmap](../../Roadmap.md) for its place in the learning sequence.
