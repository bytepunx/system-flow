#!/usr/bin/env sh
# Install flaiover dependencies with pnpm (installed into .flai-cache by scripts/install-tools.sh).
set -eu
. "$(dirname "$0")/env.sh"
cd "$ROOT/flaiover"
pnpm install --frozen-lockfile
