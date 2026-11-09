#!/usr/bin/env sh
# Build bin/flai from flai/, statically (CGO_ENABLED=0), so the one binary in a story worktree runs
# both on the host and in its agent container, whose Alpine has no glibc (S-0510).
# It builds to a temporary file beside bin/flai and renames it into place, so tiers that run at
# once (ADR-0266), each of which may build bin/flai, never run a half-written binary (S-0610).
set -eu
. "$(dirname "$0")/env.sh"
cd "$ROOT/flai"
echo "building bin/flai"
mkdir -p "$ROOT/bin"
tmp="$ROOT/bin/.flai.$$"
trap 'rm -f "$tmp"' EXIT
CGO_ENABLED=0 go build -o "$tmp" .
mv -f "$tmp" "$ROOT/bin/flai"
