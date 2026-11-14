---
id: D-0622
title: "S-0737 cannot be finalized yet: no forecast delivery and no cost of delay value until E-0034 has inputs: the epic now has a 10h penalty associated with it"
kind: answer
sources: [TH-1114, TH-1115, TH-1097]
topics: [cli, template, flai, go, flaiover, dashboard, sveltekit, conventions, code, all]
paths: [wip/kanban/epics]
date: 2026-11-14
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0622 S-0737 cannot be finalized yet: no forecast delivery and no cost of delay value until E-0034 has inputs: the epic now has a 10h penalty associated with it

E-0034's cost of delay input is `time_lost_per_cycle: 10h`, set on the epic by the operator, alex. It overrules the planner's recommended 2h. flai shares the epic's value out to its stories by their forecasts, so no story of E-0034 needs inputs of its own.

The planner asked on TH-1097, recommending 2h, because E-0034 and its stories had no inputs and `flai cod` could value none of them. The operator answered at 2026-11-14T17:42:34Z:

```text
I set the CoD on the epic to 10h
```

The orchestrator could not finalize S-0737 (TH-1114) or S-0738 (TH-1115) while the epic had no inputs. The operator answered both the same way:

```text
the epic now has a 10h penalty associated with it
the epic has a 10h per cycle CoD
```

E-0034's front matter holds `time_lost_per_cycle: 10h`. The operator called it a penalty on TH-1114, but the input is time lost per cycle, not `penalty_per_week`. A later story of E-0034 takes its value from this input, and nobody asks the operator again.

Sources: TH-1097, TH-1114, TH-1115. D-0623 and D-0624 recorded the same answers and are merged here.
