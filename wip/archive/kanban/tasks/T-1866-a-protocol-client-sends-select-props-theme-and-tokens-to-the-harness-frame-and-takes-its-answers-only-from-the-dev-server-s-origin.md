---
id: T-1866
type: task
nature: feature
title: A protocol client sends select, props, theme, and tokens to the harness frame and takes its answers only from the dev server's origin
status: done
parent: S-0439
owner: alex
created: 2026-10-16T21:34:47Z
updated: 2026-10-17T01:18:40Z
transitions:
  - to: ready
    at: 2026-10-17T01:11:15Z
    by: agent-S-0439
  - to: in-progress
    at: 2026-10-17T01:11:16Z
    by: agent-S-0439
  - to: done
    at: 2026-10-17T01:18:40Z
    by: agent-S-0439
stream: S-0439
tags: [dashboard]
touches: [flaiover/src/lib/explorer/protocol.ts, flaiover/src/lib/explorer/protocol.test.ts]
usage:
  source: log
  seconds: 444
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 21
      output: 7998
      cache_read: 1127949
      cache_write: 30174
      cost: 0.5616
---
# T-1866 A protocol client sends select, props, theme, and tokens to the harness frame and takes its answers only from the dev server's origin

## Work

`protocol.ts`: typed messages as S-0438's ADR fixes them, and `connect(frame, origin)`, which posts to the frame's window with the dev server's origin as the target and listens for `ready`, `rendered`, `error`, `event`, and `size` from that origin only. Messages sent before `ready` are queued, and only the last of each kind is kept.

Waits for nothing in this story; the protocol is S-0438's ADR.

## Done when

Tests cover the queue before `ready`, the last-of-kind rule, an answer from another origin ignored, and each answer's handler. `flai test flaiover/src/lib/explorer` passes.

## Notes

None.
