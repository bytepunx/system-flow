#!/usr/bin/env sh
# Close out a story before review: no rebase left unfinished, the branch
# contains the main branch, lint, tests, flai check --strict, the narrative,
# then the commit, stopping at the first step that fails.
# Usage: scripts/close-out.sh S-nnnn [git commit options: -m, or -F with a file
# outside the worktree, which git add -A would otherwise commit]
# The cycle before it (git.md): commit each task when it is done, run
# flai stream sync and resolve what it reports, run the task's tests, and
# commit any fix; then commit what is outstanding and sync again before this.
# Run it in the story's worktree. Add this project's other checks as steps
# of their own; keep each one gated on its exit code.
set -eu
cd "$(dirname "$0")/.."

[ $# -ge 1 ] || { echo "usage: scripts/close-out.sh S-nnnn [git commit options]" >&2; exit 2; }
story="$1"
shift

# The narrative is written in the main checkout, which a story worktree shares
# its git directory with; the main branch is the one checked out there.
common="$(git rev-parse --git-common-dir)"
main="$(cd "$common/.." && pwd)"
base="$(git -C "$main" rev-parse --abbrev-ref HEAD)"
[ "$base" != HEAD ] || { echo "close-out: the main checkout $main is not on a branch" >&2; exit 1; }

# A rebase flai stream sync stopped on is finished or undone first.
for dir in rebase-merge rebase-apply; do
  if [ -e "$(git rev-parse --git-path "$dir")" ]; then
    echo "close-out: a rebase is in progress; finish it with git rebase --continue, or undo it with git rebase --abort, then run this again" >&2
    exit 1
  fi
done

# The branch must contain the main branch: flai stream sync rebases it there,
# and refuses uncommitted changes, so with changes this is checked after the
# commit.
synced() {
  rc=0
  git merge-base --is-ancestor "$base" HEAD || rc=$?
  case "$rc" in
    0) echo "close-out: the branch contains $base" ;;
    1) echo "close-out: the branch does not contain $base; run flai stream sync $story, resolve what it reports, and run this again" >&2; exit 1 ;;
    *) echo "close-out: cannot tell whether the branch contains $base" >&2; exit 1 ;;
  esac
}
checked=no
if [ -z "$(git status --porcelain)" ]; then
  synced
  checked=yes
fi

echo "close-out: markdown lint"
scripts/lint-md.sh
echo "close-out: behavior tests"
scripts/test.sh
echo "close-out: integration tests"
scripts/integration.sh
echo "close-out: smoke tests"
scripts/smoke.sh
echo "close-out: flai check --strict"
scripts/check.sh

narrative="$main/wip/agents/$story.md"
echo "close-out: narrative $narrative"
[ -f "$narrative" ] || { echo "close-out: no narrative for $story at $narrative" >&2; exit 1; }
for section in "Current state" "Next steps"; do
  awk -v h="## $section" '
    $0 == h { on = 1; next }
    on && /^## / { exit }
    on { sub(/^[[:space:]]*([-*]|[0-9]+\.)?[[:space:]]*/, ""); if ($0 != "") { found = 1; exit } }
    END { exit !found }
  ' "$narrative" || { echo "close-out: ## $section in $narrative is empty; write it, then run this again" >&2; exit 1; }
done

status="$(git status --porcelain)"
if [ -n "$status" ]; then
  [ $# -gt 0 ] || { git status --short >&2; echo "close-out: uncommitted changes and no message; pass the message as git commit options (-m, or -F with a file outside the worktree)" >&2; exit 1; }
  echo "close-out: commit"
  git add -A
  git diff --cached --stat
  git commit "$@"
fi
[ "$checked" = yes ] || synced
status="$(git status --porcelain)"
[ -z "$status" ] || { git status --short >&2; echo "close-out: the worktree is not clean after the commit" >&2; exit 1; }
echo "close-out: $story is ready to move to review"
