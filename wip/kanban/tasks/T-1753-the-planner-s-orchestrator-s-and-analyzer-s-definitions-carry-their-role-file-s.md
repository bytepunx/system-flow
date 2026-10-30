---
id: T-1753
type: task
nature: improvement
title: The planner's, orchestrator's, and analyzer's definitions carry their role file's body, its additions included, when the harness is given the file
status: done
parent: S-0419
owner: alex
created: 2026-10-16T05:46:02Z
updated: 2026-10-30T12:59:37Z
transitions:
  - to: ready
    at: 2026-10-30T12:59:29Z
    by: agent-S-0419
  - to: in-progress
    at: 2026-10-30T12:59:30Z
    by: agent-S-0419
  - to: done
    at: 2026-10-30T12:59:37Z
    by: agent-S-0419
stream: S-0419
tags: [cli, template]
touches: [flai/internal/harness/harness.go, flai/internal/harness/adapters.go, flai/internal/harness/harness_test.go]
after: [T-1754]
usage:
  source: log
  seconds: 7
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 3
      output: 782
      cache_read: 280864
      cache_write: 7307
      cost: 0.1228
---
# T-1753 The planner's, orchestrator's, and analyzer's definitions carry their role file's body, its additions included, when the harness is given the file

## Work

- Add `Builtin *roles.Role` to `harness.Request`. It holds the role file a built-in role runs from, the project's or flai's embedded copy (`roles.Builtin`).
- In `subAgents`, a built-in role's definition still reads `.claude/agents/<name>.md` for its description, tools, and model, and keeps the refusal that names `flai upgrade`. When the request carries `Builtin`, the definition's prompt is the file's body, read through `roleBody`, so the text a project writes below the marker reaches it (ADR-0207 § 4, TH-0760).
- The start prompts, `planPrompt`, `orchestratePrompt`, and `analyzePrompt`, stay in Go and keep their text.
- Test that a planner request carrying a role file with an addition passes that addition in its `--agents` prompt.

Second layer, after T-1754, whose `roles.Builtin` the serve tasks pass here.

## Done when

`flai/internal/harness/harness_test.go` passes with `TestClaudeCodeStartsThePlanner`, `TheOrchestrator`, and `TheAnalyzer` unchanged, and the new test shows the addition in the planner's definition.

## Notes

- Rewritten by agent-S-0419 on TH-0760. A role file's body is the definition, not the start prompt (ADR-0207 § 4). Every option on the thread needs this part.
