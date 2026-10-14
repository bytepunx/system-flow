---
id: T-1582
type: task
nature: feature
title: The dashboard design and the user docs describe the undone-path cards on the review page and the badge on the board and the story page
status: done
parent: S-0388
owner: alex
created: 2026-10-14T21:09:18Z
updated: 2026-10-14T23:38:42Z
transitions:
  - to: ready
    at: 2026-10-14T23:29:45Z
    by: agent-S-0388
  - to: in-progress
    at: 2026-10-14T23:29:46Z
    by: agent-S-0388
  - to: done
    at: 2026-10-14T23:38:42Z
    by: agent-S-0388
stream: S-0388
tags: [dashboard, docs]
touches: [design/system/flaiover-dashboard.md, docs/users/flaiover.md, docs/users/flai.md]
after: [T-1579, T-1580, T-1581]
usage:
  source: log
  seconds: 536
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 84
      output: 25492
      cache_read: 3701005
      cache_write: 136882
      cost: 2.0468
---
# T-1582 The dashboard design and the user docs describe the undone-path cards on the review page and the badge on the board and the story page

## Work

- In `design/system/flaiover-dashboard.md`:
  - § Views: the review page's cards and their three actions, and the badge on the board's card and the story page.
  - § Writes (S-0075, ADR-0029): the new route `api/items/[id]/restore` behind `stream.restore`, and `keep_removals` on the accept route.
- In `docs/users/flaiover.md`:
  - § Reviewing a story: how to restore a path, keep a removal with a reason, or send the story back, and why Accept stays disabled until each card is resolved.
  - § Board, in the card table, and § The item page: the badge and where it leads.
- Say that the orchestrator never keeps a removal, and link S-0371's ADR rather than restate it.
- It waits for T-1579, T-1580, and T-1581, whose behaviour it describes.

## Done when

- [ ] Both documents describe the cards, their three actions, and the badge as built, with `updated` bumped
- [ ] `flai test design/system/flaiover-dashboard.md docs/users/flaiover.md` passes

## Notes
