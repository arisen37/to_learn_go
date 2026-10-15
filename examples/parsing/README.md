# String, JSON and CSV parsing

From the workspace root:

```sh
go run ./examples/parsing
```

## Architecture

```text
Text -> tokens -> typed JSON -> CSV Reader -> CSV Writer
```

Read `main.go`, predict the result, run it, then change one input. Identify the resource
owner and every error/cleanup path. This is a complete demonstration, not an exercise starter.
See [the roadmap](../../Roadmap.md) for its place in the learning sequence.
