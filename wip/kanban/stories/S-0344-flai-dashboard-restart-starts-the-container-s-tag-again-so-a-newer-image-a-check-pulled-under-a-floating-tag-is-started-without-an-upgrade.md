---
id: S-0344
type: story
nature: remediation
title: flai dashboard restart starts the container's tag again, so a newer image a check pulled under a floating tag is started without an upgrade
status: ready
owner: alex
created: 2026-10-08T08:08:17Z
updated: 2026-10-08T08:41:08Z
transitions:
  - to: ready
    at: 2026-10-08T08:38:55Z
    by: alex
tags: [flai, dashboard]
touches: [flai/cmd/dashboard_upgrade.go, flai/cmd/dashboard.go, flai/cmd/dashboard_test.go, flai/cmd/dashboard_watch.go, flai/cmd/dashboard_watch_test.go, docs/users/flai-reference.md, docs/users/flai.md, docs/operators/index.md, design/system/flai-cli.md, design/issues/I-0116-flai-dashboard-restart-starts-the-container-s-tag-again-so-a-newer-image-a-check-pulled-under-a-floating-tag-is-started-without-an-upgrade.md, design/issues/summary.md]
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
      seconds: 375
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 97
          output: 1614
          cache_read: 21717522
          cache_write: 33239
          cost: 5.3592
cost_of_delay:
  inputs:
    time_lost_per_cycle: 5m
    by: flai
    at: 2026-10-08T08:08:17Z
  value: 12.5
  by: planner-S-0344
  at: 2026-10-08T08:36:23Z
forecast:
  duration: 30m
  delivery: 2026-10-08T09:15:00Z
  basis: "Its own forecast of 30m; 1st in the pull order with an in-progress limit of 5, behind S-0232, S-0321, S-0322, S-0340 and S-0341."
  by: flai
  at: 2026-10-08T08:41:08Z
finalized:
  by: orchestrator
  at: 2026-10-08T08:36:52Z
---
# S-0344 flai dashboard restart starts the container's tag again, so a newer image a check pulled under a floating tag is started without an upgrade

## Goal

This story remediates [I-0116](../../../design/issues/I-0116-flai-dashboard-restart-starts-the-container-s-tag-again-so-a-newer-image-a-check-pulled-under-a-floating-tag-is-started-without-an-upgrade.md), "flai dashboard restart starts the container's tag again, so a newer image a check pulled under a floating tag is started without an upgrade". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0116 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0116 is closed with `flai issue close I-0116 --reason` saying what fixed it

## Tasks
- T-1360 flai dashboard restart starts the image the container runs, by its image ID, never the tag again
- T-1361 flai host's watch restarts the dashboard from the image ID it recorded, not the tag
- T-1362 Document that a dashboard restart keeps the image that runs, and close I-0116

## Notes

Cost of delay inputs set by flai from I-0116. time_lost_per_cycle 5m: 5m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-10-08T00:09:35Z, 0.3 days before this story; under one cycle counts as one).

### Planning

Planned by planner-S-0344 on 2026-10-08.

Proposed fix, from I-0116's instance: a restart starts from the image ID the container runs (`{{.Image}}`), as S-0316's upgrade fallback does, and keeps naming the tag in status and output. The same cause is in `flai host`'s watch, which restarts from the tag `startDashboard` recorded; ADR-0118 §1 says both restarts keep the release that runs, so T-1361 fixes the watch too.

Touches, file by file; no folder touch:

| Touch | From |
|-------|------|
| `flai/cmd/dashboard_upgrade.go` | layout: `runDashboardRestart` and `containerArgs`, named by I-0116 |
| `flai/cmd/dashboard.go` | co-change (86% of commits with the restart code): `containerInfo` |
| `flai/cmd/dashboard_test.go` | co-change (57%): the restart tests |
| `flai/cmd/dashboard_watch.go` | layout: `startDashboard`, the record, `restartDashboard` |
| `flai/cmd/dashboard_watch_test.go` | co-change (29%): the watch tests |
| `docs/users/flai-reference.md` | layout: generated from restart's help, which changes |
| `docs/users/flai.md` | design: § Run the dashboard, per the documentation convention |
| `docs/operators/index.md` | design: § The dashboard's watch names the recorded reference |
| `design/system/flai-cli.md` | design: § A chosen dashboard describes restart |
| `design/issues/I-0116-…md`, `design/issues/summary.md` | criterion 2: `flai issue close` |

`flai touches suggest` found nothing alone, since the story declared no touches; run from the two restart files it gave the co-change rows above. S-0321 in progress also claims `docs/operators/index.md` and `design/system/flai-cli.md`, and S-0322 `docs/users/flai-reference.md`; only T-1362 and T-1360's regeneration change them.

Figures:

- Forecast 30m, adjusted from flai's 17m: the median per unit of size does not count the fix the story must choose and record first, nor the two reproducing tests for restart and the watch. Delivery moves from flai's 13:49 by the 13m added.
- Cost of delay 12.50 USD a week, as `flai cod` gives it from the inputs: 5m lost per 168h cycle at 150 USD an hour. It stands: one occurrence, read from the code, not yet seen on the host.

Tasks, in three layers:

1. T-1360: restart from the image ID, with its test.
2. T-1361, after T-1360: the watch from the recorded image ID.
3. T-1362, after T-1360 and T-1361: the documents, and closing I-0116.
