---
id: T-3044
type: task
nature: remediation
title: The test dashboard reads its connection for its whole life, so it answers the client's pings
status: backlog
parent: S-0703
owner: alex
created: 2026-11-10T10:34:46Z
updated: 2026-11-10T10:34:46Z
transitions: []
stream: S-0703
tags: [flai, tests]
touches: [flai/internal/channel/channeltest/dashboard.go]
after: [T-3042]
---
# T-3044 The test dashboard reads its connection for its whole life, so it answers the client's pings

## Work

In `flai/internal/channel/channeltest/dashboard.go`, after the proof, start one reader per connection that reads until the connection closes, so coder/websocket answers each ping, and that hands what it reads to `Ask` and `Next` through a channel in place of their own `c.ws.Read`. Keep `Ask`'s and `Next`'s 5-second timeouts and the answers they return. A read error ends the reader and reaches the next `Ask` or `Next` as the error it reads today. Change no production code in `flai/internal/channel`.

Waits on T-3042: its test is what this fix must turn green, and it names the cause this fix assumes.

## Done when

- T-3042's test passes, and so do the tests of `flai/internal/serve`, `flai/internal/hostapi`, and `flai/internal/channel` with `-race`.

## Notes
