# Learn Go: from first program to backend repair interviews

Start with **[Roadmap.md](Roadmap.md)**. This workspace contains 15 theory lessons,
runnable standard-library examples, 21 graded assignments (7 easy, 7 medium, 7 hard),
and a larger original repository-repair lab. Every graded task has acceptance tests,
architecture notes, and an explained reference implementation under **[solutions/](solutions/README.md)**.

```sh
go run ./examples/basics
make test                         # completed reference solutions and examples
make check ID=easy/01-statistics   # your first assignment, initially fails
```

Use Go 1.22 or newer. The original `hello.go` TCP experiment is preserved. All new local
curriculum code uses the standard library. Ordinary `go test ./...` intentionally fails
until you implement the TODOs and repair the interview defects. `make verify` compiles
starters, runs references with race detection, checks matching test copies, and vets code.
A local Go build cache is configured by the Makefile/scripts for restricted environments.

| Directory | Purpose |
| --- | --- |
| [lessons/](lessons/README.md) | Theory, diagrams, small snippets and completion gates |
| [examples/](examples/README.md) | Complete runnable demonstrations, including HTTPS and network calls |
| [assignments/](assignments/README.md) | Starter code, contracts and active acceptance tests |
| [solutions/](solutions/README.md) | Working implementations with reasoning, complexity and tradeoffs |
| [docs/](docs/resources.md) | STL mapping, parsing/API cheat sheet, system-call notes and references |
| [interviews/](interviews/README.md) | Original Go repair lab plus researched public company repository guides |
| scripts/ | Verification, test-copy audit and upstream clone automation |

GitHub research succeeded, but shell DNS/network restrictions prevented upstream clones.
The interview catalogue records that limitation, links to the original questions and
provides `make clone-interviews` to fetch them when network access is available. The
original local Go repair lab is fully included and works offline. External repository
sizes and test results are not claimed without a successful download.

For interactive servers, your environment must permit localhost listeners. Network
acceptance tests use net.Pipe and a test-only in-memory listener, including an actual TLS
handshake with certificate validation, so they do not need open ports.
# to_learn_go
