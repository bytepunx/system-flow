---
id: T-1075
type: task
nature: improvement
title: flai verify S-nnnn prints the verification result as text or --json and exits non-zero when a step fails
status: backlog
parent: S-0270
owner: alex
created: 2026-10-06T22:52:25Z
updated: 2026-10-06T22:52:25Z
transitions: []
stream: S-0270
tags: [flai, cli]
touches: [flai/cmd/verify.go, flai/cmd/verify_test.go, flai/cmd/root.go]
after: [T-1071]
---
# T-1075 flai verify S-nnnn prints the verification result as text or --json and exits non-zero when a step fails

## Work

Criterion 1: the command.

- Add `flai verify <story>` in `flai/cmd/verify.go` and register it in `flai/cmd/root.go`. It finds the story's worktree under `.flai-cache/worktrees/` and calls `verify.Verify` from T-1071.
- Print about twenty lines as text: one line per step with its state and duration, then the step's findings, then the notes from outside the story, then one last line with the outcome and the step it stopped at. Keep the last line's form the same as the close-out's, `close-out: S-nnnn ...`, or a `verify:` line like it, so that a reader of either sees the same thing. `--json` prints the result as data.
- Exit 0 when every step passed, and non-zero when one failed.
- `--record-issues` records the notes in `design/issues` through the same code `flai check --record-issues` uses (`recordOutside` in `flai/cmd/check.go`, called, not copied), so that the close-out can keep committing them.
- Waits for T-1071, whose `Verify` it calls.

## Done when

- [ ] `flai verify S-nnnn` and `--json` print the result with each step's state and duration, and the exit status follows the outcome.
- [ ] Tests cover a passing run, a failing step with its findings, an outside-the-story note, `--record-issues`, and a story with no worktree.
- [ ] `scripts/flai-test.sh` passes.

## Notes

Drafted by the planner.
