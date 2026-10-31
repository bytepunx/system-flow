---
id: T-2349
type: task
nature: improvement
title: flai guard reads a strategic role's tools from its role file rather than from .claude/agents
status: backlog
parent: S-0553
owner: alex
created: 2026-10-30T20:47:17Z
updated: 2026-10-30T20:47:17Z
transitions: []
stream: S-0553
tags: [cli]
touches: [flai/internal/guard/writes.go, flai/internal/guard/writes_test.go]
after: [T-2347]
---
# T-2349 flai guard reads a strategic role's tools from its role file rather than from .claude/agents

## Work

`ReadOnlyAgents` in `flai/internal/guard/writes.go` parses the `tools:` line of `.claude/agents/<type>.md` on each call. For the planner, the orchestrator, the analyzer, and the curator, read the tools from the generated definition of the role file instead, the project's or the embedded baseline, so the guard and the harness judge one list. The explorer and the verifier are still read from `.claude/agents`.

Second layer: it waits for the generation task, whose definition it reads.

## Done when

`writes_test.go` and `flai/cmd/guard_test.go` pass with their existing cases unchanged. A new case shows a strategic role judged read-only, or not, by its role file's `definition.tools` when no `.claude/agents/<name>.md` exists.

## Notes

Drafted by the planner of E-0024 on 2026-10-30. `ReadOnlyAgents` is at `writes.go:26` in S-0419's worktree that day.
