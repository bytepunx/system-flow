#!/usr/bin/env sh
# Validate this repository against the system-flow standard.
set -eu
. "$(dirname "$0")/env.sh"
cd "$ROOT"
exec "$ROOT/scripts/flai.sh" check --strict "$@"
