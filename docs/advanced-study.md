# Advanced study after the core roadmap

These are follow-on study directions, not prebuilt graded assignments. The 21 graded
assignments and repair lab have tests and explained solutions; additions you create here
should come with their own contracts, diagrams and tests.

| Topic | What to understand | Suggested next experiment |
| --- | --- | --- |
| Go runtime | Scheduler, stacks, GC, escape analysis | Compare heap allocations with `go build -gcflags=-m` |
| Performance | Measure CPU, allocations, contention | Profile a pipeline, change one thing, compare benchmarks |
| Memory model | Happens-before, locks, channels, atomics | Explain result publication without relying on sleeps |
| Generics | Constraints, comparable, type inference | Generic set with clear ordering and ownership policy |
| Reflection | Runtime types/values, tags, pointer/settable values | Inspect struct tags; compare with typed parsing |
| Unsafe/cgo | ABI, pointer lifetime, portability | Read docs first; these are unnecessary for core exercises |
| Ordered trees | BST invariants, AVL/red-black balancing | Implement ordered lookup with property tests against a sorted slice |
| Graphs | DFS, topological sort, Bellman-Ford | Add cycle detection and negative-edge tests |
| Database/sql | Pooling, driver requirement, transactions, context | Separate module with an explicitly chosen SQL driver |
| Authorization | Identity versus permission, owner checks | Add per-user notes and test cross-user denial |
| Retry policy | Idempotency, exponential backoff, jitter, attempt caps | Context-aware retrying GET with injected clock/transport |
| Rate limiting | Token bucket, admission/backpressure | Deterministic clock tests; no per-request ticker |
| Durability | WAL, fsync, atomic replacement limitations | Crash-recovery log with replay and corruption tests |
| Distributed systems | Partial failure, duplicates, consistency | Idempotent job API with a documented state machine |
| Observability | slog, metrics, traces, bounded labels | Request IDs and structured access logs without secrets |
| Deployment | Config, graceful stop, TLS renewal, health/readiness | Draw reverse-proxy trust and shutdown boundaries |
| Build/dependencies | Modules, go.sum, build tags, cross-compiling | Build for another GOOS/GOARCH without running it |

Architecture should evolve with requirements. Add a queue, database, interface or service
boundary when it answers a concrete ownership, persistence, testability or scale need.
An in-memory mutex cannot coordinate multiple service processes; a database transaction
cannot cancel a remote side effect automatically. Document the failure model before coding.
