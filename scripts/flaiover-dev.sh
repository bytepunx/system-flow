#!/usr/bin/env sh
# Run the flaiover dev server against this repository (PROJECT_DIR defaults to the repo root).
set -eu
. "$(dirname "$0")/env.sh"
cd "$ROOT/flaiover"
export PROJECT_DIR="${PROJECT_DIR:-$ROOT}"
[ -x "$ROOT/bin/flai" ] || "$ROOT/scripts/flai-build.sh"
export FLAI_BIN="${FLAI_BIN:-$ROOT/bin/flai}"
exec pnpm dev "$@"
