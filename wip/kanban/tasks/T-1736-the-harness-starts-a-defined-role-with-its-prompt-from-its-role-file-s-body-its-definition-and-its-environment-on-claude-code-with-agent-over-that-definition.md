---
id: T-1736
type: task
nature: feature
title: The harness starts a defined role with its prompt from its role file's body, its definition, and its environment, on claude-code with --agent over that definition
status: done
parent: S-0417
owner: alex
created: 2026-10-16T05:44:46Z
updated: 2026-10-29T03:49:17Z
transitions:
  - to: ready
    at: 2026-10-29T03:07:00Z
    by: agent-S-0417
  - to: in-progress
    at: 2026-10-29T03:07:02Z
    by: agent-S-0417
  - to: done
    at: 2026-10-29T03:49:17Z
    by: agent-S-0417
stream: S-0417
tags: [cli]
touches: [flai/internal/harness/harness.go, flai/internal/harness/adapters.go, flai/internal/harness/harness_test.go]
usage:
  source: log
  seconds: 2534
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 79
      output: 315
      cache_read: 4200137
      cache_write: 173884
      cost: 1.4929
---
# T-1736 The harness starts a defined role with its prompt from its role file's body, its definition, and its environment, on claude-code with --agent over that definition

## Work

- `harness.Request` gains the role's name and its prompt body.
- `Prompt` gives a defined role its file's body, the item, and the run's inputs.
- `roleEnv` sets `FLAI_AGENT` to `<name>` or `<name>-<item>`, `FLAI_ROLE` to its verb, and `FLAI_ITEM` when it has one, and refuses a story for a role whose file says it works none.
- The claude-code adapter builds the role's definition from its role file's definition fields and body. It passes that definition in `--agents` beside the explorer and the verifier, and selects it with `--agent <name>`. It refuses a role whose definition fields are missing, naming `flai check`. It reads and writes nothing under `.claude/` (TH-0471, answer 6; TH-0570).
- Test a fixture role's request, its environment, its generated definition, and the refusal. Keep the planner's, the orchestrator's, and the analyzer's tests unchanged.

First layer: shares no path with T-1735 or T-1737.

## Done when

`go test ./internal/harness/...` passes.

## Notes

- The fourth Work bullet was rewritten by planner-S-0417 on TH-0570's approval: the definition is generated from the role file, not read from `.claude/agents/<name>.md`. If the operator answers TH-0567's question 2 otherwise, revisit it.
