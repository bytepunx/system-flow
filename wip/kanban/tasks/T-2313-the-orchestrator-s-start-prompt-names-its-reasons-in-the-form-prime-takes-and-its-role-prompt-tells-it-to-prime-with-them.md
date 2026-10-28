---
id: T-2313
type: task
nature: improvement
title: The orchestrator's start prompt names its reasons in the form prime takes, and its role prompt tells it to prime with them
status: backlog
parent: S-0537
owner: alex
created: 2026-10-28T22:08:01Z
updated: 2026-10-28T22:08:01Z
transitions: []
stream: S-0537
tags: [flai, orchestrator, template]
touches: [flai/internal/harness/harness.go, flai/internal/harness/harness_test.go, design/workflows/roles/orchestrator.md, template/root/design/workflows/roles/orchestrator.md, ".claude/agents/orchestrator.md", template/root/.claude/agents/orchestrator.md]
after: [T-2311, T-2312]
---
# T-2313 The orchestrator's start prompt names its reasons in the form prime takes, and its role prompt tells it to prime with them

## Work

Criterion 3. S-0536's T-2303 lists each reason with its item or thread in `orchestratePrompt` (`flai/internal/harness/harness.go`). Write each in the form T-2308 parses, and replace the prime sentence (today: "Prime your session with flai prime --role orchestrate ...") with one that primes with them: `flai prime --role orchestrate --reason <r> ...`, or the MCP tool `prime` with role `orchestrate` and `reasons`, reading every part. A run with no reasons, such as one started by hand, keeps today's sentence and the general pack.

In the role file, `design/workflows/roles/orchestrator.md` step 1, and its template copy, say the same: call `prime` with role `orchestrate` and the reasons the start prompt names. Regenerate `.claude/agents/orchestrator.md` and its template copy from the role file with flai's generator, as `flai check` reports `roles.definition` until they match. Writes under `.claude/` need the operator: when the harness refuses one, send the file's whole contents on a thread on the story for the operator to paste, then check it.

Waits for T-2311 and T-2312: the prompt names the flag and the tool's field, which must exist first.

## Done when

- A test in `harness_test.go` builds the orchestrator's prompt with three reasons, a thread, a review, and `ready`, and finds each in the form `prime` takes and the prime sentence naming them; one with no reasons keeps today's sentence.
- `TestClaudeCodeStartsTheOrchestrator` and the prompt tests that assert the prime line pass with the new line.
- The role file, its template copy, and the two definitions match, and `scripts/flai.sh check --strict` is clean.

## Notes
