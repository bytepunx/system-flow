---
id: TH-0106
title: Cost of delay inputs for E-0016's stories
anchor:
  path: wip/kanban/epics/E-0016-introduce-orchestration-planning-and-analysis-agents.md
  item: E-0016
status: resolved
participants: [planner-E-0016, alex]
created: 2026-10-04T04:36:05Z
updated: 2026-10-04T04:49:24Z
---

# TH-0106 Cost of delay inputs for E-0016's stories

On wip/kanban/epics/E-0016-introduce-orchestration-planning-and-analysis-agents.md.

## Entries

### 2026-10-04T04:36:05Z planner-E-0016
No cost of delay inputs are set on E-0016 or on any of its stories. `planning.hour_rate` is also unset in `system-flow.yaml`, so a `time_lost_per_cycle` cannot be turned into a value yet. I cannot work out a value for any story without these inputs, so every story's `cost_of_delay.value` stays empty for now. Forecasts and touches go ahead without them.

Recommended: set the inputs once, on the epic, and set the hour rate:

```bash
flai edit E-0016 --time-lost-per-cycle <hours per week the three agents would take off you>
# and in system-flow.yaml: planning.hour_rate: <what an hour of your time costs, in USD>
```

`time_lost_per_cycle` would be the time you spend each week on what the planner, orchestrator, and analyzer take over: writing and sizing stories, promoting and ordering, answering threads, accepting, publishing, and finding bottlenecks. I would then apportion the epic's value over its open stories by forecast duration, as S-0210 specifies for a story with no inputs of its own.

Alternatives:
- give `--penalty-per-week` or `--revenue-per-week` on the epic instead of time lost, and no hour rate is needed;
- give inputs per story, where some stories matter more than their size says (I would suggest the orchestrator stories S-0217 to S-0222);
- leave the cost of delay unset until S-0210 ships `flai cod`.

Reply with the figures, or with which alternative you prefer.

### 2026-10-04T04:42:54Z alex
I took the recommendation and set it to 10h

### 2026-10-04T04:43:52Z planner-E-0016
Thank you; 10h per cycle is on E-0016. One input is still missing: `planning.hour_rate` is unset in `system-flow.yaml`, so 10h cannot become an amount per week yet.

Recommended: set it in the manifest, and reply here with the figure (USD per hour; `planning.currency` defaults to USD):

```yaml
planning:
  hour_rate: <USD per hour of your time>
```

With the default `planning.cycle` of 168h, the epic's value is then 10 × hour_rate per week. I apportion it over the 17 open stories not in progress by forecast duration, and write each story's `value`.

Alternative: reply with a currency amount per week for the epic instead, and I apportion that.

### 2026-10-04T04:48:01Z alex
I set the rate in system-flow.yaml. Let me know if that's not the right file.

### 2026-10-04T04:48:29Z planner-E-0016
`system-flow.yaml` was the right file. With `planning.hour_rate: 150` and the default 168h cycle, E-0016's value is 10h × 150 = **1500 USD per week**, and I set it on the epic.

I apportioned the 1500 over the 19 stories that are not in progress, by forecast duration (1845 forecast minutes, about 0.813 USD per minute). S-0209, S-0225 and S-0255 are in progress and have no forecast, so they get no share.

| Value (USD/week) | Stories |
|------------------|---------|
| 121.95 | S-0218, S-0223 |
| 97.56 | S-0210, S-0217, S-0220, S-0228, S-0229 |
| 85.37 | S-0211, S-0214 |
| 73.17 | S-0213, S-0219, S-0221 |
| 60.98 | S-0212, S-0222, S-0224, S-0227 |
| 48.78 | S-0215, S-0216 |
| 36.59 | S-0226 |

One correction to TH-0107: I revisited 22 open stories, and the 19 not in progress got touches and forecasts, not 17. Splitting the epic's value by size ranks the stories by size only. If some stories matter more than their size, give them inputs of their own, and I will work their values out from those.

### 2026-10-04T04:49:24Z alex
Resolved.
