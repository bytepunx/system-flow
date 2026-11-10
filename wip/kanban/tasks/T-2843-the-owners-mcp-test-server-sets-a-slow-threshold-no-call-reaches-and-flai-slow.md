---
id: T-2843
type: task
nature: remediation
title: "The owners' MCP test server sets a slow threshold no call reaches, and FLAI_SLOW_REQUEST=1ns reproduces the slow log"
status: done
parent: S-0652
owner: alex
created: 2026-11-10T07:01:03Z
updated: 2026-11-10T08:10:01Z
transitions:
  - to: ready
    at: 2026-11-10T08:05:41Z
    by: system-flow
  - to: in-progress
    at: 2026-11-10T08:05:41Z
    by: system-flow
  - to: done
    at: 2026-11-10T08:10:01Z
    by: agent-S-0652
stream: S-0652
tags: [testing]
touches: [flai/internal/mcpserver/server_test.go]
usage:
  source: log
  seconds: 260
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 10
      output: 112
      cache_read: 1260457
      cache_write: 9717
      cost: 0.4674
---
# T-2843 The owners' MCP test server sets a slow threshold no call reaches, and FLAI_SLOW_REQUEST=1ns reproduces the slow log

## Work

`TestAnOperatorsReplyBecomesADigestEntry` asserts the server logged nothing. `ownersServer` leaves `Options.Slow` zero, so the timing middleware takes `perf.Slow()`, 500 ms by default, and a `thread_reply` slower than that under load logs a slow "request answered" event.

- Reproduce first: run the test with `FLAI_SLOW_REQUEST=1ns` and see it fail on the logged event.
- In `ownersServer`, set `o.Slow` far above any call's time, with a comment saying why. The test then still fails on a digest write that logs a warning, which is what it asserts.
- Run the test again with `FLAI_SLOW_REQUEST=1ns`, and the three tests that call `ownersServer`, and see them pass.

It waits for no task: it touches only the MCP server test file, which no other task of this story touches.

## Done when

- The test fails with `FLAI_SLOW_REQUEST=1ns` before the change and passes with it after.
- `flai test flai/internal/mcpserver/server_test.go` passes.

## Notes
