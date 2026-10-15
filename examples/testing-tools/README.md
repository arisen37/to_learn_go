# Unit tests, a benchmark and a fuzz property

From the workspace root:

```sh
go test ./examples/testing-tools
```

## Architecture

```text
Valid UTF-8 -> rune reversal -> double reversal property
```

Read `main.go`, predict the result, run it, then change one input. Identify the resource
owner and every error/cleanup path. This is a complete demonstration, not an exercise starter.
See [the roadmap](../../Roadmap.md) for its place in the learning sequence.
