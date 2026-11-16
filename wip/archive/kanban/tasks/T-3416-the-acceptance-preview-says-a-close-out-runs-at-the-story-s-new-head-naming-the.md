---
id: T-3416
type: task
nature: remediation
title: The acceptance preview says a close-out runs at the story's new head, naming the job and its step, beside a last report from before the head moved
status: done
parent: S-0799
owner: alex
created: 2026-11-16T01:45:56Z
updated: 2026-11-16T03:16:01Z
transitions:
  - to: ready
    at: 2026-11-16T02:43:41Z
    by: system-flow
  - to: in-progress
    at: 2026-11-16T03:06:20Z
    by: agent-S-0799
  - to: done
    at: 2026-11-16T03:16:01Z
    by: agent-S-0799
stream: S-0799
tags: [flai, acceptance, verify]
touches: [flai/internal/preview/accept.go, flai/internal/preview/accept_test.go]
after: [T-3414]
usage:
  source: log
  seconds: 581
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 25
      output: 8466
      cache_read: 1290467
      cache_write: 62558
      cost: 0.7847
---
# T-3416 The acceptance preview says a close-out runs at the story's new head, naming the job and its step, beside a last report from before the head moved

## Work

The acceptance preview (`flai/internal/preview/accept.go`) has no `verify_failed` blocker. Its `queue_tiers` blocker (`queueTiersGate`, lines 499-510, `queueTiersNotRun`, lines 517-535) and `verify_behind_main` (`verifyBehindMainGate`, lines 552-564) read the last report with no word of a close-out running at the new head.

- Read the last report with `verify.LastReportAt` from T-3414, giving the branch's head the preview already reads.
- Give the preview's result the running job as data, such as `verify_running` with the job's id, step, head, and start, set whenever a job runs at a head other than the report's, whatever the blockers are. Its text, beside `queue_tiers_why`, says the close-out runs at the head and that the preview is as of the last finished report.
- Keep every blocker as it is: a running job is not a verify that passed. The `queue_tiers` message takes T-3414's text, which names the job.
- S-0796 changes `verifyBehindMainGate` in this file. Sync before you start, and keep its gate as it lands.
- Test it in `accept_test.go`, beside `TestAHandAcceptanceWhileTheQueueBatchesListsTheQueueTiersBlocker` (line 264): a report at an old commit and a job running at the new head give `verify_running` with the job's step, and the `queue_tiers` message names the job. With no job there is no `verify_running`. The first case fails before the change: it reproduces I-0384.

This task waits for T-3414, whose reader it calls. It shares no file with the other tasks of its layer.

## Done when

- `flai accept S-nnnn --dry-run --json` and `accept.preview` carry `verify_running` while a close-out runs at a new head, and the blockers are unchanged.
- `flai test flai/internal/preview` passes.

## Notes
