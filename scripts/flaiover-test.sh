#!/usr/bin/env sh
# flaiover: format check, lint, type check, unit tests.
set -eu
. "$(dirname "$0")/env.sh"
cd "$ROOT/flaiover"
echo "prettier + eslint"; pnpm lint
echo "svelte-check"; pnpm check
echo "vitest"; pnpm test:unit
