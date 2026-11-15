---
id: D-0640
title: "flai guard reads a strategic role's moves from the workflow: the planner moves only to the initial state, the orchestrator's promote_to_ready is a move into a state a commit transition enters, its acceptance a move only an acceptance may make, and the protected-write hold asks in any active state"
kind: decision
sources: [S-0749, T-3232]
topics: [flai, go, code, guard]
paths: [flai/internal/guard/guard.go, flai/internal/guard/guard_rules.go, flai/internal/guard/guard_refusals.go, flai/internal/guard/roles.go]
date: 2026-11-15
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0640 flai guard reads a strategic role's moves from the workflow: the planner moves only to the initial state, the orchestrator's promote_to_ready is a move into a state a commit transition enters, its acceptance a move only an acceptance may make, and the protected-write hold asks in any active state

S-0749 moved flai guard off the names `backlog`, `ready`, `done`, and `in-progress` onto the workflow (ADR-0301). Its narrative's Decisions (`wip/archive/agents/S-0749.md`):

> Guard (T-3232): the planner moves only to the workflow's initial state, the orchestrator's promote_to_ready is a move into a state a commit transition enters, its acceptance is a move a role may make only as an acceptance, and the protected-write hold asks in any active state. The Backlog and Ready constants stay because guard/roles.go, S-0748's, still reads them. Refusal texts are fixed strings, word for word with the baseline.

So on the baseline each rule means what it did, `backlog`, `ready`, `done`, and `in-progress`, and on a project's own board it follows the states' kinds and the transitions' effects and `by`. The refusal texts did not change. Until S-0755 reads a project's workflow file, and while `flai/cmd/guard.go` sets no `Guard.Workflow` (I-0379), the guard runs the baseline.

Source: S-0749's narrative, `wip/archive/agents/S-0749.md`, § Decisions; I-0379.
