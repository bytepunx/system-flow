---
id: MS-0141
title: "S-0405 holds S-0400 on overlap alone: may S-0400 start beside it?"
from: S-0400
to: S-0405
about: [design/system/workflow.md]
status: closed
participants: [flai, agent-S-0405, orchestrator]
created: 2026-10-19T20:34:23Z
updated: 2026-10-19T21:52:13Z
---

# MS-0141 S-0405 holds S-0400 on overlap alone: may S-0400 start beside it?

Between S-0400 and S-0405, about `design/system/workflow.md`.

## Entries

### 2026-10-19T20:34:23Z flai S-0400
S-0405 holds S-0400 on overlap alone: may S-0400 start beside it?

S-0400, "A green batch lands each story in order through the acceptance of today, and a red batch is bisected to the story that fails, which goes back to in-progress with the failing step in its notes and a thread to its agent", is ready, and its claim overlaps S-0405's on `design/system/workflow.md`. It has no agent until it starts, so flai asks for it (ADR-0134).

Its goal:

> A verified batch ends in acceptances or in a send-back, with no person in between. When the batch passes, each of its stories is accepted in order through the acceptance of today, so that its transitions, acceptance commit, archive, and evidence read as they do now. When it fails, the queue finds the story that breaks it, lands the stories before it, and sends that story back to its agent with the failing step.

Answer in one of three ways:

1. Narrow your touches with `flai touches S-0405 --remove <path>`, or on the task that names the path, if S-0405 will not change it: the hold clears when the claims no longer overlap.
2. Share the paths with `flai message share MS-0141 --paths <path> "<split>"`, or the MCP tool `message_share`, saying who changes what, if the two stories can change them apart: the overlap on them no longer holds, and S-0400 can start.
3. Reply with `message_reply` saying why the hold stands: S-0400 waits until S-0405 moves to review, is cancelled, or is sent back.

### 2026-10-19T21:50:47Z agent-S-0405 S-0405 (settled)
S-0405 is in review now, so it no longer holds S-0400. On workflow.md, S-0405 added a paragraph at the end of the merge queue part of "What waits on an acceptance": while the queue can batch, the close-out leaves the tiers marked queue to the batch (ADR-0174). It also replaced the "So far it only reads (S-0398)" sentence. S-0400's landing should keep that paragraph true.

### 2026-10-19T21:52:13Z orchestrator
Closed: S-0405 was accepted
