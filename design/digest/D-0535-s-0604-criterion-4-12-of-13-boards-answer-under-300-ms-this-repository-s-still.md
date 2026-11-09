---
id: D-0535
title: "S-0604 criterion 4: 12 of 13 boards answer under 300 ms; this repository's still takes 342 to 761 ms, from causes 13 and 16: Recommended: take option 2, and move the full 13-board measure to a…"
kind: answer
sources: [TH-0912]
topics: [server-side, back-end, flai, cli, go, code, all]
paths: [flai/internal/release/history.go, flai/internal/mergequeue/eligible.go, flai/internal/hostapi/hostapi.go, flai/internal/release/release.go, flai/internal/release/history_test.go, flai/internal/serve/mergequeue.go, flai/internal/serve/mergequeue_test.go, flai/internal/mergequeue/eligible_test.go, flai/internal/hostapi/reads_test.go, flai/internal/storygit/undone.go, design/system/server-performance.md, design/system/flai-cli.md]
date: 2026-11-09
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0535 S-0604 criterion 4: 12 of 13 boards answer under 300 ms; this repository's still takes 342 to 761 ms, from causes 13 and 16: Recommended: take option 2, and move the full 13-board measure to a…

The operator, alex, answered TH-0912 at 2026-11-09T07:35:10Z.

What was asked:

```text
**Recommendation:** accept criterion 4 as met for what cause 14 owns, and leave it unticked with a note in S-0604's Notes. The rest of this repository's time belongs to S-0602 (cause 13) and S-0606 (cause 16).
```

The answer:

```text
**Recommended:** take option 2, and move the full 13-board measure to a follow-up. This is how you answered for S-0424 on TH-0810 (D-0465): a criterion the story cannot meet comes out, and what is left of it moves to a draft that waits.

- agent-S-0604 rewords criterion 4 to: "13 `board.get` calls at once on the operator's host, one per served project: the other 12 answer under 300 ms, this repository's remaining time is shown to be causes 13 and 16, and server-performance.md records the before and after". Then it ticks the criterion.
- It writes a draft story in E-0033 with `after: [S-0602, S-0606]`. The draft holds criterion 4's original words: all 13 boards under 300 ms, measured again.
- S-0604's Notes record the move.

**Why not option 1.** An unticked criterion blocks every acceptance but yours by hand. ADR-0093 § Decision makes "every acceptance criterion is ticked" a blocker for the orchestrator. The merge queue lists `criteria_unticked` as a reason the same way. S-0604 would sit in review until you accept it yourself.

**Why not option 3.** S-0602 is in ready and not started. S-0606 is in progress and blocked. Holding S-0604 keeps an in-progress slot for work it does not own.

**Why I recommend rather than answer.** Changing a finalized story's criterion and adding a story changes scope. That is your call. Finalizing the draft is yours too.

Source: design/digest/D-0465-s-0424-the-two-live-runs-criterion-4-t-1787-t-1788-cannot-run-before-s-0424-is.md
```

Source: TH-0912.
