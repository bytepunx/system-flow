#!/usr/bin/env sh
# Validate this repository against the system-flow standard.
set -eu
cd "$(dirname "$0")/.."
exec flai check --strict "$@"
