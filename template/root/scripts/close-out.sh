#!/usr/bin/env sh
# Close out a story before review: flai verify, then the commit, then the
# sync check, stopping at the first step that fails.
# Usage: scripts/close-out.sh S-nnnn [git commit options: -m, or -F with a file
# outside the worktree, which git add -A would otherwise commit]
# The cycle before it (git.md): close each task when it is done with
# flai task done, which commits, syncs, and checks, resolve what a stopped
# sync lists, run the task's tests, and close any fix the same way; then
# commit what is outstanding and sync again before this.
# Run it in the story's worktree. flai verify runs every check, so the
# close-out and flai verify never disagree: no rebase left unfinished, the
# branch contains the main branch, the narrative, flai check --strict scoped
# to the story, then the tiers of system-flow.yaml's tests that the branch's
# changes select, with CLOSE_OUT_STORY set to the story in their environment.
# --record-issues records the check's findings outside the story in
# design/issues, which the commit step commits. Add this project's other
# checks as tiers in system-flow.yaml's tests, not as steps here.
set -eu
cd "$(dirname "$0")/.."

[ $# -ge 1 ] || { echo "usage: scripts/close-out.sh S-nnnn [git commit options]" >&2; exit 2; }
story="$1"
shift

# Every run ends with one line on stdout naming the story, the outcome, and
# the step it stopped at, so nobody runs it again to learn why it stopped.
# Set once the story is known: a usage error prints only the usage. The trap
# reads $? before anything else can change it, and exits with it.
step="the worktree check"
hint=""
signal=""
finish() {
  code="$1"
  if [ -n "$signal" ]; then
    echo "close-out: $story stopped at $step ($signal)"
  elif [ "$code" -eq 0 ]; then
    echo "close-out: $story passed every step; ready to move to review"
  else
    echo "close-out: $story stopped at $step (exit $code)$hint"
  fi
  exit "$code"
}
trap 'finish "$?"' EXIT
# An interrupt exits, which runs the EXIT trap once, with the signal's status.
trap 'signal=interrupted; exit 130' INT
trap 'signal=terminated; exit 143' TERM

# flai verify runs in the story's worktree whatever the checkout it is called
# from, and the commit below is made in this one: they must be the same.
branch="$(git rev-parse --abbrev-ref HEAD)"
[ "$branch" = "story/$story" ] || { echo "close-out: this checkout is on $branch, not story/$story; run this in $story's worktree" >&2; exit 1; }

step="flai verify"
hint="; flai verify's last line above names the step it stopped at, or why it could not run"
echo "close-out: $step $story --record-issues"
flai verify "$story" --record-issues
hint=""

step="commit"
status="$(git status --porcelain)"
if [ -n "$status" ]; then
  [ $# -gt 0 ] || { git status --short >&2; echo "close-out: uncommitted changes and no message; pass the message as git commit options (-m, or -F with a file outside the worktree)" >&2; exit 1; }
  echo "close-out: commit"
  git add -A
  git diff --cached --stat
  git commit "$@"
fi

# The branch still contains the main branch flai verify checked it against,
# the branch the main checkout has, which may have moved while the tiers ran.
# The main checkout shares its git directory with the story's worktree.
step="the sync check"
common="$(git rev-parse --git-common-dir)"
main="$(cd "$common/.." && pwd)"
base="$(git -C "$main" rev-parse --abbrev-ref HEAD)"
[ "$base" != HEAD ] || { echo "close-out: the main checkout $main is not on a branch" >&2; exit 1; }
rc=0
git merge-base --is-ancestor "$base" HEAD || rc=$?
case "$rc" in
  0) echo "close-out: the branch contains $base" ;;
  1) echo "close-out: the branch does not contain $base; run flai stream sync $story, resolve what it reports, and run this again" >&2; exit 1 ;;
  *) echo "close-out: cannot tell whether the branch contains $base" >&2; exit 1 ;;
esac

step="the clean worktree check"
status="$(git status --porcelain)"
[ -z "$status" ] || { git status --short >&2; echo "close-out: the worktree is not clean after the commit" >&2; exit 1; }
# The EXIT trap prints the passing line.
