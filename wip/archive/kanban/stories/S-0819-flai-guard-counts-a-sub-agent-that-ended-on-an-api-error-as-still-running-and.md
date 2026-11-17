---
id: S-0819
type: story
nature: remediation
title: flai guard counts a sub-agent that ended on an API error as still running and refuses wait_for_events for the close-out
status: done
owner: alex
created: 2026-11-17T02:47:36Z
updated: 2026-11-17T07:53:28Z
transitions:
  - to: ready
    at: 2026-11-17T03:37:47Z
    by: orchestrator
  - to: in-progress
    at: 2026-11-17T07:12:21Z
    by: agent-S-0819
  - to: review
    at: 2026-11-17T07:52:27Z
    by: system-flow
  - to: done
    at: 2026-11-17T07:53:28Z
    by: alex
tags: [flai, guard]
topics: [guard, delegation]
touches: [flai/internal/guard/subagents.go, flai/internal/guard/subagents_test.go, flai/internal/guard/claudecode.go, flai/internal/guard/claudecode_test.go, flai/internal/guard/transcript.go, flai/internal/guard/transcript_test.go, flai/cmd/guard.go, flai/cmd/guard_test.go, design/adrs, design/adrs/README.md, design/system/agent-context.md, design/system/flai-cli.md, docs/users/flai.md, design/issues/I-0376-flai-guard-counts-a-sub-agent-that-ended-on-an-api-error-as-still-running-and.md, design/issues/summary.md, flai/internal/guard/call.go, design/issues/I-0374-flai-guard-counts-a-sub-agent-that-ended-on-an-api-error-as-still-running-and.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 2423
  priced_by: harness
  turns:
    - day: 2026-11-17
      work: 81
  models:
    - model: claude-opus-5-5
      input: 264
      output: 66410
      cache_read: 23120281
      cache_write: 507271
      cost: 9.4145
  strategic:
    - kind: planner
      seconds: 422
      estimated: true
      models:
        - model: claude-haiku-5-5
          input: 28
          output: 21400
          cache_read: 891169
          cache_write: 125063
          cost: 0.091
        - model: claude-opus-5-5
          input: 58
          output: 23393
          cache_read: 4006229
          cache_write: 159623
          cost: 2.5463
    - kind: orchestrator
      seconds: 196
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 58
          output: 467
          cache_read: 2061293
          cache_write: 327298
          cost: 0.737
cost_of_delay:
  inputs:
    time_lost_per_cycle: 20m
    by: orchestrator
    at: 2026-11-17T07:39:43Z
  value: 25
  by: planner-S-0819
  at: 2026-11-17T03:34:18Z
forecast:
  duration: 1h30m
  delivery: 2026-11-17T06:42:00Z
  basis: "Its own forecast of 1h30m; 1st in the pull order with an in-progress limit of 3, behind S-0750, S-0805 and S-0817."
  by: flai
  at: 2026-11-17T04:46:05Z
finalized:
  by: orchestrator
  at: 2026-11-17T03:37:13Z
---
# S-0819 flai guard counts a sub-agent that ended on an API error as still running and refuses wait_for_events for the close-out

## Goal

This story remediates [I-0376](../../../design/issues/I-0376-flai-guard-counts-a-sub-agent-that-ended-on-an-api-error-as-still-running-and.md), "flai guard counts a sub-agent that ended on an API error as still running and refuses wait_for_events for the close-out". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [x] The cause I-0376 describes no longer occurs, with a test that reproduces it where one fits
- [x] I-0376 is closed with `flai issue close I-0376 --reason` saying what fixed it

## Tasks
- T-3504 Confirm from the S-0738 and S-0739 transcripts what Claude Code records for a foreground sub-agent that ends on an API error, and write the ADR that refines ADR-0092
- T-3505 The guard finds the foreground sub-agents whose launch has its result in the session transcript and drops them from the session's record
- T-3506 flai guard prunes the session's record from its transcript before it holds a story agent's wait_for_events, with a test that reproduces I-0376
- T-3507 The design and the user guide say the guard forgets a foreground sub-agent once its result is in the session transcript
- T-3508 Close I-0376 with flai issue close, saying the guard now prunes a foreground sub-agent whose result is in the session transcript

## Notes

Cost of delay inputs set by flai from I-0376. time_lost_per_cycle 10m: 10m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-11-15T01:05:45Z, 2.1 days before this story; under one cycle counts as one).

### Planning

Planned by planner-S-0819 on 2026-11-17.

Cause, as the code stands: `flai guard` records a sub-agent at `SubagentStart` and forgets it only at `SubagentStop` (`flai/internal/guard/subagents.go`, `RecordChange`). A foreground sub-agent that ends on an API error fires no `SubagentStop`, so the record keeps it, and `holdWaits` in `flai/cmd/guard.go` refuses every later `wait_for_events` of the session. The record keeps each entry's `started` time, but nothing reads it, and no rule expires an entry. ADR-0092's Consequences accepted a stale entry as costless. I-0376, and I-0374 with the same cause on S-0738, show that it is not.

Proposed fix, to be confirmed by T-3504 before it is built: a foreground sub-agent blocks its story's agent, so once the session's transcript holds the result of the Agent call that launched it, that sub-agent has ended, error or not. At a story agent's `wait_for_events`, the guard reads the transcript the hook input names (`transcript_path`). Under the record's lock, it drops from the session's record each sub-agent whose foreground launch has its result there, and decides on what is left. A transcript it cannot read leaves the record as it is. The `PreToolUse` matcher already covers `mcp__flai__.*`, so no hook configuration changes.

Layers: T-3504; then T-3505 and T-3507 together; then T-3506; then T-3508.

Touches:

- Layout: `flai/internal/guard/subagents.go` and `subagents_test.go` (the record and its prune), `claudecode.go` and `claudecode_test.go` (reading `transcript_path` from the hook input), `transcript.go` and `transcript_test.go` (new: finding the foreground Agent results in a session transcript), `flai/cmd/guard.go` and `guard_test.go` (`holdWaits` prunes before it decides; the end-to-end reproduction). The explorer found them from the refusal's text and ADR-0092.
- Design: `design/system/agent-context.md` § Waiting for a sub-agent states the record and the missing-stop sentence this changes; `design/system/flai-cli.md`'s `flai guard` row and `docs/users/flai.md` § Sub-agents describe the record.
- Design: the new ADR refining ADR-0092 and the index `design/adrs/README.md`. The folder touch `design/adrs/` is kept because the ADR's number and slug are given only when T-3504 writes it, so no task can name the file yet.
- Layout: `design/issues/I-0376-…md` and `design/issues/summary.md`, written by `flai issue close`.
- Declared: none; `flai touches suggest` had no path to start from.

Forecast: `flai forecast` gave 6m with no touches, and 34m once the 15 touches were set (118 s a unit of size over 78 remediation stories in the large band, size 17). Adjusted to 1h30m: T-3504 measures the harness's transcripts and writes an ADR before any code, which the median remediation story does not carry. Delivery is flai's from the pull order.

Cost of delay: `flai cod` gives 25.00 USD a week from the input 10m per 168h cycle at 150 USD an hour; it stands as the value. I-0374 records a second occurrence of the same cause, which the inputs do not count. Whether to raise them is asked on the plan's thread, since they are the operator's.

No `after`: the story reads no other story's output.
