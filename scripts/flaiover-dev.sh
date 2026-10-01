#!/usr/bin/env sh
# Run the flaiover dev server against this repository (PROJECT_DIR defaults to the repo root).
set -eu
. "$(dirname "$0")/env.sh"
cd "$ROOT/flaiover"
export PROJECT_DIR="${PROJECT_DIR:-$ROOT}"
[ -x "$ROOT/bin/flai" ] || "$ROOT/scripts/flai-build.sh"
export FLAI_BIN="${FLAI_BIN:-$ROOT/bin/flai}"
# The container's token: flai creates it if missing, beside flai serve's state, and says where.
token="$(cd "$PROJECT_DIR" && "$FLAI_BIN" dashboard token --json)"
json_field() { printf '%s\n' "$token" | sed -n "s/.*\"$1\": *\"\\([^\"]*\\)\".*/\\1/p"; }
export FLAIOVER_TOKEN_FILE="${FLAIOVER_TOKEN_FILE:-$(json_field file)}"
echo "token: $FLAIOVER_TOKEN_FILE"
echo "log in with: $(json_field login_url | sed 's|http://localhost:[0-9]*|http://localhost:5173|')"
exec pnpm dev "$@"
