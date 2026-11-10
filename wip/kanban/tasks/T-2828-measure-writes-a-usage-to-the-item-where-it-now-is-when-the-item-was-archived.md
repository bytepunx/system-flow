---
id: T-2828
type: task
nature: remediation
title: Measure writes a usage to the item where it now is when the item was archived since it was read
status: backlog
parent: S-0648
owner: alex
created: 2026-11-10T06:37:26Z
updated: 2026-11-10T06:37:26Z
transitions: []
stream: S-0648
tags: [flai, serve]
touches: [flai/internal/serve/usage.go, flai/internal/serve/usage_test.go]
after: [T-2827]
---
# T-2828 Measure writes a usage to the item where it now is when the item was archived since it was read

## Work

Second layer: it waits for T-2827, whose refusal of a `Save` on a moved item, such as `ErrItemMoved`, this task handles.

- In `Measure`'s `put` (`usage.go`), when `repo.Save` is refused because the item moved, read the item again with `repo.Get`, archived items included, and save the usage on that copy. Try once. If it is refused again, return the error as `put` returns any other.
- So a task archived between `put`'s read and its save gets its usage in its archived file, and no copy is left at its old path. That is how the harness-priced usage of T-2760 was lost from its archive.
- Add a test in `usage_test.go` that archives a done task between `put`'s read and its save, through a seam or a hook the test sets. It checks that the task's file exists only under `wip/archive` and carries the measured usage.

## Done when

- The new test passes, and fails with the retry taken out.
- `go test ./internal/serve -run 'Measure|Usage' -count 10` passes in `flai/`.
- `flai test flai/internal/serve/usage.go flai/internal/serve/usage_test.go` passes.

## Notes
