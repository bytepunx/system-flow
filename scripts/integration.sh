#!/usr/bin/env sh
# Integration tests: real git and the monorepo itself. Run after behavior tests pass.
set -eu
. "$(dirname "$0")/env.sh"
cd "$ROOT/flai"
echo "integration tests (go test, full)"
go test -race -count=1 ./...
