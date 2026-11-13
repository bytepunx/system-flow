---
id: TH-1132
title: "Plan for E-0034: 19 draft stories, 59 tasks, an ADR first, then flai package by package, the timers, the dashboard, and an end-to-end test"
anchor:
  path: wip/kanban/epics/E-0034-an-operator-defines-the-board-s-states-transitions-rules-and-timers-in-a.md
  item: E-0034
status: open
participants: [planner-E-0034, orchestrator]
created: 2026-11-12T16:18:06Z
updated: 2026-11-12T16:20:29Z
---

# TH-1132 Plan for E-0034: 19 draft stories, 59 tasks, an ADR first, then flai package by package, the timers, the dashboard, and an end-to-end test

On wip/kanban/epics/E-0034-an-operator-defines-the-board-s-states-transitions-rules-and-timers-in-a.md.

## Entries

### 2026-11-12T16:18:06Z planner-E-0034
I drafted 19 stories under E-0034, each with tasks, touches, and a forecast. Cost of delay values wait on TH-1097. Every story is a draft in the backlog.

## Stories, in pull order, with what each `after` reads

1. S-0736: the ADR that fixes the workflow file's format, kinds, effects, rules, timers, and migration, plus the design. No `after`. Tasks: T-3187, then T-3193 (asks you, then accepts), then T-3200. 45m.
2. S-0737: the `flai/internal/workflow` package parses the file, flai embeds the baseline, and `workitem` reads its states and transitions from it. After S-0736 (the format). Tasks: T-3188, then T-3195 and T-3196, then T-3201. 33m.
3. S-0738: `Repo.Move`'s rules become named rules on the baseline's transitions. After S-0737 (the parser and baseline). Tasks: T-3189, then T-3197 and T-3199, then T-3202. 29m.
4. S-0739: `workitem`'s epic following, board view, task plan, holds, order, and promotion read kinds. After S-0737 (kinds, WIP limits, test hook). Tasks: T-3203, then T-3206, T-3208, and T-3210. 45m.
5. S-0741: `flai check` validates the workflow and the items against it, and its item rules read kinds. After S-0738 (the rule vocabulary). Tasks: T-3204 and T-3205, then T-3211. 29m.
6. S-0742: flai serve's agents, roles, ready-empty, replan, and reminders read kinds and the pull effect. After S-0737. Tasks: T-3213, T-3215, T-3216, and T-3217 in one layer. 1h20m.
7. S-0743: serve's kluster, forge, and merge queue read the acceptance and send-back effects. After S-0737. Tasks: T-3218, T-3219, and T-3220 in one layer. 45m.
8. S-0744: `flai move`, `flai accept`, formula stop, and the acceptance preview. After S-0737 (who may make a transition). Tasks: T-3221, T-3222, and T-3223, then T-3227. 1h.
9. S-0746: `item_move`, `story_start`, task done, sync, item edits, and the inbox. After S-0737. Tasks: T-3224, T-3225, and T-3226 in one layer. 42m.
10. S-0747: the host API's writes and reads. After S-0737. Tasks: T-3228 and T-3229, then T-3230. 38m.
11. S-0748: `to-review` and `to-backlog` become one condition naming any state, and the baseline role files use it. After S-0737 (states to check against). Tasks: T-3231, then T-3238, T-3239, and T-3240. 53m.
12. S-0749: the guard, messages, forecasts, packs, release evaluation, digest, and turns. After S-0737. Tasks: T-3232, T-3233, T-3234, and T-3235 in one layer. 42m.
13. S-0750: the metrics ADR, and `flai stats` computes the baseline's figures from timers, byte for byte. After S-0737 (`Timer`, the baseline). Tasks: T-3236 and T-3237, then T-3241 and T-3242, then T-3246. 1h.
14. S-0752: `flai stats` and `stats.get` report a workflow's own timers. After S-0750 (timer evaluation, section name). Tasks: T-3243, then T-3247 and T-3248. 33m.
15. S-0753: the dashboard's board draws columns and moves from `board.get`. After S-0739 (the board view's columns and kinds). Tasks: T-3244, then T-3249 and T-3250. 42m.
16. S-0754: the charts stack the workflow's states and chart each timer. After S-0752 (timers in `stats.get`) and S-0753 (columns in `board.get`). Tasks: T-3245, then T-3251, then T-3254. 30m.
17. S-0755: a project's own workflow file is read in place of the baseline; the template ships it, `flai upgrade` keeps it, and a runbook says how to change the board. After S-0741 (the `workflow.` findings that guard a project's file). Tasks: T-3253, then T-3259 and T-3260, then T-3262. 47m.
18. S-0758: the dashboard's other views (item page, inbox, review, cards, deployments, and others) read kinds. After S-0753. Tasks: T-3252, then T-3255, T-3256, T-3257, and T-3258. 44m.
19. S-0756: the closing end-to-end test in a project made from the template. After every story whose surface it drives (its Planning names what it reads from each). Tasks: T-3261, then T-3263, then T-3264. 45m.

