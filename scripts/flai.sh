#!/usr/bin/env sh
# Run the flai built from this tree, building it first if missing or older than the sources.
set -eu
. "$(dirname "$0")/env.sh"
if [ ! -x "$ROOT/bin/flai" ] || [ -n "$(find "$ROOT/flai" -name '*.go' -newer "$ROOT/bin/flai" | head -1)" ]; then
  "$ROOT/scripts/flai-build.sh" >&2
fi
exec "$ROOT/bin/flai" "$@"
