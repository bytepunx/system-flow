---
id: T-2406
type: task
nature: improvement
title: The planner's role file, its baseline, and its start prompt say after names a story whose output this one reads, with the reason in its Notes
status: backlog
parent: S-0551
owner: alex
created: 2026-10-31T07:44:14Z
updated: 2026-10-31T07:44:14Z
transitions: []
stream: S-0551
tags: [flai, conventions]
touches: [design/workflows/roles/planner.md, template/root/design/workflows/roles/planner.md, flai/internal/roles/baseline/planner.md, flai/internal/harness/harness.go, flai/internal/harness/harness_test.go, template/root/.claude/agents/planner.md, ".claude/agents/planner.md"]
after: [T-2388]
---
# T-2406 The planner's role file, its baseline, and its start prompt say after names a story whose output this one reads, with the reason in its Notes

## Work

Replace the planner's `after` rule, "only a story whose accepted change it needs on main, never an order of preference", with the one the ADR records: `after` names a story whose output this story reads, and the story's Notes name it with the reason. Change it in every copy, above each role file's baseline marker:

- `design/workflows/roles/planner.md` and `template/root/design/workflows/roles/planner.md`.
- `flai/internal/roles/baseline/planner.md`, flai's own baseline.
- The planner's start prompt in `flai/internal/harness/harness.go`, with the strings `harness_test.go` asserts.
- The generated definitions `template/root/.claude/agents/planner.md` and `.claude/agents/planner.md`. A write under `.claude/` asks the operator; send the file's whole text on the thread it opens.

It waits for the ADR task, whose wording it carries.

## Done when

- Every copy states the same rule, and `flai check --strict` reports no role or definition drift.
- `flai test flai/internal/harness flai/internal/roles` passes.

## Notes

S-0553 and S-0554 in the backlog move the start prompt into the role file; whichever lands second carries this rule over.
