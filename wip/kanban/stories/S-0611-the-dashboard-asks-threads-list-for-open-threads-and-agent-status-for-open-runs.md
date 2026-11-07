---
id: S-0611
type: story
nature: improvement
title: The dashboard asks threads.list for open threads and agent.status for open runs, and the request timing counts the send, so a thread list is under 100 KB, not 3.3 MB
status: backlog
parent: E-0033
owner: alex
created: 2026-11-07T22:42:08Z
updated: 2026-11-07T22:42:08Z
transitions: []
tags: [dashboard, cli]
touches: [flaiover/src/lib/server/repo.ts, flai/internal/channel/channel.go, flai/internal/hostapi/people.go]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
draft: true
---
# S-0611 The dashboard asks threads.list for open threads and agent.status for open runs, and the request timing counts the send, so a thread list is under 100 KB, not 3.3 MB

## Goal

Cause 21 in [server-performance.md](../../../design/system/server-performance.md#where-the-time-goes): the dashboard asks `threads.list` with `all: true` (`flaiover/src/lib/server/repo.ts`), 3.28 MB for 884 threads of which 12 are open, and `agent.status` sends 437 finished runs and 288 planner runs, 392 KB. Each holds its connection's one write lock while it is written (`flai/internal/channel/channel.go`), and the time is not in `duration_ms`, which is logged before the send.

## Acceptance criteria

- [ ] The dashboard asks `threads.list` for open threads, and for the resolved ones only on a page that shows them.
- [ ] `agent.status` sends the open and live runs, and the finished ones only to a page that shows them, through a parameter the host API documents.
- [ ] The `request answered` event gives the send's time in its own field, and `docs/users/flai.md` § Request timing says so.
- [ ] On this repository a board's `threads.list` is under 100 KB and `agent.status` under 50 KB, and server-performance.md records the before and after.

## Tasks

## Notes

- These change what the host API answers by default, not how flai keeps state.
