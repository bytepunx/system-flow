---
id: MS-0502
title: "S-0694 holds S-0689 on overlap alone: may S-0689 start beside it?"
from: S-0689
to: S-0694
about: [flai/internal/harness/harness.go]
status: open
participants: [flai, agent-S-0694]
created: 2026-11-11T14:23:11Z
updated: 2026-11-11T14:24:20Z
shares:
  - holder: S-0694
    held: S-0689
    paths: [flai/internal/harness/harness.go]
    split: "S-0694's change to harness.go is committed (964508846 on story/S-0694): the Kluster field of harness.SentBack and its branch in SentBack.then. It changes nothing else there. S-0689 changes the rest."
    by: agent-S-0694
    at: 2026-11-11T14:24:20Z
---

# MS-0502 S-0694 holds S-0689 on overlap alone: may S-0689 start beside it?

Between S-0689 and S-0694, about `flai/internal/harness/harness.go`.

## Entries

### 2026-11-11T14:23:11Z flai S-0689
S-0694 holds S-0689 on overlap alone: may S-0689 start beside it?

S-0689, "flai stream sync run in a story's agent container cannot rebase, since the main checkout's .git is read-only, and the git rebase --abort it names leaves HEAD detached", is ready, and its claim overlaps S-0694's on `flai/internal/harness/harness.go`. It has no agent until it starts, so flai asks for it (ADR-0134).

Its goal:

> This story remediates [I-0272](../../../design/issues/I-0272-flai-stream-sync-run-in-a-story-s-agent-container-cannot-rebase-since-the-main.md), "flai stream sync run in a story's agent container cannot rebase, since the main checkout's .git is read-only, and the git rebase --abort it names leaves HEAD detached". The issue recommends no solution yet: propose one from its instances before building it.

Answer in one of three ways:

1. Narrow your touches with `flai touches S-0694 --remove <path>`, or on the task that names the path, if S-0694 will not change it: the hold clears when the claims no longer overlap.
2. Share the paths with `flai message share MS-0502 --paths <path> "<split>"`, or the MCP tool `message_share`, saying who changes what, if the two stories can change them apart: the overlap on them no longer holds, and S-0689 can start.
3. Reply with `message_reply` saying why the hold stands: S-0689 waits until S-0694 moves to review, is cancelled, or is sent back.

### 2026-11-11T14:24:20Z agent-S-0694 S-0694
S-0694 shares `flai/internal/harness/harness.go` with S-0689, split so:

> S-0694's change to harness.go is committed (964508846 on story/S-0694): the Kluster field of harness.SentBack and its branch in SentBack.then. It changes nothing else there. S-0689 changes the rest.

The overlap between the two on them no longer holds (ADR-0134).
