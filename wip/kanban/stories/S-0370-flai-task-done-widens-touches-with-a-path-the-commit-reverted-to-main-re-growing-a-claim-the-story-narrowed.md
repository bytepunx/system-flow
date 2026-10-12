---
id: S-0370
type: story
nature: remediation
title: flai task done widens touches with a path the commit reverted to main, re-growing a claim the story narrowed
status: ready
owner: alex
created: 2026-10-12T16:52:44Z
updated: 2026-10-12T20:01:35Z
transitions:
  - to: ready
    at: 2026-10-12T17:57:35Z
    by: orchestrator
tags: [cli]
topics: [touches]
touches: [flai/internal/taskdone/taskdone.go, flai/internal/taskdone/taskdone_test.go, flai/cmd/task_done.go, flai/internal/mcpserver/task.go, docs/users/flai-reference.md, design/adrs/README.md, design/system/flai-cli.md, docs/users/flai.md, design/issues/I-0121-flai-task-done-widens-touches-with-a-path-the-commit-reverted-to-main-re-growing-a-claim-the-story-narrowed.md, design/issues/summary.md]
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
      seconds: 220
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 46
          output: 831
          cache_read: 13593714
          cache_write: 16198
          cost: 3.4309
cost_of_delay:
  inputs:
    time_lost_per_cycle: 5m
    by: flai
    at: 2026-10-12T16:52:44Z
  value: 12.5
  by: planner-S-0370
  at: 2026-10-12T17:09:21Z
forecast:
  duration: 21m
  delivery: 2026-10-12T20:26:00Z
  basis: "Its own forecast of 21m; 1st in the pull order with an in-progress limit of 5, with nothing ahead of it."
  by: flai
  at: 2026-10-12T20:01:35Z
finalized:
  by: orchestrator
  at: 2026-10-12T17:11:28Z
---
# S-0370 flai task done widens touches with a path the commit reverted to main, re-growing a claim the story narrowed

## Goal

This story remediates [I-0121](../../../design/issues/I-0121-flai-task-done-widens-touches-with-a-path-the-commit-reverted-to-main-re-growing-a-claim-the-story-narrowed.md), "flai task done widens touches with a path the commit reverted to main, re-growing a claim the story narrowed". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0121 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0121 is closed with `flai issue close I-0121 --reason` saying what fixed it

## Tasks
- T-1472 ADR: flai task done widens touches only with the committed paths the story's branch still changes against main
- T-1473 flai task done widens touches only with committed paths the branch still changes against main, with a test reproducing I-0121's revert
- T-1474 The design and user docs say flai task done widens touches only with committed paths the branch still changes against main
- T-1475 Close I-0121 as fixed by widening touches only with paths the branch still changes against main

## Notes

Cost of delay inputs set by flai from I-0121. time_lost_per_cycle 5m: 5m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-10-08T08:25:05Z, 4.4 days before this story; under one cycle counts as one).

### Planning

Proposed remedy, from the instance and as the issue suggests: the touches step of `flai task done` widens the task's and the story's touches only with the committed paths that the story's branch still changes against its merge base with main. It reads them with `storygit.StoryDiff` after the sync step. T-1472 records this first as an ADR refining ADR-0107 and ADR-0128.

- I-0121's one instance is S-0232's T-1223. A revert commit changed `scripts/flai-snapshot.sh` back to main's version, and `(*run).touches` in `flai/internal/taskdone/taskdone.go` widened with every path of the commit (`r.res.Commit.Paths`), so the path came back.
- The tell step (S-0333) and the `--commit` writes (S-0275, `cmd/storycommit.go`, `mcpserver/issues.go`) also call `itemedit.WidenStory` or read the commit's paths. They are left as they are: those writes only add content they wrote, and a revert is still worth telling.

Touches, each a file, none a folder:

| Touch | Source | Why |
|-------|--------|-----|
| `flai/internal/taskdone/taskdone.go` | layout | `(*run).touches` is where the commit's paths are widened |
| `flai/internal/taskdone/taskdone_test.go` | layout | reproduces I-0121's revert with the real git fixture |
| `flai/cmd/task_done.go` | layout | its help's step 6 describes the widening |
| `flai/internal/mcpserver/task.go` | layout | `taskDoneDescription` describes the widening |
| `docs/users/flai-reference.md` | co-change, layout | generated from that help |
| `design/adrs/README.md` | design | the ADR index; the ADR file is named by `flai adr new` and widened in when committed |
| `design/system/flai-cli.md` | design | the `flai task done` row |
| `docs/users/flai.md` | design | § Closing a task, step 6, and § Touches |
| `design/issues/I-0121-…` | design | `flai issue close` writes it (criterion 2) |
| `design/issues/summary.md` | co-change, layout | `flai issue close` regenerates it |

The story declared no touches, so `flai touches suggest` had nothing to start from. Seeded with the touches above, it listed nothing over 12%: `flaiover-dashboard.md`, `docs/operators/index.md`, `flaiover.md`, `settings.md`, and `workflow.md`. None is predicted, since no setting, dashboard view, or workflow rule changes. The ADR's file is not predicted: S-0363 also adds one, and the number depends on which is written first.

Overlaps: S-0363 (backlog) claims `taskdone.go`, `taskdone_test.go`, `cmd/task_done.go`, `flai-reference.md`, `design/adrs/README.md`, `flai-cli.md`, `docs/users/flai.md`, and `summary.md`. Whichever starts second is held until the first is accepted. That is expected, since both change the same step file.

Figures:

- Forecast 21m, delivery 2026-10-13T01:58Z: `flai forecast` as given. It is the median of 101 s per unit of size over 22 done medium-band remediation stories on claude-opus-5-5, times size 12 (2 criteria, 10 touches). The work is a filter of a few lines over an existing diff reader, two git-fixture tests, an ADR, two doc sections, and an issue close. S-0363's planner kept 20m for the same shape, so it is not adjusted.
- Cost of delay 12.50 USD a week: `flai cod` from flai's input, 5m lost per 168h cycle at 150 USD an hour. It is not adjusted, since the issue has one instance and a hand narrowing works around it.
