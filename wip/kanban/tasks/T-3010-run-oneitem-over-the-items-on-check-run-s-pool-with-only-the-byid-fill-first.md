---
id: T-3010
type: task
nature: improvement
title: Run oneItem over the items on check.Run's pool, with only the byID fill first and alone
status: done
parent: S-0677
owner: alex
created: 2026-11-10T10:31:27Z
updated: 2026-11-10T19:51:50Z
transitions:
  - to: ready
    at: 2026-11-10T19:33:56Z
    by: system-flow
  - to: in-progress
    at: 2026-11-10T19:33:56Z
    by: system-flow
  - to: done
    at: 2026-11-10T19:51:50Z
    by: agent-S-0677
stream: S-0677
tags: [flai, performance, check]
touches: [flai/internal/check/check.go, flai/internal/check/check_items.go, flai/internal/check/concurrent_test.go, flai/internal/check/items_test.go]
usage:
  source: log
  seconds: 1074
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 40
      output: 13319
      cache_read: 5152666
      cache_write: 163818
      cost: 2.2664
---
# T-3010 Run oneItem over the items on check.Run's pool, with only the byID fill first and alone

## Work

`check.Run` runs `workItems` first and alone, because it fills `byID`, the index the other rules read (`flai/internal/check/check.go`, `sequence`). `workItems` also runs `oneItem` on every item, one after another, which takes about 0.45 to 0.8 s of the critical path on a loaded host (I-0312). `oneItem` also calls `workitem.Children(c.items, it.ID)`, a scan of every item for each item, so the step grows with the square of the items: 3,349 of them when I-0312 was measured.

Split it in two:

- The first step stays alone: the duplicate-ID findings, the `byID` fill, and `c.write.pick`. It also builds an index of each item's children, once, which `oneItem` reads in place of `workitem.Children`. No rule writes it.
- The `oneItem` calls go on the pool beside the other rules, in chunks of items, each chunk into a `Result` of its own through `runRule`. The chunks merge in item order, at the place `workItems` held in `sequence`. A whole check's findings, their order, and the counts must stay what they are now. `item.parent-list` reports on the parent's file, so the order matters.

A checked write (`c.write` set, and `items.go`'s run on the items given) keeps checking only the items it picked. Waits for no task: it shares no file with the digest or conflict-marker tasks.

## Done when

- `workItems` alone only fills `byID` and the children index, reports duplicates, and picks a write's items. `oneItem` runs on the pool and reads the index.
- `TestRunAtOnceMatchesTheRulesOneAfterAnother` passes, and covers item findings from more than one chunk.
- `TestItemsMatchesRunOnTheItemsGiven` passes, and `flai test` passes on the changed paths.

## Notes

Drafted by the planner for S-0677.
