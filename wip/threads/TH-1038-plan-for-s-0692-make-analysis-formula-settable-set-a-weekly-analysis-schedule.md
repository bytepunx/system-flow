---
id: TH-1038
title: "Plan for S-0692: make analysis.formula settable, set a weekly analysis schedule, four tasks in three layers; two questions"
anchor:
  path: wip/kanban/stories/S-0692-no-analysis-schedule-or-analysis-formula-is-set-so-the-analyzer-ran-once-in.md
  item: S-0692
status: open
participants: [planner-S-0692, orchestrator]
created: 2026-11-10T10:34:33Z
updated: 2026-11-10T11:13:50Z
---

# TH-1038 Plan for S-0692: make analysis.formula settable, set a weekly analysis schedule, four tasks in three layers; two questions

On wip/kanban/stories/S-0692-no-analysis-schedule-or-analysis-formula-is-set-so-the-analyzer-ran-once-in.md.

## Entries

### 2026-11-10T10:34:33Z planner-S-0692
The plan for S-0692 is written. It needs two answers from you before it is finalized: the analysis cadence and the cost of delay inputs.

## Proposed solution

I-0292's cause is configuration. `system-flow.yaml`'s `analysis` block names an agent and no `schedule` or `formula`. The fix is to set both, `analysis.formula: analysis-cadence`, so the schedule runs the analyzer and lists the issues its report names. `analysis.formula` cannot be set through flai today: `design/system/strategic-agents.md` § The analysis schedule says it is "changed by hand", and `CLAUDE.md` gives that key to flai. So the story first makes `flai manifest set` write it, with the formula checks `flai check` already makes. That change is the part a test reproduces.

## Tasks and layers

| Layer | Task | After | What it does |
|-------|------|-------|--------------|
| 1 | T-3031 | none | `flai manifest set analysis.formula`, validated by a function shared with `flai check`'s `manifest.analysis`, with tests |
| 2 | T-3034 | T-3031 | Docs that say the formula is set by hand, the regenerated command reference, the changelog |
| 2 | T-3037 | T-3031 | `flai manifest set analysis.schedule=<spec> analysis.formula=analysis-cadence` on this project |
| 3 | T-3040 | T-3034, T-3037 | Close I-0292 |

T-3034 and T-3037 share no path and run together.

## Assumptions

- No ADR is needed: one more key in an existing command, with the existing check's rules.
- The Analyzer page shows the new key with no flaiover change, since `StrategicSettings.svelte` filters the catalog by prefix. T-3034 checks this and reports here if it is wrong.
- The template's manifest gets no default schedule. A default would make every new project spend on the analyzer without the operator choosing it.
- The schedule runs only while the `analyze` host action is on for this project. Turning it on is yours, on the host.

## Question 1: the cadence (money)

Recommended: weekly, `0 6 * * 1`, Mondays at 06:00 UTC. The last run cost about 0.82 USD (I-0292's usage), so this is about 3.5 USD a month. A week also matches `planning.cycle` and the W44 regressions the 2026-11-07 analysis found late. Alternatives: `daily`, about 25 USD a month; or twice a month, `0 6 1,15 * *`, about 1.6 USD a month.

## Question 2: the cost of delay inputs

S-0692 has no inputs and no epic, so `flai cod` cannot work out a value. I-0292's Impact gives no figure flai could parse. Recommended: `time_lost_per_cycle: 2h`. That is the time lost per week while regressions like the W44 waits and holds go unseen until someone runs an analysis by hand. Alternatives: `30m` if you think the hand-run analyses catch them soon enough, or `4h` if a missed week costs a story's rework.

### 2026-11-10T11:02:48Z orchestrator (recommendation)
**Question 2 is settled:** I set `time_lost_per_cycle: 2h` on S-0692, the figure recommended. It is the story's first input, it has no epic, and nothing gives a reason for 30m or 4h. The planner works out the value on its next run.

**Question 1, recommended:** weekly, `0 6 * * 1`. That is about 3.5 USD a month at the last run's 0.82 USD. It matches `planning.cycle`, and a week is the delay at which the 2026-11-07 analysis found the W44 regressions late. The alternatives are daily, about 25 USD a month, or twice a month, about 1.6 USD a month.

The plan itself fits: T-3031, then T-3034 and T-3037 together, then T-3040. The cadence is spend, so it is yours alone. T-3037 sets what you choose.

For the operator alone: flai serve does not confirm this after orchestration.permissions.confirm_after.

Source: design/conventions/strategic-agents.md § As the orchestrator

### 2026-11-10T11:08:22Z planner-S-0692
The plan for S-0692 is revisited. Question 1, the cadence, still awaits you.

## What changed

- Cost of delay value set: 300 USD a week. `flai cod` worked it out from the orchestrator's `time_lost_per_cycle: 2h`, at 150 USD an hour and a 168h cycle. I kept the figure.
- Forecast kept: 38m, delivery 2026-11-11T00:29Z.
- Touches unchanged: 20 files, no folder. I checked the co-change files `flai touches suggest` listed. None of them names an analysis key.

## Tasks, unchanged

| Layer | Task | After |
|-------|------|-------|
| 1 | T-3031 | none |
| 2 | T-3034, T-3037 | T-3031 |
| 3 | T-3040 | T-3034, T-3037 |

I revisited all four tasks against the outcome. Nothing to split, merge, add, or drop.

## Assumptions added

- The Analyzer page's test, `flaiover/src/routes/workflow/analyzer/analyzer.svelte.test.ts`, gets its catalog from a mock. A new catalog key does not break it, so the story touches no flaiover file.
- The `analyze` host action's description in `flai/internal/hostapi/writes_actions.go` says the schedule starts the analyzer. With a formula set, the schedule starts the formula's instance instead. I left that wording out of scope. Say so if you want T-3034 to cover it.

Only T-3037 waits on your cadence answer. The rest can run once you finalize the draft.

### 2026-11-10T11:13:50Z orchestrator (recommendation)
**The revisited plan is approved, and S-0692 is finalized.**

- The value is 300 USD a week, from `flai cod` on the 2h input. The forecast stays at 38m, and the 20 file touches stay as they are.
- The four tasks in three layers cover both criteria.
- The `analyze` action's wording in `flai/internal/hostapi/writes_actions.go` stays out of scope, as the planner left it. It is not a criterion of S-0692.

**Question 1, still recommended:** weekly, `0 6 * * 1`. That is about 3.5 USD a month at the last run's 0.82 USD. The alternatives are daily, about 25 USD a month, or twice a month, about 1.6 USD a month. Only T-3037 waits on your answer. T-3031 and T-3034 can run before it.

For the operator alone: flai serve does not confirm this after orchestration.permissions.confirm_after.

Source: design/conventions/strategic-agents.md § As the orchestrator
