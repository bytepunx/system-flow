#!/usr/bin/env sh
# Production build of flaiover (adapter-node output in flaiover/build).
set -eu
. "$(dirname "$0")/env.sh"
cd "$ROOT/flaiover"
pnpm build
