#!/usr/bin/env sh
# flaiover's format check and lint on the files given, relative to flaiover/,
# as flai test's flaiover-lint tier passes them (I-0105): prettier --check on
# every file it formats, eslint on the scripts and components among them.
# svelte-check takes no file list, so it stays in scripts/flaiover-test.sh.
# In a story worktree it first installs flaiover's dependencies when they are
# missing or stale (I-0080); in the main checkout it is skipped while
# flaiover/node_modules is missing.
set -eu
. "$(dirname "$0")/env.sh"
if [ "$CACHE_ROOT" != "$ROOT" ]; then
  "$ROOT/scripts/flaiover-install.sh" --if-needed >&2
fi
[ -d "$ROOT/flaiover/node_modules" ] || exit 0
[ "$#" -gt 0 ] || exit 0
cd "$ROOT/flaiover"
status=0
echo "prettier --check $*"
pnpm exec prettier --check --ignore-unknown "$@" || status=1
# eslint lints only what its configuration covers: pass it those files alone
for f in "$@"; do
  shift
  case "$f" in
    *.js | *.ts | *.svelte) set -- "$@" "$f" ;;
  esac
done
if [ "$#" -gt 0 ]; then
  echo "eslint $*"
  pnpm exec eslint --no-warn-ignored "$@" || status=1
fi
exit "$status"
