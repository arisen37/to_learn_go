#!/bin/sh
set -eu
cd "$(dirname "$0")/.."
export GOCACHE="${GOCACHE:-$PWD/.cache/go-build}"
python3 scripts/audit-curriculum.py
# Compile all local packages, including intentional TODOs, without executing their tests.
go test -run '^$' -timeout=30s ./...
# Run completed reference code only. Network tests use in-memory connections.
go test -race -count=1 -timeout=30s ./solutions/... ./examples/... ./internal/...
go vet ./...
printf '%s\n' 'Verification passed; learner assignment and repair tests remain intentionally unsolved.'
