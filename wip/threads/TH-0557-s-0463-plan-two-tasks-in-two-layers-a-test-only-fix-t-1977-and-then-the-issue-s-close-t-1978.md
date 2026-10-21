---
id: TH-0557
title: "S-0463 plan: two tasks in two layers, a test-only fix (T-1977) and then the issue's close (T-1978)"
anchor:
  path: wip/kanban/stories/S-0463-testachildisstoppedonlyonceitsprocesshasended-fails-under-load-when-the-stubborn-child-dies-on-sigterm-before-it-ignores-it.md
  item: S-0463
status: resolved
participants: [planner-S-0463, orchestrator]
created: 2026-10-21T07:33:07Z
updated: 2026-10-21T07:33:19Z
---

# TH-0557 S-0463 plan: two tasks in two layers, a test-only fix (T-1977) and then the issue's close (T-1978)

On wip/kanban/stories/S-0463-testachildisstoppedonlyonceitsprocesshasended-fails-under-load-when-the-stubborn-child-dies-on-sigterm-before-it-ignores-it.md.

## Entries

### 2026-10-21T07:33:07Z planner-S-0463
The plan for S-0463 has two tasks, worked one at a time.

| Layer | Task | Touches | After |
|-------|------|---------|-------|
| 1 | T-1977: the `stubborn` test child writes a file once it ignores SIGTERM. Both host tests wait for that file instead of sleeping 100 ms. A delay knob reproduces the race. | `flai/internal/host/host_test.go` | none |
| 2 | T-1978: close I-0173 with a reason naming the fix | the I-0173 issue file, `design/issues/summary.md` | T-1977 |

Figures written:

- Forecast: 25m. flai gave 11m; I raised it for a race stress run and a close-out whose integration tier took 5m56s.
- Delivery: 2026-10-22T01:16:00Z.
- Cost of delay value: 25 USD a week, from flai's input of 10m lost per cycle.

Assumptions:

- The fix is in the tests only. The host's stop logic is right; the test assumed the child was ready without checking. So none of the co-change files `flai touches suggest` listed is taken (`host.go`, `child.go`, `api.go`, `client.go`, `flai/cmd/host.go`).
- `TestAChildThatWillNotStopIsKilled` has the same 100 ms sleep, so T-1977 fixes it too.
- The repro is a delay knob on the child that pushes its ignore past the old sleep. T-1977 keeps the knob only if a test uses it.

I propose no splits, merges, or drops. Nothing waits on an answer here.

### 2026-10-21T07:33:15Z orchestrator
Approved. T-1977 meets criterion 1: the stubborn child signals readiness, and the delay knob reproduces the race. T-1978 meets criterion 2. The touches are the files the fix and the close change, and leaving out the host's production files is reasoned. The forecast of 25m is reasoned, and the 25 USD value comes from flai's input.

Source: design/conventions/strategic-agents.md

### 2026-10-21T07:33:19Z orchestrator
Resolved: The plan for S-0463 is approved.
