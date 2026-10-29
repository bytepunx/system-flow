---
id: T-2313
type: task
nature: improvement
title: The orchestrator's start prompt names its reasons in the form prime takes, and its role prompt tells it to prime with them
status: backlog
parent: S-0537
owner: alex
created: 2026-10-28T22:08:01Z
updated: 2026-10-29T15:24:14Z
transitions: []
stream: S-0537
tags: [flai, orchestrator, template]
touches: [flai/internal/harness/harness.go, flai/internal/harness/harness_test.go, flai/internal/serve/orchestrate_reasons.go, flai/internal/serve/orchestrate_reasons_test.go, flai/internal/serve/orchestrate.go, design/workflows/roles/orchestrator.md, template/root/design/workflows/roles/orchestrator.md]
after: [T-2311, T-2312]
---
# T-2313 The orchestrator's start prompt names its reasons in the form prime takes, and its role prompt tells it to prime with them

## Work

Criterion 3, ADR-0217 §6. S-0536 lists each reason in `orchestratePrompt` (`flai/internal/harness/harness.go`) from `orchestrateReason.line()` (`flai/internal/serve/orchestrate_reasons.go`).

- In serve, give each reason the form `prime` takes:

  | Reason | Form |
  |---|---|
  | moved to review | `review:S-nnnn` |
  | moved to backlog | `draft:S-nnnn` |
  | a draft's edit | `draft:S-nnnn` |
  | planned | `planned:<id>` |
  | thread | `thread:TH-nnnn` |
  | release-due | `release` |
  | ready-empty | `ready` |
  | reordered | `ready` |

  `asked`, `schedule`, and an event with no form have none. Pass the forms to the harness in a new `harness.Request` field beside `Reasons`, deduplicated and in order, and only when every reason has a form. A run with an `asked` or `schedule` reason, or one with no form, primes with the general pack.
- In the prompt, write each reason's form beside its line. Replace the prime sentence with one that primes with them: `flai prime --role orchestrate --reason <r> ...`, or the MCP tool `prime` with role `orchestrate` and `reasons`, reading every part. A run with no forms keeps today's sentence and the general pack.
- In the role file, `design/workflows/roles/orchestrator.md` step 1, and its template copy, say the same: call `prime` with role `orchestrate` and the reasons the start prompt names, else the general pack.

The two `.claude/agents/orchestrator.md` copies are T-2317's, since they are protected (D-0411).

Waits for T-2311 and T-2312: the prompt names the flag and the tool's field, which must exist first.

## Done when

- A test in `harness_test.go` builds the orchestrator's prompt with three reasons, a thread, a review, and `ready`, and finds each in the form `prime` takes and the prime sentence naming them. One with no forms keeps today's sentence.
- A test in `orchestrate_reasons_test.go` gives each reason kind its form, and checks that `context` parses every form serve writes.
- `TestClaudeCodeStartsTheOrchestrator` and the prompt tests that assert the prime line pass with the new line.
- The role file and its template copy match, and `flai test` passes on the changed paths.

## Notes