## Assumptions

- The file is `design/workflows/board.md`, and flai's copy is `flai/internal/workflow/baseline/board.md`. Where S-0736's ADR names otherwise, each story's agent renames the touches.
- The state kinds and effects are named as the epic names them (queue, active work, wait, accepted, cancelled; pull, acceptance, send back). The ADR may rename them.
- Who may make a transition is a field of the transition. S-0737 adds the check, which S-0744, S-0746, S-0747, and S-0749 all use, so none of them waits on another.
- Stories 4 to 12 change no behaviour on the baseline. Each is tested on a workflow with a state renamed and one added between `in-progress` and `review`, through S-0737's test hook, before S-0755 reads a project's own file.
- `flai check` advises `wip.after-hub` on S-0737 (10 stories) and on S-0753 (3). Each of those stories reads that story's output, so I kept them.

## Proposals

- **Added: S-0758.** The epic's order named the board and the charts. A search found about fifteen more flaiover views that compare a status with a state's name, so I drafted this story for them.
- **Split S-0742 if you prefer smaller stories.** At 1h20m and 38 files it is the largest. Its four tasks share no path, so it splits cleanly into the agents' starts (T-3213, T-3215) and the strategic roles (T-3216, T-3217).
- **Merge S-0752 into S-0750 if you want fewer stories.** I kept them apart: S-0750's change must leave every figure unchanged byte for byte, and S-0752 adds new output.
- I would drop none.

## Observed

`flai forecast S-0756` gave a delivery before S-0758's, though S-0756 waits for S-0758. I wrote S-0758's delivery plus 45m instead.

### 2026-11-12T16:18:13Z planner-E-0034
Correction to the title: the plan has 69 tasks, not 59.

### 2026-11-12T16:20:29Z orchestrator (recommendation)
**Recommended:** approve the plan as drafted: 19 stories, S-0736 to S-0758, 69 tasks. Keep S-0758. Keep S-0742 whole. Keep S-0752 apart from S-0750.

**Why each:**

- **Keep S-0758.** E-0034's Outcome is done when an added state is followed by the board and the charts "with no change to flai or flaiover". Any flaiover view that still compares a status with a state's name breaks that. So the fifteen views the planner found are in the epic's scope.
- **Keep S-0742 whole.** `work-management.md` sizes a story at one to a few sessions. 1h20m fits that. Its four tasks share no path, so they run together in one layer. Its agent can still split it at pull time if it grows.
- **Keep S-0752 apart.** S-0750 must leave every figure byte for byte as today, which the Notes and `metrics.md` require through an ADR. S-0752 adds new output. Kept apart, a regression in either shows up on its own.
- **The `after` hub on S-0737** (10 stories) is justified. Each of those stories reads its parser, kinds, or transition check on main, as "As the planner" in `strategic-agents.md` requires of an `after`.

**What approval does not do:** the drafts stay drafts. Each still needs a cost of delay value before it can be finalized or promoted. That waits on your answer to TH-1097, the epic's inputs, which is yours alone.

Source: design/conventions/strategic-agents.md
