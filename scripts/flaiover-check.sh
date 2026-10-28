#!/usr/bin/env sh
# flaiover's type check, as flai test's svelte-check tier runs it: pnpm check,
# which is svelte-kit sync then svelte-check on the whole project (I-0206).
# svelte-check takes no file list, so this script ignores the files it is
# given, on purpose: a type error a change causes in another file, such as a
# caller of a changed prop, is caught too, outside --all and before the
# close-out's flaiover tier. --incremental keeps svelte2tsx's output on disk
# under .svelte-kit, so a repeat run takes about a third of the time.
# In a story worktree it first installs flaiover's dependencies when they are
# missing or stale (I-0080); in the main checkout it is skipped while
# flaiover/node_modules is missing.
set -eu
. "$(dirname "$0")/env.sh"
if [ "$CACHE_ROOT" != "$ROOT" ]; then
  "$ROOT/scripts/flaiover-install.sh" --if-needed >&2
fi
[ -d "$ROOT/flaiover/node_modules" ] || exit 0
cd "$ROOT/flaiover"
echo "svelte-check --incremental"
pnpm check --incremental
