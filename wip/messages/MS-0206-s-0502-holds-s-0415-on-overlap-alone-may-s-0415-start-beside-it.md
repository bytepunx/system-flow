---
id: MS-0206
title: "S-0502 holds S-0415 on overlap alone: may S-0415 start beside it?"
from: S-0415
to: S-0502
about: [flai/internal/mcpserver/server.go, flai/internal/mcpserver/folder.go]
status: open
participants: [flai]
created: 2026-10-28T14:47:44Z
updated: 2026-10-28T14:47:44Z
---

# MS-0206 S-0502 holds S-0415 on overlap alone: may S-0415 start beside it?

Between S-0415 and S-0502, about `flai/internal/mcpserver/server.go`, `flai/internal/mcpserver/folder.go`.

## Entries

### 2026-10-28T14:47:44Z flai S-0415
S-0502 holds S-0415 on overlap alone: may S-0415 start beside it?

S-0415, "flai prime --role gives a role a project defines its pack within the budget, from the topics and the item kind its role file names, and the three built-in roles prime from their files unchanged", is ready, and its claim overlaps S-0502's on `flai/internal/mcpserver/server.go`, `flai/internal/mcpserver/folder.go`. It has no agent until it starts, so flai asks for it (ADR-0134).

Its goal:

> A role a project defines primes as the planner, the orchestrator, and the analyzer do. `flai prime --role <role>` and the MCP tool `prime` read the role's file for the topics its pack selects and the item it takes, and give it the conventions whose `roles` are empty or list it, briefs of the design its topics select, what its item names, and a catalog, within the budget (ADR-0049) and in parts (ADR-0104). The built-in roles' topics and item kinds move from Go into their role files, and their packs stay exactly as they are.

Answer in one of three ways:

1. Narrow your touches with `flai touches S-0502 --remove <path>`, or on the task that names the path, if S-0502 will not change it: the hold clears when the claims no longer overlap.
2. Share the paths with `flai message share MS-0206 --paths <path> "<split>"`, or the MCP tool `message_share`, saying who changes what, if the two stories can change them apart: the overlap on them no longer holds, and S-0415 can start.
3. Reply with `message_reply` saying why the hold stands: S-0415 waits until S-0502 moves to review, is cancelled, or is sent back.
