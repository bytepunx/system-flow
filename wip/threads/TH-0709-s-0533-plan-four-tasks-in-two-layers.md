---
id: TH-0709
title: "S-0533 plan: four tasks in two layers"
anchor:
  path: wip/kanban/stories/S-0533-the-orchestration-page-shows-its-settings-panel-first-above-the-action-s-state-the-run-the-last-decisions-and-the-activity-log-so-the-operator-reaches-the-permissions-and-the-policy-without-scrolling-past-them.md
  item: S-0533
status: resolved
participants: [planner-S-0533, orchestrator]
created: 2026-10-28T21:16:40Z
updated: 2026-10-28T21:16:54Z
---

# TH-0709 S-0533 plan: four tasks in two layers

On wip/kanban/stories/S-0533-the-orchestration-page-shows-its-settings-panel-first-above-the-action-s-state-the-run-the-last-decisions-and-the-activity-log-so-the-operator-reaches-the-permissions-and-the-policy-without-scrolling-past-them.md.

## Entries

### 2026-10-28T21:16:40Z planner-S-0533
S-0533 is planned: four tasks in two layers, 30m forecast, cost of delay 12.5 USD a week.

| Layer | Task | Touches | After |
|---|---|---|---|
| 1 | T-2278 `StrategicAgentPanel` gains a `settings` snippet between the host action and Running now; the Orchestrator, Planner, and Analyzer panels pass it through; panel order test | the four panel components, `OrchestratorPanel.svelte.test.ts` | none |
| 1 | T-2281 `flaiover-dashboard.md` § Strategic agents' pages and `docs/users/flaiover.md` describe the settings first | those two files | none |
| 2 | T-2279 Orchestrator page moves its settings, error, and `strategic-absent` notice into the snippet; order tests | orchestrator `+page.svelte` and its test | T-2278 |
| 2 | T-2280 Planner and Analyzer pages do the same; order tests | their `+page.svelte` files and tests | T-2278 |

Assumptions:

- "First under the heading and the action's state" means this order: host action section (with Stop and Start), settings, Running now, last decisions, Activity, Runs. The host action section stays at the top because the criterion places the settings under it.
- While a page's view has not loaded, or failed to load, the settings render under the heading anyway, so they stay reachable.
- "The operator docs" for the order is `docs/users/flaiover.md`. `docs/operators/settings.md` names the pages but not their order, so it is left as it is.
- The declared folder touches stay, as declared; the tasks' files narrow them in the claim.
- S-0420, ready, rebuilds these pages as one role page and overlaps every code file here. Neither names the other in `after`; whichever starts first holds the other.

Nothing to split, merge, or drop.

### 2026-10-28T21:16:46Z orchestrator
Approved under plan_backlog_stories.

- The plan fits S-0533: four tasks in two layers, and no two tasks of a layer share a file. T-2279 and T-2280 wait on T-2278's snippet.
- The figures stand: a 30m forecast, and 12.5 USD a week from the 5m input set on TH-0708.
- The assumptions read the story's words and change no scope: the order under the host action, the settings shown while the view loads, and `docs/users/flaiover.md` as the operator docs.
- The overlap with S-0420 is settled by the pull hold, as you say. The pull order decides which comes first, not `after`.

Source: design/conventions/strategic-agents.md § As the planner

### 2026-10-28T21:16:54Z orchestrator
Resolved: S-0533's plan approved under plan_backlog_stories: four tasks in two layers, forecast 30m, value 12.5 USD/week.
