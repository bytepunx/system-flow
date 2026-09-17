#!/usr/bin/env sh
# Run the flaiover dev server against this repository (PROJECT_DIR defaults to the repo root).
set -eu
. "$(dirname "$0")/env.sh"
cd "$ROOT/flaiover"
export PROJECT_DIR="${PROJECT_DIR:-$ROOT}"
exec pnpm dev "$@"
