---
id: S-0795
type: story
nature: remediation
title: TestContainerStrategicStartUnderALayerWaitsForTheHarnessImageAlone fails under load in the race run
status: done
owner: alex
created: 2026-11-16T00:01:08Z
updated: 2026-11-16T01:53:29Z
transitions:
  - to: ready
    at: 2026-11-16T00:43:21Z
    by: orchestrator
  - to: in-progress
    at: 2026-11-16T00:43:35Z
    by: agent-S-0795
  - to: review
    at: 2026-11-16T01:36:49Z
    by: agent-S-0795
  - to: done
    at: 2026-11-16T01:53:29Z
    by: alex
tags: [flai, serve, test]
touches: [flai/internal/serve/agents_state.go, flai/internal/serve/agents_start.go, flai/internal/serve/agents_test.go, design/issues/I-0369-testcontainerstrategicstartunderalayerwaitsfortheharnessimagealone-fails-under.md, design/issues/summary.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 3215
  priced_by: harness
  turns:
    - day: 2026-11-16
      ceremony: 2
      test_runs: 3
      hand_edits: 2
      work: 51
  models:
    - model: claude-opus-5-5
      input: 126
      output: 24175
      cache_read: 13403938
      cache_write: 258402
      cost: 5.232
  strategic:
    - kind: orchestrator
      seconds: 152
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 33
          output: 256
          cache_read: 1171910
          cache_write: 133437
          cost: 0.3964
cost_of_delay:
  inputs:
    time_lost_per_cycle: 10m
    by: flai
    at: 2026-11-16T00:01:08Z
  value: 25
  by: planner-S-0795
  at: 2026-11-16T00:40:34Z
forecast:
  duration: 40m
  delivery: 2026-11-16T03:30:00Z
  basis: "flai forecast's 16m (137 s per unit over 67 done medium remediation stories, size 7), raised to 40m because a close-out alone took about 9m on S-0790 and S-0755 and a load-only race flake needs repeated -race runs under load to show it gone; delivery moved by the same 24m from flai's 03:03 as 7th in the pull order."
  by: planner-S-0795
  at: 2026-11-16T00:40:34Z
finalized:
  by: orchestrator
  at: 2026-11-16T00:43:03Z
---
# S-0795 TestContainerStrategicStartUnderALayerWaitsForTheHarnessImageAlone fails under load in the race run

## Goal

This story remediates [I-0369](../../../design/issues/I-0369-testcontainerstrategicstartunderalayerwaitsfortheharnessimagealone-fails-under.md), "TestContainerStrategicStartUnderALayerWaitsForTheHarnessImageAlone fails under load in the race run". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [x] The cause I-0369 describes no longer occurs, with a test that reproduces it where one fits
- [x] I-0369 is closed with `flai issue close I-0369 --reason` saying what fixed it

## Tasks
- T-3392 Replace an image-held run with its start in one write, so it is never kept among the past runs
- T-3393 Show the held-start container tests pass under load in the race run, and close I-0369

## Notes

Cost of delay inputs set by flai from I-0369. time_lost_per_cycle 10m: 10m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-11-13T20:34:57Z, 2.1 days before this story; under one cycle counts as one).

### Planning

Proposed cause, read from the code at planning. A planner start held for its image job goes on in `goOn` (`flai/internal/serve/agents_start.go`). Its start records the new run with `AgentState.put`. `keep` (`flai/internal/serve/agents_state.go`) then pushes the held run, PID 0, into `PastPlans`, since `same` (`flai/internal/serve/agents.go`) finds the PIDs differ. A second write, `dropPast`, takes it out again. The test waits for the new run's PID and reads `PastPlans` between those two writes when the host is loaded. That is the failure at `container_test.go:1618` in both instances. Proposed fix: make the replacement one write, with `keep` leaving out a replaced run that is `imageHeld()`, and a unit test of `put` that fails on the code as it was. T-3392 confirms the cause and builds the fix; T-3393 shows it holds under a loaded `-race` run and closes I-0369.

Touches:

| Path | Source | Why |
|------|--------|-----|
| `flai/internal/serve/agents_state.go` | design | `put` and `keep`, where the held run enters the past runs |
| `flai/internal/serve/agents_start.go` | design | `goOn` and `dropPast`, the second write |
| `flai/internal/serve/agents.go` | co-change (74%) | `same` and `imageHeld`, which the fix may use or change |
| `flai/internal/serve/agents_test.go` | layout | the unit test beside `TestEveryRolesEarlierRunsAreKept` |
| `flai/internal/serve/container_test.go` | design | the failing test, if its check must change |
| `design/issues/I-0369-…-fails-under.md` | design | closed by criterion 2 |

`flai touches suggest` found no touches at first. With the files above it listed `agents.go` at 74%, which is kept; the rest (`flai/cmd/serve_actions.go` at 24% and below) belong to other changes to serve, not to this fix.

Forecast: `flai forecast` gave 16m (137 s per unit over 67 done medium remediation stories, size 7). Raised to 40m. A close-out alone took about 9m on S-0790 and S-0755. A load-only race flake also needs repeated `-race` runs under load to show it gone. Delivery 03:30Z is flai's 03:03Z, as 7th in the pull order, moved by the same 24m.

Cost of delay: `flai cod` gives 25.00 USD a week, 10m lost per 168h cycle at 150 USD an hour. It stands. I-0369 now counts two occurrences, which would double the input; the inputs are the operator's, so they are left as set.

No `after`: the story reads nothing another open story writes.
