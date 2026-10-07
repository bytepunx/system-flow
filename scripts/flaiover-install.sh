#!/usr/bin/env sh
# Install flaiover dependencies with pnpm (installed into .flai-cache by scripts/install-tools.sh).
# With --if-needed it installs only when flaiover/node_modules is missing or
# older than pnpm-lock.yaml, as in a story worktree, which git worktree add
# makes without the git-ignored node_modules (I-0080). It says so, and pnpm's
# output goes, on stderr, so that a caller's stdout stays its own.
set -eu
. "$(dirname "$0")/env.sh"
cd "$ROOT/flaiover"
case "${1:-}" in
  "") exec pnpm install --frozen-lockfile ;;
  --if-needed) ;;
  *) echo "usage: $0 [--if-needed]" >&2; exit 2 ;;
esac
if [ -f node_modules/.modules.yaml ] && [ -z "$(find pnpm-lock.yaml -newer node_modules/.modules.yaml)" ]; then
  exit 0
fi
echo "installing flaiover's dependencies in $ROOT/flaiover" >&2
pnpm install --frozen-lockfile >&2
# pnpm leaves .modules.yaml alone when nothing changed: mark the install current
touch node_modules/.modules.yaml
