---
id: T-2697
type: task
nature: feature
title: item.get copies only the item and its children it answers, and reads the rest of the list in place
status: done
parent: S-0606
owner: alex
created: 2026-11-08T01:41:02Z
updated: 2026-11-08T01:56:06Z
transitions:
  - to: ready
    at: 2026-11-08T01:42:04Z
    by: agent-S-0606
  - to: in-progress
    at: 2026-11-08T01:42:11Z
    by: agent-S-0606
  - to: done
    at: 2026-11-08T01:56:06Z
    by: agent-S-0606
stream: S-0606
tags: []
touches: [flai/internal/workitem/view.go, flai/internal/workitem/view_test.go, flai/internal/hostapi/hostapi.go, flai/internal/hostapi/hostapi_test.go]
usage:
  source: log
  seconds: 835
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 76
      output: 24308
      cache_read: 12491609
      cache_write: 579163
      cost: 6.3489
---

# T-2697 item.get copies only the item and its children it answers, and reads the rest of the list in place

## Work

T-2619's profile found `item.get` still spends about 27 ms of CPU a call with the snapshot on: `Repo.List` deep-copies all 3,347 items (46%) and sorts them (20%). In a new file `flai/internal/workitem/view.go`, add a read of the items in place: one that calls a function with each kept item and whether it is archived, never handing out a copy, and one that returns copies of an item's children in ID order. `item.get` in `flai/internal/hostapi/hostapi.go` gets the item with `Repo.Get`, its children with the new call, and works out the cost per agent hour from the items read in place, so it copies only what it answers. `store.go` and `repo.go` are not changed: S-0606 agreed with S-0600 and S-0599 (MS-0354, MS-0352) to change them no further.

Waits for no open task: T-2615 and T-2617, whose code it builds on, are done.

## Done when

- [ ] `item.get`'s answers are the same as before, field for field, in `hostapi_test.go`'s tests.
- [ ] Tests in `flai/internal/workitem/view_test.go` show the children come in ID order, as copies, with `Archived` set, and that the in-place read sees every item once.

## Notes
