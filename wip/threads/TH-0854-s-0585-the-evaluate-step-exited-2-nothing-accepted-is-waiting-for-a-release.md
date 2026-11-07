---
id: TH-0854
title: "S-0585: the evaluate step exited 2, nothing accepted is waiting for a release"
anchor:
  path: wip/kanban/stories/S-0585-publish-what-is-accepted-and-not-yet-released.md
  item: S-0585
status: resolved
participants: [agent-S-0585, agent-S-0562, alex, system-flow]
created: 2026-11-07T08:06:55Z
updated: 2026-11-07T08:16:33Z
---

# TH-0854 S-0585: the evaluate step exited 2, nothing accepted is waiting for a release

On wip/kanban/stories/S-0585-publish-what-is-accepted-and-not-yet-released.md.

## Entries

### 2026-11-07T08:06:55Z agent-S-0585
The publish formula's step 2, T-2548 `flai release --evaluate --require-met`, exited 2 in the main checkout. Step 1, the fetch, passed (exit 0). Step 3, `flai release --pending`, was not run, and nothing was tagged or pushed.

The evaluation says:

```text
release policy judgement: not met
the release policy is judgement, which is never met by itself: whether to release is the orchestrator's or the operator's call, with no accepted story waiting for a release
value  0 USD/week
count  0
nothing accepted is waiting for a release
```

T-2548 is blocked with this output in its Notes, and S-0585's Notes record the reason.

What should S-0585 do?

1. Recommended: cancel S-0585. There is nothing to publish, and the formula's stop is the right result.
2. Unblock T-2548 and run the steps again later, once a story is accepted and waiting for a release.
3. Something else: say what.

### 2026-11-07T08:07:08Z agent-S-0562 (recommendation)
**Recommended:** do not cancel S-0585. Take it to review and accept it with the stop recorded:

1. Unblock T-2548 and close it with its Notes as they stand. With nothing pending, the evaluation stopping is the right result. It is not a refusal.
2. Cancel T-2549 with the reason "step 2 stopped the run: nothing accepted is waiting for a release". Nothing is tagged or pushed.
3. Move S-0585 to review with the evaluation's reason in its Notes, and accept it.

**Why:** S-0585 is the instance S-0562 runs. S-0562's criterion 1 needs the instance "worked to review and accepted", with its Notes naming "the evaluation that stopped it" when nothing was pending. Cancelling S-0585 would leave criterion 1 unmet, and the run would have to be started again. The formula's body says the same: "When step 2 stopped the run, record the evaluation's reason instead". Its blocking rule covers the refusals of step 3: a `conflict:` or the `push` gate.

**What S-0562 does with it:** the formula does not say how an instance stopped at step 2 ends. That is why the agent blocked and asked. S-0562's T-2506 will fix the formula's body to say it, or record it as an issue.

**Alternative:** cancel S-0585, as its agent recommends. S-0562 then needs a criterion reworded, or a second run once a story is pending.

Source: design/workflows/formulas/publish.md

### 2026-11-07T08:15:34Z alex
Confirmed the recommendation of 2026-11-07T08:07:08Z agent-S-0562.

Source: design/workflows/formulas/publish.md

### 2026-11-07T08:16:33Z system-flow
Resolved: the operator confirmed agent-S-0562's recommendation: T-2548 closed as it stands, T-2549 cancelled, S-0585 to review
