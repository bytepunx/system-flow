---
id: T-2292
type: task
nature: feature
title: flai appends the events it emits to a log under .flai-cache, and inbox.CatchUp reads them back once, described, under the 50-event cap, waking wait_for_events
status: done
parent: S-0535
owner: alex
created: 2026-10-28T21:37:37Z
updated: 2026-10-28T22:44:49Z
transitions:
  - to: ready
    at: 2026-10-28T22:31:51Z
    by: system-flow
  - to: in-progress
    at: 2026-10-28T22:31:52Z
    by: system-flow
  - to: done
    at: 2026-10-28T22:44:49Z
    by: agent-S-0535
stream: S-0535
tags: [flai, events, inbox, mcp]
touches: [flai/internal/inbox/events.go, flai/internal/inbox/events_test.go, flai/internal/inbox/inbox.go, flai/internal/inbox/inbox_test.go, flai/internal/mcpserver/server.go, flai/internal/mcpserver/server_test.go]
after: [T-2291]
usage:
  source: log
  seconds: 777
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 47
      output: 285
      cache_read: 1910186
      cache_write: 108138
      cost: 0.6875
---
# T-2292 flai appends the events it emits to a log under .flai-cache, and inbox.CatchUp reads them back once, described, under the 50-event cap, waking wait_for_events

## Work

`planned`, `ready-empty`, and `release-due` happen in flai serve and in `flai accept`, not in a work item file, and `workitem.Changes` derives nothing from them. They need a log, as overlaps do (`itemedit.Overlaps`, `.flai-cache/overlaps.jsonl`).

- Add `flai/internal/inbox/events.go`: one function that appends an event (`at`, `by`, `kind`, `id`, `to`, and the kind's fields: the IDs a plan wrote, the release figures that changed) to `.flai-cache/events.jsonl` in the main checkout (`workitem.Repo.CacheDir`), safe for several writers, and one that reads the entries after a time.
- In `inbox.go`, add the three kinds beside `Edited`, `Overlapped`, and `Message`. Make `CatchUp` read the log after the reader's cursor and merge it with the other events under `MaxEvents`. Give `Describe` a case for each kind.
- In `mcpserver/server.go`, add the log to `noticeLogs()` so that an append wakes `wait_for_events`.
- If S-0532's agent event log (`flai/internal/inbox/agentevents.go`) is on main when this story starts, append these kinds to that log through its helper instead of adding a second one, and narrow this task's touches.

It waits for T-2291, whose ADR settles the names and the fields.

## Done when

- A test appends one event of each kind and reads each back through `CatchUp` once, with its fields, and not again after the cursor moves.
- A test shows the new events counted under the 50-event cap with the others, newest kept.
- An append wakes `wait_for_events`, with a test.
- `flai test` on the changed paths passes.

## Notes

The writers come in the next layer: T-2294 (`planned`), T-2295 (`ready-empty`), and T-2296 (`release-due`).
