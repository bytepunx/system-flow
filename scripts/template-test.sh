#!/usr/bin/env sh
# Render template/ into a temporary directory and check the result. The standard's smoke test.
set -eu
. "$(dirname "$0")/env.sh"
DEST="$(mktemp -d)"
trap 'rm -rf "$DEST"' EXIT
"$ROOT/scripts/flai.sh" new "$DEST/sample" --template "$ROOT/template" --defaults --no-git
"$ROOT/scripts/flai.sh" check "$DEST/sample" --strict
echo "template renders and passes check"
