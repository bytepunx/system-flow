---
id: D-0498
title: "S-0585: the evaluate step exited 2, nothing accepted is waiting for a release: Recommended: do not cancel S-0585"
kind: answer
sources: [TH-0854]
topics: [all]
paths: [wip/kanban/stories/S-0585-publish-what-is-accepted-and-not-yet-released.md]
date: 2026-11-07
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0498 S-0585: the evaluate step exited 2, nothing accepted is waiting for a release: Recommended: do not cancel S-0585

The operator, alex, answered TH-0854 at 2026-11-07T08:15:34Z.

What was asked:

```text
The publish formula's step 2, T-2548 `flai release --evaluate --require-met`, exited 2 in the main checkout. Step 1, the fetch, passed (exit 0). Step 3, `flai release --pending`, was not run, and nothing was tagged or pushed.
```

The answer:

```text
**Recommended:** do not cancel S-0585. Take it to review and accept it with the stop recorded:

1. Unblock T-2548 and close it with its Notes as they stand. With nothing pending, the evaluation stopping is the right result. It is not a refusal.
2. Cancel T-2549 with the reason "step 2 stopped the run: nothing accepted is waiting for a release". Nothing is tagged or pushed.
3. Move S-0585 to review with the evaluation's reason in its Notes, and accept it.

**Why:** S-0585 is the instance S-0562 runs. S-0562's criterion 1 needs the instance "worked to review and accepted", with its Notes naming "the evaluation that stopped it" when nothing was pending. Cancelling S-0585 would leave criterion 1 unmet, and the run would have to be started again. The formula's body says the same: "When step 2 stopped the run, record the evaluation's reason instead". Its blocking rule covers the refusals of step 3: a `conflict:` or the `push` gate.

**What S-0562 does with it:** the formula does not say how an instance stopped at step 2 ends. That is why the agent blocked and asked. S-0562's T-2506 will fix the formula's body to say it, or record it as an issue.

**Alternative:** cancel S-0585, as its agent recommends. S-0562 then needs a criterion reworded, or a second run once a story is pending.

Source: design/workflows/formulas/publish.md
```

Source: TH-0854.
