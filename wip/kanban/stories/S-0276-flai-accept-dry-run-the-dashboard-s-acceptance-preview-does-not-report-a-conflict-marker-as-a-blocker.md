---
id: S-0276
type: story
nature: remediation
title: flai accept --dry-run, the dashboard's acceptance preview, does not report a conflict marker as a blocker
status: ready
owner: alex
created: 2026-10-05T03:48:28Z
updated: 2026-10-05T05:46:20Z
transitions:
  - to: ready
    at: 2026-10-05T04:05:04Z
    by: alex
tags: []
touches: [flai/internal/preview, flai/cmd/branch.go, design/system/flai-cli.md, docs/users/flai.md, flai/internal/conflictmark, flai/cmd/accept_conflict_test.go]
after: [S-0253]
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
    - kind: planner
      seconds: 431
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 76
          output: 13860
          cache_read: 2095308
          cache_write: 64013
          cost: 1.2087
cost_of_delay:
  inputs:
    time_lost_per_cycle: 5m
    by: planner-S-0276
    at: 2026-10-05T04:11:51Z
  value: 12.5
  by: planner-S-0276
  at: 2026-10-05T04:12:07Z
forecast:
  duration: 30m
  delivery: 2026-10-05T06:50:00Z
  basis: "Its own forecast of 30m; 1st in the pull order with an in-progress limit of 3, with nothing ahead of it."
  by: flai
  at: 2026-10-05T05:46:20Z
finalized:
  by: alex
  at: 2026-10-05T04:05:03Z
---
# S-0276 flai accept --dry-run, the dashboard's acceptance preview, does not report a conflict marker as a blocker

## Goal

The acceptance preview names every reason the acceptance would refuse, so that the operator sees a story branch carrying a conflict marker before pressing Accept rather than when the acceptance refuses it.

## Acceptance criteria
- [ ] `flai accept --dry-run`, and the dashboard's preview that runs it, lists a blocker naming each file and line on the story branch that carries a conflict marker, as `flai accept` refuses it (S-0253's `refuseConflictMarkers`).
- [ ] The preview and the acceptance find markers by one check, `flai/internal/conflictmark` over the files the branch adds or changes, so they cannot disagree.
- [ ] A test shows the preview reporting the blocker for a branch with a marker and none for a clean branch.
- [ ] `design/system/flai-cli.md` and `docs/users/flai.md` say the preview reports it.

## Tasks
- T-0871 conflictmark reads the conflict markers a story branch adds or changes, and flai accept refuses by it
- T-0872 flai accept --dry-run lists a blocker naming each conflict marker the story branch carries
- T-0873 flai-cli.md and the users' guide say the acceptance preview reports conflict markers

## Notes

Raised in S-0253: its acceptance gate refuses a branch with a conflict marker, but `flai/internal/preview` was outside that story, so the preview still offers Accept. The operator asked for this story in S-0253's Decisions (2026-10-05).

### Planning

Planned by planner-S-0276 on 2026-10-05. The plan thread lists the tasks and their layers.

Touches:

- `flai/internal/preview`, `flai/cmd/branch.go`, `design/system/flai-cli.md`, `docs/users/flai.md`: declared, kept.
- `flai/internal/conflictmark`: from the criteria and the layout. The branch scan, `(*app).conflictMarkers`, sits in `flai/cmd/branch.go`, which `preview` cannot import. Criterion 2 needs it moved into `conflictmark`.
- `flai/cmd/accept_conflict_test.go`: from the layout. `flai/internal/preview` has no tests of its own. The acceptance tests live in `flai/cmd`, and S-0253 put its conflict marker tests in this file.
- Not taken from `flai touches suggest`:
  - `design/system/flaiover-dashboard.md` and `docs/users/flaiover.md` (co-change 33% and 26%). The dashboard's Accept dialog shows `blockers[]` as they come, so neither the dashboard nor its docs change.
  - `docs/users/flai-reference.md` (25%). No command or flag changes.
  - `design/system/workflow.md` (16%). S-0253 already stated the acceptance rule there, and the preview adds no rule.

Forecast: 30m, delivery 2026-10-05T05:19Z.

- `flai forecast` gave 20m: 116 s per unit of size, times size 10 from 4 criteria and 6 touches.
- Raised for three serial tasks: a move of git-reading code across packages, with a test that builds a git repository, then the preview and its tests, then the docs. S-0253, the two-task story this one follows, spent 51m in progress.
- Delivery is flai's 05:09Z, fifth in the pull order, plus the extra 10m.

Cost of delay: 12.50 USD/week, from `flai cod`, kept as computed.

- The input `time_lost_per_cycle: 5m` is the operator's, given on TH-0127: "take recommendation".
- 5m is the cost recorded on I-0066. Without this story, a branch with a marker passes the preview, the acceptance refuses it, and the operator has to send the story back and accept it again.
