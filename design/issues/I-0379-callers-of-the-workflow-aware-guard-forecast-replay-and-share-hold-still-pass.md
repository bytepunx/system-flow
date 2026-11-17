---
id: I-0379
title: Callers of the workflow-aware guard, forecast replay, and share hold still pass the baseline, so a project's own workflow file will not reach them
class: defect
status: open
count: 1
cost: 15m
first_reported: 2026-11-15T20:04:03Z
last_reported: 2026-11-15T20:04:03Z
updated: 2026-11-17T08:32:54Z
---

# I-0379 Callers of the workflow-aware guard, forecast replay, and share hold still pass the baseline, so a project's own workflow file will not reach them

## Description
Callers of the workflow-aware guard, forecast replay, and share hold still pass the baseline, so a project's own workflow file will not reach them

## Instances

### 2026-11-15T20:04:03Z
Story: S-0749.
S-0749 made the guard, the forecast, and the messages read the workflow, with a baseline wrapper for each caller. Four callers outside its touches still run the baseline: flai/cmd/guard.go builds guard.Guard with no Workflow; flai/internal/serve/replan.go calls planning.Replay rather than planning.ReplayIn(repo.Flow(), ...); flai/internal/metrics/costorder.go builds its puller with planning.NewPuller (S-0750's); and flai/internal/workitem/hold.go still names InProgress for the story that holds a ready one, so a share's hold does not follow a second active state. Today repo.Flow() is always the baseline, so nothing behaves differently yet. Once S-0755 reads a project's own workflow file, each caller must pass repo.Flow(), or the guard, replan, and holds will run the baseline under a renamed board. Also, digest.Decisions now reads the accepted kind but still filters on Item.CompletedAt, which I-0377 says reads the baseline's done and cancelled names.

## Remediation

Story S-0843 remediates this issue, created from it at 2026-11-17T08:32:54Z.
