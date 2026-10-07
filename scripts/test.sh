#!/usr/bin/env sh
# Behavior tests: fast, no external dependencies. Run on every iteration.
# Keep in step with the tiers under tests in system-flow.yaml, which flai test
# runs for the paths a change touches: go-test and vitest here.
set -eu
. "$(dirname "$0")/env.sh"
cd "$ROOT/flai"
echo "behavior tests (go test -short)"
go test -race -short -count=1 ./...
"$ROOT/scripts/flaiover-unit.sh"
