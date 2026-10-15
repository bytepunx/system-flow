---
id: T-1547
type: task
nature: feature
title: The dashboard's design and user docs say how tokens are changed on the Design page and saved
status: done
parent: S-0383
owner: alex
created: 2026-10-13T17:47:46Z
updated: 2026-10-15T02:52:37Z
transitions:
  - to: ready
    at: 2026-10-15T02:49:44Z
    by: agent-S-0383
  - to: in-progress
    at: 2026-10-15T02:49:45Z
    by: agent-S-0383
  - to: done
    at: 2026-10-15T02:52:37Z
    by: agent-S-0383
stream: S-0383
tags: [dashboard]
touches: [design/system/flaiover-dashboard.md, docs/users/flaiover.md]
after: [T-1546]
usage:
  source: log
  seconds: 172
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 36
      output: 25085
      cache_read: 2284809
      cache_write: 89489
      cost: 1.4518
---
# T-1547 The dashboard's design and user docs say how tokens are changed on the Design page and saved

## Work

Describe the Edit mode, the live preview of unsaved changes, the contrast marks, the save through `doc.save` and its refusals, and the leave guard in `design/system/flaiover-dashboard.md` and `docs/users/flaiover.md`. It waits for T-1546, whose behaviour it describes.

## Done when

- Both documents describe it, and the markdown lint passes on them.

## Notes

- Layer 4.
