---
id: TH-0961
title: "S-0636 criterion 1: 12 of 13 boards answer under 300 ms; this repository's takes 169 to 404 ms, 5 of 10 calls at 300 ms or more"
anchor:
  path: wip/kanban/stories/S-0636-13-board-get-calls-at-once-answer-under-300-ms-on-every-served-project-measured.md
  item: S-0636
status: resolved
participants: [agent-S-0636, orchestrator, alex]
created: 2026-11-09T23:09:52Z
updated: 2026-11-10T08:48:59Z
---

# TH-0961 S-0636 criterion 1: 12 of 13 boards answer under 300 ms; this repository's takes 169 to 404 ms, 5 of 10 calls at 300 ms or more

On wip/kanban/stories/S-0636-13-board-get-calls-at-once-answer-under-300-ms-on-every-served-project-measured.md.

## Entries

### 2026-11-09T23:09:52Z agent-S-0636
**Recommended:** reword criterion 1 to what S-0602 and S-0606 delivered against `main` before them on the same host, and leave the 300 ms target to I-0326. Then I tick it and move S-0636 to review.

Proposed wording: "13 `board.get` calls at once on the operator's host, one per served project: the other 12 answer under 300 ms in every call, this repository's p50 falls against `main` before S-0602 (8b29cf861) on the same host, its remaining time is shown to be the merge queue's git with stories in review, and server-performance.md records the before and after."

**What was measured** (2026-11-09, 23:04 to 23:05 UTC, load 96 to 123 on 24 cores; `design/system/server-performance.md` § 13 boards at once, measured again after S-0602 and S-0606 (S-0636)):

| 13 at once, a call | Before, 8b29cf861 | After, f667efc2d |
|---|---|---|
| The other 12, p50 (max) | 66 ms (359) | 9 ms (216) |
| The other 12, under 300 ms | 115 of 120 | 120 of 120 |
| This repository, p50 (range) | 405 ms (244 to 657) | 297 ms (169 to 404) |
| This repository, under 300 ms | 3 of 10 | 5 of 10 |

- S-0602 and S-0606 took out what S-0604 left them: this repository's `repo.list` fell from 103 to 32 ms at p50, `release.pending` from 41 to 20 ms, and the thread read no longer shows in `merge.queue`.
- What is left is the merge queue's git, which runs on every board call while stories are in review (5 were): `git status --porcelain` on the main checkout, 52 to 148 ms, and five `git rev-parse` one after another, 21 to 81 ms. No story owns that; I-0326 records it with a candidate fix (keep those reads while the main-refs watcher saw no move, as S-0604 did for `show-ref`).

**Why reword:** this story changes no code, so it cannot reach 300 ms itself. D-0535 (TH-0912) and D-0532 (TH-0911) answered the same case: the criterion is reworded to what was measured, and the absolute target goes to an issue when no story measures it.

**Alternatives:**

- Leave criterion 1 unticked and move S-0636 to review for you to accept by hand. An unticked criterion blocks the orchestrator's acceptance.
- Keep S-0636 open until a story fixes I-0326 and measure again. It holds an in-progress slot for work it does not own.

Source: design/digest/D-0535-s-0604-criterion-4-12-of-13-boards-answer-under-300-ms-this-repository-s-still.md

### 2026-11-09T23:12:01Z orchestrator (recommendation)
**Recommended:** take the agent's option. Reword S-0636's criterion 1 to the proposed wording, tick it, and leave the absolute 300 ms target for this repository's board to I-0326.

This is the operator's answer last time, D-0535 (TH-0912): when a performance story misses an absolute target only because of a cause it does not own, its criterion is reworded to what it measured, and the target moves on. D-0532 (TH-0911) and D-0537 (TH-0910) make the same move: the target goes to an issue when no story will measure it. S-0636 changes no code. What is left is the merge queue's git on every board call while stories are in review. No story owns that, and I-0326 records it with a candidate fix.

The two alternatives are worse for the same reasons D-0535 gives. An unticked criterion blocks the orchestrator's acceptance and keeps the story out of the merge queue. Holding S-0636 until I-0326 is fixed keeps an in-progress slot for work it does not own.

One check for the story's agent before review: on `main`, `I-0326` was already used and retired. S-0622 recorded it in 6d666eab1 and folded it into I-0323 in 943417945. The `I-0326` added on `story/S-0636` (9a97a3cc2) reuses that ID. Make sure `flai check --strict` accepts it after the rebase on `main`. If it does not, re-record the issue under the next free ID and name that ID in the reworded criterion and in server-performance.md.

This is a change to the story's scope, so it is the operator's to confirm. Under `confirm_after` (48h), flai serve confirms it if no other answer comes.

Source: design/digest/D-0535-s-0604-criterion-4-12-of-13-boards-answer-under-300-ms-this-repository-s-still.md

### 2026-11-10T08:40:50Z alex
Confirmed the recommendation of 2026-11-09T23:12:01Z orchestrator.

Source: design/digest/D-0535-s-0604-criterion-4-12-of-13-boards-answer-under-300-ms-this-repository-s-still.md

### 2026-11-10T08:48:59Z alex
Resolved: S-0636 was accepted
