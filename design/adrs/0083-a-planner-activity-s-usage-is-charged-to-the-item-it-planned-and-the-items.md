---
id: ADR-0083
title: "A planner activity's usage is charged to the item it planned and the items above it, under usage.strategic per agent kind, apart from the agents' figures"
status: accepted
date: 2026-10-04
supersedes: []
superseded_by: []
refines: [ADR-0051, ADR-0079]
topics: [cli, dashboard, planning]
---

# ADR-0083 A planner activity's usage is charged to the item it planned and the items above it, under usage.strategic per agent kind, apart from the agents' figures

## Context

ADR-0051 puts on each item what the story agents spent on it, and sums it up the hierarchy. ADR-0079 logs what the planner, the orchestrator, and the analyzer spend in their activity documents, and on no item. The planner runs for one epic or one story and ends (ADR-0082), so what it spends is spent on that item, and without this it is missing from the item's totals and from the charts. The designer asked (S-0225) that it land on the item it planned, marked as strategic, so that the item's own agent cost stays readable, and that an item the planner forecast show what it is expected to cost before any agent works it.

## Decision

An item's `usage` may carry `strategic`: one entry per strategic agent kind, with `kind`, `seconds`, `estimated: true`, and `models` in the shape the agents' models have.

1. **Charged when a planner activity is logged.** When flai serve logs a planner activity, as its run ends or through `activity_log` during it, the usage apportioned to the activity's span, the same share its entry's cost is (ADR-0079), is added to the planner's entry on the item the run was started for, as `serve/agents.json` records it under `plans`. An activity outside a run flai serve logged charges nothing.
2. **Summed up the hierarchy when it is charged.** The same charge is added to every item above the planned one, up to its epic, at once. Roll-up, which sums the agents' figures from the children, and a story's measurement from its logs keep `strategic` as it is. An item with no usage gets `source: sum` with no seconds and no models beside it.
3. **Apart from the agents' figures.** An item's tokens, cost, seconds, and models are the agents' alone, and so is every per-model figure in `flai stats` and the dashboard. `flai stats` reports strategic usage per item, in totals, and in spend as figures of its own, and the cost charts draw it as a series of its own.
4. **Expected cost.** An item with a `forecast.duration`, or else an `estimate`, has an expected cost: that duration in hours times the project's mean cost per agent hour, the agents' cost over their hours across the stories measured from their logs, archived ones included. It is marked estimated and shown beside the forecast or estimate in `flai stats`, `flai show`, and the item's page, and it is absent while no story has been measured.

## Consequences

- A planner run's cost is in two places: on the items it planned and in `wip/agents/planner.md`. They are two views of one spend; nothing adds them together.
- A story moved to another epic after it was planned leaves its planning on the old epic: the charge was summed when it was made.
- An item may carry usage with nothing but `strategic`; every agent aggregate leaves it out, as it leaves out an empty usage.
- The keys of `usage` are listed in `front-matter-fields.txt`, so publishing the flai release that carries `strategic` raises `flai.minimum`. A flai older than it reads past the key and writes usage back without it.
- The orchestrator and the analyzer have no item of their own yet; the entry's `kind` leaves room for them (S-0226, S-0227).
- `design/system/metrics.md` and `design/system/work-hierarchy.md` define the entry and its figures.

## Alternatives considered

- **Add the planner's models to the item's `models`.** Every per-model figure would mix planning with delivery, which is what the designer asked to keep readable.
- **Sum `strategic` from the children at roll-up.** A parent's own planning cannot be told from its children's once summed, so each roll-up would lose it, unless every item carried its own planning in a second block.
- **Keep the charge on the planned item only and sum at read time.** Every reader would need the hierarchy to total an epic, and the item's usage would not say what its children's planning cost, as its agent figures do.
- **Price the expected cost at `planning.hour_rate`.** That rate is what an hour of a person's work costs, for time lost; an agent hour costs what the logs say.
