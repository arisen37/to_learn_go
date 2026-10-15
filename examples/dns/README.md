# Context-aware DNS resolution

From the workspace root:

```sh
go run ./examples/dns localhost
```

## Architecture

```text
Hostname -> context-aware resolver -> address list
```

Read `main.go`, predict the result, run it, then change one input. Identify the resource
owner and every error/cleanup path. This is a complete demonstration, not an exercise starter.
See [the roadmap](../../Roadmap.md) for its place in the learning sequence.
