---
id: D-0479
title: "E-0032 cost of delay inputs, so its stories can be valued: I entered my own change since holding up confident releases of signet and authstar is preventing me from shipping commercial software with…"
kind: recommendation
sources: [TH-0829, E-0032]
topics: [all, planning, kluster]
paths: [wip/kanban/epics/E-0032-integrate-kluster-for-demo-and-testing.md]
date: 2026-11-06
status: active
absorbed_by: ""
asked_again: []
cited: [TH-0829, TH-0979, TH-0980, TH-0981, TH-0982, TH-0983]
---

# D-0479 E-0032 cost of delay inputs, so its stories can be valued: I entered my own change since holding up confident releases of signet and authstar is preventing me from shipping commercial software with…

E-0032 (integrate kluster for demo and testing) is valued as revenue, not as time lost: the operator set `revenue_per_week: 1250` USD on the epic themselves, at 2026-11-06T21:36:25Z, in place of the 4h `time_lost_per_cycle` (600 USD/week) that planner-E-0032 and the orchestrator recommended on TH-0829.

The operator's reason, on TH-0829 at 2026-11-06T21:38:41Z:

```text
I entered my own change since holding up confident releases of signet and authstar is preventing me from shipping commercial software with waiting subscribers.
```

So the epic's delay costs revenue the operator's commercial projects, signet and authstar, cannot earn while their releases wait on kluster-backed integration and smoke tiers. It is not only the hand loop's hours. A later planner or orchestrator that values E-0032's stories, or an epic that unblocks signet's or authstar's releases, should not offer the time-lost figure again. Ask for the revenue those releases hold up instead.

planner-E-0032 wrote each story's value, S-0575 to S-0581, as `flai cod`'s share of 1250 USD/week by forecast duration (TH-0830). The orchestrator then finalized them.

The recommendation it overruled cited D-0342 (E-0024's 4h) and D-0383 (E-0029's flat `penalty_per_week: 150`). Neither applies here.

Sources: TH-0829, E-0032.
