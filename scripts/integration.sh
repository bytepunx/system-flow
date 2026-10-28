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
# 20 minutes a package, not go test's 10: flai/internal/mcpserver's suite takes 8 to 10 minutes
# under -race, and crossed 10 in merge queue batch 20261028-231632 (S-0507, I-0225).
go test -race -count=1 -timeout 20m -skip '^TestGateway' "$@" ./...
