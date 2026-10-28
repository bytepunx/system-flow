---
id: T-2145
type: task
nature: improvement
title: The orchestrator's prompt and its agent definition tell it to mark a recommendation on money or a host change for the operator alone, and that flai serve confirms the rest after confirm_after
status: in-progress
parent: S-0502
owner: alex
created: 2026-10-26T04:45:16Z
updated: 2026-10-28T13:05:51Z
transitions:
  - to: ready
    at: 2026-10-28T13:05:50Z
    by: agent-S-0502
  - to: in-progress
    at: 2026-10-28T13:05:51Z
    by: agent-S-0502
stream: S-0502
tags: [flai, template]
touches: [flai/internal/harness/harness.go, flai/internal/harness/harness_test.go, ".claude/agents/orchestrator.md", template/root/.claude/agents/orchestrator.md]
after: [T-2143]
usage:
  source: log
  seconds: 2399
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 153
      output: 52700
      cache_read: 10778085
      cache_write: 333737
      cost: 5.2297
---
# T-2145 The orchestrator's prompt and its agent definition tell it to mark a recommendation on money or a host change for the operator alone, and that flai serve confirms the rest after confirm_after

## Work

Point 3 of T-2141's ADR, on the orchestrator's side:

- `flai/internal/harness/harness.go`, `orchestratePrompt`: after the sentences on `recommend` and `autonomous`, one on the mark: when a recommendation concerns money (a cost of delay input, an estimate, spend) or a host change (an install, a key, a host action), give `operator_only: true` (`--operator-only` with `--recommend`), so that flai serve never confirms it by timeout; otherwise, while `orchestration.permissions.confirm_after` is set, flai serve confirms it as the answer once it has waited that long with no entry from the operator, and the orchestrator still never confirms one itself. `harness_test.go` asserts the prompt names `operator_only` and `confirm_after`.
- `.claude/agents/orchestrator.md`, point 9, and `template/root/.claude/agents/orchestrator.md`: the same rule in one sentence each, the two copies kept identical. A write under `.claude/` is a protected path: the story's agent makes it itself, `permission_prompt` asks the operator on a thread (ADR-0086), and a story whose branch changes such a path is accepted by the operator only (ADR-0106), so the story's agent says so in the narrative and the close-out.

## Done when

- The prompt test passes and names both terms.
- Both agent definitions carry the sentence, and `flai check --strict` scoped to the story is clean.

## Notes

Layer 3, waits for T-2143 (the parameter's name). The docs task waits for it. It shares no path with the serve look or the metrics task.
