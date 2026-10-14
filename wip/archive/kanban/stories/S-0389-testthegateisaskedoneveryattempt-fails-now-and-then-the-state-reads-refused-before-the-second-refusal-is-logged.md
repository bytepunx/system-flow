---
id: S-0389
type: story
nature: remediation
title: "TestTheGateIsAskedOnEveryAttempt fails now and then: the state reads refused before the second refusal is logged"
status: done
owner: alex
created: 2026-10-14T23:09:14Z
updated: 2026-10-14T23:36:28Z
transitions:
  - to: ready
    at: 2026-10-14T23:12:03Z
    by: orchestrator
  - to: in-progress
    at: 2026-10-14T23:12:11Z
    by: agent-S-0389
  - to: review
    at: 2026-10-14T23:35:42Z
    by: agent-S-0389
  - to: done
    at: 2026-10-14T23:36:28Z
    by: orchestrator
tags: [cli]
touches: [flai/internal/channel/channel.go, flai/internal/channel/channel_test.go, design/issues/I-0146-testthegateisaskedoneveryattempt-fails-now-and-then-the-state-reads-refused-before-the-second-refusal-is-logged.md, design/issues/summary.md, design/issues/I-0141-the-mcp-test-tool-of-the-installed-flai-1-40-1-refuses-main-s-manifest-so-a-story-s-agent-runs-flai-test-from-the-tree.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 1425
  priced_by: harness
  turns:
    - day: 2026-10-14
      ceremony: 3
      test_runs: 3
      hand_edits: 3
      work: 51
  models:
    - model: claude-opus-5-5
      input: 126
      output: 21942
      cache_read: 8454296
      cache_write: 157954
      cost: 3.3938
  strategic:
    - kind: orchestrator
      seconds: 1741
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 182
          output: 3378
          cache_read: 37971307
          cache_write: 46361
          cost: 9.5309
        - model: claude-sonnet-5-5
          input: 4
          output: 40
          cache_read: 12244
          cache_write: 28440
          cost: 0.0333
cost_of_delay:
  inputs:
    time_lost_per_cycle: 16m
    by: flai
    at: 2026-10-14T23:09:14Z
  value: 40
  by: planner-S-0389
  at: 2026-10-14T23:11:34Z
forecast:
  duration: 30m
  delivery: 2026-10-15T01:43:00Z
  basis: "flai's 14m (131 s per unit over 8 done remediation stories, size 6), raised to 30m for a race reproduction, a -race loop, and a close-out whose integration tier holds the test; delivery is flai's 9th-in-pull-order date moved by the 16m added"
  by: planner-S-0389
  at: 2026-10-14T23:11:34Z
finalized:
  by: orchestrator
  at: 2026-10-14T23:11:59Z
---
# S-0389 TestTheGateIsAskedOnEveryAttempt fails now and then: the state reads refused before the second refusal is logged

## Goal

This story remediates [I-0146](../../../design/issues/I-0146-testthegateisaskedoneveryattempt-fails-now-and-then-the-state-reads-refused-before-the-second-refusal-is-logged.md), "TestTheGateIsAskedOnEveryAttempt fails now and then: the state reads refused before the second refusal is logged". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [x] The cause I-0146 describes no longer occurs, with a test that reproduces it where one fits
- [x] I-0146 is closed with `flai issue close I-0146 --reason` saying what fixed it

## Tasks
- T-1587 The channel client logs an attempt's end before it publishes the state, with a test that makes the race certain
- T-1588 Close I-0146 as fixed by the channel client logging before it publishes the state
- T-1589 The channel client logs a proven connection before it publishes the state, with a test that makes the race certain

## Notes

Cost of delay inputs set by flai from I-0146. time_lost_per_cycle 16m: 8m per occurrence × 2 occurrences ÷ 1 cycle of 168h (first reported 2026-10-13T19:01:05Z, 1.2 days before this story; under one cycle counts as one).

### Planning

Proposed solution, from the issue's first instance: `Client.Run` in `flai/internal/channel/channel.go` calls `setState` with the attempt's refusal before `c.Logger.Error` writes it (lines 417 to 468). The fix logs each attempt's end before the state publishes it, so any reader that sees the state refused finds the line logged. The issue's other remedy, polling the log in the test, would leave the order wrong for every other reader and is not taken.

Touches, all files, no folder touch:

| Touch | Source |
|-------|--------|
| `flai/internal/channel/channel.go` | layout: the `Run` loop that sets the state before it logs |
| `flai/internal/channel/channel_test.go` | layout: the flaky test and the new test that makes the race certain |
| `design/issues/I-0146-…md` | design: criterion 2 closes it |
| `design/issues/summary.md` | design: `flai issue close` regenerates it |

`flai touches suggest` from the two channel files listed serve, hostapi, and flaiover files at 13 to 47%. Those commits were the features that built the channel; none bears on the order of a log line, so none was kept. No user or design document describes the order, so none is touched.

Forecast 30m, from flai's 14m (131 s per unit of size over 8 done remediation stories, size 6). Raised because the story writes a test that forces the race, loops it under `-race`, and its close-out runs the integration tier that holds the test. Delivery is flai's 2026-10-15T01:27:00Z, 9th in the pull order, moved by the 16m added.

Cost of delay 40 USD a week, as `flai cod` gives it from the operator's 16m lost per 168h cycle at 150 USD an hour. It stands: the two instances each cost a close-out run, which the input already counts.

Tasks: T-1587 (layer 1, the reorder and the test), then T-1588 (layer 2, closing I-0146, after T-1587).

### Accepted by the orchestrator

- Verified: adddb8fce3aaf26b5c3f8b6750bfb860a0a91ceb
- At: 2026-10-14T23:36:28Z

Verdict: meets both criteria; flai verify passed every step at adddb8fc; every changed path is within the story's touches, none under .claude/. Run and serveOnce log before setState publishes, and two slow-writer tests make each race certain.

- 1: flai/internal/channel/channel.go, flai/internal/channel/channel_test.go
- 2: design/issues/I-0146-testthegateisaskedoneveryattempt-fails-now-and-then-the-state-reads-refused-before-the-second-refusal-is-logged.md, design/issues/summary.md
