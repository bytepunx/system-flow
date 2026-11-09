---
id: TH-0918
title: "Plan for S-0612: flai guard lets every agent run flai hostapi's read methods"
anchor:
  path: wip/kanban/stories/S-0612-flai-guard-refuses-a-story-s-task-sub-agent-flai-hostapi-so-a-measurement-task.md
  item: S-0612
status: resolved
participants: [planner-S-0612, orchestrator]
created: 2026-11-09T08:00:06Z
updated: 2026-11-09T08:02:47Z
---

# TH-0918 Plan for S-0612: flai guard lets every agent run flai hostapi's read methods

On wip/kanban/stories/S-0612-flai-guard-refuses-a-story-s-task-sub-agent-flai-hostapi-so-a-measurement-task.md.

## Entries

### 2026-11-09T08:00:06Z planner-S-0612
S-0612 is planned. It has 4 tasks in 3 layers, 8 file touches, a 45m forecast, and a cost of delay of 50 USD a week. Nothing needs your answer. The story is still a draft for you to finalize.

## Proposed fix

`flai guard` treats `flai hostapi` as a read when it names no method or a method on a list of the host API's read methods kept in the guard. It refuses every other method and names it. This follows S-0568, which did the same for `flai role`.

## Tasks and layers

| Layer | Task | Waits for | Touches |
|-------|------|-----------|---------|
| 1 | T-2702: guard rule, read list, tests that reproduce I-0307 | none | `flai/internal/guard/guard.go`, `flai/internal/guard/hostapi_test.go`, `flai/cmd/guard_test.go` |
| 2 | T-2703: test that keeps the guard's list in step with the host API's method table | T-2702 (the list) | `flai/internal/hostapi/guardreads_test.go` |
| 2 | T-2704: docs, § Request timing and the `flai guard` row | T-2702 (the rule) | `docs/users/flai.md`, `design/system/flai-cli.md` |
| 3 | T-2705: close I-0307 | T-2703, T-2704 | `design/issues/I-0307-…md`, `design/issues/summary.md` |

T-2703 and T-2704 share no path, so they can run together.

## Assumptions

- The read list lives in `flai/internal/guard`. That package cannot import `flai/internal/hostapi`, because hostapi imports guard. T-2703's test is what stops the two lists drifting apart.
- A method that writes anything is left off the list, a cursor or kept state under `.flai-cache` included. T-2702 checks each candidate's code before it lists it.
- `reads` is shared, so the planner's, the orchestrator's, and the analyzer's own calls get these reads too, as they got `flai role list`.
- No ADR is needed: allowing a read fits within ADR-0060, as S-0568 and S-0491 did.
- The `git clone` refusal that I-0307's instance also mentions is out of scope. The issue is about `flai hostapi`.
- Forecast: flai's 18m is raised to 45m. S-0568, the closest earlier fix, took 1927 s, and T-2703 has to class every host API method.
- Cost of delay: 50 USD a week, taken as flai worked it out from the 20m input.

I would not split, merge, or drop any task.

### 2026-11-09T08:02:18Z orchestrator
The plan for S-0612 is approved as drafted.

- **Tasks fit.** There are 4 tasks in 3 layers. T-2702 adds the guard's hostapi read rule and the tests that reproduce I-0307. T-2703 keeps the read list in step with the method table. T-2704 does the docs. T-2705 closes I-0307. Criterion 1 is met by T-2702 and T-2703, and criterion 2 by T-2705.
- **Touches fit.** The 8 file touches cover every task's touches, with no folder touch. T-2703 and T-2704 share no path, so they can run together.
- **Figures fit.** `flai promote --drafts` lists S-0612 as complete. Its forecast is 45m against flai's 18m, with the reason stated (S-0568's 1927 s). Its value of 50 USD a week is flai's, worked out from the issue's 20m input.
- **Scope holds.** It allows a read within ADR-0060, as S-0568 did, so no ADR is needed. The `git clone` refusal stays out of scope.

I finalize S-0612 next under `finalize_drafts`.

Source: design/adrs/0119-with-plan-backlog-stories-on-the-orchestrator-asks-the-planner-for-each-backlog.md § Decision

### 2026-11-09T08:02:47Z orchestrator
Resolved: Plan approved under plan_backlog_stories: tasks, touches, and figures fit S-0612.
