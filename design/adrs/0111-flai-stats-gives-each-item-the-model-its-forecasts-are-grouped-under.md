---
id: ADR-0111
title: flai stats gives each item the model its forecasts are grouped under
status: accepted
date: 2026-10-07
supersedes: []
superseded_by: []
refines: [ADR-0081]
---

# ADR-0111 flai stats gives each item the model its forecasts are grouped under

## Context

S-0212 charts forecast error per story agent model and filters the forecast and delivery error by model. `flai stats --json` groups the errors by model under `forecasts.*.by_model`, from each story's `agent.model`, `(none)` without one ([ADR-0081](0081-flai-stats-reports-forecast-error-cost-of-delay-waiting-holds-touches-drift-and.md)). Its `items` carry the errors but not the model, so the dashboard cannot place a story under a model, nor work out a model's error per bucket, the same way flai does. `usage.models` names every model that spent on an item, sub-agents included, which is not the story agent's model. `design/system/metrics.md` defines each value `flai stats` reports and changes only with an ADR.

## Decision

Each item in `flai stats --json` carries `model`: its agent's `agent.model`, or `(none)` when it has no agent or its agent names no model. It is the model `forecasts.*.by_model` groups the item under.

## Consequences

- The dashboard filters the planning charts by model and works out a model's percentiles per bucket from `items` alone, matching `forecasts` to the second.
- Every item grows by one short field. A dashboard reading a flai older than this treats every item as `(none)`.

## Alternatives considered

- Taking the model with the most spend in `usage.models`: it counts sub-agents' models and items without usage, and would disagree with `by_model`.
- Reading `agent.model` from the work items in the dashboard: the charts read `/api/stats` alone, and a second source could disagree with the report's window.
- Setting `model` only on items with a forecast or an estimate: the field would mean different things on different items.
