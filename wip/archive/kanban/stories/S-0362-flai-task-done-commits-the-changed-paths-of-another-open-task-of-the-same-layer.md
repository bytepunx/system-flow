---
id: S-0362
type: story
nature: remediation
title: flai task done commits the changed paths of another open task of the same layer
status: done
owner: alex
created: 2026-10-12T14:12:29Z
updated: 2026-10-12T15:10:25Z
transitions:
  - to: ready
    at: 2026-10-12T14:32:20Z
    by: orchestrator
  - to: in-progress
    at: 2026-10-12T14:32:27Z
    by: agent-S-0362
  - to: review
    at: 2026-10-12T14:46:35Z
    by: agent-S-0362
  - to: done
    at: 2026-10-12T15:10:25Z
    by: alex
tags: [cli]
touches: [flai/internal/taskdone/taskdone_test.go, flai/internal/mcpserver/task_test.go, design/issues/I-0126-flai-task-done-commits-the-changed-paths-of-another-open-task-of-the-same-layer.md, design/issues/summary.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 862
  turns:
    - day: 2026-10-12
      ceremony: 1
      test_runs: 4
      hand_edits: 2
      work: 21
  models:
    - model: claude-opus-5-5
      input: 94
      output: 23329
      cache_read: 4767165
      cache_write: 245824
      cost: 3.1112
  strategic:
    - kind: planner
      seconds: 358
      estimated: true
      models:
        - model: claude-haiku-5-5
          input: 60
          output: 243
          cache_read: 2247389
          cache_write: 115906
          cost: 0
        - model: claude-opus-5-5
          input: 57
          output: 1508
          cache_read: 2573286
          cache_write: 108683
          cost: 1.2435
    - kind: orchestrator
      seconds: 1238
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 63
          output: 899
          cache_read: 7867556
          cache_write: 36414
          cost: 1.9926
        - model: claude-sonnet-5-5
          input: 4
          output: 40
          cache_read: 12362
          cache_write: 28753
          cost: 0.0333
cost_of_delay:
  inputs:
    time_lost_per_cycle: 3m
    by: flai
    at: 2026-10-12T14:12:29Z
  value: 7.5
  by: planner-S-0362
  at: 2026-10-12T14:19:58Z
forecast:
  duration: 12m
  delivery: 2026-10-13T02:11:00Z
  basis: "Its own forecast of 12m; 25th in the pull order with an in-progress limit of 5, behind S-0232, S-0337, S-0233, S-0234, S-0235, S-0236, S-0237, S-0238, S-0239, S-0241, S-0304, S-0305, S-0306, S-0313, S-0349, S-0350, S-0351, S-0352, S-0353, S-0354, S-0355, S-0356, S-0357, S-0358, S-0359 and S-0360."
  by: flai
  at: 2026-10-12T14:31:40Z
finalized:
  by: orchestrator
  at: 2026-10-12T14:32:13Z
---
# S-0362 flai task done commits the changed paths of another open task of the same layer

## Goal

This story remediates [I-0126](../../../design/issues/I-0126-flai-task-done-commits-the-changed-paths-of-another-open-task-of-the-same-layer.md), "flai task done commits the changed paths of another open task of the same layer". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [x] The cause I-0126 describes no longer occurs, with a test that reproduces it where one fits
- [x] I-0126 is closed with `flai issue close I-0126 --reason` saying what fixed it

## Tasks
- T-1437 Tests reproduce I-0126's layer through task_done and with a ready or backlog sibling, and pass on the current code
- T-1438 Close I-0126 as fixed by S-0322 in flai 1.39.14, its instance having run on 1.39.13

## Notes

Cost of delay inputs set by flai from I-0126. time_lost_per_cycle 3m: 3m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-10-08T10:36:24Z, 4.2 days before this story; under one cycle counts as one).

### Planning

Proposed remedy, from the instance: the current code fixes the cause, so the story adds tests and closes the issue.

- I-0126's one instance closed T-1377 at 2026-10-08T10:35:55Z through the MCP tool `task_done` on the installed flai 1.39.13.
- S-0322 implemented [ADR-0128](../../../design/adrs/0128-flai-task-done-commits-the-paths-the-closing-task-covers-and-those-no-other.md) (4040be86, 774e6784) and was published in flai 1.39.14 at 09:30Z the same day.
- The current `commit` and `coverage` in `flai/internal/taskdone/taskdone.go` leave a path that only another open task covers. Ready, backlog, in-progress, and review siblings all count as open.
- The MCP tool (`flai/internal/mcpserver/task.go`) and the host method `task.done` run the same code.
- Missing tests: an MCP-level check that the closing task's touches are not widened, and a sibling in `ready` or `backlog`.
- I-0127 has the same cause, an MCP flai older than the project's history (S-0338, also on 1.39.13). Its guard is S-0363's, not this story's.

Touches, each a file, none a folder:

| Touch | Source | Why |
|-------|--------|-----|
| `flai/internal/taskdone/taskdone_test.go` | design, layout | ADR-0128's tests live here; T-1437 adds the sibling-state table |
| `flai/internal/mcpserver/task_test.go` | design, layout | the MCP `task_done` test I-0126's instance went through |
| `design/issues/I-0126-flai-task-done-commits-the-changed-paths-of-another-open-task-of-the-same-layer.md` | design | `flai issue close` writes it (criterion 2) |
| `design/issues/summary.md` | co-change, layout | `flai issue close` regenerates it |

`flai touches suggest` declared nothing to start from. Seeded with the two test files, it listed `taskdone.go` (57%), `mcpserver/task.go` (43%), `cmd/task_done.go`, its test, and `docs/users/flai-reference.md`. None is predicted. The code already meets ADR-0128, and no behaviour or command changes. If a T-1437 test fails, `task_done` widens the touches to the fix in `taskdone.go`.

Figures:

- Forecast 12m, delivery 2026-10-13T02:01Z: `flai forecast` as given. It is the median of 119 s per unit of size over 7 done small remediation stories on claude-opus-5-5, times size 6 (2 criteria, 4 touches). Two tests and an issue close fit that, so it is not adjusted.
- Cost of delay 7.50 USD a week: `flai cod` from flai's input, 3m lost per 168h cycle at 150 USD an hour. Not adjusted: one instance, on a flai no longer installed.
