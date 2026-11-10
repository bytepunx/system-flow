---
id: MS-0463
title: "S-0682 holds S-0665 on overlap alone: may S-0665 start beside it?"
from: S-0665
to: S-0682
about: [flai/internal/harness/harness.go]
status: closed
participants: [flai, agent-S-0682, alex]
created: 2026-11-10T09:54:28Z
updated: 2026-11-10T10:10:48Z
---

# MS-0463 S-0682 holds S-0665 on overlap alone: may S-0665 start beside it?

Between S-0665 and S-0682, about `flai/internal/harness/harness.go`.

## Entries

### 2026-11-10T09:54:28Z flai S-0665
S-0682 holds S-0665 on overlap alone: may S-0665 start beside it?

S-0665, "flai serve records no image, container name, or mounts for a container harness run, so a run cannot be shown afterwards to have used the project's layer image", is ready, and its claim overlaps S-0682's on `flai/internal/harness/harness.go`. It has no agent until it starts, so flai asks for it (ADR-0134).

Its goal:

> This story remediates [I-0274](../../../design/issues/I-0274-flai-serve-records-no-image-container-name-or-mounts-for-a-container-harness.md), "flai serve records no image, container name, or mounts for a container harness run, so a run cannot be shown afterwards to have used the project's layer image". The issue recommends no solution yet: propose one from its instances before building it.

Answer in one of three ways:

1. Narrow your touches with `flai touches S-0682 --remove <path>`, or on the task that names the path, if S-0682 will not change it: the hold clears when the claims no longer overlap.
2. Share the paths with `flai message share MS-0463 --paths <path> "<split>"`, or the MCP tool `message_share`, saying who changes what, if the two stories can change them apart: the overlap on them no longer holds, and S-0665 can start.
3. Reply with `message_reply` saying why the hold stands: S-0665 waits until S-0682 moves to review, is cancelled, or is sent back.

### 2026-11-10T10:08:08Z agent-S-0682 S-0682 (settled)
S-0682's work is done and its close-out passed; it moves to review now, which holds nothing, so S-0665 can start. S-0682 adds `Probe bool` (JSON `probe`) to harness.Request and a probe case at the top of Prompt; a change elsewhere in harness.go should merge beside it.

### 2026-11-10T10:10:48Z alex
Closed: S-0682 was accepted
