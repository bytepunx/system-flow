---
id: D-0531
title: "S-0602 criterion 3: stats.get's reads spend 260 ms less CPU but only 46 ms less wall time, since S-0607 runs them at once: Recommended: answer \"met\""
kind: answer
sources: [TH-0913, S-0602]
topics: [server-side, back-end, flai, cli, go, code, all]
paths: [flai/internal/threads/threads.go, flai/internal/messages/messages.go, flai/internal/digest/digest.go, flai/internal/issues/issues.go, flai/internal/threads/threads_test.go, flai/internal/messages/messages_test.go, flai/internal/workitem/share.go, flai/internal/workitem/share_test.go, flai/internal/digest/digest_test.go, flai/internal/issues/issues_test.go, flai/internal/workitem/activity.go, flai/internal/workitem/activity_test.go, design/system/server-performance.md, design/system/flai-cli.md, flai/internal/workitem/keptfiles.go, flai/internal/workitem/keptfiles_test.go]
date: 2026-11-09
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0531 S-0602 criterion 3: stats.get's reads spend 260 ms less CPU but only 46 ms less wall time, since S-0607 runs them at once: Recommended: answer "met"

The operator, alex, answered TH-0913 at 2026-11-09T07:14:59Z.

What was asked:

```text
**Recommendation:** answer "met". Criterion 3 reads as the work the reads do, and that falls by 260 ms of CPU. The wall-time shortfall comes from S-0607, accepted after S-0602 was planned, which already overlaps the reads.
```

The answer:

```text
**Recommended:** answer "met". Tick criterion 3, and say in S-0602's Notes that `stats.get` was judged on the reads' CPU and on the kept reads taken one at a time. The wall-time figure is bounded by reads this story does not keep.

**Why "met":**

- The design measured cause 13 per read, one read at a time, before S-0607 overlapped them: threads 80 to 110 ms, digest 60 to 85 ms, and so on. Its estimate for `stats.get` is "about 180 to 200 ms", the work the parses cost. Measured the same way, the five kept reads fell from 292 ms to 18 ms, and the reads' CPU fell 260 ms.
- `threads.list` and `inbox.designer` meet 80 ms in wall time too: 248 ms and 220 ms less.
- The wall shortfall comes from cause 17. S-0607, accepted after S-0602 was planned, runs `statsread.Read`'s twelve reads at once. Its wall time is now the slowest read's, and that read is one cause 13 does not cover: the commits, the board, the conflict log, or the queue's events.
- Meeting 80 ms of wall time would mean keeping those reads as well. That is beyond cause 13 and would widen the story's scope.

**What makes it hold:** `design/system/server-performance.md` § Cause 13 fixed records both the CPU and the wall figures, and the method. Criterion 3 asks for that record.

**Alternative:** answer "not met". The agent leaves criterion 3 unticked, explains it in the Notes, and moves S-0602 to review. The story then waits in review for your call, as S-0598 and S-0601 do on TH-0907 and TH-0910.

Source: design/system/server-performance.md § Where the time goes
```

The operator confirmed the recommendation, and S-0602 recorded how the criterion was judged (`wip/archive/agents/S-0602.md`):

```text
The operator, alex, answered TH-0913 on 2026-11-09T07:14:59Z by confirming the orchestrator's recommendation: criterion 3 is met. `stats.get` is judged on its reads' CPU and on the kept reads taken one at a time, since its wall time is bounded by reads cause 13 does not cover. Criterion 3 is ticked and the story's Notes say so.
```

So a criterion that asks a method to spend less time in its reads is met by the reads' CPU and by the story's own reads timed one at a time, when the method runs its reads at once (S-0607) and its wall time is bounded by a read the story does not change. The story records both the CPU and the wall figures in its design section.

Sources: TH-0913, S-0602.
