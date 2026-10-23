---
id: MS-0187
title: "S-0430 holds S-0491 on overlap alone: may S-0491 start beside it?"
from: S-0491
to: S-0430
about: [flai/internal/mcpserver/folder.go, flai/internal/mcpserver/server_test.go, design/conventions/strategic-agents.md, template/root/design/conventions/strategic-agents.md]
status: closed
participants: [flai, agent-S-0430, orchestrator]
created: 2026-10-23T20:37:05Z
updated: 2026-10-23T21:00:06Z
---

# MS-0187 S-0430 holds S-0491 on overlap alone: may S-0491 start beside it?

Between S-0491 and S-0430, about `flai/internal/mcpserver/folder.go`, `flai/internal/mcpserver/server_test.go`, `design/conventions/strategic-agents.md`, `template/root/design/conventions/strategic-agents.md`.

## Entries

### 2026-10-23T20:37:05Z flai S-0491
S-0430 holds S-0491 on overlap alone: may S-0491 start beside it?

S-0491, "Add sub-project management", is ready, and its claim overlaps S-0430's on `flai/internal/mcpserver/folder.go`, `flai/internal/mcpserver/server_test.go`, `design/conventions/strategic-agents.md`, `template/root/design/conventions/strategic-agents.md`. It has no agent until it starts, so flai asks for it (ADR-0134).

Its goal:

> No flai command registers a sub-project, and no document in the project or the template gives the schema of projects. The operator is expected to add the entry to system-flow.yaml by hand. flai's own messages show that an entry has name, path, and kind, but there's little else to go off.

Answer in one of three ways:

1. Narrow your touches with `flai touches S-0430 --remove <path>`, or on the task that names the path, if S-0430 will not change it: the hold clears when the claims no longer overlap.
2. Share the paths with `flai message share MS-0187 --paths <path> "<split>"`, or the MCP tool `message_share`, saying who changes what, if the two stories can change them apart: the overlap on them no longer holds, and S-0491 can start.
3. Reply with `message_reply` saying why the hold stands: S-0491 waits until S-0430 moves to review, is cancelled, or is sent back.

### 2026-10-23T20:47:05Z agent-S-0430 S-0430 (settled)
S-0430 is done and moving to review now, which clears the hold, so S-0491 can start then. On the overlapping paths, S-0430 added the digest step to the `prime` tool's description and the `thread_reply` source in `folder.go`. It added `TestAThreadReplyMayCiteADigestEntry` to `server_test.go`. And it added one sentence on citing digest entries to the `answer_threads` bullet of `strategic-agents.md`, template and here.

### 2026-10-23T21:00:06Z orchestrator
Closed: S-0430 was accepted
