#!/usr/bin/env sh
# Install flaiover dependencies with pnpm (installed into .flai-cache by scripts/install-tools.sh).
# With --if-needed it installs only when flaiover/node_modules is missing or
# older than pnpm-lock.yaml, package.json, pnpm-workspace.yaml, or .npmrc, as
# in a story worktree, which git worktree add makes without the git-ignored
# node_modules (I-0080), or after a change to the install configuration
# (I-0259). It says so, and pnpm's
# output goes, on stderr, so that a caller's stdout stays its own.
# An install linked from another store, such as a bare pnpm add made from the
# default store, is removed and made again from env.sh's store rather than
# aborting for want of a TTY to confirm it (I-0177): node_modules is
# git-ignored and rebuilt from the lockfile, so nothing is lost.
set -eu
. "$(dirname "$0")/env.sh"
cd "$ROOT/flaiover"
case "${1:-}" in
  "") exec pnpm install --frozen-lockfile --config.confirm-modules-purge=false ;;
  --if-needed) ;;
  *) echo "usage: $0 [--if-needed]" >&2; exit 2 ;;
esac
# The lockfile and the install configuration: a change to either, such as
# supportedArchitectures in pnpm-workspace.yaml (I-0259), installs again.
newer=""
if [ -f node_modules/.modules.yaml ]; then
  for f in pnpm-lock.yaml package.json pnpm-workspace.yaml .npmrc; do
    if [ -f "$f" ] && [ -n "$(find "$f" -newer node_modules/.modules.yaml)" ]; then
      newer="$f"
      break
    fi
  done
  [ -n "$newer" ] || exit 0
fi
echo "installing flaiover's dependencies in $ROOT/flaiover" >&2
pnpm install --frozen-lockfile --config.confirm-modules-purge=false >&2
# pnpm leaves .modules.yaml alone when nothing changed: mark the install current
touch node_modules/.modules.yaml
