---
id: TH-0912
title: "S-0604 criterion 4: 12 of 13 boards answer under 300 ms; this repository's still takes 342 to 761 ms, from causes 13 and 16"
anchor:
  path: wip/kanban/stories/S-0604-board-get-takes-one-release-lock-per-repository-runs-show-ref-only-when-a-ref.md
  item: S-0604
status: resolved
participants: [agent-S-0604, orchestrator, alex]
created: 2026-11-08T02:38:13Z
updated: 2026-11-09T07:41:17Z
---

# TH-0912 S-0604 criterion 4: 12 of 13 boards answer under 300 ms; this repository's still takes 342 to 761 ms, from causes 13 and 16

On wip/kanban/stories/S-0604-board-get-takes-one-release-lock-per-repository-runs-show-ref-only-when-a-ref.md.

## Entries

### 2026-11-08T02:38:13Z agent-S-0604
**Recommendation:** accept criterion 4 as met for what cause 14 owns, and leave it unticked with a note in S-0604's Notes. The rest of this repository's time belongs to S-0602 (cause 13) and S-0606 (cause 16).

## What was measured

13 `board.get` calls at once, one per served project, in one process with serve's main-refs watcher, 130 calls per build, load 80 to 104 on 24 cores (`design/system/server-performance.md` § Cause 14 fixed).

| A call | Before (`main`) | After (S-0604) |
|---|---|---|
| Other 12 projects, p50 (max) | 36 ms (190) | 10 ms (293) |
| This repository, p50 (range) | 415 ms (274 to 529) | 442 ms (342 to 761) |
| Calls that ran `git show-ref` | 70 of 130 | 0 of 130 |
| Calls that ran `git status` | 130 of 130 | 40 of 130 |

## Why this repository still misses 300 ms

None of its remaining time is git, and none is cause 14:

- Four of its stories are in review, so `merge.queue` reads all its threads, about 200 ms. That is cause 13, S-0602.
- Its `repo.list` and `board.view` (the item list and the holds) take 60 to 270 ms. That is cause 16, S-0606.
- `release.pending` takes 20 to 76 ms with no git process, working out plans from every item.

The shared lock and the git processes the story removes are what made the serve log's 800 ms p50 queue. A one-process harness does not reproduce that contention: `main` measured 41 ms p50 there.

## Options

1. **Accept as is (recommended).** I leave criterion 4 unticked and say why in Notes; you accept with the shortfall recorded.
2. **Reword criterion 4** to "the other 12 projects answer under 300 ms, and this repository's remaining time is causes 13 and 16", and I tick it.
3. **Hold S-0604** until S-0602 and S-0606 land, then measure again.

Source: design/system/server-performance.md § Cause 14 fixed: board.get's release lock per root, the watcher's word, and no threads with nothing in review (S-0604)

### 2026-11-08T02:39:49Z orchestrator (recommendation)
**Recommended:** take option 2, and move the full 13-board measure to a follow-up. This is how you answered for S-0424 on TH-0810 (D-0465): a criterion the story cannot meet comes out, and what is left of it moves to a draft that waits.

- agent-S-0604 rewords criterion 4 to: "13 `board.get` calls at once on the operator's host, one per served project: the other 12 answer under 300 ms, this repository's remaining time is shown to be causes 13 and 16, and server-performance.md records the before and after". Then it ticks the criterion.
- It writes a draft story in E-0033 with `after: [S-0602, S-0606]`. The draft holds criterion 4's original words: all 13 boards under 300 ms, measured again.
- S-0604's Notes record the move.

**Why not option 1.** An unticked criterion blocks every acceptance but yours by hand. ADR-0093 § Decision makes "every acceptance criterion is ticked" a blocker for the orchestrator. The merge queue lists `criteria_unticked` as a reason the same way. S-0604 would sit in review until you accept it yourself.

**Why not option 3.** S-0602 is in ready and not started. S-0606 is in progress and blocked. Holding S-0604 keeps an in-progress slot for work it does not own.

**Why I recommend rather than answer.** Changing a finalized story's criterion and adding a story changes scope. That is your call. Finalizing the draft is yours too.

Source: design/digest/D-0465-s-0424-the-two-live-runs-criterion-4-t-1787-t-1788-cannot-run-before-s-0424-is.md

### 2026-11-09T07:35:10Z alex
Confirmed the recommendation of 2026-11-08T02:39:49Z orchestrator.

Source: design/digest/D-0465-s-0424-the-two-live-runs-criterion-4-t-1787-t-1788-cannot-run-before-s-0424-is.md

### 2026-11-09T07:41:17Z agent-S-0604
Resolved: Option 2 done: S-0604's criterion 4 reworded and ticked, its original words moved to the draft S-0636 (after S-0602, S-0606), and the move recorded in S-0604's Notes and server-performance.md's cause 14 section.
