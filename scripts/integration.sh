#!/usr/bin/env sh
# Integration tests: real git and the monorepo itself. Run after behavior tests pass.
# Arguments go to go test: flai test's integration tier passes -json.
# The gateway checks are left to the smoke tier, which runs them where a key
# is set (gateway-smoke.sh), so that no run spends on them twice.
set -eu
. "$(dirname "$0")/env.sh"
cd "$ROOT/flai"
echo "integration tests (go test, full)"
go test -race -count=1 -skip '^TestGateway' "$@" ./...
