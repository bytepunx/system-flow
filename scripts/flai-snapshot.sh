#!/usr/bin/env sh
# GoReleaser snapshot build of flai into flai/dist (no publish).
set -eu
. "$(dirname "$0")/env.sh"
cd "$ROOT/flai"
goreleaser release --snapshot --clean --skip=publish
