---
id: T-1414
type: task
nature: feature
title: The user guide says what gateway, harness, and estimated mean beside a cost
status: backlog
parent: S-0358
owner: alex
created: 2026-10-08T08:51:52Z
updated: 2026-10-08T08:51:52Z
transitions: []
stream: S-0358
tags: [cli]
touches: [docs/users/flai.md]
after: [T-1412]
---
# T-1414 The user guide says what gateway, harness, and estimated mean beside a cost

## Work

- `docs/users/flai.md`, where usage and `flai stats` are described: the three sources, which runs get each, why a gateway's figure can differ from Claude Code's, and that the sum shows the least certain.

It waits for T-1412, so that it quotes the output as built. It runs together with T-1413, which touches other files.

## Done when

- `flai test docs/users/flai.md` passes.

## Notes

Layer 3 of S-0358.
