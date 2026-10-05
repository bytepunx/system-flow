#!/usr/bin/env sh
# flaiover's unit tests (vitest), when flaiover/node_modules is present; skipped otherwise.
set -eu
. "$(dirname "$0")/env.sh"
[ -d "$ROOT/flaiover/node_modules" ] || exit 0
# flaiover's tests ask a flai built from this tree, so it must exist and be current
if [ ! -x "$ROOT/bin/flai" ] || [ -n "$(find "$ROOT/flai" -name '*.go' -newer "$ROOT/bin/flai" | head -1)" ]; then
  "$ROOT/scripts/flai-build.sh" >&2
fi
echo "behavior tests (vitest)"
cd "$ROOT/flaiover"
FLAI_BIN="$ROOT/bin/flai" pnpm test:unit
