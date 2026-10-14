---
id: TH-0452
title: "S-0389 plan: the channel client logs an attempt's end before it publishes the state, with a test that forces the race; two tasks in two layers"
anchor:
  path: wip/kanban/stories/S-0389-testthegateisaskedoneveryattempt-fails-now-and-then-the-state-reads-refused-before-the-second-refusal-is-logged.md
  item: S-0389
status: resolved
participants: [planner-S-0389, orchestrator]
created: 2026-10-14T23:11:49Z
updated: 2026-10-14T23:11:57Z
---

# TH-0452 S-0389 plan: the channel client logs an attempt's end before it publishes the state, with a test that forces the race; two tasks in two layers

On wip/kanban/stories/S-0389-testthegateisaskedoneveryattempt-fails-now-and-then-the-state-reads-refused-before-the-second-refusal-is-logged.md.

## Entries

### 2026-10-14T23:11:49Z planner-S-0389
**Plan:** two tasks in two layers. The channel client will write each attempt's log line before it publishes that attempt's state. A new test forces the race. Then I-0146 is closed.

## Cause

`Client.Run` in `flai/internal/channel/channel.go` calls `setState` with the refusal, then calls `c.Logger.Error`. `TestTheGateIsAskedOnEveryAttempt` waits until the state reads refused, then reads the log at once (`channel_test.go:883`). It can read the log before the second line is written.

## Tasks

| Layer | Task | Touches | After |
|-------|------|---------|-------|
| 1 | T-1587: reorder `Run` so each attempt's log line comes before its state; add a test whose logger writes through a slow writer, so the race is certain on the old code | `flai/internal/channel/channel.go`, `flai/internal/channel/channel_test.go` | none |
| 2 | T-1588: `flai issue close I-0146 --reason` naming the reorder and the test | `design/issues/I-0146-…md`, `design/issues/summary.md` | T-1587 |

## Figures

- **Forecast:** 30m, delivery 2026-10-15T01:43:00Z. flai gave 14m. I raised it for the race reproduction, a `-race -count=50` loop, and a close-out whose integration tier holds the test.
- **Cost of delay:** 40 USD a week, flai's value from the 16m input. Unchanged.

## Assumptions

- **The fix goes in the client, not the test.** The issue's other remedy polls the log in the test. That would leave the order wrong for any other reader of `State`.
- **`proven` can be read before the log.** It is the attempt's own `Since`, and nothing else changes it once `attempt` returns.
- **No document changes.** No design or user document describes the order of the state and the log line.
- **No split, merge, or drop to propose.**

### 2026-10-14T23:11:55Z orchestrator
Approved as drafted.

- **The tasks cover both criteria.**
  - T-1587 removes the cause: `Run` logs each attempt before `setState` publishes it. A test with a slow writer makes the race certain on the old code (1).
  - T-1588 closes I-0146 (2).
- **The fix belongs in the client.** Fixing it there, not by polling in the test, keeps the order right for every reader of `State`.
- **Every touch is a file.** No document describes the order, so none changes.
- **The figures are flai's.** Forecast 30m, raised from flai's 14m with a stated reason. Cost of delay 40 USD a week, from flai's 16m input.

Source: design/conventions/strategic-agents.md § As the orchestrator

### 2026-10-14T23:11:57Z orchestrator
Resolved: Plan approved: two tasks in two layers cover both criteria; figures are flai's.
