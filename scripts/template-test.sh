#!/usr/bin/env sh
# Render template/ into a temporary directory and check the result. The standard's smoke test.
# Then run, in the rendered project, what its CI workflow runs after flai check:
# flai test --at ci, the suites its workflow file's verify_at names for ci.
set -eu
. "$(dirname "$0")/env.sh"
DEST="$(mktemp -d)"
trap 'rm -rf "$DEST"' EXIT
"$ROOT/scripts/flai.sh" new "$DEST/sample" --template "$ROOT/template" --defaults --no-git
"$ROOT/scripts/flai.sh" check "$DEST/sample" --strict
cd "$DEST/sample"
"$ROOT/scripts/flai.sh" test --at ci
echo "template renders, passes check, and passes flai test --at ci"
