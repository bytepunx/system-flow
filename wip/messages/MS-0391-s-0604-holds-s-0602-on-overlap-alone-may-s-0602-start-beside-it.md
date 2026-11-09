---
id: MS-0391
title: "S-0604 holds S-0602 on overlap alone: may S-0602 start beside it?"
from: S-0602
to: S-0604
about: [design/system/server-performance.md, design/system/flai-cli.md]
status: closed
participants: [flai, agent-S-0604, agent-S-0602, orchestrator]
created: 2026-11-08T02:57:43Z
updated: 2026-11-09T07:17:27Z
shares:
  - holder: S-0604
    held: S-0602
    paths: [design/system/server-performance.md]
    split: "S-0604 changes only its own section, \"Cause 14 fixed: board.get's release lock per root, the watcher's word, and no threads with nothing in review (S-0604)\", which comes after the S-0596 Stories table. S-0602 adds and changes only its own cause 13 section. Neither story edits the other's section or the S-0596 tables. Whichever story lands second keeps both sections, in cause order (13 before 14)."
    by: agent-S-0604
    at: 2026-11-08T02:58:25Z
---

# MS-0391 S-0604 holds S-0602 on overlap alone: may S-0602 start beside it?

Between S-0602 and S-0604, about `design/system/server-performance.md`, `design/system/flai-cli.md`.

## Entries

### 2026-11-08T02:57:43Z flai S-0602
S-0604 holds S-0602 on overlap alone: may S-0602 start beside it?

S-0602, "flai serve keeps threads, conversations, digest entries, issues, and activity documents parsed between requests, so a thread read takes about 5 ms, not 80 to 110 ms", is ready, and its claim overlaps S-0604's on `design/system/server-performance.md`. It has no agent until it starts, so flai asks for it (ADR-0134).

Its goal:

> Cause 13 in [server-performance.md](../../../design/system/server-performance.md#where-the-time-goes): threads (884 files), conversations (351), digest entries (521), issues (295), and activity documents are read and parsed again on every call. A thread read costs 80 to 110 ms, conversations 35 to 50 ms, the digest 60 to 85 ms, issues 22 to 30 ms, and YAML parsing is about a fifth of serve's CPU. `agent.state`, `inbox.designer`, `threads.list`, `messages.list`, `board.get`'s `merge.queue`, and `stats.get` all pay it; `threads.read` averages about 600 ms in the serve log.

Answer in one of three ways:

1. Narrow your touches with `flai touches S-0604 --remove <path>`, or on the task that names the path, if S-0604 will not change it: the hold clears when the claims no longer overlap.
2. Share the paths with `flai message share MS-0391 --paths <path> "<split>"`, or the MCP tool `message_share`, saying who changes what, if the two stories can change them apart: the overlap on them no longer holds, and S-0602 can start.
3. Reply with `message_reply` saying why the hold stands: S-0602 waits until S-0604 moves to review, is cancelled, or is sent back.

### 2026-11-08T02:58:25Z agent-S-0604 S-0604
S-0604 shares `design/system/server-performance.md` with S-0602, split so:

> S-0604 changes only its own section, "Cause 14 fixed: board.get's release lock per root, the watcher's word, and no threads with nothing in review (S-0604)", which comes after the S-0596 Stories table. S-0602 adds and changes only its own cause 13 section. Neither story edits the other's section or the S-0596 tables. Whichever story lands second keeps both sections, in cause order (13 before 14).

The overlap between the two on them no longer holds (ADR-0134).

### 2026-11-08T02:59:11Z agent-S-0602 S-0602 (settled)
Agreed. S-0602 adds and changes only its own cause 13 section in design/system/server-performance.md, leaves S-0604's cause 14 section and the S-0596 tables alone, and keeps both sections in cause order (13 before 14) if it lands second.

### 2026-11-08T03:30:37Z agent-S-0602 S-0602 (notice)
T-2684 of S-0602 changed paths S-0604's claim covers.

T-2684, Measure warm threads.list, inbox.designer, and stats.get before and after, and record them in the design, committed 3afbba7 on story/S-0602, `docs: [S-0602] record cause 13's fix and the reads' before and after`, changing `design/system/flai-cli.md`, `design/system/server-performance.md`. It reaches the main branch when S-0602 is accepted; `git show 3afbba7` shows it until then. Reply here if it breaks your work, or adjust to it early.

### 2026-11-08T03:30:37Z flai S-0602
story/S-0602 and story/S-0604 conflict when merged.

A trial merge of the two at flai stream sync conflicts in:

- `design/system/server-performance.md`

Whichever of S-0602 and S-0604 is accepted second will stop on these paths when it rebases. Agree here who changes what: one narrows its change, or names the other in `after:` and waits for it. The next sync that finds the two merging cleanly closes this conversation. When you do not agree, either of you asks the operator with `flai message escalate` on this conversation, or the MCP tool `message_escalate`, saying what you could not agree.

### 2026-11-08T03:31:07Z agent-S-0604 S-0604 (settled)
T-2684 does not break S-0604. The conflict in `design/system/server-performance.md` is both stories adding a section after the S-0596 Stories table; as agreed, whichever lands second keeps both, in cause order (13 before 14). S-0602's two added lines in `design/system/flai-cli.md` do not touch the `release.Pending` paragraph S-0604 changes.

### 2026-11-08T03:31:09Z agent-S-0602 S-0602 (settled)
The conflict is the one the share foresaw: both stories add their own section after the S-0596 Stories table. Whichever lands second keeps both, cause 13 before cause 14. Neither changes the other's section. In design/system/flai-cli.md, S-0602 adds one paragraph after the S-0156 store's paragraph in Internal structure; keep both there too.

### 2026-11-09T07:17:27Z orchestrator
Closed: S-0602 was accepted
