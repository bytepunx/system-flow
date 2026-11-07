---
id: S-0602
type: story
nature: improvement
title: flai serve keeps threads, conversations, digest entries, issues, and activity documents parsed between requests, so a thread read takes about 5 ms, not 80 to 110 ms
status: backlog
parent: E-0033
owner: alex
created: 2026-11-07T22:40:05Z
updated: 2026-11-07T22:40:05Z
transitions: []
tags: [cli]
touches: [flai/internal/threads/threads.go, flai/internal/messages/messages.go, flai/internal/digest/digest.go, flai/internal/issues/issues.go]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
draft: true
---
# S-0602 flai serve keeps threads, conversations, digest entries, issues, and activity documents parsed between requests, so a thread read takes about 5 ms, not 80 to 110 ms

## Goal

Cause 13 in [server-performance.md](../../../design/system/server-performance.md#where-the-time-goes): threads (884 files), conversations (351), digest entries (521), issues (295), and activity documents are read and parsed again on every call. A thread read costs 80 to 110 ms, conversations 35 to 50 ms, the digest 60 to 85 ms, issues 22 to 30 ms, and YAML parsing is about a fifth of serve's CPU. `agent.state`, `inbox.designer`, `threads.list`, `messages.list`, `board.get`'s `merge.queue`, and `stats.get` all pay it; `threads.read` averages about 600 ms in the serve log.

## Acceptance criteria

- [ ] flai serve and flai mcp keep each of the five parsed between requests, as S-0156 keeps items: each read stats every file and parses again only those whose size or modification time changed, or that are new, and drops those that are gone.
- [ ] Tests show a changed, added, and removed file of each kind is seen on the next read, and that two readers at once get the same answer.
- [ ] A warm `threads.list`, `inbox.designer`, and `stats.get` on this repository, timed with `flai hostapi --timing` in one process or through flai serve, each spend at least 80 ms less in their reads, and server-performance.md records the before and after.

## Tasks

## Notes

- Proposed fix: one kept-files store shared by the five packages, on the model of `flai/internal/workitem/store.go`.
- The stories for causes 12 and 14 gain from this one but do not wait for it.
