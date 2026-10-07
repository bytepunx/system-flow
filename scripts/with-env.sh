#!/usr/bin/env sh
# Run a command with this tree's environment: bin/ and the main checkout's bin/
# on PATH, and the caches in .flai-cache, as env.sh sets them. flai test's
# tiers in system-flow.yaml run a tool through it, such as golangci-lint, that
# is installed in bin/ rather than on the PATH flai was started with.
set -eu
[ "$#" -gt 0 ] || { echo "with-env: give the command to run, such as scripts/with-env.sh golangci-lint run ./..." >&2; exit 2; }
. "$(dirname "$0")/env.sh"
exec "$@"
