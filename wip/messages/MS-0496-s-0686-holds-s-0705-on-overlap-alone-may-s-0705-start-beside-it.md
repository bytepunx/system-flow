---
id: MS-0496
title: "S-0686 holds S-0705 on overlap alone: may S-0705 start beside it?"
from: S-0705
to: S-0686
about: [flai/internal/harness/container.go]
status: open
participants: [flai]
created: 2026-11-11T07:23:13Z
updated: 2026-11-11T07:23:13Z
---

# MS-0496 S-0686 holds S-0705 on overlap alone: may S-0705 start beside it?

Between S-0705 and S-0686, about `flai/internal/harness/container.go`.

## Entries

### 2026-11-11T07:23:13Z flai S-0705
S-0686 holds S-0705 on overlap alone: may S-0705 start beside it?

S-0705, "A role a project defines in a role file of its own cannot run on the container harness, since the start request does not carry its file to pilot", is ready, and its claim overlaps S-0686's on `flai/internal/harness/container.go`. It has no agent until it starts, so flai asks for it (ADR-0134).

Its goal:

> This story remediates [I-0258](../../../design/issues/I-0258-a-role-a-project-defines-in-a-role-file-of-its-own-cannot-run-on-the-container.md), "A role a project defines in a role file of its own cannot run on the container harness, since the start request does not carry its file to pilot". The issue recommends no solution yet: propose one from its instances before building it.

Answer in one of three ways:

1. Narrow your touches with `flai touches S-0686 --remove <path>`, or on the task that names the path, if S-0686 will not change it: the hold clears when the claims no longer overlap.
2. Share the paths with `flai message share MS-0496 --paths <path> "<split>"`, or the MCP tool `message_share`, saying who changes what, if the two stories can change them apart: the overlap on them no longer holds, and S-0705 can start.
3. Reply with `message_reply` saying why the hold stands: S-0705 waits until S-0686 moves to review, is cancelled, or is sent back.
