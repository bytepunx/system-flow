---
id: T-1736
type: task
nature: feature
title: The harness starts a defined role with its prompt from its role file's body, its definition, and its environment, on claude-code with --agent over that definition
status: backlog
parent: S-0417
owner: alex
created: 2026-10-16T05:44:46Z
updated: 2026-10-16T05:44:46Z
transitions: []
stream: S-0417
tags: [cli]
touches: [flai/internal/harness/harness.go, flai/internal/harness/adapters.go, flai/internal/harness/harness_test.go]
---
# T-1736 The harness starts a defined role with its prompt from its role file's body, its definition, and its environment, on claude-code with --agent over that definition

## Work

- `harness.Request` gains the role's name and its prompt body.
- `Prompt` gives a defined role its file's body, the item, and the run's inputs.
- `roleEnv` sets `FLAI_AGENT` to `<name>` or `<name>-<item>`, `FLAI_ROLE` to its verb, and `FLAI_ITEM` when it has one, and refuses a story for a role whose file says it works none.
- The claude-code adapter passes `--agent <definition>` over the `.claude/agents/<name>.md` its file names, beside the explorer and the verifier, and refuses a missing definition, naming `flai upgrade` or the file.
- Test a fixture role's request, its environment, and the refusal. Keep the planner's, the orchestrator's, and the analyzer's tests unchanged.

First layer: shares no path with T-1735 or T-1737.

## Done when

`go test ./internal/harness/...` passes.

## Notes
