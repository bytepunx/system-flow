---
id: T-1656
type: task
nature: feature
title: flai board, the MCP board, and the host read board.get carry the queue
status: backlog
parent: S-0398
owner: alex
created: 2026-10-16T05:01:53Z
updated: 2026-10-16T05:01:53Z
transitions: []
stream: S-0398
tags: [cli]
touches: [flai/cmd/board.go, flai/cmd/board_test.go, flai/internal/mcpserver/server.go, flai/internal/mcpserver/server_test.go, flai/internal/hostapi/hostapi.go, flai/internal/hostapi/reads_test.go]
after: [T-1639]
---
# T-1656 flai board, the MCP board, and the host read board.get carry the queue

## Work

- Add `queue` to the board's JSON, from `mergequeue`: the eligible stories in order, the others' reasons, the batch, and its step.
- Fill it where the board is built for `flai board --json`, the MCP `board` (`server.go`), and `board.get` (`hostapi.go`), so all three give the same.
- `flai board` as text adds one line under the review column naming the queue's order and the batch.
- Test each of the three on one fixture.
- Waits for T-1639, whose package it reads.

## Done when

The three surfaces carry the same `queue`, and the tests pass.

## Notes
