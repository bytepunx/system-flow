---
id: T-0570
type: task
nature: feature
title: Measure flai version's start on this host and record it in the design
status: done
parent: S-0160
owner: alex
created: 2026-09-29T19:34:26Z
updated: 2026-09-29T19:40:16Z
transitions:
  - to: ready
    at: 2026-09-29T19:34:42Z
    by: agent-S-0160
  - to: in-progress
    at: 2026-09-29T19:39:18Z
    by: agent-S-0160
  - to: done
    at: 2026-09-29T19:40:16Z
    by: agent-S-0160
stream: S-0160
tags: []
touches: [design/system/server-performance.md]
usage:
  source: log
  seconds: 58
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 18
      output: 5042
      cache_read: 896288
      cache_write: 17086
      cost: 0.4169
---

# T-0570 Measure flai version's start on this host and record it in the design

## Work

Build flai and, with this host's own `PATH`, run `GODEBUG=inittrace=1 flai version` and time `flai version` several times. Check a prompt in a terminal (`flai new` or `flai move ... cancelled` on a scratch project under a pseudo-terminal). Record the numbers under cause 5 of `design/system/server-performance.md`.

## Done when

- No package init takes over 5 ms and `flai version` takes under 30 ms with the host's `PATH`, both recorded in the design with the date.
- A prompt answered at a terminal works.

## Notes
