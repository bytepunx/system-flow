---
id: MS-0207
title: "S-0502 holds S-0416 on overlap alone: may S-0416 start beside it?"
from: S-0416
to: S-0502
about: [design/system/strategic-agents.md, flai/internal/harness/harness_test.go]
status: open
participants: [flai]
created: 2026-10-28T14:47:44Z
updated: 2026-10-28T14:47:44Z
---

# MS-0207 S-0502 holds S-0416 on overlap alone: may S-0416 start beside it?

Between S-0416 and S-0502, about `design/system/strategic-agents.md`, `flai/internal/harness/harness_test.go`.

## Entries

### 2026-10-28T14:47:44Z flai S-0416
S-0502 holds S-0416 on overlap alone: may S-0416 start beside it?

S-0416, "flai guard reads a role's permitted calls from its role file and refuses the rest, and the planner's, the orchestrator's, and the analyzer's rules move into their files with their guard tests unchanged", is ready, and its claim overlaps S-0502's on `design/system/strategic-agents.md`, `flai/internal/harness/harness_test.go`. It has no agent until it starts, so flai asks for it (ADR-0134).

Its goal:

> The guard holds a role a project defines as it holds the three built-in roles today. With `FLAI_ROLE` naming a defined role, `flai guard` reads the role's file for the MCP tools, flai commands, and file edits it may make, each under the conditions the file names, and for the permissions that gate them, read from the manifest at each call; it refuses the rest, naming the role and the rule. The planner's, the orchestrator's, and the analyzer's rules move from Go into their role files, so one path decides for every role, and their guard tests pass unchanged.

Answer in one of three ways:

1. Narrow your touches with `flai touches S-0502 --remove <path>`, or on the task that names the path, if S-0502 will not change it: the hold clears when the claims no longer overlap.
2. Share the paths with `flai message share MS-0207 --paths <path> "<split>"`, or the MCP tool `message_share`, saying who changes what, if the two stories can change them apart: the overlap on them no longer holds, and S-0416 can start.
3. Reply with `message_reply` saying why the hold stands: S-0416 waits until S-0502 moves to review, is cancelled, or is sent back.
