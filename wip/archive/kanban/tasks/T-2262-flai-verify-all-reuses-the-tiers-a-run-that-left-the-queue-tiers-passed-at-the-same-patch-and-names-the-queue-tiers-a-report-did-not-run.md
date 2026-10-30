---
id: T-2262
type: task
nature: remediation
title: flai verify --all reuses the tiers a run that left the queue tiers passed at the same patch, and names the queue tiers a report did not run
status: done
parent: S-0529
owner: alex
created: 2026-10-28T18:16:46Z
updated: 2026-10-30T00:23:33Z
transitions:
  - to: ready
    at: 2026-10-30T00:18:05Z
    by: agent-S-0529
  - to: in-progress
    at: 2026-10-30T00:18:06Z
    by: agent-S-0529
  - to: done
    at: 2026-10-30T00:23:33Z
    by: agent-S-0529
stream: S-0529
tags: [cli]
touches: [flai/internal/verify/story.go, flai/internal/verify/story_test.go]
usage:
  source: log
  seconds: 327
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 22
      output: 6623
      cache_read: 1768252
      cache_write: 72992
      cost: 0.96
---
# T-2262 flai verify --all reuses the tiers a run that left the queue tiers passed at the same patch, and names the queue tiers a report did not run

## Work

In `flai/internal/verify/story.go`, the report's `fingerprint` (around `story.go:867-891`) holds the names of the tiers left to the queue (`story.go:392`), so `reusable` (`story.go:707-780`) reuses nothing when a `--all` run follows a close-out at the same patch: every tier runs again. Make the reuse rule pass over that difference: a run that leaves no tier to the queue reuses each tier a stored report at the same patch passed, and runs the tiers it left `queued`. Keep a fingerprint change from any other cause, such as a tier's command, refusing reuse.

Add a helper the acceptance reads, such as `QueueTiersNotRun(rep *Report, head string) []string`: the names of the steps a report left `queued`, or a reason when the report is missing, not at `head`, or did not pass.

It waits for nothing: the acceptance task builds on the helper.

## Done when

- A test in `story_test.go`, on `newVerifyFixture` and `fakeProc`, runs a verify that leaves the queue tiers, then one with `LeaveToQueue` false at the same patch, and asserts the second reuses the other tiers and runs only the queue tiers.
- A test asserts a changed tier command still refuses reuse.
- The helper is tested for a report at the head with queued steps, with none, at another commit, and failed.
- `flai test flai/internal/verify` passes.

## Notes
