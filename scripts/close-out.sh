#!/usr/bin/env sh
# Close out a story before review: flai verify, then the commit, then the
# sync check, stopping at the first step that fails (I-0012). When main moved
# during the run, it syncs and runs them again, at most three rounds
# (ADR-0156).
# Usage: scripts/close-out.sh S-nnnn [git commit options: -m, or -F with a file
# outside the worktree, which git add -A would otherwise commit]
# The cycle before it (git.md): close each task when it is done with
# flai task done, which commits, syncs, and checks, resolve what a stopped
# sync lists, run the task's tests, and close any fix the same way; then
# commit what is outstanding and sync again before this.
# Run it in the story's worktree. flai verify runs every check (S-0270), so
# the close-out and flai verify never disagree: no rebase left unfinished, the
# branch contains the main branch, the narrative, flai check --strict scoped
# to the story, then the tiers of system-flow.yaml's tests that the branch's
# changes select, with CLOSE_OUT_STORY set to the story in their environment
# (which scripts/check.sh in smoke and TestMonorepoIsClean in integration read
# to scope themselves, ADR-0085). --record-issues records the check's findings
# outside the story in design/issues, which the commit step commits (S-0249).
# While the merge queue can batch for the project (the merge-queue host action
# on for it in the configuration flai reads, which scripts/flai.sh sets to
# .flai-cache/config.json, accept_reviews on, acceptance local), flai verify
# leaves the tiers marked queue: true, here integration and smoke, to the
# queue's batch, which runs them, and lists each as queued (ADR-0174). Turn
# the host action on in that configuration with scripts/flai.sh serve enable
# merge-queue: flai serve's own configuration does not count here. While it is
# off, flai verify runs them and says why on a line before its last (S-0466).
set -eu
. "$(dirname "$0")/env.sh"
cd "$ROOT"

[ $# -ge 1 ] || { echo "usage: scripts/close-out.sh S-nnnn [git commit options]" >&2; exit 2; }
story="$1"
shift

# Every run ends with one line on stdout naming the story, the outcome, and
# the step it stopped at, so nobody runs it again to learn why it stopped
# (S-0266), and the round it ended on when it ran more than one (ADR-0156).
# Set once the story is known: a usage error prints only the usage.
# The trap reads $? before anything else can change it, and exits with it.
# A pass or a stop at flai verify carries the tiers verify reused (S-0341).
step="the worktree check"
hint=""
signal=""
reused=""
tmp=""
round=1
rounds=3
finish() {
  code="$1"
  [ -z "$tmp" ] || rm -rf "$tmp" || :
  [ "$step" = "flai verify" ] || [ "$code" -eq 0 ] || reused=""
  in=""
  [ "$round" -eq 1 ] || in=" in round $round of $rounds"
  if [ -n "$signal" ]; then
    echo "close-out: $story stopped at $step$in ($signal)"
  elif [ "$code" -eq 0 ]; then
    echo "close-out: $story passed every step$in; ready to move to review$reused"
  else
    echo "close-out: $story stopped at $step$in (exit $code)$reused$hint"
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

# run runs a flai verify: its output shows as it comes and is kept, to read
# its last line into $last, and its status goes through a file into $rc,
# which the pipe would otherwise hide.
tmp="$(mktemp -d)"
run() {
  { rc=0; "$ROOT/scripts/flai.sh" verify "$@" || rc=$?; echo "$rc" >"$tmp/status"; } | tee "$tmp/out"
  rc="$(cat "$tmp/status")"
  last="$(tail -n 1 "$tmp/out")"
}

while :; do
  step="flai verify"
  hint="; flai verify's last line above names the step it stopped at, or why it could not run"
  reused=""
  echo "close-out: $step $story --record-issues"
  # Without --fresh, a second run at an unchanged head and base resumes at the
  # tier that failed, and one after a clean rebase reuses the tiers main's new
  # paths do not select (ADR-0156); verify's last line then ends
  # "(reused N tiers from <time>)".
  run "$story" --record-issues
  case "$last" in
    "verify: $story "*" (reused "[0-9]*" tier"*")") reused=" (${last##*" ("}" ;;
  esac
  [ "$rc" -eq 0 ] || exit "$rc"
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
  # which may have moved while the tiers ran: flai verify --sync-only asks again
  # at the new head, and passes over commits on the main branch that change only
  # wip paths the branch does not change, such as flai's replans (ADR-0135).
  step="the sync check"
  hint="; run flai stream sync $story, resolve what it reports, and run this again"
  echo "close-out: flai verify $story --sync-only"
  run "$story" --sync-only
  [ "$rc" -ne 0 ] || break
  # Main moved, such as by an acceptance, when it stopped at the sync step; at
  # any other step, or after the last round, the close-out stops here. Else it
  # syncs and runs every step again (I-0136, ADR-0156).
  case "$last" in
    "verify: $story stopped at sync ("*) ;;
    *) exit "$rc" ;;
  esac
  [ "$round" -lt "$rounds" ] || exit "$rc"
  round=$((round + 1))
  step="flai stream sync"
  hint="; resolve each conflicting path flai stream sync lists above in the worktree, git add it, run git rebase --continue there, and run this again"
  echo "close-out: main moved during the run; flai stream sync $story (round $round of $rounds)"
  "$ROOT/scripts/flai.sh" stream sync "$story"
done
hint=""

step="the clean worktree check"
status="$(git status --porcelain)"
[ -z "$status" ] || { git status --short >&2; echo "close-out: the worktree is not clean after the commit" >&2; exit 1; }
# The EXIT trap prints the passing line.
