---
id: T-3044
type: task
nature: remediation
title: The test dashboard reads its connection for its whole life, so it answers the client's pings
status: done
parent: S-0703
owner: alex
created: 2026-11-10T10:34:46Z
updated: 2026-11-11T23:12:54Z
transitions:
  - to: ready
    at: 2026-11-11T23:07:37Z
    by: agent-S-0703
  - to: in-progress
    at: 2026-11-11T23:07:37Z
    by: agent-S-0703
  - to: done
    at: 2026-11-11T23:12:54Z
    by: agent-S-0703
stream: S-0703
tags: [flai, tests]
touches: [flai/internal/channel/channeltest/dashboard.go, flai/internal/hostapi/components_test.go, flai/internal/serve/serve_test.go]
usage:
  source: log
  seconds: 317
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 14
      output: 1836
      cache_read: 1507521
      cache_write: 18736
      cost: 0.4882
---
# T-3044 The test dashboard reads its connection for its whole life, so it answers the client's pings

## Work

In `flai/internal/channel/channeltest/dashboard.go`, after the proof, start one reader per connection that reads until the connection closes, so coder/websocket answers each ping, and that hands what it reads to `Ask` and `Next` through a channel in place of their own `c.ws.Read`. Keep `Ask`'s and `Next`'s 5-second timeouts and the answers they return. A read error ends the reader and reaches the next `Ask` or `Next` as the error it reads today. Change no production code in `flai/internal/channel`.

Waits on T-3042: its test is what this fix must turn green, and it names the cause this fix assumes.

## Done when

- T-3042's test passes, and so do the tests of `flai/internal/serve`, `flai/internal/hostapi`, and `flai/internal/channel` with `-race`.

## Notes
