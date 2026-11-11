---
id: MS-0503
title: "S-0694 holds S-0705 on overlap alone: may S-0705 start beside it?"
from: S-0705
to: S-0694
about: [flai/internal/harness/harness.go, design/system/strategic-agents.md]
status: open
participants: [flai, agent-S-0694]
created: 2026-11-11T14:23:11Z
updated: 2026-11-11T14:24:21Z
shares:
  - holder: S-0694
    held: S-0705
    paths: [flai/internal/harness/harness.go, design/system/strategic-agents.md]
    split: "S-0694's change to harness.go is committed (964508846 on story/S-0694): the Kluster field of harness.SentBack and its branch in SentBack.then, nothing else. In strategic-agents.md, S-0694 changes only the orchestrator's acceptance blockers, to add kluster. S-0705 changes the role-file and container-harness passages."
    by: agent-S-0694
    at: 2026-11-11T14:24:21Z
---

# MS-0503 S-0694 holds S-0705 on overlap alone: may S-0705 start beside it?

Between S-0705 and S-0694, about `flai/internal/harness/harness.go`, `design/system/strategic-agents.md`.

## Entries

### 2026-11-11T14:23:11Z flai S-0705
S-0694 holds S-0705 on overlap alone: may S-0705 start beside it?

S-0705, "A role a project defines in a role file of its own cannot run on the container harness, since the start request does not carry its file to pilot", is ready, and its claim overlaps S-0694's on `flai/internal/harness/harness.go`, `design/system/strategic-agents.md`. It has no agent until it starts, so flai asks for it (ADR-0134).

Its goal:

> This story remediates [I-0258](../../../design/issues/I-0258-a-role-a-project-defines-in-a-role-file-of-its-own-cannot-run-on-the-container.md), "A role a project defines in a role file of its own cannot run on the container harness, since the start request does not carry its file to pilot". The issue recommends no solution yet: propose one from its instances before building it.

Answer in one of three ways:

1. Narrow your touches with `flai touches S-0694 --remove <path>`, or on the task that names the path, if S-0694 will not change it: the hold clears when the claims no longer overlap.
2. Share the paths with `flai message share MS-0503 --paths <path> "<split>"`, or the MCP tool `message_share`, saying who changes what, if the two stories can change them apart: the overlap on them no longer holds, and S-0705 can start.
3. Reply with `message_reply` saying why the hold stands: S-0705 waits until S-0694 moves to review, is cancelled, or is sent back.

### 2026-11-11T14:24:21Z agent-S-0694 S-0694
S-0694 shares `flai/internal/harness/harness.go`, `design/system/strategic-agents.md` with S-0705, split so:

> S-0694's change to harness.go is committed (964508846 on story/S-0694): the Kluster field of harness.SentBack and its branch in SentBack.then, nothing else. In strategic-agents.md, S-0694 changes only the orchestrator's acceptance blockers, to add kluster. S-0705 changes the role-file and container-harness passages.

The overlap between the two on them no longer holds (ADR-0134).
