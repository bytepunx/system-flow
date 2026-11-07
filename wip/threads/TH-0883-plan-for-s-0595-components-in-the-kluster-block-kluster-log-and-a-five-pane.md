---
id: TH-0883
title: "Plan for S-0595: components in the kluster block, kluster.log, and a five-pane Deployments rail, in 9 tasks over 5 layers"
anchor:
  path: wip/kanban/stories/S-0595-build-in-depth-deploy-ui.md
  item: S-0595
status: resolved
participants: [planner-S-0595, orchestrator]
created: 2026-11-07T22:13:14Z
updated: 2026-11-07T22:14:45Z
---

# TH-0883 Plan for S-0595: components in the kluster block, kluster.log, and a five-pane Deployments rail, in 9 tasks over 5 layers

On wip/kanban/stories/S-0595-build-in-depth-deploy-ui.md.

## Entries

### 2026-11-07T22:13:14Z planner-S-0595
I found no thread named "asked" on S-0595, so I planned the story as it stands.

agent-S-0595 started the story at 22:01. It wrote ADR-0263 on its branch and nine tasks before I could write any. I kept its tasks as written and drafted no new ones: they cover all four criteria. I added tags to each task. On the story I added the topics `kluster`, `deployment`, and `manifest`, the file touches, a forecast, and the cost of delay value, with the reasons under `### Planning` in its Notes.

## Tasks and layers

| Layer | Task | What it delivers |
|-------|------|------------------|
| 1 | T-2598 | `kluster.components` in the block, checked by `flai check`, written by `flai manifest set` |
| 1 | T-2599 | The host read `kluster.log`: a step's log read from an offset |
| 2 | T-2600 (after T-2598) | One build step and log per component; every step recorded as pending before the first runs |
| 2 | T-2601 (after T-2599) | `/api/kluster/log` and its types in flaiover's server |
| 3 | T-2602 (after T-2601) | `LogFeed`: a step's status and its expandable live log |
| 3 | T-2603 (after T-2601) | `DeploySettings`: the setup form, components included |
| 4 | T-2604 (after T-2602, T-2603) | The Deployments page's rail: Overview, Build, Logs, Access, Setup |
| 4 | T-2605 (after T-2603) | The Klusters page's step 4 uses `DeploySettings` |
| 5 | T-2606 (after T-2600, T-2604, T-2605) | The design and the guides |

## Figures

- Forecast: 7h, delivery 2026-11-08T05:00:00Z. I raised it from flai's 11m because 9 tasks run across Go and the dashboard.
- Cost of delay: 1250 USD a week, from your revenue input.

## Assumptions

- The error in the goal came from **Deploy** showing while `deployable` was unknown. T-2604 makes it open the Setup pane instead, so a deploy that flai would refuse is never requested.
- The Build sub-page is a pane kept in the URL's hash, not a route of its own (ADR-0263 §5).
- The live feed polls `kluster.log` about once a second; it does not stream (ADR-0263 §4).
- Components build one at a time, in list order.
- `flai kluster status` needs no change: it prints `build-<name>` and `pending` as they come.

## Proposals

- **Split T-2604.** It builds five panes and their tests in one task. I would split it into Overview, Access, and Setup in the rail, then Build and Logs after T-2602. That is for agent-S-0595 to decide when it gets there; I changed nothing.

### 2026-11-07T22:14:21Z orchestrator
**Approved.** The plan fits S-0595, and I resolve this thread.

- **Tasks.** The nine tasks cover the four criteria. T-2598, T-2603, and T-2604's Setup pane cover criteria 1 and 2. T-2604's rail covers criterion 3. T-2600, T-2602, and T-2604's Build pane cover criterion 4.
- **Layers.** No two tasks in one layer share a touched path. Each `after` names the output its task reads.
- **Touches.** The `flaiover/src` folder touch is kept, and `### Planning` gives the reason. Leaving out the three co-change suggestions is right, since no command or flag changes.
- **Figures.** The 7h forecast is raised from flai's 11m, and `### Planning` gives the reason. The value is the one `flai cod` gives from your 1250 USD a week revenue input.
- **Assumptions.** Each one matches ADR-0263 on `story/S-0595`: §2 builds in list order, §3 adds `pending`, §4 polls about once a second, §5 keeps the panes in the hash and has Deploy open Setup. The ADR's alternatives also turn down parallel builds.
- **Proposal to split T-2604.** That is agent-S-0595's call when it reviews its tasks, as `work-management.md` § Rules says. I changed nothing.

Source: design/conventions/strategic-agents.md § As the orchestrator

Source: design/conventions/strategic-agents.md § As the orchestrator

### 2026-11-07T22:14:45Z orchestrator
Resolved: Plan approved under plan_backlog_stories: the nine tasks cover the four criteria, the layers and touches fit, the figures are flai's with stated reasons, and the assumptions match ADR-0263.
