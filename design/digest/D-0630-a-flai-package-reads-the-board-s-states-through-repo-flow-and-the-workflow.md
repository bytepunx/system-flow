---
id: D-0630
title: A flai package reads the board's states through repo.Flow() and the workflow package's lookups, Kind, OfKind, Find, WithEffect, and Permits, not through a workitem helper
kind: decision
sources: [S-0746]
topics: [flai, go, code]
paths: [flai/internal/workflow, flai/internal/workitem/repo.go, flai/internal/mcpserver, flai/internal/storystart, flai/internal/taskdone, flai/internal/inbox, flai/internal/itemedit, flai/internal/storygit]
date: 2026-11-15
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0630 A flai package reads the board's states through repo.Flow() and the workflow package's lookups, Kind, OfKind, Find, WithEffect, and Permits, not through a workitem helper

S-0746 moved the MCP tools, story start, task done, the sync, item edits, and the inbox off state names and onto a state's kind and a transition's effect (ADR-0301). Its narrative's Decisions (`wip/archive/agents/S-0746.md`) record how each package reaches the workflow:

> Every package reads its states through repo.Flow() and the workflow package's lookups (Kind, OfKind, Find, WithEffect, Permits) rather than a workitem helper, so that it does not wait on S-0739, which changes workitem in parallel.

So an E-0034 story that turns another package's state-name comparison into a kind or an effect, as S-0747, S-0749, and S-0750 do, calls `repo.Flow()` and those lookups in `flai/internal/workflow`. It does not add a helper to `flai/internal/workitem`, which a parallel story may be changing. ADR-0301 and `design/system/workflow.md` say that code reads a state's kind, not which functions it calls.

Source: S-0746's narrative, `wip/archive/agents/S-0746.md`, § Decisions.
