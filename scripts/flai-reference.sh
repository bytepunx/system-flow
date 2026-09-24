#!/usr/bin/env sh
# Regenerate docs/users/flai-reference.md, and the flag index in
# docs/operators/settings.md, from the command help in flai/cmd.
# Runs the tree's source with go run, so bin/flai is left as it is.
set -eu
. "$(dirname "$0")/env.sh"
cd "$ROOT/flai"
go run . reference --write "$ROOT/docs/users/flai-reference.md" --settings "$ROOT/docs/operators/settings.md"
