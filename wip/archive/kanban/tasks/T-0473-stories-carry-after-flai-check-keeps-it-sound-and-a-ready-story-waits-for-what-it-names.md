---
id: T-0473
type: task
nature: feature
title: "Stories carry after:, flai check keeps it sound, and a ready story waits for what it names"
status: done
parent: S-0130
owner: alex
created: 2026-09-26T17:55:02Z
updated: 2026-09-26T17:59:29Z
transitions:
  - to: ready
    at: 2026-09-26T17:55:12Z
    by: agent-S-0130
  - to: in-progress
    at: 2026-09-26T17:55:13Z
    by: agent-S-0130
  - to: done
    at: 2026-09-26T17:59:29Z
    by: agent-S-0130
stream: S-0130
tags: []
touches: [flai/internal/workitem, flai/internal/check, flai/internal/serve]
---
# T-0473 Stories carry after:, flai check keeps it sound, and a ready story waits for what it names

## Work

- Add `after: [S-nnnn, …]` to the item schema: read, written, and kept in the round trip.
- Add the `after` hold code to `workitem.Holds`: a ready story is held while a story it names is not done; the reason names each with its state, says a cancelled one was cancelled, and says what clears it. Named stories are looked up in the archive too, so a done story that was archived clears the hold.
- Combine with an overlap hold: the reason carries both.
- `flai check` reports an `after:` entry that names no story, names the story itself, names a non-story, or forms a cycle, and `after:` on an epic or task.
- flai serve's launcher and `agent.status`, `wait_for_work`, `flai board`, and `inbox` go through `Holds`, so they hold on `after` as on overlap; tests prove each.

## Done when

- Behaviour tests cover the hold (open, done, archived done, cancelled, missing, combined with overlap) and each check finding; the launcher and `wait_for_work` skip a story held by `after`.
- `make test` and lint pass.

## Notes
