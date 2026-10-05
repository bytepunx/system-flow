---
id: TH-0144
title: "S-0228 plan: five tasks in four layers, and an orchestrator Stop that holds it stopped"
anchor:
  path: wip/kanban/stories/S-0228-the-workflow-menu-has-orchestrator-and-analyzer-pages-showing-their-status-activity-log-and-runs.md
  item: S-0228
status: open
participants: [planner-S-0228]
created: 2026-10-05T05:46:00Z
updated: 2026-10-05T05:46:00Z
---

# TH-0144 S-0228 plan: five tasks in four layers, and an orchestrator Stop that holds it stopped

On wip/kanban/stories/S-0228-the-workflow-menu-has-orchestrator-and-analyzer-pages-showing-their-status-activity-log-and-runs.md.

## Entries

### 2026-10-05T05:46:00Z planner-S-0228
Plan for S-0228, by planner-S-0228.

Tasks and layers:

- Layer 1, run together, sharing no path:
  - T-0917: flai serve reads an orchestrator or analyzer run's stream by `role`, and adds the writes `orchestrate.stop` and `orchestrate.start`.
  - T-0921: the dashboard API `/api/orchestrator` and `/api/analyzer` (GET for state, activity, and runs; POST for stop, start, and run with focus), the stream route by role, and a shared `strategic.ts` taken from `planner.ts`.
- Layer 2: T-0928, the Orchestrator page, after T-0921. It takes the shared panel out of `PlannerPanel.svelte` and adds the menu entry.
- Layer 3: T-0931, the Analyzer page, after T-0928, because it builds on that shared panel and also edits `sitemenu.ts`.
- Layer 4: T-0938, the design, user guide, and flai reference, after all four.

Assumptions. Correct me on any of them:

1. Stop holds the orchestrator stopped. S-0218 has flai serve restart the orchestrator whenever a run ends and `orchestrate` is still on, so a plain stop would be undone at once. I planned for Stop to end the run and hold it, and for a Start button (the goal's "start and stop") to lift the hold. The action itself stays on, because host actions are turned on only on the host. Turning the action off and on also lifts the hold. If you would rather have Stop do nothing more than turn `orchestrate` off, T-0917 shrinks to the stream read.
2. The orchestrator's "last decisions with their reasons" are the newest entries of `wip/agents/orchestrator.md`, shown above the full log. S-0218 has every decision and refusal logged there with its reason.
3. The analyzer's Run calls S-0223's `analyze.run` with focus `bottlenecks`, `intent`, `risk`, or none. Its runs are typed as S-0223 records them in `agent.status`. T-0917 and T-0921 take their shapes from S-0218 and S-0223 as those are built.
4. The story's own Notes are unchanged apart from the `### Planning` section. Topics `orchestration` and `analysis` were added to the story.

Figures: forecast 1h, against flai's 14m on the declared touches and 1h4m on the predicted ones. The reason: S-0259 was one such page and took 29m, and this story is two pages plus the stop and start. Cost of delay is 60.81 USD a week, its share of E-0016's value. The reasons are under `### Planning` in the story's Notes.

No split, merge, or drop proposed.
