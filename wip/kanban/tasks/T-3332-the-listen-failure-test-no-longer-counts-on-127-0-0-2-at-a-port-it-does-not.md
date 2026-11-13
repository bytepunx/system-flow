---
id: T-3332
type: task
nature: remediation
title: The listen-failure test no longer counts on 127.0.0.2 at a port it does not hold staying free
status: done
parent: S-0779
owner: alex
created: 2026-11-13T20:33:50Z
updated: 2026-11-13T20:44:05Z
transitions:
  - to: ready
    at: 2026-11-13T20:39:52Z
    by: agent-S-0779
  - to: in-progress
    at: 2026-11-13T20:39:53Z
    by: agent-S-0779
  - to: done
    at: 2026-11-13T20:44:05Z
    by: agent-S-0779
stream: S-0779
tags: [flai]
touches: [flai/cmd/mcp_http_test.go]
usage:
  source: log
  seconds: 252
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 14
      output: 148
      cache_read: 1518636
      cache_write: 7866
      cost: 0.573
---
# T-3332 The listen-failure test no longer counts on 127.0.0.2 at a port it does not hold staying free

## Work

`TestMCPListenFailureClosesWhatItOpened` (`flai/cmd/mcp_http_test.go`) holds `127.0.0.1:<port>` and counts on `127.0.0.2:<port>` being free before `listenMCP` runs and again after it closes what it opened. Nothing ensures either.

Make the test decide for itself. One way: when `listenMCP` fails on the first address rather than the held one, or the probe after it finds the first address taken, try again with a new port, a bounded number of times, and fail only when every try loses. Keep `listenMCP` listening on real addresses here: this test is one of those that keep what `listenMCP` does with the operating system covered. Should the agent give `listenMCP` its listen as a value instead, the task's touches widen to `flai/cmd/mcp_http.go`, and the narrative's `## Decisions` says why.

Waits for nothing: it shares no path with the other tasks.

## Done when

- The test no longer fails when another process takes `127.0.0.2:<port>` before or after `listenMCP` runs.
- It still shows that a failure on any address closes the listeners already opened and names the address that failed.
- `flai test` on `flai/cmd/mcp_http_test.go` passes.

## Notes
