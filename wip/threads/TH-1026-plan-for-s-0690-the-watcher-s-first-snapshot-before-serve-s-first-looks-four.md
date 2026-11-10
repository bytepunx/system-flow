---
id: TH-1026
title: "Plan for S-0690: the watcher's first snapshot before serve's first looks, four tasks in three layers"
anchor:
  path: wip/kanban/stories/S-0690-flai-serve-takes-a-project-s-watcher-s-first-snapshot-after-its-first-looks-so.md
  item: S-0690
status: open
participants: [planner-S-0690]
created: 2026-11-10T10:30:53Z
updated: 2026-11-10T10:30:53Z
---

# TH-1026 Plan for S-0690: the watcher's first snapshot before serve's first looks, four tasks in three layers

On wip/kanban/stories/S-0690-flai-serve-takes-a-project-s-watcher-s-first-snapshot-after-its-first-looks-so.md.

## Entries

### 2026-11-10T10:30:53Z planner-S-0690
S-0690 is planned. It has four tasks in three layers, 7 file touches, a 40m forecast, and a cost of delay of 12.50 USD a week.

## The fix

`watch.Watcher` gets a call that takes its first snapshot straight away. `reconcile` in `flai/internal/serve/serve.go` makes that call before `workitem.SnapshotItems` and before the first looks. `Run` then compares its first tick against that snapshot. A change written while the first looks run is reported within a tick, instead of waiting for the once-a-minute look.

## Tasks

| Layer | Task | Touches |
|-------|------|---------|
| 1 | T-2981 the watcher takes its first snapshot through a call of its own before `Run` | `flai/internal/watch/watch.go`, `watch_test.go` |
| 2 | T-2984 serve makes that call before the item snapshot and the first looks, with a test that reproduces I-0273; the probe wait in `TestServeStartsADefinedRoleOnAnEventAndOnASchedule` is removed (after T-2981) | `flai/internal/serve/serve.go`, `serve_test.go` |
| 3 | T-2988 `design/system/flai-cli.md` says when the first snapshot is taken (after T-2984) | `design/system/flai-cli.md` |
| 3 | T-2989 close I-0273 (after T-2984) | the issue file, `design/issues/summary.md` |

## Assumptions

- `Run` with no prior call still takes its own first snapshot, so the merge queue's `watchRefs` is unchanged.
- The once-a-minute `ItemsChanged` and `starter.look` stay as a safety net.
- The reproducing test may not fail every run, because the window is a race (5 of 25 runs under load). T-2984 asks for at least one failure in 20 runs on the old code, recorded in its Notes.

## Proposed, not drafted

`watchRefs` in `flai/internal/serve/mergequeue.go` works around the same window with a timer (`armed`) that holds the release word back for one tick. It could make the new call before `release.Watch` and drop the timer. I left it out of scope because I-0273 names only the project watcher. If you want it, say so and I will add it as a task after T-2981.
