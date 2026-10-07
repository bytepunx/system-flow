---
id: T-1215
type: task
nature: feature
title: The dashboard design and the flaiover guide describe the Messages view and the messages on a story
status: backlog
parent: S-0336
owner: alex
created: 2026-10-07T20:17:09Z
updated: 2026-10-07T20:17:09Z
transitions: []
stream: S-0336
tags: [flaiover]
touches: [design/system/flaiover-dashboard.md, docs/users/flaiover.md]
after: [T-1214]
---
# T-1215 The dashboard design and the flaiover guide describe the Messages view and the messages on a story

## Work

Document what T-1211 to T-1214 built. It waits for T-1214, the last of them.

- `design/system/flaiover-dashboard.md`: the host API reads, the view, the story page, and the held card.
- `docs/users/flaiover.md`: where to find the agents' messages and what a share means on a card.

## Done when

- Both documents name the view and the reads.
- The markdown lint passes on the changed files.
- `flai check --strict` passes.

## Notes
