---
id: TH-1081
title: "S-0717 plan: three tasks in two layers, after S-0716, 1h, 12.50 USD a week; proposal to fold it into S-0716"
anchor:
  path: wip/kanban/stories/S-0717-wait-for-events-fails-with-no-front-matter-while-another-story-s-checked.md
  item: S-0717
status: open
participants: [planner-S-0717, orchestrator]
created: 2026-11-11T09:08:33Z
updated: 2026-11-11T09:11:02Z
---

# TH-1081 S-0717 plan: three tasks in two layers, after S-0716, 1h, 12.50 USD a week; proposal to fold it into S-0716

On wip/kanban/stories/S-0717-wait-for-events-fails-with-no-front-matter-while-another-story-s-checked.md.

## Entries

### 2026-11-11T09:08:33Z planner-S-0717
Plan for S-0717. The story had no tasks and no earlier plan, and no thread on it held a request from the operator. So this is a first plan.

## Finding

I-0345 and I-0344 record one event: T-3079's file held only `PLACEHOLDER` from about 11:04:00Z to 11:04:25Z on 2026-11-10. The checked creation of ADR-0272 did not write it. `itemnew.Create` writes a whole item through `atomicfile.WriteFile`, and no Go code under `flai/` holds that string. T-3079 was created at 11:03:52Z, so the file was overwritten after it was created and then restored. The writer is not known yet. The failure itself is `Repo.List` failing the whole read on one file that does not parse, and S-0716 fixes that.

## Tasks

| Layer | Task | Nature | Touches | After |
|-------|------|--------|---------|-------|
| 1 | T-3120 Find what wrote PLACEHOLDER into T-3079's file and correct I-0345's account of it | research | I-0345 issue file | none |
| 1 | T-3121 wait_for_events, in one project and across a folder of projects, answers past a PLACEHOLDER task file, with tests | remediation | `flai/internal/mcpserver/cursor.go`, `folder.go`, `server_test.go`, `folder_test.go` | none |
| 2 | T-3122 Close I-0345 with what fixed it | remediation | I-0345 issue file, `design/issues/summary.md` | T-3120, T-3121 |

T-3120 and T-3121 share no file, so they run together.

## Figures

- `after` S-0716: T-3121's tests pass only on S-0716's change to `Repo.List`.
- Forecast 1h, delivery 2026-11-12T00:07:00Z. flai gave 18m. The extra time covers the log search, a fixture for a folder of two projects, and the integration and smoke tiers at close-out.
- Cost of delay 12.50 USD a week, from the operator's 5m per cycle, unchanged.

## Assumptions

- S-0716's T-3109 covers `inbox` and single-project `wait_for_events` in `inbox_test.go`. This story adds the held wait and the folder of projects only.
- `cursor.go` and `folder.go` change only if a test still fails once S-0716 is on main.
- If T-3120 finds that flai code wrote the file, a fix of five minutes or less becomes a task here. A larger one goes to a thread.

## Proposal

I recommend folding S-0717 into S-0716 and cancelling S-0717. The two issues record the same event. S-0716's T-3110 could close I-0345 next to I-0344, and T-3121's folder-of-projects test could join T-3109. That saves one close-out and one pass through review. If you agree, say so here. I have cancelled and changed nothing in either story.

### 2026-11-11T09:11:02Z orchestrator (recommendation)
**Recommended:** keep S-0717 as its own story, after S-0716, and approve this plan as written. Do not fold it into S-0716.

- S-0716 is in ready, finalized, with four tasks in three layers, 17 touches, and a 1h30m forecast. Folding adds T-3120's search for an unknown writer, which has no known size, to a story about to be pulled. A story is sized before it starts, not after.
- The `after` on S-0716 is the right link: T-3121's tests read S-0716's change to `Repo.List`. S-0716 starts unchanged, and S-0717 follows once it is done.
- The saving is one close-out and one pass through review. Against that, S-0716's fix would wait on log research that may find nothing.

The plan itself fits the story: T-3120 and T-3121 in layer 1 share no file, T-3122 closes I-0345 after both, the touches stay within the story's, and the 1h forecast and 12.50 USD a week value stand on the operator's 5m input.

**Alternative:** say "fold it", and S-0716's agent takes T-3121's folder-of-projects test into T-3109 and I-0345's close into T-3110, with S-0717 cancelled. That changes a finalized story's scope, so it is yours to call.

Source: design/conventions/work-management.md
