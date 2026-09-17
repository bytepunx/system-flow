#!/usr/bin/env sh
# Run the flaiover dev server against this repository (PROJECT_DIR defaults to the repo root).
set -eu
. "$(dirname "$0")/env.sh"
cd "$ROOT/flaiover"
export PROJECT_DIR="${PROJECT_DIR:-$ROOT}"
[ -x "$ROOT/bin/flai" ] || "$ROOT/scripts/flai-build.sh"
export FLAI_BIN="${FLAI_BIN:-$ROOT/bin/flai}"
# Same token as the container: create it if missing and point the dev server at it.
(cd "$PROJECT_DIR" && "$FLAI_BIN" dashboard token >/dev/null)
export FLAIOVER_TOKEN_FILE="${FLAIOVER_TOKEN_FILE:-$PROJECT_DIR/.flai-cache/dashboard.token}"
echo "log in with: $(cd "$PROJECT_DIR" && "$FLAI_BIN" dashboard token --json | sed -E 's/.*"login_url": *"([^"]+)".*/\1/' | sed 's|http://localhost:[0-9]*|http://localhost:5173|')"
exec pnpm dev "$@"
