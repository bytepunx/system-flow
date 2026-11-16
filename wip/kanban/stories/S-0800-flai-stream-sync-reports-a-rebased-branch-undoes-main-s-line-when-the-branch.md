---
id: S-0800
type: story
nature: remediation
title: flai stream sync reports a rebased branch undoes main's line when the branch removes a closing brace elsewhere in a Go file main added a closing brace to, and the merge queue then never batches the story
status: review
owner: alex
created: 2026-11-16T01:39:43Z
updated: 2026-11-16T02:42:42Z
transitions:
  - to: ready
    at: 2026-11-16T01:45:11Z
    by: orchestrator
  - to: in-progress
    at: 2026-11-16T01:45:20Z
    by: agent-S-0800
  - to: review
    at: 2026-11-16T02:42:42Z
    by: agent-S-0800
tags: [storygit, sync, merge-queue]
topics: [git, acceptance]
touches: [flai/internal/storygit/undone.go, flai/internal/storygit/undone_test.go, flai/internal/storygit/sync_test.go, design/adrs/README.md, design/system/workflow.md, design/system/flai-cli.md, docs/users/flai.md, design/issues/I-0385-flai-stream-sync-reports-a-rebased-branch-undoes-main-s-line-when-the-branch.md, design/issues/summary.md, design/adrs/0309-a-line-main-added-is-undone-only-when-the-rebased-branch-removes-it-at-the-line.md, flai/internal/hostapi/hostapi_test.go, flai/internal/hostapi/reads_test.go, flai/internal/mergequeue/eligible_test.go]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 3455
  priced_by: harness
  turns:
    - day: 2026-11-16
      ceremony: 4
      hand_edits: 1
      work: 53
  models:
    - model: claude-opus-5-5
      input: 230
      output: 59615
      cache_read: 16645810
      cache_write: 615606
      cost: 8.5032
  strategic:
    - kind: planner
      seconds: 63
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 13
          output: 2358
          cache_read: 1018838
          cache_write: 25172
          cost: 0.4411
    - kind: orchestrator
      seconds: 105
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 28
          output: 237
          cache_read: 1105110
          cache_write: 121358
          cost: 0.3735
        - model: claude-sonnet-5-5
          input: 8
          output: 830
          cache_read: 80699
          cache_write: 32579
          cost: 0.0906
cost_of_delay:
  inputs:
    time_lost_per_cycle: 45m
    by: flai
    at: 2026-11-16T01:39:43Z
  value: 112.5
  by: planner-S-0800
  at: 2026-11-16T01:43:54Z
forecast:
  duration: 50m
  delivery: 2026-11-16T05:01:00Z
  basis: "flai forecast's 28m (median 137 s per unit over 67 done medium remediation stories on claude-opus-5-5, size 12) raised to 50m, the 49m S-0767 took for the same shape of work: an ADR refining ADR-0155, the undone check's code and tests, a real-git sync test, three design and docs passages, and an issue close; delivery moved by the same 22m from flai forecast's 04:39Z, 10th in the pull order."
  by: planner-S-0800
  at: 2026-11-16T01:43:54Z
finalized:
  by: orchestrator
  at: 2026-11-16T01:44:44Z
---
# S-0800 flai stream sync reports a rebased branch undoes main's line when the branch removes a closing brace elsewhere in a Go file main added a closing brace to, and the merge queue then never batches the story

## Goal

This story remediates [I-0385](../../../design/issues/I-0385-flai-stream-sync-reports-a-rebased-branch-undoes-main-s-line-when-the-branch.md), "flai stream sync reports a rebased branch undoes main's line when the branch removes a closing brace elsewhere in a Go file main added a closing brace to, and the merge queue then never batches the story". The issue recommends this solution:

The undone check does not count a line of main's that carries no change of its own. Proposed, for the story to settle in an ADR refining ADR-0155: a line main added is undone only when the branch removes it where main added it, by comparing the hunks' positions in main's file rather than the file's lines as one multiset; or, more narrowly, a line that holds only brackets, braces, or white space is never counted alone. The revert ADR-0155 was written for (I-0124) must still be caught.

## Acceptance criteria
- [x] The cause I-0385 describes no longer occurs, with a test that reproduces it where one fits
- [x] I-0385 is closed with `flai issue close I-0385 --reason` saying what fixed it

## Tasks
- T-3405 Write the ADR refining ADR-0155 and ADR-0297 on which of main's lines a rebased branch undoes
- T-3407 Count a line main added as undone only by the ADR's rule in FindUndone and RecordedUndone, with unit tests
- T-3410 Reproduce I-0385 with a real-git sync whose branch removes a closing brace elsewhere in a Go file main added one to
- T-3412 Say in the design and the user guide which of main's lines a rebased branch undoes
- T-3413 Close I-0385 saying what fixed it

## Notes

Cost of delay inputs set by flai from I-0385. time_lost_per_cycle 45m: 45m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-11-16T01:39:37Z, 0 days before this story; under one cycle counts as one).

### Planning

Touches. The story declared none, so `flai touches suggest` was run from `flai/internal/storygit/undone.go`, the file holding `FindUndone` and `RecordedUndone`.

- `flai/internal/storygit/undone.go`, `undone_test.go` (design and layout): the undone check of ADR-0155 and ADR-0297 and its unit tests; co-change gives `undone_test.go` 67%.
- `flai/internal/storygit/sync_test.go` (co-change, 33%): the real-git reproduction, beside S-0767's I-0360 test.
- `flai/internal/storygit/sync.go` (co-change, 50%) is left out: the fix keeps `FindUndone`'s signature, so its caller does not change. So are `flai/internal/preview/accept.go` and `flai/internal/mergequeue/eligible.go`, which read `RecordedUndone` as it is. Task done widens the touches if one changes.
- `design/adrs/0309-a-line-main-added-is-undone-only-when-the-rebased-branch-removes-it-at-the-line.md` and `design/adrs/README.md` (design): the ADR T-3405 wrote, refining ADR-0155 and ADR-0297. The planner's `design/adrs/` folder touch was narrowed to this file once it was written.
- `design/system/workflow.md`, `design/system/flai-cli.md`, `docs/users/flai.md` (design): the three passages S-0767 changed for ADR-0297, which state the rule.
- `design/issues/I-0385-…md`, `design/issues/summary.md` (design): what `flai issue close` writes, for criterion 2.
- The rest of what co-change lists are hub files, each under 10%: none bears on the check.

Forecast 50m, delivery 2026-11-16T05:01:00Z. `flai forecast` gave 28m (137 s per unit of size over 67 done medium remediation stories on claude-opus-5-5, size 12) and delivery 04:39Z, 10th in the pull order. Raised to 50m because S-0767, the same shape of work for I-0360 (an ADR refining ADR-0155, the check's code and tests, a real-git sync test, three design and docs passages, an issue close), was forecast 40m and took 49m. Delivery moved by the same 22m.

Cost of delay 112.50 USD a week. `flai cod` works it from the operator's input, 45m lost per 168h cycle, at 150 USD an hour. It stands: the input is the operator's, and one occurrence so far is what it counts. The issue's record also keeps S-0739 out of the merge queue until the fix, which the input does not price.

After: none. The story reads no output of another open story.
