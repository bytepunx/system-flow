#!/usr/bin/env sh
# Integration tests: real git and the monorepo itself. Run after behavior tests pass.
# Arguments go to go test: flai test's integration tier passes -json.
# The gateway checks are left to the manual gateway tier, which runs
# gateway-smoke.sh only when named (flai test gateway, S-0548), so that no
# close-out or batch spends on them.
# In a worktree it first installs flaiover's dependencies when they are missing
# or stale, so that the components tests run on this tree's lockfile rather
# than on the main checkout's install (I-0178).
#
# A close-out runs it on one core and without the race detector (S-0722):
# -p 1 builds and tests one package at a time, and GOMAXPROCS=1 holds the go
# command, the compiler, and each test binary, with the go processes the tests
# start, to one thread of Go code, and -parallel to its default of 1, so that
# several stories closing out at once do not saturate the host. With --race first it runs the race detector on every core
# instead, as the manual race tier (flai test race) and make flai-test do.
set -eu
. "$(dirname "$0")/env.sh"
race=""
if [ "${1:-}" = "--race" ]; then
  race=1
  shift
fi
if [ "$CACHE_ROOT" != "$ROOT" ]; then
  "$ROOT/scripts/flaiover-install.sh" --if-needed
fi
cd "$ROOT/flai"
if [ -n "$race" ]; then
  # 20 minutes a package under -race: flai/internal/mcpserver's suite took 8 to 10
  # minutes before S-0548 ran it in parallel, and crossed 10 in merge queue batch
  # 20261028-231632 (S-0507, I-0225); the slowest package took 252 s on every core
  # with the host at a load of 69 (S-0722).
  echo "integration tests (go test -race, full, every core)"
  go test -race -count=1 -timeout 20m -skip '^TestGateway' "$@" ./...
else
  # 10 minutes a package at one core: the slowest, ./cmd and ./internal/mcpserver,
  # took 130 to 133 s and 121 s in three runs (S-0722), so a host loaded three
  # times over still leaves them inside it, and a hung test is named well before
  # the whole run would end.
  echo "integration tests (go test, full, one core)"
  GOMAXPROCS=1 go test -p 1 -count=1 -timeout 10m -skip '^TestGateway' "$@" ./...
fi
