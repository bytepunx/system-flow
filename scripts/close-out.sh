#!/usr/bin/env sh
# Close out a story before review: lint, tests, flai check --strict, the
# narrative, then the commit, stopping at the first step that fails (I-0012).
# Usage: scripts/close-out.sh S-nnnn [git commit options: -m, or -F with a file
# outside the worktree, which git add -A would otherwise commit]
# Run it in the story's worktree. The tests follow what the branch changes
# against main (CLOSE_OUT_BASE): flai/ runs flai-test.sh (lint, every tier,
# the check, and the markdown lint); otherwise template/ runs template-test.sh,
# then the markdown lint and the check; flaiover/ adds flaiover-test.sh.
set -eu
. "$(dirname "$0")/env.sh"
cd "$ROOT"

[ $# -ge 1 ] || { echo "usage: scripts/close-out.sh S-nnnn [git commit options]" >&2; exit 2; }
story="$1"
shift

base="${CLOSE_OUT_BASE:-main}"
# Both sides of a rename, and paths unquoted, so a prefix match sees them all.
committed="$(git -c core.quotePath=false diff --no-renames --name-only "$base...HEAD")"
pending="$(git -c core.quotePath=false diff --no-renames --name-only HEAD)"
untracked="$(git -c core.quotePath=false ls-files --others --exclude-standard)"
touched() { printf '%s\n%s\n%s\n' "$committed" "$pending" "$untracked" | grep -q "^$1/"; }

if touched flai; then
  echo "close-out: flai lint and every test tier"
  "$ROOT/scripts/flai-test.sh"
else
  if touched template; then
    echo "close-out: template render and check"
    "$ROOT/scripts/template-test.sh"
  fi
  echo "close-out: markdown lint"
  "$ROOT/scripts/lint-md.sh"
  echo "close-out: flai check --strict"
  "$ROOT/scripts/check.sh"
fi
if touched flaiover; then
  echo "close-out: flaiover lint, types, and unit tests"
  "$ROOT/scripts/flaiover-test.sh"
fi

# The narrative is written in the main checkout, not on the story branch.
narrative="$CACHE_ROOT/wip/agents/$story.md"
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
