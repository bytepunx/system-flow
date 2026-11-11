---
id: MS-0496
title: "S-0686 holds S-0705 on overlap alone: may S-0705 start beside it?"
from: S-0705
to: S-0686
about: [flai/internal/harness/container.go]
status: closed
participants: [flai, agent-S-0686, alex]
created: 2026-11-11T07:23:13Z
updated: 2026-11-11T08:46:40Z
shares:
  - holder: S-0686
    held: S-0705
    paths: [flai/internal/harness/container.go]
    split: "S-0686 is done with container.go: its only change, committed, is two lines of safeMain's doc comment (go build's VCS stamping needs the main checkout's safe.directory too, I-0255). It changes no code there and nothing more is planned. S-0705 may change anything else in the file; if it touches safeMain's comment, keep the I-0255 sentence. The scripts/README.md trial-merge conflict is gone too: story/S-0686 was just rebased onto main."
    by: agent-S-0686
    at: 2026-11-11T07:43:19Z
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

### 2026-11-11T07:43:19Z agent-S-0686 S-0686
S-0686 shares `flai/internal/harness/container.go` with S-0705, split so:

> S-0686 is done with container.go: its only change, committed, is two lines of safeMain's doc comment (go build's VCS stamping needs the main checkout's safe.directory too, I-0255). It changes no code there and nothing more is planned. S-0705 may change anything else in the file; if it touches safeMain's comment, keep the I-0255 sentence. The scripts/README.md trial-merge conflict is gone too: story/S-0686 was just rebased onto main.

The overlap between the two on them no longer holds (ADR-0134).

### 2026-11-11T08:46:40Z alex
Closed: S-0686 was accepted
