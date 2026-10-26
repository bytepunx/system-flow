---
id: D-0383
title: "Cost of delay inputs for E-0029: Recommended for the operator: take the planner's figure, `penalty_per_week: 150` USD on E-0029"
kind: answer
sources: [TH-0657]
topics: [all]
paths: []
date: 2026-10-26
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0383 Cost of delay inputs for E-0029: Recommended for the operator: take the planner's figure, `penalty_per_week: 150` USD on E-0029

The operator, alex, answered TH-0657 at 2026-10-26T12:16:13Z.

What was asked:

```text
E-0029 has no cost of delay inputs, so `flai cod` refuses a value for it and for each story I draft under it (a story without inputs takes a share of its epic's value). The inputs are yours.
```

The answer:

````text
**Recommended for the operator:** take the planner's figure, `penalty_per_week: 150` USD on E-0029.

```sh
flai edit E-0029 --penalty-per-week 150
```

- **Why 150.** It is three times the 50 USD/week you gave S-0496 and S-0497, the research behind this epic, for a six-story delivery. Nothing on the thread or the epic argues for 50 (not urgent) or 300 (needed before E-0023).
- **Why it is yours.** An epic's input is money and the operator's, as on TH-0467 (D-0330). Under plan_backlog_stories I give inputs only to a backlog story, never to an epic. Each story's value is the epic's input shared by duration, so setting it once here values all six.
- **What follows.** The planner writes each story's value on its next run. I then judge S-0505 to S-0510 for finalizing under finalize_drafts.

Source: design/digest/D-0330-s-0407-plan-five-tasks-in-four-layers-two-touches-added-needs-a-cost-of-delay.md
````

The operator, alex, answered TH-0657 again at 2026-10-26T12:27:56Z.

The answer:

```text
TH-0659 is settled
```

Source: TH-0657.
