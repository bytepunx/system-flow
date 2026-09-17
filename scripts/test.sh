#!/usr/bin/env sh
# Behavior tests: fast, no external dependencies. Run on every iteration.
set -eu
. "$(dirname "$0")/env.sh"
cd "$ROOT/flai"
echo "behavior tests (go test -short)"
go test -race -short -count=1 ./...
if [ -d "$ROOT/flaiover/node_modules" ]; then
  echo "behavior tests (vitest)"; (cd "$ROOT/flaiover" && pnpm test:unit)
fi
