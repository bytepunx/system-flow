---
id: T-2031
type: task
nature: improvement
title: flai close-out runs the close-out's rounds in Go, verify, commit, sync check, stream sync, at most three rounds, with the script's last line
status: cancelled
parent: S-0466
owner: alex
created: 2026-10-22T20:57:59Z
updated: 2026-10-23T03:22:00Z
transitions:
  - to: cancelled
    at: 2026-10-23T03:22:00Z
    by: agent-S-0466
stream: S-0466
tags: [cli]
touches: [flai/cmd/closeout.go, flai/cmd/closeout_test.go, flai/internal/closeout/closeout.go, flai/internal/closeout/closeout_test.go]
after: [T-2030]
---
# T-2031 flai close-out runs the close-out's rounds in Go, verify, commit, sync check, stream sync, at most three rounds, with the script's last line

## Work

Add `flai close-out <story> [-m <message> | -F <file>] [--json]` in `flai/cmd/closeout.go`, on a package `flai/internal/closeout` that runs, in the story's worktree, what `scripts/close-out.sh` runs today, step by step and in the same order: the branch check (`story/<S-nnnn>`), `verify.Verify` with `RecordIssues` as `flai verify --record-issues` runs it, the commit of what is outstanding (`git add -A`, `git commit` with the message or file given, refused with the script's words when there is something to commit and no message), `verify.SyncOnly`, and, when that stops at `sync`, `storygit.Sync` as `flai stream sync` runs it, then the next round, at most three rounds (ADR-0156), and the clean worktree check at the end. Keep the script's hints per step and its last line, `close-out: <story> passed every step[ in round r of 3]; ready to move to review[ (reused N tiers from <time>)]` or `close-out: <story> stopped at <step>[ in round r of 3] (exit n)<hint>`, so that `docs/users/flai.md` and the agents' prompts keep reading it the same way. Print each step as it starts, as the script does, and stream the verify's own output. This task runs the rounds in the foreground; T-2033 detaches them.

Waits for T-2030, which settles the command and its flags.

## Done when

- `flai close-out S-nnnn -m "<message>"` in a fixture worktree, as the cmd tests build one for `flai verify`, passes every step and prints the script's passing line; a run whose verify stops prints the stopped line with the step and the hint; a run with uncommitted changes and no message stops at the commit with the script's words.
- A test where main gains a commit outside `wip/` between the verify and the sync check shows the second round: `flai stream sync` ran, verify ran again, and the last line names round 2 of 3; a test where main moves in every round stops at the sync check in round 3.
- `flai test` on the changed paths passes and `flai check --strict` is clean.

## Notes

The second layer. `scripts/close-out.sh` is not changed here: T-2046 makes it call this command, beside T-2033, which detaches the run.
- 2026-10-23T03:22:00Z: moved to cancelled: TH-0594 settled on option 3: S-0466 makes the close-out see the merge queue; the detached close-out is left to a backlog story from I-0148
