---
id: T-1775
type: task
nature: feature
title: A formula's role step starts its role with the step's prompt and inputs through the defined-role path and is done when the run ends worked, its cost charged to the instance
status: backlog
parent: S-0422
owner: alex
created: 2026-10-16T05:47:45Z
updated: 2026-10-16T05:47:45Z
transitions: []
stream: S-0422
tags: [cli]
touches: [flai/internal/serve/formula_role.go, flai/internal/serve/formula_role_test.go, flai/internal/serve/usage.go, flai/internal/serve/usage_test.go]
after: [T-1774]
---
# T-1775 A formula's role step starts its role with the step's prompt and inputs through the defined-role path and is done when the run ends worked, its cost charged to the instance

## Work

- A role step calls `RunRole` with the step's prompt as the run's input, the instance story as its item when the role takes a story, and the trigger `formula <name> <S-nnnn>`.
- Wait for the run to settle. `worked` moves the task to done; `failed` or `asked` hands the step to the failure path.
- Measure each command and script step's seconds onto its task, and charge each role run's usage to the instance story, so the story's usage is measured as any story's is.
- Test with a fake harness.

Waits for the command step runner it extends.

## Done when

`go test ./internal/serve/...` passes.

## Notes
