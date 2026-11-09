---
id: TH-0913
title: "S-0602 criterion 3: stats.get's reads spend 260 ms less CPU but only 46 ms less wall time, since S-0607 runs them at once"
anchor:
  path: wip/kanban/stories/S-0602-flai-serve-keeps-threads-conversations-digest-entries-issues-and-activity.md
  item: S-0602
status: resolved
participants: [agent-S-0602, orchestrator, alex]
created: 2026-11-08T03:31:34Z
updated: 2026-11-09T07:15:44Z
---

# TH-0913 S-0602 criterion 3: stats.get's reads spend 260 ms less CPU but only 46 ms less wall time, since S-0607 runs them at once

On wip/kanban/stories/S-0602-flai-serve-keeps-threads-conversations-digest-entries-issues-and-activity.md.

## Entries

### 2026-11-08T03:31:34Z agent-S-0602
**Recommendation:** answer "met". Criterion 3 reads as the work the reads do, and that falls by 260 ms of CPU. The wall-time shortfall comes from S-0607, accepted after S-0602 was planned, which already overlaps the reads.

Criterion 3 asks that a warm `threads.list`, `inbox.designer`, and `stats.get` each spend at least 80 ms less in their reads. Warm medians on this repository, main at 17097cea1 against story/S-0602, are below. `design/system/server-performance.md` § Cause 13 fixed has the full tables and the method.

| Method | Reads before | Reads after | Less |
|--------|-------------:|------------:|-----:|
| `threads.list`, `threads.read` | 262 ms | 14 ms | 248 ms |
| `inbox.designer`, `threads.read` | 233 ms | 13 ms | 220 ms |
| `stats.get`, `statsread.Read` wall time | 96 ms | 50 ms | 46 ms |
| `stats.get`, `statsread.Read` CPU time | 389 ms | 129 ms | 260 ms |

- **Wall time falls only 46 ms for stats.get.** S-0607 runs `statsread.Read`'s twelve reads at once, so its wall time is that of the slowest. Taken one at a time, the five reads this story keeps fell from 292 ms to 18 ms. What is left is the reads this story does not keep: the stories' commits, the board, the conflict log, and the merge queue's events.
- **`stats.get`'s `stats.compute` phase fell 135 ms,** from 942 to 807 ms. That phase holds both the reads and the computation, and the host was loaded (load average 37 to 49), so the figure is noisy.

## Why this needs you

The criteria are yours. "Spend 80 ms less in their reads" could mean CPU or wall time. If you mean wall time, the stats.get part is not met, and meeting it would need the remaining reads kept too, which is beyond cause 13.

What I will do with each answer:

- **"met":** tick criterion 3 and move S-0602 to review.
- **"not met":** leave it unticked, say why in the story's Notes, and move S-0602 to review for you to decide at acceptance, or follow what you say instead.

### 2026-11-08T03:32:48Z orchestrator (recommendation)
**Recommended:** answer "met". Tick criterion 3, and say in S-0602's Notes that `stats.get` was judged on the reads' CPU and on the kept reads taken one at a time. The wall-time figure is bounded by reads this story does not keep.

**Why "met":**

- The design measured cause 13 per read, one read at a time, before S-0607 overlapped them: threads 80 to 110 ms, digest 60 to 85 ms, and so on. Its estimate for `stats.get` is "about 180 to 200 ms", the work the parses cost. Measured the same way, the five kept reads fell from 292 ms to 18 ms, and the reads' CPU fell 260 ms.
- `threads.list` and `inbox.designer` meet 80 ms in wall time too: 248 ms and 220 ms less.
- The wall shortfall comes from cause 17. S-0607, accepted after S-0602 was planned, runs `statsread.Read`'s twelve reads at once. Its wall time is now the slowest read's, and that read is one cause 13 does not cover: the commits, the board, the conflict log, or the queue's events.
- Meeting 80 ms of wall time would mean keeping those reads as well. That is beyond cause 13 and would widen the story's scope.

**What makes it hold:** `design/system/server-performance.md` § Cause 13 fixed records both the CPU and the wall figures, and the method. Criterion 3 asks for that record.

**Alternative:** answer "not met". The agent leaves criterion 3 unticked, explains it in the Notes, and moves S-0602 to review. The story then waits in review for your call, as S-0598 and S-0601 do on TH-0907 and TH-0910.

Source: design/system/server-performance.md § Where the time goes

### 2026-11-09T07:14:59Z alex
Confirmed the recommendation of 2026-11-08T03:32:48Z orchestrator.

Source: design/system/server-performance.md § Where the time goes

### 2026-11-09T07:15:44Z agent-S-0602
Resolved: Criterion 3 is met by the operator's confirmation. It is ticked, the story's Notes record that stats.get was judged on its reads' CPU and on the kept reads taken one at a time, and S-0602 moves to review.
