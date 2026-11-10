#!/usr/bin/env sh
# Run the flai built from this tree, building it first if missing or older than the sources.
# Without go (as in the agent image), run bin/flai as it is, else the installed flai, with a
# warning; with neither, exit 127 so the guard hook refuses the call (ADR-0235).
# When the build fails, guard does the same, so the guard hooks never pass a call for a failed
# build; every other command exits with the build's status (ADR-0280).
set -eu
. "$(dirname "$0")/env.sh"
stale() {
  [ -n "$(find "$ROOT/flai" -name '*.go' -newer "$ROOT/bin/flai" | head -1)" ]
}
# fall_back runs bin/flai as it is, else the installed flai, warning why, or exits 127.
fall_back() {
  why=$1
  shift
  if [ -x "$ROOT/bin/flai" ]; then
    if stale; then
      echo "flai.sh: $why; running bin/flai, older than the sources" >&2
    fi
    exec "$ROOT/bin/flai" "$@"
  fi
  if installed=$(command -v flai); then
    echo "flai.sh: $why; running the installed flai, $installed" >&2
    exec "$installed" "$@"
  fi
  echo "flai.sh: $why and no flai on PATH" >&2
  exit 127
}
if ! command -v go >/dev/null 2>&1; then
  fall_back "go not found" "$@"
fi
if [ ! -x "$ROOT/bin/flai" ] || stale; then
  status=0
  "$ROOT/scripts/flai-build.sh" >&2 || status=$?
  if [ "$status" -ne 0 ]; then
    [ "${1:-}" = guard ] || exit "$status"
    fall_back "building bin/flai failed" "$@"
  fi
fi
exec "$ROOT/bin/flai" "$@"
