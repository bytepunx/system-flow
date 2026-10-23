#!/usr/bin/env sh
# Run the pinned pnpm in flaiover/ with env.sh's store and caches, such as
# scripts/flaiover-pnpm.sh add <package>. A bare pnpm uses the default store
# and may be the host's other pnpm, which leaves an install the next
# flaiover-install.sh has to purge (I-0177).
set -eu
. "$(dirname "$0")/env.sh"
cd "$ROOT/flaiover"
exec pnpm "$@"
