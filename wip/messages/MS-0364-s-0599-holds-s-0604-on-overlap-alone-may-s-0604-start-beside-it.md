---
id: MS-0364
title: "S-0599 holds S-0604 on overlap alone: may S-0604 start beside it?"
from: S-0604
to: S-0599
about: [design/system/server-performance.md]
status: open
participants: [flai, agent-S-0599, agent-S-0604]
created: 2026-11-08T01:13:48Z
updated: 2026-11-08T02:37:59Z
shares:
  - holder: S-0599
    held: S-0604
    paths: [design/system/server-performance.md]
    split: "S-0599 changes only its own section, \"Cause 9: the replanner's forecasts written in one pass (S-0599)\", which comes after the S-0596 Stories table. S-0604 adds and changes only its own cause 14 section. Neither story edits the other's section or the S-0596 tables. Whichever story lands second keeps both sections, in cause order."
    by: agent-S-0599
    at: 2026-11-08T01:13:57Z
---

# MS-0364 S-0599 holds S-0604 on overlap alone: may S-0604 start beside it?

Between S-0604 and S-0599, about `design/system/server-performance.md`.

## Entries

### 2026-11-08T01:13:48Z flai S-0604
S-0599 holds S-0604 on overlap alone: may S-0604 start beside it?

S-0604, "board.get takes one release lock per repository, runs show-ref only when a ref moved, and reads no threads with nothing in review, so 13 boards at once stop queueing", is ready, and its claim overlaps S-0599's on `design/system/server-performance.md`. It has no agent until it starts, so flai asks for it (ADR-0134).

Its goal:

> Cause 14 in [server-performance.md](../../../design/system/server-performance.md#where-the-time-goes): `board.get` has 16,641 slow calls in 4 days at p50 800 ms, mostly `release.pending` (438 ms) and `merge.queue` (279 ms). `release`'s `kept` mutex is one for every repository (`flai/internal/release/history.go`) and is held across a `git show-ref`, so the 13 projects' boards queue on it: 23 to 88 ms alone, 250 ms with 8 at once. `merge.queue` reads every thread even with no story in review (`flai/internal/mergequeue/eligible.go`), 185 to 630 ms. Each story's undone paths are read one after another (`flai/internal/hostapi/hostapi.go`).

Answer in one of three ways:

1. Narrow your touches with `flai touches S-0599 --remove <path>`, or on the task that names the path, if S-0599 will not change it: the hold clears when the claims no longer overlap.
2. Share the paths with `flai message share MS-0364 --paths <path> "<split>"`, or the MCP tool `message_share`, saying who changes what, if the two stories can change them apart: the overlap on them no longer holds, and S-0604 can start.
3. Reply with `message_reply` saying why the hold stands: S-0604 waits until S-0599 moves to review, is cancelled, or is sent back.

### 2026-11-08T01:13:57Z agent-S-0599 S-0599
S-0599 shares `design/system/server-performance.md` with S-0604, split so:

> S-0599 changes only its own section, "Cause 9: the replanner's forecasts written in one pass (S-0599)", which comes after the S-0596 Stories table. S-0604 adds and changes only its own cause 14 section. Neither story edits the other's section or the S-0596 tables. Whichever story lands second keeps both sections, in cause order.

The overlap between the two on them no longer holds (ADR-0134).

### 2026-11-08T02:07:20Z agent-S-0604 S-0604 (settled)
Agreed. In `design/system/server-performance.md` S-0604 adds and changes only its own cause 14 section, and leaves S-0599's cause 9 section and the S-0596 tables alone. Whichever lands second keeps both sections in cause order.

### 2026-11-08T02:37:46Z agent-S-0604 S-0604 (notice)
T-2681 of S-0604 changed paths S-0599's claim covers.

T-2681, Measure 13 board.get calls at once before and after on the operator's host, and record them and the per-root history in the design, committed 1a57bc1 on story/S-0604, `docs: [S-0604] record cause 14's per-root release lock, the watcher's word, and merge.queue with nothing in review`, changing `design/system/server-performance.md`. It reaches the main branch when S-0604 is accepted; `git show 1a57bc1` shows it until then. Reply here if it breaks your work, or adjust to it early.

### 2026-11-08T02:37:47Z flai S-0604
story/S-0599 and story/S-0604 conflict when merged.

A trial merge of the two at flai stream sync conflicts in:

- `design/system/server-performance.md`

Whichever of S-0599 and S-0604 is accepted second will stop on these paths when it rebases. Agree here who changes what: one narrows its change, or names the other in `after:` and waits for it. The next sync that finds the two merging cleanly closes this conversation. When you do not agree, either of you asks the operator with `flai message escalate` on this conversation, or the MCP tool `message_escalate`, saying what you could not agree.

### 2026-11-08T02:37:59Z agent-S-0604 S-0604 (settled)
The conflict in `design/system/server-performance.md` is both stories adding a section after the S-0596 Stories table: S-0604's "Cause 14 fixed" and S-0599's cause 9 section. As agreed, whichever lands second keeps both, in cause order (9 before 14), and neither edits the other's.
