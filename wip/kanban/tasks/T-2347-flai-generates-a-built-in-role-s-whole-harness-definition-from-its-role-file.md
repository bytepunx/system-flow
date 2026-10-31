---
id: T-2347
type: task
nature: improvement
title: flai generates a built-in role's whole harness definition from its role file, and subAgents passes it with no .claude/agents read
status: backlog
parent: S-0553
owner: alex
created: 2026-10-30T20:47:00Z
updated: 2026-10-30T20:47:00Z
transitions: []
stream: S-0553
tags: [cli]
touches: [flai/internal/roles/definition.go, flai/internal/roles/definition_test.go, flai/internal/harness/adapters.go, flai/internal/harness/harness_test.go]
---
# T-2347 flai generates a built-in role's whole harness definition from its role file, and subAgents passes it with no .claude/agents read

## Work

Add `flai/internal/roles/definition.go`, which builds a role's harness definition from its file: `name`, `definition.description`, `definition.tools`, `definition.model` with the role's agent's model over it, and the body with its project additions as the prompt, without the marker line. `roleDefinition` in `flai/internal/harness/adapters.go` does this for a defined role today; move it there so the harness and `flai role definition` share it.

In `subAgents`, give the four built-in strategic roles the generated definition from `r.Builtin`, the project's file or the embedded baseline, in place of `readDefinition` on `.claude/agents/<name>.md`. Remove the refusal that names `flai upgrade` for a missing strategic definition. The explorer's and the verifier's definitions, and the sub-agents' models, are still read from `.claude/agents`.

First layer: it waits for nothing. The guard, check, and command tasks wait for it, since they use the generation it adds.

## Done when

`TestClaudeCodeStartsThePlanner`, `TestClaudeCodeStartsTheOrchestrator`, `TestClaudeCodeStartsTheAnalyzer`, and `TestClaudeCodeStartsTheCurator` pass with their `.claude/agents/<name>.md` fixtures removed and the same start prompt text. A case shows a definition's description and tools taken from the role file when `.claude/agents` holds a different copy. `TestClaudeCodeStartsADefinedRole` passes unchanged, and `definition_test.go` covers the model override and the stripped marker.

## Notes

Drafted by the planner of E-0024 on 2026-10-30. The lines named are S-0419's worktree as read that day: `subAgents` at `adapters.go:330`, `readDefinition` at `:423`, the refusal at `:343`.
