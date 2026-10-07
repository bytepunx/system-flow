#!/usr/bin/env sh
# flaiover's unit tests (vitest), when flaiover/node_modules is present; skipped otherwise.
# With arguments it runs vitest with them instead of the whole suite, such as
# flai test's vitest tier does (related --run --reporter=json <files>), and
# prints its banner on stderr so that the reporter's output alone is on stdout.
set -eu
. "$(dirname "$0")/env.sh"
[ -d "$ROOT/flaiover/node_modules" ] || exit 0
# flaiover's tests ask a flai built from this tree, so it must exist and be current
if [ ! -x "$ROOT/bin/flai" ] || [ -n "$(find "$ROOT/flai" -name '*.go' -newer "$ROOT/bin/flai" | head -1)" ]; then
  "$ROOT/scripts/flai-build.sh" >&2
fi
cd "$ROOT/flaiover"
if [ "$#" -gt 0 ]; then
  echo "behavior tests (vitest $*)" >&2
  FLAI_BIN="$ROOT/bin/flai" exec pnpm exec vitest "$@"
fi
echo "behavior tests (vitest)"
FLAI_BIN="$ROOT/bin/flai" pnpm test:unit
