---
id: T-1414
type: task
nature: feature
title: The user guide says what gateway, harness, and estimated mean beside a cost
status: done
parent: S-0358
owner: alex
created: 2026-10-08T08:51:52Z
updated: 2026-10-13T21:45:52Z
transitions:
  - to: ready
    at: 2026-10-13T21:45:04Z
    by: system-flow
  - to: in-progress
    at: 2026-10-13T21:45:05Z
    by: system-flow
  - to: done
    at: 2026-10-13T21:45:52Z
    by: system-flow
stream: S-0358
tags: [cli]
touches: [docs/users/flai.md]
after: [T-1412]
usage:
  source: log
  seconds: 47
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 12
      output: 3909
      cache_read: 634263
      cache_write: 22599
      cost: 0.3543
---
# T-1414 The user guide says what gateway, harness, and estimated mean beside a cost

## Work

- `docs/users/flai.md`, where usage and `flai stats` are described: the three sources, which runs get each, why a gateway's figure can differ from Claude Code's, and that the sum shows the least certain.

It waits for T-1412, so that it quotes the output as built. It runs together with T-1413, which touches other files.

## Done when

- `flai test docs/users/flai.md` passes.

## Notes

Layer 3 of S-0358.
