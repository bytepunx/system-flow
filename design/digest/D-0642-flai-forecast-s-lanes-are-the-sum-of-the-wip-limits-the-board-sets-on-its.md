---
id: D-0642
title: flai forecast's lanes are the sum of the WIP limits the board sets on its active states, a story in any active state holds one, and its pull order is the queue states, nearest the pull first
kind: decision
sources: [S-0749, T-3234, S-0755, S-0742]
topics: [flai, go, code, planning]
paths: [flai/internal/planning/forecast.go, flai/cmd/forecast.go, flai/internal/statsread/statsread.go, flai/internal/metrics/costorder.go, flai/internal/serve/replan.go]
date: 2026-11-15
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0642 flai forecast's lanes are the sum of the WIP limits the board sets on its active states, a story in any active state holds one, and its pull order is the queue states, nearest the pull first

S-0749 moved flai forecast's play-out off `in-progress` and `ready` onto the workflow's kinds (ADR-0301). Its narrative's Decisions (`wip/archive/agents/S-0749.md`):

> Forecast (T-3234): the lanes are the sum of the limits the board sets on the active states with limit, and a story in any active state holds one; the pull order is the queue states, nearest the pull first. flai forecast now passes repo.Flow() (flai/cmd/forecast.go added to T-3234's touches).
  
> statsread's WIP limit is the limit of the state a pull moves a story into.

So on the baseline the lanes are the in-progress limit, as before. On a board with two active states each with a limit, such as one that adds a testing state after in-progress, the lanes are the two limits summed. The forecast's lanes and statsread's WIP limit are not the same figure on such a board. `flai/internal/metrics/costorder.go` still builds its puller with `planning.NewPuller`, the baseline form, so the cost of delay projection of `flai stats` plays out on the baseline (I-0379). `flai/internal/serve/replan.go`'s replay no longer does: S-0755 made it call `planning.ReplayIn(flow, …)` with the project's workflow, and S-0742 kept that line on its rebase (`wip/archive/agents/S-0742.md`, § Decisions).

Sources: S-0749's narrative, `wip/archive/agents/S-0749.md`, § Decisions; I-0379; S-0755; S-0742's narrative, `wip/archive/agents/S-0742.md`, § Decisions.
