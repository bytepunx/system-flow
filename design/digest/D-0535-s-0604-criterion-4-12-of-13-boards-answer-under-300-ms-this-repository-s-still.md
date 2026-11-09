---
id: D-0535
title: "S-0604 criterion 4: 12 of 13 boards answer under 300 ms; this repository's still takes 342 to 761 ms, from causes 13 and 16: Recommended: take option 2, and move the full 13-board measure to a…"
kind: answer
sources: [TH-0912, S-0604, S-0636]
topics: [server-side, back-end, flai, cli, go, code, all]
paths: [design/system/server-performance.md, flai/internal/release/history.go, flai/internal/mergequeue/eligible.go, flai/internal/hostapi/hostapi.go]
date: 2026-11-09
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0535 S-0604 criterion 4: 12 of 13 boards answer under 300 ms; this repository's still takes 342 to 761 ms, from causes 13 and 16: Recommended: take option 2, and move the full 13-board measure to a…

The operator, alex, answered TH-0912 at 2026-11-09T07:35:10Z by confirming the orchestrator's recommendation, option 2.

What was asked: S-0604's criterion 4 wanted all 13 `board.get` calls at once under 300 ms. After S-0604, the other 12 projects' boards answered at 10 ms p50, but this repository's still took 342 to 761 ms. None of that time was git or cause 14: about 200 ms was `merge.queue` reading every thread with stories in review (cause 13, S-0602), and 60 to 270 ms the item list and holds (cause 16, S-0606).

The answer:

```text
**Recommended:** take option 2, and move the full 13-board measure to a follow-up. This is how you answered for S-0424 on TH-0810 (D-0465): a criterion the story cannot meet comes out, and what is left of it moves to a draft that waits.
```

Why not leave the criterion unticked for the operator to accept by hand: an unticked criterion is a blocker for the orchestrator's acceptance (ADR-0093 § Decision) and `criteria_unticked` keeps a story out of the merge queue, so the story would sit in review. Why not hold S-0604 until S-0602 and S-0606 land: it keeps an in-progress slot for work it does not own.

Applied: S-0604's criterion 4 now reads "the other 12 answer under 300 ms, this repository's remaining time is shown to be causes 13 and 16, and server-performance.md records the before and after", and is ticked. Its original words moved to the draft S-0636 in E-0033, `after: [S-0602, S-0606]`, and S-0604's Notes record the move.

So a performance story whose absolute target is missed only because of causes other stories own has its criterion reworded to what its own cause owns, and the absolute target moves to a draft story whose `after` names the stories that own the rest, where it is measured again. D-0532 is the same answer for S-0603, where the target went to an issue because no story would measure it; D-0465 and D-0466 are the same move for a run that cannot happen before acceptance.

Sources: TH-0912, S-0604, S-0636.
