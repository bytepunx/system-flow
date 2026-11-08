---
id: MS-0391
title: "S-0604 holds S-0602 on overlap alone: may S-0602 start beside it?"
from: S-0602
to: S-0604
about: [design/system/server-performance.md]
status: open
participants: [flai]
created: 2026-11-08T02:57:43Z
updated: 2026-11-08T02:57:43Z
---

# MS-0391 S-0604 holds S-0602 on overlap alone: may S-0602 start beside it?

Between S-0602 and S-0604, about `design/system/server-performance.md`.

## Entries

### 2026-11-08T02:57:43Z flai S-0602
S-0604 holds S-0602 on overlap alone: may S-0602 start beside it?

S-0602, "flai serve keeps threads, conversations, digest entries, issues, and activity documents parsed between requests, so a thread read takes about 5 ms, not 80 to 110 ms", is ready, and its claim overlaps S-0604's on `design/system/server-performance.md`. It has no agent until it starts, so flai asks for it (ADR-0134).

Its goal:

> Cause 13 in [server-performance.md](../../../design/system/server-performance.md#where-the-time-goes): threads (884 files), conversations (351), digest entries (521), issues (295), and activity documents are read and parsed again on every call. A thread read costs 80 to 110 ms, conversations 35 to 50 ms, the digest 60 to 85 ms, issues 22 to 30 ms, and YAML parsing is about a fifth of serve's CPU. `agent.state`, `inbox.designer`, `threads.list`, `messages.list`, `board.get`'s `merge.queue`, and `stats.get` all pay it; `threads.read` averages about 600 ms in the serve log.

Answer in one of three ways:

1. Narrow your touches with `flai touches S-0604 --remove <path>`, or on the task that names the path, if S-0604 will not change it: the hold clears when the claims no longer overlap.
2. Share the paths with `flai message share MS-0391 --paths <path> "<split>"`, or the MCP tool `message_share`, saying who changes what, if the two stories can change them apart: the overlap on them no longer holds, and S-0602 can start.
3. Reply with `message_reply` saying why the hold stands: S-0602 waits until S-0604 moves to review, is cancelled, or is sent back.
