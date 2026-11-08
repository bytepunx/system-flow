---
id: T-2622
type: task
nature: improvement
title: itemnew.Create checks a creation with the item-write rule set before and after, in place of check.Run
status: done
parent: S-0598
owner: alex
created: 2026-11-07T23:47:47Z
updated: 2026-11-08T00:34:22Z
transitions:
  - to: ready
    at: 2026-11-08T00:08:26Z
    by: system-flow
  - to: in-progress
    at: 2026-11-08T00:08:31Z
    by: system-flow
  - to: done
    at: 2026-11-08T00:34:22Z
    by: agent-S-0598
stream: S-0598
tags: [cli]
touches: [flai/internal/itemnew/itemnew.go, flai/internal/itemnew/itemnew_test.go]
after: [T-2620]
usage:
  source: log
  seconds: 1551
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 53
      output: 18412
      cache_read: 8933805
      cache_write: 182748
      cost: 3.2665
---
# T-2622 itemnew.Create checks a creation with the item-write rule set before and after, in place of check.Run

## Work

In `flai/internal/itemnew/itemnew.go`, `Create` runs `check.Run` before and after it writes the item. Run the item-write rule set T-2620 adds in both places instead. The "before" run takes the item's parent and the paths the write will make, so it covers the same rules on the same files as the "after" run. Keep the scoping with `check.ScopeToItem`, the known set, the advisory rules that come back in `Warnings`, and the undo as they are (ADR-0213).

This task waits for T-2620, whose rule set it calls. It shares no file with the edit's task, so the two run together.

## Done when

- `Create` no longer calls `check.Run`.
- Tests in `itemnew_test.go` create a story and a task that each bring a finding of each kept rule. Each is refused as before, with nothing left behind, and the advisory rules still come back as warnings.
- A test shows that a finding of a skipped rule already in the repository, such as a documentation finding, refuses nothing.
- `flai test` passes on the paths this task changed.

## Notes
