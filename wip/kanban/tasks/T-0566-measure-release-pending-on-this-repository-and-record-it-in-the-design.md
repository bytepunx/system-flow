---
id: T-0566
type: task
nature: feature
title: Measure release.pending on this repository and record it in the design
status: done
parent: S-0157
owner: alex
created: 2026-09-29T19:23:06Z
updated: 2026-09-29T19:37:29Z
transitions:
  - to: ready
    at: 2026-09-29T19:23:31Z
    by: agent-S-0157
  - to: in-progress
    at: 2026-09-29T19:31:51Z
    by: agent-S-0157
  - to: done
    at: 2026-09-29T19:37:29Z
    by: agent-S-0157
stream: S-0157
tags: []
touches: [design/system/server-performance.md, design/system/flai-cli.md]
usage:
  source: log
  seconds: 338
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 89
      output: 46296
      cache_read: 7204037
      cache_write: 110014
      cost: 3.2472
---

# T-0566 Measure release.pending on this repository and record it in the design

## Work

- Time `board.get` with `flai hostapi --timing` on this repository, cold and warm, and MCP `board`.
- Record the result under cause 2 in `design/system/server-performance.md`, and describe the kept history where the design describes `release.pending`.

## Done when

- `release.pending` is under 15 ms with at most 3 `exec.git` steps on a warm `board.get`, measured and written down.
- `flai check --strict` passes.

## Notes
