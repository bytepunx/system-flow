---
id: T-3415
type: task
nature: remediation
title: flai queue names the close-out job running at a story's new head, and says a verify_failed report is from before the head moved
status: done
parent: S-0799
owner: alex
created: 2026-11-16T01:45:44Z
updated: 2026-11-16T03:15:53Z
transitions:
  - to: ready
    at: 2026-11-16T02:43:41Z
    by: system-flow
  - to: in-progress
    at: 2026-11-16T03:06:20Z
    by: agent-S-0799
  - to: done
    at: 2026-11-16T03:15:53Z
    by: agent-S-0799
stream: S-0799
tags: [flai, merge-queue, verify]
touches: [flai/internal/mergequeue/eligible.go, flai/internal/mergequeue/eligible_test.go, flai/cmd/queue_test.go]
after: [T-3414]
usage:
  source: log
  seconds: 573
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 26
      output: 8971
      cache_read: 1367447
      cache_write: 66289
      cost: 0.8315
---
# T-3415 flai queue names the close-out job running at a story's new head, and says a verify_failed report is from before the head moved

## Work

`readLastVerify` (`flai/internal/mergequeue/eligible.go`, lines 404-419) gives a failed report the reason `verify_failed` with the report's step and commit, and never looks at the head. This is what I-0384's operator read on S-0739 and S-0747 as the sync having done nothing.

- Read the story's last report with `verify.LastReportAt` from T-3414, giving the branch's head. `ComputeWith` (lines 208-296) reads heads only when some report passed (`anyPassed`, lines 242-253, S-0698). Widen that to the stories whose last report failed too, and keep the S-0698 saving for the rest.
- While a close-out job runs at another head, give the reason a new kind, `verify_running` (beside `KindVerifyFailed`, line 36), in place of `verify_failed` or `verify_stale`. Its text names the job, the step it is on, and the head, and says the story is eligible once it passes. Give `Reason` (lines 149-158) the job as data, such as `job` with `id`, `step`, `head`, and `started`, so the dashboard need not parse the text.
- A `verify_failed` whose report is `moved`, with no job running, says the report is from before the head moved. Keep its kind, and name both commits and the command that verifies the new head.
- Print the new kind and text in `flai/cmd/queue.go` (lines 304-305).
- Test it: in `eligible_test.go`, extend `TestEachOtherStoryInReviewHasEveryReasonItIsNotEligible` (line 222) with a failed report at an old commit and a job running at the new head, which gives `verify_running` with the job's id and step, and the same with no job, which gives `verify_failed` saying the report is from before the head moved. The first case fails before the change: it reproduces I-0384. In `flai/cmd/queue_test.go`, `TestQueueListsTheEligibleThenEveryReasonOfTheOthers` (line 64) prints the new reason.

This task waits for T-3414, whose reader it calls. It runs beside T-3416 and T-3417, which change none of its files.

## Done when

- `flai queue` and `flai queue --json` list a story whose close-out runs at a new head as `verify_running`, naming the job and its step, and an old failed report as from before the head moved.
- `flai test flai/internal/mergequeue flai/cmd/queue.go` passes.

## Notes
