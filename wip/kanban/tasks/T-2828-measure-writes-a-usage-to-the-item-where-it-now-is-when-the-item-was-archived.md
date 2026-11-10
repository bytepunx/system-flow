---
id: T-2828
type: task
nature: remediation
title: Measure writes a usage to the item where it now is when the item was archived since it was read
status: done
parent: S-0648
owner: alex
created: 2026-11-10T06:37:26Z
updated: 2026-11-10T08:35:36Z
transitions:
  - to: ready
    at: 2026-11-10T08:28:02Z
    by: agent-S-0648
  - to: in-progress
    at: 2026-11-10T08:28:04Z
    by: agent-S-0648
  - to: done
    at: 2026-11-10T08:35:36Z
    by: agent-S-0648
stream: S-0648
tags: [flai, serve]
touches: [flai/internal/serve/usage.go, flai/internal/serve/usage_test.go]
after: [T-2827]
usage:
  source: log
  seconds: 450
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 32
      output: 284
      cache_read: 1785591
      cache_write: 60046
      cost: 0.679
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
