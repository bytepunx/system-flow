#!/usr/bin/env sh
# Run the flai built from this tree, building it first if missing or older than the sources.
# Without go (as in the agent image), run bin/flai as it is, else the installed flai, with a
# warning; with neither, exit 127 so the guard hook refuses the call (ADR-0235).
set -eu
. "$(dirname "$0")/env.sh"
stale() {
  [ -n "$(find "$ROOT/flai" -name '*.go' -newer "$ROOT/bin/flai" | head -1)" ]
}
if ! command -v go >/dev/null 2>&1; then
  if [ -x "$ROOT/bin/flai" ]; then
    if stale; then
      echo "flai.sh: go not found; running bin/flai, older than the sources" >&2
    fi
    exec "$ROOT/bin/flai" "$@"
  fi
  if installed=$(command -v flai); then
    echo "flai.sh: go not found; running the installed flai, $installed" >&2
    exec "$installed" "$@"
  fi
  echo "flai.sh: go not found and no flai on PATH" >&2
  exit 127
fi
if [ ! -x "$ROOT/bin/flai" ] || stale; then
  "$ROOT/scripts/flai-build.sh" >&2
fi
exec "$ROOT/bin/flai" "$@"
