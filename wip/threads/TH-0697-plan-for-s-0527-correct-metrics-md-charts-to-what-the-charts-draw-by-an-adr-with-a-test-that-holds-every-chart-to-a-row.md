---
id: TH-0697
title: "Plan for S-0527: correct metrics.md § Charts to what the charts draw, by an ADR, with a test that holds every chart to a row"
anchor:
  path: wip/kanban/stories/S-0527-metrics-md-charts-says-burn-up-draws-a-forecast-line-from-throughput-and-cycle-time-filters-by-nature-but-neither-chart-does.md
  item: S-0527
status: open
participants: [planner-S-0527, orchestrator]
created: 2026-10-28T15:02:12Z
updated: 2026-10-28T15:11:48Z
---

# TH-0697 Plan for S-0527: correct metrics.md § Charts to what the charts draw, by an ADR, with a test that holds every chart to a row

On wip/kanban/stories/S-0527-metrics-md-charts-says-burn-up-draws-a-forecast-line-from-throughput-and-cycle-time-filters-by-nature-but-neither-chart-does.md.

## Entries

### 2026-10-28T15:02:12Z planner-S-0527
The plan, by planner-S-0527 on 2026-10-28. The story's Planning note records the touches, the forecast, and the cost of delay.

**The solution proposed from I-0193's instance.** The charts in `flaiover/src/lib/viz/charts.ts`, their explainers (S-0485), `docs/users/flaiover.md`, and `design/system/flaiover-dashboard.md` all agree with each other; only `design/system/metrics.md` § Charts differs, and no ADR decided a Burn-up forecast line or a Cycle Time nature filter. Recommended: correct the contract, not the charts. An ADR refining ADR-0054 and ADR-0057 settles that § Charts says what the dashboard draws: Cycle Time chosen per epic and coloured by nature, no nature filter; Burn-up scope and done, no forecast line; the dashed all-models line on every chart per model, Tokens / Min and Tokens / $ included; rows for Agent Waiting, Strategic Cost, and Strategic Use, which the table lacks, Strategic Use's naming the mean waiting as the dashboard's sum of `wait_threads_seconds` and `wait_review_seconds` (TH-0140); and a row for every chart kind, held by a test. The alternative, building the forecast line from throughput and the nature filter, is a feature no story asked for; if the operator wants it, it is a feature story of its own, and this remediation still corrects the rows that do not match. Say so on this thread before T-2255 is worked if you choose otherwise.

**Tasks, in four serial layers.**

1. T-2255: the ADR, with `flai adr new` and `flai adr accept`. Waits for nothing. Touches `design/adrs/README.md`; the ADR's own path is widened in at `flai task done`.
2. T-2256, after T-2255: `design/system/metrics.md` § Charts rewritten to what each chart draws, the three missing rows added, and § Strategic use per day saying the mean waiting is the dashboard's. Touches `design/system/metrics.md`.
3. T-2257, after T-2256: a test in `flaiover/src/lib/viz/explain.test.ts` that every kind in the catalog has a § Charts row under its title. It reproduces I-0193: on main it fails on Agent Waiting, Strategic Cost, and Strategic Use. It waits for T-2256 because `flai task done` runs the tier and the test fails until the rows exist.
4. T-2258, after T-2256 and T-2257: `flai issue close I-0193 --reason` in the worktree, and the criteria ticked.

**Assumptions.**

- A correction of § Charts to what the dashboard already draws still needs an ADR, as ADR-0054 and `CLAUDE.md` say metrics.md changes only with one; the ADR is cheap and records that the forecast line was never built.
- The test that reproduces the issue is the row-per-kind test: the wording of a row cannot be tested against the chart, but its presence can, and the three missing rows are the part of the issue a test can hold.
- `docs/users/flaiover.md` and `design/system/flaiover-dashboard.md` are not touched: they already describe the charts as drawn, and no chart changes. `flai stats` changes nothing, so no Go file is touched.
- The folder touch `design/adrs` stays on the story for the ADR's unallotted number, as the Planning note records; the folder is in the manifest's shared paths and holds nothing.
- Forecast raised from flai's 16m to 40m for four serial layers; cost of delay 25.00 USD a week stands from the input flai set from I-0193.
- S-0502, in progress, also touches `design/system/metrics.md`: the story is held until S-0502 reaches review, and the pull order settles the rest. No `after`.

### 2026-10-28T15:11:48Z orchestrator (recommendation)
**Recommended for the operator:** approve the plan as written. Correct `design/system/metrics.md` § Charts to what the dashboard draws, by T-2255's ADR, and do not build the Burn-up forecast line or the Cycle Time nature filter in this story.

## The plan fits

- **Tasks:** four serial layers. T-2255 (ADR), then T-2256 (§ Charts rewritten, three rows added), then T-2257 (a row-per-kind test that reproduces I-0193), then T-2258 (closes I-0193).
- **The ADR is required.** The decisions convention's project additions say `metrics.md` is the contract between flai and flaiover, changed only with an ADR. The plan keeps that.
- **Figures:** forecast 40m over flai's 16m with the reason stated, value 25 USD a week from the issue's input.
- **No `after`:** right. The overlap with S-0502 on `metrics.md` holds it until S-0502 reaches review.

## Why it is yours

Two features the contract describes would be taken out of it rather than built: a Burn-up forecast line from throughput, and a Cycle Time nature filter. That is a change of scope, so it is your call, not mine.

## Why I recommend correcting the contract

- No ADR ever decided those two features. The charts, their explainers, `docs/users/flaiover.md` and `design/system/flaiover-dashboard.md` already agree with each other. Only § Charts differs.
- The story is a remediation, and building two chart features is new work no story has asked for.
- The new test holds every chart to a § Charts row from now on, so the contract cannot drift from the charts again.

**The alternative:** keep the two rows and build them as a feature story of their own. This story still adds the three missing rows either way.

The thread stays open for your answer.

Source: design/conventions/decisions.md § Project additions
