# Files, metadata, Readers and a child process

From the workspace root:

```sh
go run ./examples/system
```

## Architecture

```text
Temp file -> write/seek/read/stat -> close/remove; context -> child process -> output
```

Read `main.go`, predict the result, run it, then change one input. Identify the resource
owner and every error/cleanup path. This is a complete demonstration, not an exercise starter.
See [the roadmap](../../Roadmap.md) for its place in the learning sequence.
