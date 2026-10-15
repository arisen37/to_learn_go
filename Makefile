SHELL := /bin/sh
GOCACHE ?= $(CURDIR)/.cache/go-build
export GOCACHE

.PHONY: help test check repair verify fmt clone-interviews
help:
	@printf '%s\n' 'make test                     Reference solutions and example tests' 'make check ID=easy/01-statistics  One learner assignment, with race detector' 'make repair                   Learner repair lab (initially fails)' 'make verify                   Compile starters, audit links/tests, run references, vet' 'make fmt                      Format local Go code' 'make clone-interviews         Clone researched GitHub questions (network required)'
test:
	go test -count=1 -timeout=30s ./solutions/... ./examples/... ./internal/...
check:
	@test -n "$(ID)" || { printf '%s\n' 'Use make check ID=easy/01-statistics'; exit 2; }
	go test -race -count=1 -timeout=30s "./assignments/$(ID)"
repair:
	go test -race -count=1 -timeout=30s ./interviews/repair-shop
verify:
	sh scripts/verify.sh
fmt:
	gofmt -w assignments solutions examples internal interviews/repair-shop
clone-interviews:
	python3 scripts/clone-interviews.py
