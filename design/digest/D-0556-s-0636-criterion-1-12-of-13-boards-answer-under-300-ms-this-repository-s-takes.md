---
id: D-0556
title: "S-0636 criterion 1: 12 of 13 boards answer under 300 ms; this repository's takes 169 to 404 ms, 5 of 10 calls at 300 ms or more: Recommended: take the agent's option"
kind: answer
sources: [TH-0961, S-0636, I-0333]
topics: [server-side, back-end, flai, cli, go, code, all]
paths: [design/system/server-performance.md, design/issues/summary.md, design/issues/I-0333-this-repository-s-board-get-takes-300-ms-or-more-in-half-its-calls-with-13.md]
date: 2026-11-10
status: active
absorbed_by: ""
asked_again: []
cited: [T-3007, TH-1028]
---

# D-0556 S-0636 criterion 1: 12 of 13 boards answer under 300 ms; this repository's takes 169 to 404 ms, 5 of 10 calls at 300 ms or more: Recommended: take the agent's option

The operator, alex, answered TH-0961 at 2026-11-10T08:40:50Z.

What was asked:

```text
**Recommended:** reword criterion 1 to what S-0602 and S-0606 delivered against `main` before them on the same host, and leave the 300 ms target to I-0326. Then I tick it and move S-0636 to review.
```

The answer:

```text
**Recommended:** take the agent's option. Reword S-0636's criterion 1 to the proposed wording, tick it, and leave the absolute 300 ms target for this repository's board to I-0326.

This is the operator's answer last time, D-0535 (TH-0912): when a performance story misses an absolute target only because of a cause it does not own, its criterion is reworded to what it measured, and the target moves on. D-0532 (TH-0911) and D-0537 (TH-0910) make the same move: the target goes to an issue when no story will measure it. S-0636 changes no code. What is left is the merge queue's git on every board call while stories are in review. No story owns that, and I-0326 records it with a candidate fix.

The two alternatives are worse for the same reasons D-0535 gives. An unticked criterion blocks the orchestrator's acceptance and keeps the story out of the merge queue. Holding S-0636 until I-0326 is fixed keeps an in-progress slot for work it does not own.

One check for the story's agent before review: on `main`, `I-0326` was already used and retired. S-0622 recorded it in 6d666eab1 and folded it into I-0323 in 943417945. The `I-0326` added on `story/S-0636` (9a97a3cc2) reuses that ID. Make sure `flai check --strict` accepts it after the rebase on `main`. If it does not, re-record the issue under the next free ID and name that ID in the reworded criterion and in server-performance.md.

This is a change to the story's scope, so it is the operator's to confirm. Under `confirm_after` (48h), flai serve confirms it if no other answer comes.

Source: design/digest/D-0535-s-0604-criterion-4-12-of-13-boards-answer-under-300-ms-this-repository-s-still.md
```

Applied, from S-0636's Decisions (`wip/archive/agents/S-0636.md`): "TH-0961 (2026-11-10): the operator confirmed the orchestrator's recommendation, reworded criterion 1 to "each answer moving towards 300 ms, and server-performance.md records the before and after", and ticked it. The diff meets it: this repository's p50 fell from 405 to 297 ms, and the other 12 from 66 to 9 ms. The 300 ms target is left to I-0333." The `I-0326` the answer names was a reused ID, recorded again as I-0333 (D-0560). TH-0961 asked what D-0535 had answered.

Source: TH-0961, S-0636, I-0333.
