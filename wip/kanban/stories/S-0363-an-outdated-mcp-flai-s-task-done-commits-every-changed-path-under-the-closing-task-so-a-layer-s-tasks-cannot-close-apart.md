---
id: S-0363
type: story
nature: remediation
title: An outdated MCP flai's task_done commits every changed path under the closing task, so a layer's tasks cannot close apart
status: backlog
owner: alex
created: 2026-10-12T14:12:30Z
updated: 2026-10-12T16:51:13Z
transitions: []
tags: [cli]
touches: [flai/internal/taskdone/taskdone.go, flai/internal/taskdone/taskdone_test.go, flai/internal/mcpserver/task_test.go, flai/cmd/task_done.go, docs/users/flai-reference.md, design/adrs/0136-flai-task-done-refuses-to-close-a-task-while-the-running-flai-is-older-than-the.md, design/adrs/README.md, design/system/flai-cli.md, docs/users/flai.md, docs/operators/index.md, design/issues/I-0127-an-outdated-mcp-flai-s-task-done-commits-every-changed-path-under-the-closing-task-so-a-layer-s-tasks-cannot-close-apart.md, design/issues/summary.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: sum
  seconds: 0
  models: []
  strategic:
    - kind: orchestrator
      seconds: 226
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 35
          output: 591
          cache_read: 3897612
          cache_write: 23346
          cost: 0.9885
draft: true
cost_of_delay:
  inputs:
    time_lost_per_cycle: 5m
    by: flai
    at: 2026-10-12T14:12:30Z
  value: 12.5
  by: planner-S-0363
  at: 2026-10-12T14:35:38Z
forecast:
  duration: 20m
  delivery: 2026-10-13T02:18:00Z
  basis: "Its own forecast of 20m; 21st in the pull order with an in-progress limit of 5, behind S-0367, S-0233, S-0236, S-0237, S-0238, S-0239, S-0241, S-0304, S-0305, S-0306, S-0350, S-0351, S-0352, S-0353, S-0354, S-0355, S-0356, S-0357, S-0358, S-0359 and S-0360."
  by: flai
  at: 2026-10-12T16:51:13Z
---
# S-0363 An outdated MCP flai's task_done commits every changed path under the closing task, so a layer's tasks cannot close apart

## Goal

This story remediates [I-0127](../../../design/issues/I-0127-an-outdated-mcp-flai-s-task-done-commits-every-changed-path-under-the-closing-task-so-a-layer-s-tasks-cannot-close-apart.md), "An outdated MCP flai's task_done commits every changed path under the closing task, so a layer's tasks cannot close apart". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0127 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0127 is closed with `flai issue close I-0127 --reason` saying what fixed it

## Tasks
- T-1439 ADR-0136: flai task done refuses to close a task while the running flai is older than the project's history
- T-1440 taskdone.Run refuses, changing nothing, while the running flai is older than the project's history, tested through Run and the MCP tool
- T-1441 The design, user, and operator docs say flai task done refuses while the running flai is older than the project's history
- T-1442 Close I-0127 as fixed by task_done's refusal while the running flai is older than the project's history

## Notes

Cost of delay inputs set by flai from I-0127. time_lost_per_cycle 5m: 5m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-10-08T10:41:45Z, 4.1 days before this story; under one cycle counts as one).

### Planning

Proposed remedy, from the instance: `taskdone.Run` refuses, changing nothing, while `release.FlaiOutdated` says the running flai is older than the newest `flai/vX.Y.Z` tag in the main checkout's history. It is recorded first as ADR-0136 (T-1439).

- I-0127's one instance is S-0338's T-1324, closed through the MCP tool `task_done` on the installed flai 1.39.13, below 1.40.0 in the project's history. The inbox reported `flai_outdated` then, and the close went ahead.
- The specific defect, the commit of another open task's paths, was fixed by S-0322 (ADR-0128) in 1.39.14. S-0362 adds its tests. This story removes the cause behind both issues: an outdated flai running `task_done` by older rules. The host runs 1.40.1 today, with 1.40.3 in the history, so the cause still holds.
- Today nothing refuses on an outdated flai. `flai_outdated` only warns, in `inbox`, `project.info`, `flai serve`, and after `flai task done`'s inbox. ADR-0117 has the minimum warn, and ADR-0136 must say why `task_done` differs.
- The MCP tool (`mcpserver/task.go`), `flai task done` (`cmd/task_done.go`), and the host's `task.done` (which runs the CLI) all reach `taskdone.Run`, so the check goes there once.
- A dev build, such as `scripts/flai.sh`'s, is never refused, since `FlaiOutdated` returns nil on a version that is not X.Y.Z.

Touches, each a file, none a folder:

| Touch | Source | Why |
|-------|--------|-----|
| `flai/internal/taskdone/taskdone.go` | layout | `Run` and `start` are the one entry all three paths reach |
| `flai/internal/taskdone/taskdone_test.go` | layout | reproduces I-0127's layer and the refusal |
| `flai/internal/mcpserver/task_test.go` | layout | the MCP `task_done` the instance went through |
| `flai/cmd/task_done.go` | layout | its help says it refuses |
| `docs/users/flai-reference.md` | co-change, layout | generated from that help |
| `design/adrs/0136-flai-task-done-refuses-to-close-a-task-while-the-running-flai-is-older-than-the.md` | design | the decision, which departs from ADR-0117's warn-only precedent |
| `design/adrs/README.md` | co-change | the ADR index |
| `design/system/flai-cli.md` | design | § Versions: the host's flai and the tree |
| `docs/users/flai.md` | design | § Closing a task |
| `docs/operators/index.md` | design | § Keeping the host's flai current |
| `design/issues/I-0127-…` | design | `flai issue close` writes it (criterion 2) |
| `design/issues/summary.md` | co-change, layout | `flai issue close` regenerates it |

The story declared no touches, so `flai touches suggest` had nothing to start from. Seeded with the touches above, it listed nothing over 12%: `flaiover-dashboard.md`, `flaiover.md`, `settings.md`, and `workflow.md`. None is predicted, since no setting, dashboard view, or workflow rule changes. The ADR's number and slug are predicted from `flai adr new`'s 80-character slug. If they differ, `task_done` widens the touches.

Overlaps: S-0362 (in progress) touches both test files, and S-0337 (in progress) touches `design/adrs`, `flai-cli.md`, `docs/users/flai.md`, and `flai-reference.md`. S-0363 is held behind them once ready. That is expected: the tests and the docs must change.

Figures:

- Forecast 20m, delivery 2026-10-13T02:24Z: `flai forecast` as given. It is the median of 82 s per unit of size over 23 done large-band remediation stories on claude-opus-5-5, times size 14 (2 criteria, 12 touches). An ADR, a check of a few lines over an existing function, two tests, three doc sections, and an issue close fit that, so it is not adjusted.
- Cost of delay 12.50 USD a week: `flai cod` from flai's input, 5m lost per 168h cycle at 150 USD an hour. It is not adjusted. The issue has one instance, and its specific defect is fixed in the installed flai.
