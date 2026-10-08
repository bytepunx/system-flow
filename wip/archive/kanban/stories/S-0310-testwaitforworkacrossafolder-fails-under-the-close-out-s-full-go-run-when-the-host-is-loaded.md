---
id: S-0310
type: story
nature: remediation
title: TestWaitForWorkAcrossAFolder fails under the close-out's full Go run when the host is loaded
status: done
owner: alex
created: 2026-10-07T06:48:45Z
updated: 2026-10-08T00:00:19Z
transitions:
  - to: ready
    at: 2026-10-07T23:31:30Z
    by: orchestrator
  - to: in-progress
    at: 2026-10-07T23:31:39Z
    by: agent-S-0310
  - to: review
    at: 2026-10-07T23:59:45Z
    by: agent-S-0310
  - to: done
    at: 2026-10-08T00:00:19Z
    by: orchestrator
tags: [flai, mcp, tests]
touches: [flai/internal/mcpserver/folder_test.go, design/issues/I-0102-testwaitforworkacrossafolder-fails-under-the-close-out-s-full-go-run-when-the-host-is-loaded.md, design/issues/summary.md, design/issues/I-0079-testroundtriprepositoryitems-reads-the-live-main-checkout-and-fails-a-close-out-when-another-agent-edits-a-story-mid-run.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 1700
  turns:
    - day: 2026-10-07
      ceremony: 3
      hand_edits: 1
      work: 39
  models:
    - model: claude-opus-5-5
      input: 88
      output: 24161
      cache_read: 6081113
      cache_write: 196594
      cost: 3.2725
  strategic:
    - kind: orchestrator
      seconds: 481
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 82
          output: 1337
          cache_read: 17418325
          cache_write: 41291
          cost: 4.3024
        - model: claude-sonnet-5-5
          input: 6
          output: 43
          cache_read: 39077
          cache_write: 35713
          cost: 0.0645
cost_of_delay:
  inputs:
    time_lost_per_cycle: 10m
    by: flai
    at: 2026-10-07T06:48:45Z
  value: 25
  by: planner-S-0310
  at: 2026-10-07T23:30:59Z
forecast:
  duration: 1h
  delivery: 2026-10-08T06:15:00Z
  basis: "flai forecast's 12m (134 s per unit of size over 4 done remediation stories, size 5) raised to 1h: reproducing a load-dependent failure takes repeated runs of the mcpserver package, which took 178s under load in I-0102, before the fix and after it; delivery is flai's 05:25 behind 19 stories, moved by the 48m added."
  by: planner-S-0310
  at: 2026-10-07T23:30:59Z
finalized:
  by: orchestrator
  at: 2026-10-07T23:31:27Z
---
# S-0310 TestWaitForWorkAcrossAFolder fails under the close-out's full Go run when the host is loaded

## Goal

This story remediates [I-0102](../../../design/issues/I-0102-testwaitforworkacrossafolder-fails-under-the-close-out-s-full-go-run-when-the-host-is-loaded.md), "TestWaitForWorkAcrossAFolder fails under the close-out's full Go run when the host is loaded". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [x] The cause I-0102 describes no longer occurs, with a test that reproduces it where one fits
- [x] I-0102 is closed with `flai issue close I-0102 --reason` saying what fixed it

## Tasks
- T-1279 Surface the tool's error in TestWaitForWorkAcrossAFolder and confirm the cause under load
- T-1280 Write the story file atomically in readyStoryIn, with a test that reproduces the truncated read
- T-1281 Close I-0102 with what fixed it

## Notes

Cost of delay inputs set by flai from I-0102. time_lost_per_cycle 10m: 10m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-10-07T03:13:07Z, 0.1 days before this story; under one cycle counts as one).

### Planning

Proposed cause, from the instance and the code: `handed over: map[]` is a tool error, not a timeout. `readyStoryIn` in `flai/internal/mcpserver/folder_test.go` rewrites the story with `os.WriteFile`, which truncates first. The wait polls every 20ms, and a poll that reads the truncated file fails the parse in `workitem`'s `store.list`. The hold returns that error, and the test drops it. T-1279 confirms or refutes this before T-1280 fixes it.

Touches, all files, no folder touch:

| Touch | Source | Why |
|-------|--------|-----|
| `flai/internal/mcpserver/folder_test.go` | design (the issue's instance) | The failing test and its helper `readyStoryIn` |
| I-0102's file under `design/issues/` | design (criterion 2) | `flai issue close` writes it |
| `design/issues/summary.md` | layout | `flai issue close` regenerates it |

`flai touches suggest` listed `flai/internal/mcpserver/folder.go` (96% co-change) and `server.go`. They are left out: the fix proposed is in the test helper. Making the held waits tolerate a half-written file is proposed on the plan's thread as a story of its own.

Figures:

- Forecast: flai gave 12m (134 s per unit of size over 4 done remediation stories, size 5). Raised to 1h, because reproducing a failure that needs a loaded host takes repeated runs of the mcpserver package, 178s under load in I-0102, before the fix and after it. Delivery 2026-10-08T06:15:00Z is flai's 05:25, 19th in the pull order, moved by the 48m added.
- Cost of delay value: 25 USD a week, as `flai cod` gives it from the operator's input, 10m lost per 168h cycle at 150 USD an hour. It stands: one occurrence so far, each costing a re-run of the close-out.

### Accepted by the orchestrator

- Verified: 0118f54c7c7b5d624c38a7f4aa8169737fc3cda9
- At: 2026-10-08T00:00:19Z

Verdict: meets both criteria (verifier at 0118f54c7c7b5d624c38a7f4aa8169737fc3cda9; flai verify passed every step at that commit). The reproduction test characterizes the server's failure on a truncated story file; a race test of the helper itself does not fit, as the narrative records. The I-0079 bump is the close-out's own record.

- 1: flai/internal/mcpserver/folder_test.go
- 2: design/issues/I-0102-testwaitforworkacrossafolder-fails-under-the-close-out-s-full-go-run-when-the-host-is-loaded.md, design/issues/summary.md
