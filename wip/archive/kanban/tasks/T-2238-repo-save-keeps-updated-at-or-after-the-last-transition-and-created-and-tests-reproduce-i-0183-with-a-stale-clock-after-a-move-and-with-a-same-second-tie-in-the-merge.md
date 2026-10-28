---
id: T-2238
type: task
nature: remediation
title: Repo.Save keeps updated at or after the last transition and created, and tests reproduce I-0183 with a stale clock after a move and with a same-second tie in the merge
status: done
parent: S-0522
owner: alex
created: 2026-10-28T14:21:44Z
updated: 2026-10-28T14:38:31Z
transitions:
  - to: ready
    at: 2026-10-28T14:25:05Z
    by: agent-S-0522
  - to: in-progress
    at: 2026-10-28T14:25:06Z
    by: agent-S-0522
  - to: done
    at: 2026-10-28T14:38:31Z
    by: agent-S-0522
stream: S-0522
tags: [cli]
touches: [flai/internal/workitem/repo.go, flai/internal/workitem/itemlock.go, flai/internal/workitem/itemlock_test.go]
usage:
  source: log
  seconds: 805
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 20
      output: 8267
      cache_read: 881872
      cache_write: 46969
      cost: 0.6471
---
# T-2238 Repo.Save keeps updated at or after the last transition and created, and tests reproduce I-0183 with a stale clock after a move and with a same-second tie in the merge

## Work

First write the tests in `flai/internal/workitem/itemlock_test.go`, beside `TestSaveKeepsAMoveSavedSinceTheCopyWasRead`, and see them fail. One reproduces I-0183's shape: a story is moved to `ready` at a time T with `Transition` or `TransitionAll`; a second writer then reads the item with `repo.Get`, so its copy carries the transition, changes another key, such as `touches` or `forecast`, sets `Updated` to a stamp before T, as a writer whose clock was read before the move does, and saves; after the save the file's `updated` is T, not the earlier stamp, and `flai check`'s `history` reports nothing on it. The other covers the merge's tie: a copy read before the move, whose `Updated` the writer sets to a stamp before T, saved after a move whose stamp equals the base's `updated` to the second, so that `mergeParts` sees `updated` changed on one side only; the file's `updated` is at or after the last transition's `at` afterwards. A third checks `created`: a stamp before `created` is raised to it.

Then make them pass in `Repo.Save` (`flai/internal/workitem/repo.go`) and `mergeSaved` (`flai/internal/workitem/itemlock.go`): after the merge and before the write, when the document's `updated` is before the last transition's `at` or before `created`, raise it to the later of the two, in `workitem.TimeFormat`, and set it on the item too, so that the item is what was written, as `Save`'s contract says. Keep the merge's `max` for the case where both sides changed `updated`. Do not refuse the save: the writer cannot do better than write again with a later clock, and ADR-0184's `*ConflictError` stays for a key both sides changed. Keep the change inside the lock, and leave `Transition` and the other writers as they are: this task is the invariant, the task on the writers' clocks is the honest stamp.

This task waits for nothing. It shares no path with the task on the writers' clocks, so the two run together.

## Done when

The three tests exist in `flai/internal/workitem/itemlock_test.go`, failed before the change and pass after it, `Repo.Save` leaves `updated` at or after the last transition's `at` and `created` on every write, `go test ./flai/internal/workitem/...` and `flai test` on the files changed pass, and `flai check --strict` is clean on them.

## Notes

The check's rule is `history` in `flai/internal/check/check.go`: whole seconds, `workitem.TimeFormat`; equal stamps pass, only an earlier `updated` fails. The fixture `flai/internal/check/testdata/bad/wip/kanban/stories/S-003-wrong-name.md` already exercises the finding, so the check itself needs no new test.
