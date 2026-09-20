#!/usr/bin/env sh
# Behavior tests: fast, no external dependencies. Run on every iteration.
set -eu
. "$(dirname "$0")/env.sh"
cd "$ROOT/flai"
echo "behavior tests (go test -short)"
go test -race -short -count=1 ./...
if [ -d "$ROOT/flaiover/node_modules" ]; then
  # flaiover's tests ask a flai built from this tree, so it must exist and be current
  if [ ! -x "$ROOT/bin/flai" ] || [ -n "$(find "$ROOT/flai" -name '*.go' -newer "$ROOT/bin/flai" | head -1)" ]; then
    "$ROOT/scripts/flai-build.sh" >&2
  fi
  echo "behavior tests (vitest)"; (cd "$ROOT/flaiover" && FLAI_BIN="$ROOT/bin/flai" pnpm test:unit)
fi
