#!/usr/bin/env sh
# Integration tests: real git and the monorepo itself. Run after behavior tests pass.
# Arguments go to go test: flai test's integration tier passes -json.
# The gateway checks are left to the smoke tier, which runs them where a key
# is set (gateway-smoke.sh), so that no run spends on them twice.
# In a worktree it first installs flaiover's dependencies when they are missing
# or stale, so that the components tests run on this tree's lockfile rather
# than on the main checkout's install (I-0178).
set -eu
. "$(dirname "$0")/env.sh"
if [ "$CACHE_ROOT" != "$ROOT" ]; then
  "$ROOT/scripts/flaiover-install.sh" --if-needed
fi
cd "$ROOT/flai"
echo "integration tests (go test, full)"
go test -race -count=1 -skip '^TestGateway' "$@" ./...
