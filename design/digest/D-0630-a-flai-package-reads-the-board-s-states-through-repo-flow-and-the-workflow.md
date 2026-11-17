---
id: D-0630
title: A flai package reads the board's states through repo.Flow() and the workflow package's lookups, Kind, OfKind, Find, WithEffect, and Permits, not through a workitem helper
kind: decision
sources: [S-0746, S-0749, S-0750]
topics: [flai, go, code]
paths: [flai/internal/workflow, flai/internal/workitem/repo.go, flai/internal/mcpserver, flai/internal/storystart, flai/internal/taskdone, flai/internal/inbox, flai/internal/itemedit, flai/internal/storygit, flai/internal/guard, flai/internal/messages, flai/internal/planning, flai/internal/metrics]
date: 2026-11-15
status: active
absorbed_by: ""
asked_again: []
cited: [S-0749, S-0750, S-0752]
---

# D-0630 A flai package reads the board's states through repo.Flow() and the workflow package's lookups, Kind, OfKind, Find, WithEffect, and Permits, not through a workitem helper

S-0746 moved the MCP tools, story start, task done, the sync, item edits, and the inbox off state names and onto a state's kind and a transition's effect (ADR-0301). Its narrative's Decisions (`wip/archive/agents/S-0746.md`) record how each package reaches the workflow:

> Every package reads its states through repo.Flow() and the workflow package's lookups (Kind, OfKind, Find, WithEffect, Permits) rather than a workitem helper, so that it does not wait on S-0739, which changes workitem in parallel.

So an E-0034 story that turns another package's state-name comparison into a kind or an effect, as S-0747, S-0749, and S-0750 do, calls `repo.Flow()` and those lookups in `flai/internal/workflow`. It does not add a helper to `flai/internal/workitem`, which a parallel story may be changing. ADR-0301 and `design/system/workflow.md` say that code reads a state's kind, not which functions it calls.

S-0749 followed this entry for the guard, messages, the forecast, prime packs, the release evaluation, the digest, and statsread, and named the idiom that keeps every caller unchanged (`wip/archive/agents/S-0749.md`, § Decisions):

> Each package follows the idiom S-0744 and S-0746 set: a function XIn(flow workflow.Workflow, ...) holds the logic and reads kinds, effects, and who may move from the workflow, and the old function calls it with workflow.Baseline() or repo.Flow(), so every existing caller and test is unchanged.
  
> Followed D-0630: each package reads states through repo.Flow() and the workflow's lookups (Kind, State, Initial, Permits, transitions' effects), with a baseline wrapper keeping every exported signature.

So an E-0034 story moving a package onto kinds writes `XIn(flow workflow.Workflow, ...)` with the logic and keeps the old exported function as a wrapper that passes `workflow.Baseline()` or `repo.Flow()`, rather than changing each caller's signature.

S-0750 followed it for the metrics in `flai/internal/metrics` (`wip/archive/agents/S-0750.md`, § Decisions):

> Following D-0630, the metrics read states through the workflow package's lookups (Kind, OfKind, WithEffect), not through a new workitem helper. committed and started are the first moves into the states the commit and pull transitions move to (ADR-0305 §5).

Sources: S-0746's, S-0749's, and S-0750's narratives, `wip/archive/agents/S-0746.md`, `wip/archive/agents/S-0749.md`, and `wip/archive/agents/S-0750.md`, § Decisions.
