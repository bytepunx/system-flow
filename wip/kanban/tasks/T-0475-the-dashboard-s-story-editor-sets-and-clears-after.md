---
id: T-0475
type: task
nature: feature
title: The dashboard's story editor sets and clears after
status: in-progress
parent: S-0130
owner: alex
created: 2026-09-26T17:55:02Z
updated: 2026-09-26T18:12:45Z
transitions:
  - to: ready
    at: 2026-09-26T17:55:12Z
    by: agent-S-0130
  - to: in-progress
    at: 2026-09-26T18:12:45Z
    by: agent-S-0130
stream: S-0130
tags: []
touches: [flaiover]
---
# T-0475 The dashboard's story editor sets and clears after

## Work

- Add an `after` field to the story editor in flaiover, next to touches, sending `after` to the host API's item edit.
- Start only when S-0129 and S-0133 are accepted (TH-0022); sync first.

## Done when

- The editor loads, sets, and clears `after`; flaiover's tests and check pass.

## Notes
