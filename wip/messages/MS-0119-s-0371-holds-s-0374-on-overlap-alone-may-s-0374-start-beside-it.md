---
id: MS-0119
title: "S-0371 holds S-0374 on overlap alone: may S-0374 start beside it?"
from: S-0374
to: S-0371
about: [design/system/workflow.md]
status: open
participants: [flai, agent-S-0371]
created: 2026-10-14T21:33:36Z
updated: 2026-10-14T21:34:11Z
shares:
  - holder: S-0371
    held: S-0374
    paths: [design/system/workflow.md]
    split: "S-0374 changes only what workflow.md says about the close-out and flai verify (in § Branches and collisions, the sentences from \"Before moving the story to review\" to the close-out's ADR-0110 reference). S-0371 changes what it says about flai stream sync's rebase and flai accept (the sentences on `flai stream sync` and `flai accept` in that paragraph), and adds its own paragraph on a rebased branch that undoes main's change. Neither edits the other's sentences; the second to be accepted rebases over the first."
    by: agent-S-0371
    at: 2026-10-14T21:34:11Z
---

# MS-0119 S-0371 holds S-0374 on overlap alone: may S-0374 start beside it?

Between S-0374 and S-0371, about `design/system/workflow.md`.

## Entries

### 2026-10-14T21:33:36Z flai S-0374
S-0371 holds S-0374 on overlap alone: may S-0374 start beside it?

S-0374, "An acceptance on main during a close-out's tiers fails its final sync check, so the whole close-out runs again", is ready, and its claim overlaps S-0371's on `design/system/workflow.md`. It has no agent until it starts, so flai asks for it (ADR-0134).

Its goal:

> This story remediates [I-0136](../../../design/issues/I-0136-an-acceptance-on-main-during-a-close-out-s-tiers-fails-its-final-sync-check-so-the-whole-close-out-runs-again.md), "An acceptance on main during a close-out's tiers fails its final sync check, so the whole close-out runs again". The issue recommends no solution yet: propose one from its instances before building it.

Answer in one of three ways:

1. Narrow your touches with `flai touches S-0371 --remove <path>`, or on the task that names the path, if S-0371 will not change it: the hold clears when the claims no longer overlap.
2. Share the paths with `flai message share MS-0119 --paths <path> "<split>"`, or the MCP tool `message_share`, saying who changes what, if the two stories can change them apart: the overlap on them no longer holds, and S-0374 can start.
3. Reply with `message_reply` saying why the hold stands: S-0374 waits until S-0371 moves to review, is cancelled, or is sent back.

### 2026-10-14T21:34:11Z agent-S-0371 S-0371
S-0371 shares `design/system/workflow.md` with S-0374, split so:

> S-0374 changes only what workflow.md says about the close-out and flai verify (in § Branches and collisions, the sentences from "Before moving the story to review" to the close-out's ADR-0110 reference). S-0371 changes what it says about flai stream sync's rebase and flai accept (the sentences on `flai stream sync` and `flai accept` in that paragraph), and adds its own paragraph on a rebased branch that undoes main's change. Neither edits the other's sentences; the second to be accepted rebases over the first.

The overlap between the two on them no longer holds (ADR-0134).
