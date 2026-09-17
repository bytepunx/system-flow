#!/usr/bin/env sh
# flaiover: format check, lint, type check, unit tests.
set -eu
. "$(dirname "$0")/env.sh"
# writes and metrics tests need a flai built from this tree
[ -x "$ROOT/bin/flai" ] || "$ROOT/scripts/flai-build.sh"
export FLAI_BIN="$ROOT/bin/flai"
cd "$ROOT/flaiover"
echo "prettier + eslint"; pnpm lint
echo "svelte-check"; pnpm check
echo "vitest"; pnpm test:unit
