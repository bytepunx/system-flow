---
id: T-3428
type: task
nature: feature
title: flai verify reports each step and tier as it starts and ends, and the closeout package reads them, the resolve request, and the resolve job's outcome
status: done
parent: S-0804
owner: alex
created: 2026-11-16T23:19:11Z
updated: 2026-11-16T23:57:02Z
transitions:
  - to: ready
    at: 2026-11-16T23:20:38Z
    by: agent-S-0804
  - to: in-progress
    at: 2026-11-16T23:20:39Z
    by: agent-S-0804
  - to: done
    at: 2026-11-16T23:57:02Z
    by: agent-S-0804
stream: S-0804
tags: []
touches: [flai/internal/verify/story.go, flai/internal/verify/story_test.go, flai/cmd/verify.go, flai/cmd/verify_test.go, flai/internal/closeout/closeout.go, flai/internal/closeout/closeout_test.go, flai/internal/verify/run.go, flai/internal/verify/verify.go, flai/internal/verify/run_test.go]
usage:
  source: log
  seconds: 2183
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 126
      output: 53179
      cache_read: 9594096
      cache_write: 360507
      cost: 4.9942
---
# T-3428 flai verify reports each step and tier as it starts and ends, and the closeout package reads them, the resolve request, and the resolve job's outcome

## Work

First layer: it waits for no task. It owns every change to `flai/internal/closeout`, so that the serve job (the next task) builds on a settled contract.

1. `flai/internal/verify`: give `StoryOptions` (or what `Verify` takes) progress hooks, called as each step of a story's verify starts and ends (the rebase check, sync, narrative, check, each tier), reusing `TierOptions.Started` and `Ended` for the tiers. A skipped or reused tier ends with its outcome (skipped, reused, queued) without a start where it does not run.
2. `flai/cmd/verify.go`: when the environment sets `FLAI_VERIFY_PROGRESS=1`, `flai verify` writes one line to stderr as each step and tier starts and ends, never to stdout, whose last line the close-out script reads. The lines are exactly:
   - `verify: S-nnnn <RFC3339 UTC> step started <name>`
   - `verify: S-nnnn <RFC3339 UTC> step ended <outcome> <name>`
   - `verify: S-nnnn <RFC3339 UTC> tier started <name>`
   - `verify: S-nnnn <RFC3339 UTC> tier ended <outcome> <name>`
   where outcome is one word (passed, failed, skipped, reused, queued). Without the variable nothing changes.
3. `flai/internal/closeout`:
   - `KindResolve = "resolve"`; `Request.Kind` and `Job.Kind` (`json:"kind,omitempty"`, empty for a close-out).
   - `Job` gains the resolve job's end: `Outcome` (`OutcomeStayed = "stayed in review"`, `OutcomeSentBack = "sent back"`, empty otherwise, and for an interrupted job), `Thread` (the send-back thread), `Conflicts []string`, `Reason` (why it was sent back, or what it did).
   - `StepRebase = "rebase"`; the resolve job writes lines `resolve: S-nnnn <RFC3339 UTC> step started rebase`, `... step ended <outcome> rebase`, and the same for `close-out`, in the same shape as verify's.
   - `Progress` gains `Kind`, `Tier` (the tier running, "" when none), and `Steps []StepProgress` (`Name`, `Kind` step or tier, `Source` verify or resolve, `Started`, `Ended` zero while running, `Outcome`), parsed from those lines in the log in order; `Step` is `StepRebase` while the resolve job's rebase runs, else as `stepOf` says today.
   - `Status(cacheDir, story) (Status, bool, error)`: the story's last job, running or ended, with its progress parsed from the log and its log path, for the dashboard; `Running` stays as it is for its callers.
4. Tests for each: the hooks' order, the stderr lines with and without the variable, and the parse of a log holding a resolve job's and verify's lines.

## Done when

- `flai test` passes on the changed paths.
- `flai verify` prints nothing new without `FLAI_VERIFY_PROGRESS=1`, and its stdout is unchanged with it.
- `closeout.Status` and `closeout.Running` report a resolve job's step, tier, and steps from a log, tested.

## Notes
