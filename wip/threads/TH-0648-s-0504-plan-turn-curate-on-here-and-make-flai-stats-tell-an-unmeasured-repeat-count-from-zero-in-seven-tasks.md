---
id: TH-0648
title: "S-0504 plan: turn curate on here and make flai stats tell an unmeasured repeat count from zero, in seven tasks"
anchor:
  path: wip/kanban/stories/S-0504-the-curator-has-never-run-on-system-flow-so-no-thread-is-in-any-entry-s-asked-again-and-the-digest-s-repeat-measure-reads-0-whatever-was-asked.md
  item: S-0504
status: resolved
participants: [planner-S-0504, orchestrator]
created: 2026-10-26T04:11:04Z
updated: 2026-10-26T04:11:16Z
---

# TH-0648 S-0504 plan: turn curate on here and make flai stats tell an unmeasured repeat count from zero, in seven tasks

On wip/kanban/stories/S-0504-the-curator-has-never-run-on-system-flow-so-no-thread-is-in-any-entry-s-asked-again-and-the-digest-s-repeat-measure-reads-0-whatever-was-asked.md.

## Entries

### 2026-10-26T04:11:04Z planner-S-0504
Plan for S-0504 by planner-S-0504, 2026-10-26.

**The cause has two parts.** The `curate` host action was never turned on for system-flow, so no curator ran and `wip/agents/curator.md` does not exist; and `flai stats` reports `digest.weeks[].asked_again` as 0 whether no repeat was counted or no curator ever ran, because `metrics.digestSection` reads threads and entries alone. The plan remedies both. Only the operator can turn a host action on; the stats change is code, with an ADR because `metrics.md` changes only with one.

**Tasks, in layers:**

1. T-2120 (research): confirm the cause on the host and propose the remedy to the operator on a thread, recommendation first, as the story's goal asks before anything is built.
2. T-2121 (ADR and `metrics.md` § Digest: `digest.weeks[].curator_runs`, `digest.curator` with `runs` and `last_run`, `asked_again` absent in a week no curator run has reached) and T-2122 (the operator turns `curate` on, the first run is recorded), both after T-2120. No path in common.
3. T-2123 (`flai stats`: `digest.go`, `metrics.go`, `stats.go` and their tests, with the test that reproduces I-0207: threads, entries, no curator activity) and T-2124 (`docs/users/flai.md`, the analyzer's convention bullet and its template copy, `template/CHANGELOG.md`), both after T-2121. No path in common.
4. T-2125 (the Digest chart: `charts.ts`, `explain.ts`, the page's table and their tests, `flaiover-dashboard.md`, `docs/users/flaiover.md`), after T-2123.
5. T-2126 (verify with `flai stats`, close I-0207 with `flai issue close --reason`, tick the criteria), after T-2122, T-2124, T-2125.

**Figures.** Forecast 2h30m of agent time, delivery 2026-10-26T21:00Z: `flai forecast` gave 38m from size 25, raised because the size counts touches, not an ADR, a Go section with its tests, a chart with its tests, and four documents; the wait for the operator's enabling and the first run is wall-clock, not agent time. Cost of delay 50 USD a week, as `flai cod` gives it from the issue's input; it stands. Touches: 23, including the folder `design/adrs` for the ADR's file, which `flai adr new` names; it is a shared path and holds nothing. `after`: none.

**Assumptions.**

- The operator agrees to the code half. The alternative, turning the action on alone, leaves every project's measure reading 0 until its curator runs and the analyzer and the dashboard reading that 0 as a result; T-2120 puts both to the operator and the story's agent changes the tasks if the answer differs.
- `asked_again` is left out of a week no curator run has reached, as `per_story` is left out of a week with no completed story, rather than reported with a flag. The ADR may choose otherwise; T-2123 to T-2125 follow it.
- The first curator run needs an acceptance, `digest.schedule`, or the operator's `flai curate`; T-2122 asks for `flai curate` on T-2120's thread when none is due. An agent `flai serve` started cannot start it.
- Not this story: the re-run of S-0435's measurement once W44 and W45 are whole, from 2026-11-09. It needs two whole weeks after the curator's first run, so I propose it as an experiment story of its own, drafted when the operator wants it, not now.

Nothing to split, merge, or drop.

### 2026-10-26T04:11:12Z orchestrator
Approved. The plan fits S-0504.

- **The tasks fit.** There are seven tasks in five layers. T-2120 proposes the remedy to the operator before anything is built, as the goal asks. Tasks in the same layer share no path: T-2121 with T-2122, and T-2123 with T-2124. T-2123 carries the test that reproduces I-0207, and T-2126 closes the issue.
- **The touches fit.** There are 23. The one folder touch, `design/adrs`, is kept only for the ADR file, is named under `### Planning`, and lies in `claims.shared`. The ADR follows the rule that `metrics.md` changes only with an ADR.
- **The figures fit.** The forecast is 2h30m of agent time, raised from flai's 38m with the reason stated. The operator's enabling of `curate` is kept out of agent time. The value is 50 USD/week from `flai cod` on the issue's input, and I set no input.
- **The assumptions stand.** Enabling `curate` and the first `flai curate` are the operator's, and T-2122 asks for them on T-2120's thread. The re-run of S-0435's measurement after 2026-11-09 stays out of this story, to be drafted when the operator asks for it.

Source: design/adrs/0119-with-plan-backlog-stories-on-the-orchestrator-asks-the-planner-for-each-backlog.md

### 2026-10-26T04:11:16Z orchestrator
Resolved: Plan approved under plan_backlog_stories: the tasks, touches, forecast, and value fit S-0504.
