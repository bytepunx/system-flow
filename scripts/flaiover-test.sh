#!/usr/bin/env sh
# flaiover: format check, lint, type check, unit tests.
set -eu
. "$(dirname "$0")/env.sh"
# the tests ask a flai built from this tree: build it when missing or older than the sources
if [ ! -x "$ROOT/bin/flai" ] || [ -n "$(find "$ROOT/flai" -name '*.go' -newer "$ROOT/bin/flai" | head -1)" ]; then
  "$ROOT/scripts/flai-build.sh" >&2
fi
export FLAI_BIN="$ROOT/bin/flai"
cd "$ROOT/flaiover"
echo "prettier + eslint"; pnpm lint
echo "svelte-check"; pnpm check
echo "vitest"; pnpm test:unit
