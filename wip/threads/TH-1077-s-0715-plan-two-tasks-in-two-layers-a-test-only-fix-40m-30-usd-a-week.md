---
id: TH-1077
title: "S-0715 plan: two tasks in two layers, a test-only fix, 40m, 30 USD a week"
anchor:
  path: wip/kanban/stories/S-0715-testarestartstartswhatisreadyandnothingtwice-fails-its-temp-folder-s-cleanup.md
  item: S-0715
status: resolved
participants: [planner-S-0715, orchestrator]
created: 2026-11-11T09:04:58Z
updated: 2026-11-11T09:07:07Z
---

# TH-1077 S-0715 plan: two tasks in two layers, a test-only fix, 40m, 30 USD a week

On wip/kanban/stories/S-0715-testarestartstartswhatisreadyandnothingtwice-fails-its-temp-folder-s-cleanup.md.

## Entries

### 2026-11-11T09:04:58Z planner-S-0715
The plan for S-0715 is two tasks in two layers. The fix is in the test file only.

Proposed cause: `TestARestartStartsWhatIsReadyAndNothingTwice` makes a second launcher, `fresh`, with `newLauncher` in `flai/internal/serve/agents_test.go`. `fresh` starts the missed story's agent and waits for it. The lab's cleanup, `settle`, waits only for `lab.l`'s `waiting` and for each stub's `ended-<pid>` file. The stub writes that file before it exits, so `settle` can return while `fresh` still records the run and its usage under the project root. The root is the `.../001` folder in the failure. This is I-0048's cause, a launcher nobody waits for, on a launcher S-0122's fix did not reach.

| Layer | Task | Touches | After |
|-------|------|---------|-------|
| 1 | T-3112 The lab keeps every launcher it makes, and `settle` waits for all of them; the restart test takes `fresh` from the lab; a test checks it | `flai/internal/serve/agents_test.go` | none |
| 2 | T-3113 Close I-0343 with the fix's commit and race-run counts | I-0343's file, `design/issues/summary.md` | T-3112 |

Figures:

- Forecast: 40m. flai gave 6m; I raised it for repeated `-race` runs and an 8m integration close-out.
- Cost of delay: 30 USD a week, from the operator's 12m per cycle. It stands as given.

Assumptions:

- The fix needs no change to the launcher. If T-3112 finds otherwise, its agent widens the touches.
- A test of the lab's `settle` is the reproduction that fits. A deterministic reproduction of the race would need a hook in production code, which this plan does not add.

The start prompt asked me to read a thread named `asked`. No thread exists on S-0715, so I planned the story from its issue.

### 2026-11-11T09:06:42Z orchestrator
Approved. The plan fits S-0715.

- The cause is I-0048's: a launcher `settle` does not wait for. The fix stays in `flai/internal/serve/agents_test.go`, which the story touches.
- T-3112, then T-3113, is the right order: the issue's close reason cites T-3112's commit and race counts.
- Each task's touches are files, and the story's touches cover them.
- The 40m forecast and the 30 USD a week value, from flai's 12m input, stand.

If T-3112 finds the cause in `flai/internal/serve/agents.go`, its agent widens the touches, as the plan says.

Source: design/conventions/strategic-agents.md § As the orchestrator

### 2026-11-11T09:07:07Z orchestrator
Resolved: Plan approved by the orchestrator under plan_backlog_stories: T-3112 then T-3113, 40m, 30 USD a week.
