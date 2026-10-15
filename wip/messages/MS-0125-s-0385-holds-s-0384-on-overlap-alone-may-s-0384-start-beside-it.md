---
id: MS-0125
title: "S-0385 holds S-0384 on overlap alone: may S-0384 start beside it?"
from: S-0384
to: S-0385
about: [flai/cmd/design.go, flai/cmd/design_test.go]
status: open
participants: [flai]
created: 2026-10-15T02:55:08Z
updated: 2026-10-15T02:55:08Z
---

# MS-0125 S-0385 holds S-0384 on overlap alone: may S-0384 start beside it?

Between S-0384 and S-0385, about `flai/cmd/design.go`, `flai/cmd/design_test.go`.

## Entries

### 2026-10-15T02:55:08Z flai S-0384
S-0385 holds S-0384 on overlap alone: may S-0384 start beside it?

S-0384, "An agent's proposal for the design system is a story whose branch changes it, and the Design page previews that branch beside the main branch's before the operator accepts it", is ready, and its claim overlaps S-0385's on `flai/cmd/design.go`, `flai/cmd/design_test.go`. It has no agent until it starts, so flai asks for it (ADR-0134).

Its goal:

> The operator works with agents on the design system through the board's own workflow: they ask for a change from the Design page, which becomes a story; the story's agent changes the design system on its branch; and the Design page previews that branch's tokens on the sample components beside the main branch's, so the operator sees the proposal respond in the browser before accepting it, and accepting it adopts it.

Answer in one of three ways:

1. Narrow your touches with `flai touches S-0385 --remove <path>`, or on the task that names the path, if S-0385 will not change it: the hold clears when the claims no longer overlap.
2. Share the paths with `flai message share MS-0125 --paths <path> "<split>"`, or the MCP tool `message_share`, saying who changes what, if the two stories can change them apart: the overlap on them no longer holds, and S-0384 can start.
3. Reply with `message_reply` saying why the hold stands: S-0384 waits until S-0385 moves to review, is cancelled, or is sent back.
