---
id: ADR-0219
title: "metrics.md § Charts says what the dashboard draws, with a row under its catalogue title for every chart, and Burn-up has no forecast line nor Cycle Time a nature filter"
status: accepted
date: 2026-10-30
supersedes: []
superseded_by: []
refines: [ADR-0054, ADR-0057]
---

# ADR-0219 metrics.md § Charts says what the dashboard draws, with a row under its catalogue title for every chart, and Burn-up has no forecast line nor Cycle Time a nature filter

## Context

[ADR-0054](0054-every-chart-spans-the-window-chosen-its-time-axis-runs-from-the-window-s-start.md) makes `design/system/metrics.md` § Charts the contract for what each chart plots, changed only with an ADR. [ADR-0057](0057-the-dashboard-charts-agent-time-and-cost-per-item-per-model-over-time-instead.md) decided the dashed line for all models on Avg. Time / Model and Avg. Cost / Model.

Writing the chart explainers in S-0485 found § Charts and the charts apart ([I-0193](../issues/I-0193-metrics-md-charts-says-burn-up-draws-a-forecast-line-from-throughput-and-cycle-time-filters-by-nature-but-neither-chart-does.md)):

| § Charts says | The dashboard draws |
|---|---|
| Burn-up: a forecast line from throughput | Scope and done only |
| Cycle Time: filter by nature | Chosen per epic or whole, each point coloured by nature |
| Nothing on Tokens / Min and Tokens / $'s dashed line | A dashed line for all models with two models or more, as on the ADR-0057 charts |
| No row for Agent Waiting, Strategic Cost, or Strategic Use | All three |

No ADR decided the forecast line or the nature filter, and no story asked for them. The charts, their explainers, `docs/users/flaiover.md`, and `design/system/flaiover-dashboard.md` agree with each other; only § Charts differs. Nothing kept § Charts in step with the charts the catalogue lists.

## Decision

**§ Charts says what the dashboard draws, with a row under its catalogue title for every chart, and Burn-up has no forecast line nor Cycle Time a nature filter.** This refines ADR-0054 and ADR-0057; the rest of each stands.

1. **Burn-up** draws scope and done, per epic or the whole repository. It has no forecast line. A forecast from throughput would be a feature story of its own.
2. **Cycle Time** is chosen per epic or whole and colours each point by its nature. It has no nature filter.
3. **The all-models line.** Every chart per model drawn from the spend buckets, Tokens / Min, Tokens / $, Avg. Time / Model, and Avg. Cost / Model, adds a dashed line for all models together when two or more models worked in the window. ADR-0057's line on its two charts is the same line.
4. **Every chart has a row.** Each chart kind in `flaiover/src/lib/viz/catalog.ts` has a row in § Charts whose first cell is its title there. Agent Waiting, Strategic Cost, and Strategic Use gain theirs. A chart titled by its bucket, Tokens / Day and $ / Day, is named by its day title. Strategic Use's mean waiting is the dashboard's: the mean over the items completed that day of `wait_threads_seconds` plus `wait_review_seconds`, which `flai stats` does not send as a field of its own.
5. **A test holds it.** `flaiover/src/lib/viz/explain.test.ts` fails when a kind has no row under its title, so a story that adds a chart adds its row in the same change.

## Consequences

- § Charts describes no chart feature the dashboard lacks, and lacks no chart it draws.
- A chart's title change in the catalogue fails the test until § Charts follows, under this ADR.
- The test checks a row's presence, not its wording: a row that says the wrong thing is still found by review.
