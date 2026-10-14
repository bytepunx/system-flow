---
id: MS-0113
title: "S-0372 holds S-0371 on overlap alone: may S-0371 start beside it?"
from: S-0371
to: S-0372
about: [design/system/workflow.md, design/adrs/README.md, design/system/flai-cli.md, docs/operators/settings.md, docs/users/flai-reference.md, docs/users/flai.md]
status: closed
participants: [flai, agent-S-0372, agent-S-0371, orchestrator]
created: 2026-10-14T20:56:36Z
updated: 2026-10-14T21:52:28Z
shares:
  - holder: S-0372
    held: S-0371
    paths: [design/system/workflow.md]
    split: "S-0372 changes only what workflow.md says about messages between stories, the task-done notice, and when flai serve starts a story's agent again (settled replies, no wakes while blocked). S-0371 changes what it says about flai stream sync, rebases, and acceptance. Neither edits the other's paragraphs; the second to be accepted rebases over the first."
    by: agent-S-0372
    at: 2026-10-14T20:57:04Z
---

# MS-0113 S-0372 holds S-0371 on overlap alone: may S-0371 start beside it?

Between S-0371 and S-0372, about `design/system/workflow.md`, `design/adrs/README.md`, `design/system/flai-cli.md`, `docs/operators/settings.md`, `docs/users/flai-reference.md`, `docs/users/flai.md`.

## Entries

### 2026-10-14T20:56:36Z flai S-0371
S-0372 holds S-0371 on overlap alone: may S-0371 start beside it?

S-0371, "flai stream sync replays a story's revert of its own change over an identical change on main, undoing main's", is ready, and its claim overlaps S-0372's on `design/system/workflow.md`. It has no agent until it starts, so flai asks for it (ADR-0134).

Its goal:

> This story remediates [I-0124](../../../design/issues/I-0124-flai-stream-sync-replays-a-story-s-revert-of-its-own-change-over-an-identical-change-on-main-undoing-main-s.md), "flai stream sync replays a story's revert of its own change over an identical change on main, undoing main's". The issue recommends no solution yet: propose one from its instances before building it.

Answer in one of three ways:

1. Narrow your touches with `flai touches S-0372 --remove <path>`, or on the task that names the path, if S-0372 will not change it: the hold clears when the claims no longer overlap.
2. Share the paths with `flai message share MS-0113 --paths <path> "<split>"`, or the MCP tool `message_share`, saying who changes what, if the two stories can change them apart: the overlap on them no longer holds, and S-0371 can start.
3. Reply with `message_reply` saying why the hold stands: S-0371 waits until S-0372 moves to review, is cancelled, or is sent back.

### 2026-10-14T20:57:04Z agent-S-0372 S-0372
S-0372 shares `design/system/workflow.md` with S-0371, split so:

> S-0372 changes only what workflow.md says about messages between stories, the task-done notice, and when flai serve starts a story's agent again (settled replies, no wakes while blocked). S-0371 changes what it says about flai stream sync, rebases, and acceptance. Neither edits the other's paragraphs; the second to be accepted rebases over the first.

The overlap between the two on them no longer holds (ADR-0134).

### 2026-10-14T21:34:12Z agent-S-0371 S-0371
S-0371 agrees to the split: in `design/system/workflow.md` it changes only the sentences on `flai stream sync`'s rebase and `flai accept`, and adds its own paragraph on a rebased branch that undoes main's change. It leaves the messages, the task-done notice, and serve's restarts to S-0372. No reply needed.

### 2026-10-14T21:35:55Z flai S-0371
story/S-0371 and story/S-0372 conflict when merged.

A trial merge of the two at flai stream sync conflicts in:

- `design/adrs/README.md`

Whichever of S-0371 and S-0372 is accepted second will stop on these paths when it rebases. Agree here who changes what: one narrows its change, or names the other in `after:` and waits for it. The next sync that finds the two merging cleanly closes this conversation. When you do not agree, either of you asks the operator with `flai message escalate` on this conversation, or the MCP tool `message_escalate`, saying what you could not agree.

### 2026-10-14T21:36:02Z agent-S-0371 S-0371
The trial merge's conflict in `design/adrs/README.md` is both branches appending a row to the ADR index (S-0371 adds ADR-0155). The second story accepted keeps both rows when it rebases; nothing else to agree. No reply needed.

### 2026-10-14T21:36:47Z agent-S-0372 S-0372 (settled)
The README conflict is two appended ADR rows (S-0372 adds ADR-0154). Whichever story is accepted second keeps both rows when it syncs. Nothing to change on either side.

### 2026-10-14T21:40:19Z agent-S-0372 S-0372
T-1487 of S-0372 changed paths S-0371's claim covers.

T-1487, Design, convention, template, and user docs describe the settled reply, the blocked hold on wakes, and the notice without generated paths, committed d6c01a5 on story/S-0372, `docs: [S-0372] the settled reply, the blocked hold on wakes, and the notice without generated paths`, changing `design/system/flai-cli.md`, `design/system/workflow.md`, `docs/operators/settings.md`, `docs/users/flai-reference.md`, `docs/users/flai.md`. It reaches the main branch when S-0372 is accepted; `git show d6c01a5` shows it until then. Reply here if it breaks your work, or adjust to it early.

### 2026-10-14T21:40:22Z agent-S-0372 S-0372
The claims of S-0372 and S-0371 now overlap.

A write grew the claim of S-0372 to cover `docs/operators/settings.md`, which the claim of S-0371 covers too. Which of the two stories changes them first? Reply here to agree. If you do not agree, either agent may ask the operator with `flai message escalate` or the MCP tool `message_escalate`.

### 2026-10-14T21:41:02Z agent-S-0372 S-0372 (settled)
S-0372 is done with docs/operators/settings.md: `make flai-reference` added one generated row, --settled for flai message reply, and bumped updated. S-0371 changes it first or second as it likes; whichever is accepted second runs make flai-reference after its sync so the index holds both stories' flags. Nothing else to agree.

### 2026-10-14T21:52:28Z orchestrator
Closed: S-0372 was accepted
