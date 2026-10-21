---
id: S-0457
type: story
nature: improvement
title: The merge queue's send-back thread quotes only the first line of a failing plain tier's finding and keeps no log of the batch, so the step that failed cannot be told
status: in-progress
owner: alex
created: 2026-10-21T02:13:46Z
updated: 2026-10-21T07:02:49Z
transitions:
  - to: ready
    at: 2026-10-21T06:12:25Z
    by: orchestrator
  - to: in-progress
    at: 2026-10-21T06:12:35Z
    by: agent-S-0457
  - to: review
    at: 2026-10-21T07:01:56Z
    by: alex
  - to: in-progress
    at: 2026-10-21T07:02:37Z
    by: alex
tags: [flai]
topics: [cli]
touches: [flai/internal/verify/text.go, flai/internal/verify/text_test.go, flai/cmd/verify.go, flai/cmd/verify_test.go, flai/internal/verify/run.go, flai/internal/verify/run_test.go, flai/internal/verify/story.go, flai/internal/mergequeue/batch.go, flai/internal/mergequeue/batch_test.go, flai/internal/mergequeue/state.go, flai/internal/mergequeue/bisect.go, flai/internal/mergequeue/bisect_test.go, design/system/flai-cli.md, docs/users/flai.md, docs/operators/index.md, docs/operators/runbooks/backup.md, design/issues/I-0172-the-merge-queue-s-send-back-thread-quotes-only-the-first-line-of-a-failing-plain-tier-s-finding-and-keeps-no-log-of-the-batch-so-the-step-that-failed-cannot-be-told.md, design/issues/summary.md, flai/internal/verify/verify.go, flai/internal/mergequeue/land_test.go, flai/internal/serve/mergequeue_test.go]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 2662
  estimated: true
  priced_by: estimate
  turns:
    - day: 2026-10-21
      ceremony: 1
      hand_edits: 1
      work: 35
  models:
    - model: claude-opus-5-5
      input: 262
      output: 1631
      cache_read: 10722082
      cache_write: 588843
      cost: 4.1801
  strategic:
    - kind: orchestrator
      seconds: 470
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 36
          output: 944
          cache_read: 16195196
          cache_write: 28579
          cost: 4.0672
cost_of_delay:
  inputs:
    time_lost_per_cycle: 30m
    by: flai
    at: 2026-10-21T02:13:46Z
  value: 75
  by: planner-S-0457
  at: 2026-10-21T06:11:22Z
forecast:
  duration: 1h
  delivery: 2026-10-21T23:05:00Z
  basis: "flai forecast's 30m (89 s a unit over 81 done large improvement stories, size 20) doubled for six tasks in five layers and a Verifier signature change that reaches every queue test stub; delivery moved by the same 30m."
  by: planner-S-0457
  at: 2026-10-21T06:11:22Z
finalized:
  by: orchestrator
  at: 2026-10-21T06:12:20Z
---
# S-0457 The merge queue's send-back thread quotes only the first line of a failing plain tier's finding and keeps no log of the batch, so the step that failed cannot be told

## Goal

This story remediates [I-0172](../../../design/issues/I-0172-the-merge-queue-s-send-back-thread-quotes-only-the-first-line-of-a-failing-plain-tier-s-finding-and-keeps-no-log-of-the-batch-so-the-step-that-failed-cannot-be-told.md), "The merge queue's send-back thread quotes only the first line of a failing plain tier's finding and keeps no log of the batch, so the step that failed cannot be told". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [x] The cause I-0172 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0172 is closed with `flai issue close I-0172 --reason` saying what fixed it

## Tasks
- T-1947 Move flai verify's step findings formatter into internal/verify so the merge queue can quote a failing tier as flai verify prints it
- T-1948 Let verify.Run write each tier's whole output to a log file under a folder its options name
- T-1949 Keep each merge-queue batch's and bisection step's tier output under .flai-cache/merge-queue/batches and record it on the batch's result
- T-1950 Quote the failing tier's findings in full, as flai verify prints them, and name the batch's log in the merge queue's send-back thread
- T-1951 Document the merge queue's kept batch logs and the findings its send-back thread quotes, for builders, users, and operators
- T-1952 Close I-0172 with the reason that the send-back thread now quotes the failing tier's findings in full and the queue keeps each batch's output

## Notes

Cost of delay inputs set by flai from I-0172. time_lost_per_cycle 30m: 30m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-10-21T01:54:16Z, 0 days before this story; under one cycle counts as one).

### Planning

Proposed fix, from the instance on S-0454 (TH-0545): the cause is `failedSummary` in `flai/internal/mergequeue/batch.go`, which keeps only the first line of each finding's `Text()`; for a plain tier that is the script's first output line. `SendBack` and `failingText` in `flai/internal/mergequeue/bisect.go` quote that summary, and nothing keeps the tiers' output once the batch worktree is cleaned. The fix keeps each batch's and bisection step's tier output under `.flai-cache/merge-queue/batches/<batch id>/` (T-1948, T-1949) and quotes the failing step's findings in the thread as `flai verify` prints them (T-1947, T-1950).

Touches, all files, no folder touch:

| Touch | Source |
|-------|--------|
| `flai/internal/mergequeue/batch.go`, `bisect.go`, `state.go` | layout: the summary, the send-back thread, and `Result` (explorer) |
| `flai/internal/mergequeue/batch_test.go`, `bisect_test.go` | co-change (38%, 15%) and layout: the existing summary and send-back tests |
| `flai/internal/verify/run.go`, `run_test.go`, `story.go` | layout: where a tier's output lives; co-change 31% and 23% |
| `flai/internal/verify/text.go`, `text_test.go`, `flai/cmd/verify.go`, `flai/cmd/verify_test.go` | layout: `stepFindings` moves out of package `cmd`, which `mergequeue` cannot import; co-change 15% |
| `design/system/flai-cli.md` | design: `### The merge queue pass` describes the state and the thread |
| `docs/users/flai.md`, `docs/operators/index.md`, `docs/operators/runbooks/backup.md` | design: the user, operator, and backup texts on the queue and `.flai-cache` |
| the I-0172 issue file, `design/issues/summary.md` | criterion 2: `flai issue close` writes both |

`flai touches suggest` found nothing at first, since the story declared no touches; it was run again from the four code files the explorer named. `flai/internal/serve/mergequeue.go` is left out: it calls `CheckoutVerifier(git, version)`, whose signature need not change.

Forecast: `flai forecast` gave 30m (89 s a unit over 81 done large improvement stories, size 20 from 2 criteria and 18 touches). Set to 1h: six tasks in five layers, each closed with its tests, and a `Verifier` signature change that reaches every stub in the queue's tests; delivery moved by the same 30m, to 2026-10-21T23:05:00Z.

Cost of delay: `flai cod` gives 75.00 USD a week, 30m lost per 168h cycle at 150 USD an hour from flai's inputs. It stands: one occurrence so far, and the queue runs only while its host action is on.
- 2026-10-21T07:02:37Z: moved to in-progress: unchecked criteria
