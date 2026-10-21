---
id: T-2007
type: task
nature: remediation
title: Repo.Save merges an item copy read before a concurrent write with the file on disk, key by key, under one lock, and refuses a key both changed
status: done
parent: S-0479
owner: alex
created: 2026-10-21T12:12:53Z
updated: 2026-10-21T12:40:28Z
transitions:
  - to: ready
    at: 2026-10-21T12:13:36Z
    by: system-flow
  - to: in-progress
    at: 2026-10-21T12:13:36Z
    by: system-flow
  - to: done
    at: 2026-10-21T12:40:28Z
    by: agent-S-0479
stream: S-0479
tags: []
touches: [flai/internal/workitem/itemlock.go, flai/internal/workitem/itemlock_test.go, flai/internal/workitem/item.go, flai/internal/workitem/repo.go, flai/cmd/accept_resume_test.go, flai/cmd/check_stats_test.go, flai/cmd/issue_test.go]
usage:
  source: log
  seconds: 1612
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 92
      output: 32838
      cache_read: 6985500
      cache_write: 240874
      cost: 3.7058
---
# T-2007 Repo.Save merges an item copy read before a concurrent write with the file on disk, key by key, under one lock, and refuses a key both changed

## Work

Every item writer reads the item, changes its copy, and saves it with `Repo.Save`, which writes the whole file. A writer that holds its copy for long, as `itemedit.Apply` does across a whole-repository `flai check`, writes back the fields another writer changed meanwhile (I-0175: the replanner's forecast write undid `flai story start`'s move to in-progress). An item read from a file keeps the canonical form it was read in; `Repo.Save` takes one lock for item writes under the cache folder, reads the file again, and when it changed since the copy was read, merges by top-level front-matter key and the body: a key only the copy changed takes the copy's value, one only the file changed keeps the file's, `updated` takes the later, and a key both changed differently refuses the save with a conflict error naming the item and the keys, writing nothing. The copy is then what was written. First in the plan: the other tasks build on it.

## Done when

- A test reads one item twice, moves one copy to in-progress and saves it, sets the other's forecast and saves it, and finds both in the file; it fails on main.
- A test where both copies change the status finds the save refused with the conflict error and the first write kept.
- A test saves concurrently from several goroutines, each a different key, and loses none.
- `flai test` on the changed paths passes.

## Notes
