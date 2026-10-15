# Concurrent TCP echo with connection ownership

From the workspace root:

```sh
go run ./examples/tcp-echo
```

## Architecture

```text
Accept loop -> bounded handler slots -> newline frames -> echo; cancel -> close all -> join
```

Read `main.go`, predict the result, run it, then change one input. Identify the resource
owner and every error/cleanup path. This is a complete demonstration, not an exercise starter.
See [the roadmap](../../Roadmap.md) for its place in the learning sequence.
