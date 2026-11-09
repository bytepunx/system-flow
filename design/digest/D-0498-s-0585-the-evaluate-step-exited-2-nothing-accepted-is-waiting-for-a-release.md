---
id: D-0498
title: "S-0585: the evaluate step exited 2, nothing accepted is waiting for a release: Recommended: do not cancel S-0585"
kind: answer
sources: [TH-0854, S-0585]
topics: [all]
paths: [wip/kanban/stories/S-0585-publish-what-is-accepted-and-not-yet-released.md]
date: 2026-11-07
status: absorbed
absorbed_by: "design/adrs/0271-a-formula-step-may-name-the-exits-that-stop-its-instance-cleanly-flai-formula.md#decision"
asked_again: []
cited: [S-0621, T-2743, TH-0934]
---

# D-0498 S-0585: the evaluate step exited 2, nothing accepted is waiting for a release: Recommended: do not cancel S-0585

A publish formula instance whose step 2, `flai release --evaluate --require-met`, stops because nothing accepted is waiting for a release is not cancelled: it is taken to review and accepted with the stop recorded.

The operator, alex, confirmed agent-S-0562's recommendation on TH-0854 at 2026-11-07T08:15:34Z:

```text
1. Unblock T-2548 and close it with its Notes as they stand. With nothing pending, the evaluation stopping is the right result. It is not a refusal.
2. Cancel T-2549 with the reason "step 2 stopped the run: nothing accepted is waiting for a release". Nothing is tagged or pushed.
3. Move S-0585 to review with the evaluation's reason in its Notes, and accept it.
```

Why: S-0562's criterion 1 needs the instance "worked to review and accepted", with its Notes naming the evaluation that stopped it when nothing was pending. The formula's body says "When step 2 stopped the run, record the evaluation's reason instead" (`design/workflows/formulas/publish.md`).

S-0585's narrative records it as its decision: T-2548 closed with its Notes, T-2549 cancelled because step 2 stopped the run, and S-0585 sent to review with the evaluation's reason in its Notes. This departs from the formula instance rule to cancel no step; the operator's word takes precedence. Its criterion 3 stays unticked, since no publish ran and nothing was waiting. S-0562's T-2506 is to make the formula say how an instance stopped at step 2 ends.

Sources: TH-0854, S-0585 (`wip/archive/agents/S-0585.md`).
