---
id: MS-0230
title: "S-0536 holds S-0515 on overlap alone: may S-0515 start beside it?"
from: S-0515
to: S-0536
about: [docs/operators/index.md]
status: open
participants: [flai]
created: 2026-10-29T07:28:28Z
updated: 2026-10-29T07:28:28Z
---

# MS-0230 S-0536 holds S-0515 on overlap alone: may S-0515 start beside it?

Between S-0515 and S-0536, about `docs/operators/index.md`.

## Entries

### 2026-10-29T07:28:28Z flai S-0515
S-0536 holds S-0515 on overlap alone: may S-0515 start beside it?

S-0515, "flai's HTTP MCP server refuses an agent in a container with 403 invalid Host header, since the MCP SDK's localhost protection admits only a loopback Host on a loopback listener", is ready, and its claim overlaps S-0536's on `docs/operators/index.md`. It has no agent until it starts, so flai asks for it (ADR-0134).

Its goal:

> This story remediates [I-0213](../../../design/issues/I-0213-flai-s-http-mcp-server-refuses-an-agent-in-a-container-with-403-invalid-host-header-since-the-mcp-sdk-s-localhost-protection-admits-only-a-loopback-host-on-a-loopback-listener.md), "flai's HTTP MCP server refuses an agent in a container with 403 invalid Host header, since the MCP SDK's localhost protection admits only a loopback Host on a loopback listener". The issue recommends no solution yet: propose one from its instances before building it.

Answer in one of three ways:

1. Narrow your touches with `flai touches S-0536 --remove <path>`, or on the task that names the path, if S-0536 will not change it: the hold clears when the claims no longer overlap.
2. Share the paths with `flai message share MS-0230 --paths <path> "<split>"`, or the MCP tool `message_share`, saying who changes what, if the two stories can change them apart: the overlap on them no longer holds, and S-0515 can start.
3. Reply with `message_reply` saying why the hold stands: S-0515 waits until S-0536 moves to review, is cancelled, or is sent back.
