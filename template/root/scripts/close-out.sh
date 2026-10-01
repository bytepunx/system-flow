#!/usr/bin/env sh
# Close out a story before review: lint, tests, flai check --strict, the
# narrative, then the commit, stopping at the first step that fails.
# Usage: scripts/close-out.sh S-nnnn [git commit options: -m, or -F with a file
# outside the worktree, which git add -A would otherwise commit]
# Run it in the story's worktree. Add this project's other checks as steps
# of their own; keep each one gated on its exit code.
set -eu
cd "$(dirname "$0")/.."

[ $# -ge 1 ] || { echo "usage: scripts/close-out.sh S-nnnn [git commit options]" >&2; exit 2; }
story="$1"
shift

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

# The narrative is written in the main checkout, which a story worktree shares
# its git directory with.
common="$(git rev-parse --git-common-dir)"
main="$(cd "$common/.." && pwd)"
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
status="$(git status --porcelain)"
[ -z "$status" ] || { git status --short >&2; echo "close-out: the worktree is not clean after the commit" >&2; exit 1; }
echo "close-out: $story is ready to move to review"
