---
id: T-3179
type: task
nature: feature
title: The close-out container and the checks container wait on an image job when their image is missing or stale
status: done
parent: S-0751
owner: alex
created: 2026-11-12T16:02:33Z
updated: 2026-11-13T15:43:00Z
transitions:
  - to: ready
    at: 2026-11-13T15:33:11Z
    by: agent-S-0751
  - to: in-progress
    at: 2026-11-13T15:33:12Z
    by: agent-S-0751
  - to: done
    at: 2026-11-13T15:43:00Z
    by: agent-S-0751
stream: S-0751
tags: [cli]
touches: [flai/internal/serve/verifyjob_container.go, flai/internal/serve/verifyjob_test.go, flai/internal/serve/checks.go, flai/internal/serve/checks_test.go, flai/internal/serve/kluster.go]
after: [T-3177, T-3178]
usage:
  source: log
  seconds: 588
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 38
      output: 19491
      cache_read: 3385662
      cache_write: 83009
      cost: 1.5534
---
# T-3179 The close-out container and the checks container wait on an image job when their image is missing or stale

## Work

`containerCloseOut.run` (`flai/internal/serve/verifyjob_container.go`, line 96, its image read at line 115) and the checks container (`checks.go`, line 250) both call `containerImage`. Make them wait on T-3177's job when the image is missing or stale, rather than failing the close-out or the check at once:

- the close-out's verdict waits, and the `verified` event comes after the job and the verify;
- a failed job fails the close-out as `failed` at its first step, naming the image job's log;
- the checks fail with the same cause.

Test the close-out's wait and its failure in `verifyjob_test.go`.

It waits for T-3177, whose runner it calls, and for T-3178, which changes `containerImage` to request the job, so this task builds on its new signature.

## Done when

- `go test ./internal/serve/ -run 'Verify|Check'` passes, under `-race` too.
- A close-out whose image was missing ends `passed` after the job in a test over a fake ensure.

## Notes
