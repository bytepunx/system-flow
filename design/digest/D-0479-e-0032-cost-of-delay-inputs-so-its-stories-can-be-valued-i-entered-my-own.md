---
id: D-0479
title: "E-0032 cost of delay inputs, so its stories can be valued: I entered my own change since holding up confident releases of signet and authstar is preventing me from shipping commercial software with…"
kind: recommendation
sources: [TH-0829]
topics: [all]
paths: []
date: 2026-11-06
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0479 E-0032 cost of delay inputs, so its stories can be valued: I entered my own change since holding up confident releases of signet and authstar is preventing me from shipping commercial software with…

The operator, alex, answered TH-0829 at 2026-11-06T21:38:41Z.

What was asked:

```text
E-0032 has no cost of delay inputs, and none of its stories will have their own, so `flai cod` refuses each story's value. The inputs are yours.
```

The recommendation it overruled:

````text
**Recommended for the operator:** take the planner's figure, `time_lost_per_cycle: 4h` on E-0032 (600 USD/week at `planning.hour_rate` 150).

```sh
flai edit E-0032 --time-lost-per-cycle 4h
```

- **Why 4h.** E-0032's outcome is removing a hand loop: build an image, load it into a kluster, deploy, and check, per changeset of signet and authstar. Time lost is the input that fits that, not a flat penalty. You gave E-0024 the same 4h when its planner recommended it (D-0342). Nothing on the thread or the epic argues for 2h (signet and authstar changing less than weekly) or for E-0029's flat `penalty_per_week: 150` (D-0383).
- **Why it is yours.** It is an epic's cost of delay input, which is money. `set_cost_of_delay` is off, and `plan_backlog_stories` lets me give inputs only to a backlog story, never to an epic. So you set it, or confirm this and turn `set_cost_of_delay` on for me to set it.
- **What follows.** planner-E-0032 writes each story's value from the epic's input on its next run. I then judge its drafts under `finalize_drafts`.

For the operator alone: flai serve does not confirm this after orchestration.permissions.confirm_after.

Source: design/digest/D-0342-e-0024-needs-a-cost-of-delay-input-and-five-design-answers-recommend-4h-lost.md
````

The operator's answer:

```text
I entered my own change since holding up confident releases of signet and authstar is preventing me from shipping commercial software with waiting subscribers.
```

Source: TH-0829.
