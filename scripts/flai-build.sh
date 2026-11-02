#!/usr/bin/env sh
# Build bin/flai from flai/, statically (CGO_ENABLED=0), so the one binary in a story worktree runs
# both on the host and in its agent container, whose Alpine has no glibc (S-0510).
set -eu
. "$(dirname "$0")/env.sh"
cd "$ROOT/flai"
echo "building bin/flai"
CGO_ENABLED=0 go build -o "$ROOT/bin/flai" .
