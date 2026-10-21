---
id: T-1933
type: task
nature: remediation
title: Name the threads the orchestrator answers by awaiting operator in its instructions, its agent definition, and the strategic-agents convention
status: done
parent: S-0453
owner: alex
created: 2026-10-20T22:21:02Z
updated: 2026-10-21T00:54:27Z
transitions:
  - to: ready
    at: 2026-10-21T00:41:50Z
    by: system-flow
  - to: in-progress
    at: 2026-10-21T00:41:51Z
    by: system-flow
  - to: done
    at: 2026-10-21T00:54:27Z
    by: system-flow
stream: S-0453
tags: [flai, template]
touches: [flai/internal/harness/harness.go, flai/internal/harness/harness_test.go, design/conventions/strategic-agents.md, template/root/design/conventions/strategic-agents.md, design/system/strategic-agents.md, ".claude/agents/orchestrator.md", template/root/.claude/agents/orchestrator.md, template/CHANGELOG.md]
after: [T-1927]
usage:
  source: log
  seconds: 756
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 30
      output: 9654
      cache_read: 1901964
      cache_write: 79995
      cost: 1.0643
---
# T-1933 Name the threads the orchestrator answers by awaiting operator in its instructions, its agent definition, and the strategic-agents convention

## Work

The orchestrator is told today to "take from inbox the threads awaiting the operator, those whose status is open and whose last entry is by a story's agent". Name them by the new field instead, as criterion 3 asks:

- `orchestratePrompt` in `flai/internal/harness/harness.go`, and the comment above it: take from `inbox` the threads with `awaiting: "operator"`, leaving out those you opened and those with a `pending_recommendation`. Say that `inbox` lists them without `all: true`.
- The S-0220 test in `flai/internal/harness/harness_test.go` that asserts the prompt's text.
- The `answer_threads` rule under § As the orchestrator in `design/conventions/strategic-agents.md`, its template copy, and `design/system/strategic-agents.md` § Answering threads.
- Step 9 of `.claude/agents/orchestrator.md` and of its template copy. Writes under `.claude/` ask the operator through `permission_prompt`.
- A line in `template/CHANGELOG.md` for the convention and the agent definition.

It waits for T-1927, whose ADR it links. It does not wait for the code: it changes only words and the prompt's test, and runs beside T-1929.

## Done when

- The prompt, the convention, the design, and both agent definitions name the threads by `awaiting: "operator"`, and none says "whose last entry is by a story's agent".
- `flai test flai/internal/harness` and the markdown lint pass.

## Notes

The convention's rule is above its baseline marker. This repository develops the template, so the change lands in `design/conventions` and `template/` in the same story, as CLAUDE.md says.
