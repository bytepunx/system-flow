---
id: TH-1075
title: "S-0720 plan: six tasks in four layers, 1h30m, 375 USD a week"
anchor:
  path: wip/kanban/stories/S-0720-a-usage-limit-on-one-claude-login-holds-every-start-on-claude-code-and-changing.md
  item: S-0720
status: resolved
participants: [planner-S-0720, orchestrator]
created: 2026-11-11T08:56:30Z
updated: 2026-11-11T09:00:21Z
---

# TH-1075 S-0720 plan: six tasks in four layers, 1h30m, 375 USD a week

On wip/kanban/stories/S-0720-a-usage-limit-on-one-claude-login-holds-every-start-on-claude-code-and-changing.md.

## Entries

### 2026-11-11T08:56:30Z planner-S-0720
S-0720 is planned: 21 file touches, six tasks in four layers, a forecast of 1h30m, and a cost of delay of 375 USD a week.

## Tasks and layers

| Layer | Task | What | After |
|-------|------|------|-------|
| 1 | T-3095 | Key a budget state by provider, else harness with its resolved program and login; read the old shape | none |
| 1 | T-3096 | The dashboard's banner names `flai serve budget clear <provider>` | none |
| 2 | T-3098 | The clear: journaled under `serve.budget` with why, sent as a notice; `budgetHeldWhy` names the command | T-3095 |
| 3 | T-3104 | `flai serve budget` and `flai serve budget clear [<provider>]`, with the reference regenerated | T-3098 |
| 3 | T-3105 | `flai serve agent harness` clears the old identity's state when program, arguments, or login change | T-3098 |
| 4 | T-3106 | `design/system/flai-cli.md`, the operator and user docs, the dashboard design | T-3104, T-3105, T-3096 |

The tasks of each layer share no path.

## Figures

- Forecast 1h30m, against flai's 46m. I raised it for four sequential layers over the budget file's lock and probe code, and a close-out that runs the integration and smoke tiers.
- Cost of delay 375 USD a week, from `time_lost_per_cycle` 2h30m (TH-1070).

## Assumptions

- The dashboard's harness setting runs `flai serve agent harness` as a subprocess (`flai/internal/hostapi/settings.go`). T-3105 covers it, and nothing under `flai/internal/hostapi/` is touched.
- The resolved program and login are recorded on `AgentRun` at start. So `budgetEnded`'s six callers need no change.
- T-3095 decides whether an old-shape entry keyed `claude-code` holds every program on the harness until it clears, and records the choice.
- I added the tags `flai` and `flaiover`. `story.component-tag` asks for them now that the story touches both.

## Proposals

None. I would not split, merge, add, or drop a story.

### 2026-11-11T08:59:55Z orchestrator
Approved: the plan of S-0720 fits the story.

- Each criterion has a task. Criterion 1 is T-3105. Criterion 2 is T-3095. Criterion 3 is T-3098, T-3104, and T-3096. Criterion 4 is the tests in each task, plus T-3106.
- The layers hold. T-3095 and T-3096 share no path. T-3098 waits for T-3095 on `budget.go`. T-3104 and T-3105 share no path and call T-3098's clear. T-3106 describes what was built.
- The touches are 21 files and no folder. Leaving `flai/internal/hostapi/` out is sound, because `settings.harness` runs `flai serve agent harness`.
- The figures stand: 1h30m over flai's 46m, with the reason stated, and 375 USD a week from the 2h30m input set on TH-1070.

No input is asked for and no scope changes, so I resolve the thread.

Source: design/adrs/0119-with-plan-backlog-stories-on-the-orchestrator-asks-the-planner-for-each-backlog.md

### 2026-11-11T09:00:21Z orchestrator
Resolved: Plan of S-0720 approved by the orchestrator under plan_backlog_stories: tasks, touches, and figures fit the story.
