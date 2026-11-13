---
id: T-3142
type: task
nature: improvement
title: flai unblock closes an interval left open on a closed item at the item's close, archived items included, and item.blocked-closed names that fix
status: done
parent: S-0723
owner: alex
created: 2026-11-12T15:53:47Z
updated: 2026-11-12T17:05:31Z
transitions:
  - to: ready
    at: 2026-11-12T16:54:47Z
    by: agent-S-0723
  - to: in-progress
    at: 2026-11-12T16:54:47Z
    by: agent-S-0723
  - to: done
    at: 2026-11-12T17:05:31Z
    by: agent-S-0723
stream: S-0723
tags: [cli]
touches: [flai/internal/workitem/rules.go, flai/internal/workitem/rules_test.go, flai/cmd/move.go, flai/cmd/move_test.go, flai/internal/check/check_items.go, flai/internal/check/items_test.go, docs/users/flai-reference.md]
after: [T-3141]
usage:
  source: log
  seconds: 644
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 32
      output: 7909
      cache_read: 3218033
      cache_write: 76026
      cost: 1.4101
---
# T-3142 flai unblock closes an interval left open on a closed item at the item's close, archived items included, and item.blocked-closed names that fix

## Work

T-3141 stops new closed items from keeping an open interval. This task gives the records that already have one, S-0678 in the archive and any an older flai or a hand edit leaves, a repair the operator runs once in the main checkout.

- `UnblockItem` in `flai/internal/workitem/rules.go`: on a closed item, `done` or `cancelled`, close the open interval at the item's close, `CompletedAt()`, not at now, so that `flai stats` counts the block only up to the close. An open item keeps closing at now.
- `flai unblock` in `flai/cmd/move.go` reaches an archived item: `repo.Get` reads the archive, so check that `repo.Save` writes the item back to its archived path and does not refuse it, and fix what stops it. Say in the command's output that the interval was closed at the item's close when it was.
- `item.blocked-closed` in `flai/internal/check/check_items.go` names the fix in its message: `flai unblock <ID> in the main checkout`.

Tests: `flai/internal/workitem/rules_test.go` for the close time on a closed and an open item; `flai/cmd/move_test.go` for `flai unblock` on an archived done story with an open interval, which ends closed at its done transition and leaves `flai check` with no `item.blocked-closed`; `flai/internal/check/items_test.go` for the message.

Waits for T-3141: both change `flai/internal/workitem/rules.go` and its test file.

## Done when

- `flai unblock` on an archived done story closes its open interval at its done transition's time, and `flai check` then reports no `item.blocked-closed` for it.
- `item.blocked-closed` names `flai unblock` as its fix.
- `flai test` on the six paths passes.

## Notes
