---
id: ADR-0057
title: "The dashboard charts agent time and cost per item per model over time, instead of completion, aging, and estimates, and charts per item draw item types only"
status: accepted
date: 2026-09-30
supersedes: []
superseded_by: []
refines: [ADR-0053]
topics: [cli, dashboard]
---

# ADR-0057 The dashboard charts agent time and cost per item per model over time, instead of completion, aging, and estimates, and charts per item draw item types only

## Context

The designer found that several charts do not show what was asked of them (S-0169):

- Aging work in progress and Estimate versus actual are not wanted in the flow analysis. Items here carry no `estimate`, so the second chart was always empty.
- Completion over time and Completion against cost (ADR-0051, kept by ADR-0053) plot cumulative items done per model. They do not say how long a model takes over an item, or what an item costs with each model, which is what the designer wants to compare.

Each bucket of spend over time (ADR-0053) already carries, per model, the items the model worked on, their cost, and their seconds of agent work. Tokens per item and cost per item could already draw one line per model for one type, behind a "compare" control.

The designer chose (TH-0040) agent time over cycle time for the time a model takes. Cycle time includes the time a story waits in review for the operator, which is not the model's doing. They also chose to have the new chart per model replace the "compare" control.

`design/system/metrics.md` is the contract between `flai stats` and the dashboard and changes only with an ADR.

## Decision

The dashboard charts the agent time and the cost an item takes per model over time, in place of completion, aging, and estimates, and its charts per item draw item types only.

1. **Removed.** The dashboard no longer draws Aging work in progress, Estimate versus actual, Completion over time, or Completion against cost. `flai stats` still reports aging items, estimates, and the completion series (`usage.done`, `usage.by_model`), for scripts.
2. **Avg. Time / Model.** Per model, one point per bucket in which the model worked on items of the report's type that took agent time: the mean agent minutes of those items. With two models or more, a dashed line gives the same over every item of the bucket.
3. **Avg. Cost / Model.** Per model, one point per bucket in which the model worked on items of the report's type: its cost over those items. With two models or more, a dashed line gives the cost per item of the bucket.
4. **Minutes per item.** Every set of items in spend over time gains `minutes_per_item`: agent seconds over items, in minutes. It is absent when there are no items or none took agent time. For a model it is the mean over the items it worked on, as its seconds are.
5. **Per item, per type.** Tokens per item and cost per item always draw one line per item type. The "compare" control is gone, since the charts per model show the models.
6. **Titles.** Cycle Time, Cumulative Flow, Time in State, Tokens / Min, Tokens / Day, Tokens / $, $ / Day, $ / Work Type, $ / Item. The charts per bucket follow the bucket: Tokens / Hour, $ / Week.

## Consequences

- The operator compares models by the agent time and the dollars an item of one type takes, as each changes over time.
- A model that worked on part of an item is charted with that item's whole agent time. The usage does not split agent time by model, as ADR-0053 already has it for the token rate.
- A dashboard newer than the host's flai works minutes per item out from the seconds, so it draws Avg. Time / Model against flai 1.26.
- A bookmark of a removed chart opens Cycle Time.

## Alternatives considered

- Cycle time per model for Avg. Time / Model: it counts the time in review, which the operator spends, not the model (TH-0040).
- Keeping the "compare" control beside Avg. Cost / Model: the same view could then be reached two ways (TH-0040).
- Taking the completion series and aging items out of `flai stats --json`: it breaks scripts written against it, and dropping a chart does not require it.
