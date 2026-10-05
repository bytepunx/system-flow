#!/usr/bin/env sh
# Close out a story before review: no rebase left unfinished, the branch
# contains main, lint, tests, flai check --strict, the narrative, then the
# commit, stopping at the first step that fails (I-0012).
# Usage: scripts/close-out.sh S-nnnn [git commit options: -m, or -F with a file
# outside the worktree, which git add -A would otherwise commit]
# The cycle before it (git.md): commit each task when it is done, run
# flai stream sync and resolve what it reports, run the task's tests, and
# commit any fix; then commit what is outstanding and sync again before this.
# Run it in the story's worktree. The tests follow what the branch changes
# against main (CLOSE_OUT_BASE): flai/ runs flai-test.sh (lint, vitest, the
# full Go tests once, and smoke, which holds the check and the markdown lint);
# otherwise template/ runs template-test.sh, then the markdown lint and the
# check; flaiover/ adds flaiover-test.sh.
set -eu
. "$(dirname "$0")/env.sh"
cd "$ROOT"

[ $# -ge 1 ] || { echo "usage: scripts/close-out.sh S-nnnn [git commit options]" >&2; exit 2; }
story="$1"
shift
# check.sh scopes flai check to the story, on whichever path reaches it: the
# findings outside it are notes, recorded in design/issues, which the commit
# step below commits (S-0249). TestMonorepoIsClean scopes itself the same way.
CLOSE_OUT_STORY="$story"
export CLOSE_OUT_STORY

# Every run ends with one line on stdout naming the story, the outcome, and
# the step it stopped at, so nobody runs it again to learn why it stopped
# (S-0266). Set once the story is known: a usage error prints only the usage.
# The trap reads $? before anything else can change it, and exits with it.
step="the rebase check"
signal=""
finish() {
  code="$1"
  if [ -n "$signal" ]; then
    echo "close-out: $story stopped at $step ($signal)"
  elif [ "$code" -eq 0 ]; then
    echo "close-out: $story passed every step; ready to move to review"
  else
    echo "close-out: $story stopped at $step (exit $code)"
  fi
  exit "$code"
}
trap 'finish "$?"' EXIT
# An interrupt exits, which runs the EXIT trap once, with the signal's status.
trap 'signal=interrupted; exit 130' INT
trap 'signal=terminated; exit 143' TERM

base="${CLOSE_OUT_BASE:-main}"

# A rebase flai stream sync stopped on is finished or undone first.
for dir in rebase-merge rebase-apply; do
  if [ -e "$(git rev-parse --git-path "$dir")" ]; then
    echo "close-out: a rebase is in progress; finish it with git rebase --continue, or undo it with git rebase --abort, then run this again" >&2
    exit 1
  fi
done

# The branch must contain main: flai stream sync rebases it there, and refuses
# uncommitted changes, so with changes this is checked after the commit.
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
step="the sync check"
if [ -z "$(git status --porcelain)" ]; then
  synced
  checked=yes
fi

# Both sides of a rename, and paths unquoted, so a prefix match sees them all.
step="the paths the branch changes against $base"
committed="$(git -c core.quotePath=false diff --no-renames --name-only "$base...HEAD")"
pending="$(git -c core.quotePath=false diff --no-renames --name-only HEAD)"
untracked="$(git -c core.quotePath=false ls-files --others --exclude-standard)"
touched() { printf '%s\n%s\n%s\n' "$committed" "$pending" "$untracked" | grep -q "^$1/"; }

if touched flai; then
  step="flai lint, vitest, the full Go tests, and smoke"
  echo "close-out: $step"
  "$ROOT/scripts/flai-test.sh"
else
  if touched template; then
    step="template render and check"
    echo "close-out: $step"
    "$ROOT/scripts/template-test.sh"
  fi
  step="markdown lint"
  echo "close-out: $step"
  "$ROOT/scripts/lint-md.sh"
  step="flai check --strict"
  echo "close-out: $step"
  "$ROOT/scripts/check.sh"
fi
if touched flaiover; then
  step="flaiover lint, types, and unit tests"
  echo "close-out: $step"
  "$ROOT/scripts/flaiover-test.sh"
fi

# The narrative is written in the main checkout, not on the story branch.
step="the narrative"
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

step="commit"
status="$(git status --porcelain)"
if [ -n "$status" ]; then
  [ $# -gt 0 ] || { git status --short >&2; echo "close-out: uncommitted changes and no message; pass the message as git commit options (-m, or -F with a file outside the worktree)" >&2; exit 1; }
  echo "close-out: commit"
  git add -A
  git diff --cached --stat
  git commit "$@"
fi
step="the sync check"
[ "$checked" = yes ] || synced
step="the clean worktree check"
status="$(git status --porcelain)"
[ -z "$status" ] || { git status --short >&2; echo "close-out: the worktree is not clean after the commit" >&2; exit 1; }
# The EXIT trap prints the passing line.
