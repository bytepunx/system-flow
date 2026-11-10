---
id: TH-0973
title: "Plan for S-0652: two test fixes, the dev server's deadline and the owners' MCP slow threshold, then I-0329 closes"
anchor:
  path: wip/kanban/stories/S-0652-teststartdevserversayswhatavitethatdoesnotanswer200answered-and.md
  item: S-0652
status: resolved
participants: [planner-S-0652, orchestrator]
created: 2026-11-10T07:01:20Z
updated: 2026-11-10T07:07:19Z
---

# TH-0973 Plan for S-0652: two test fixes, the dev server's deadline and the owners' MCP slow threshold, then I-0329 closes

On wip/kanban/stories/S-0652-teststartdevserversayswhatavitethatdoesnotanswer200answered-and.md.

## Entries

### 2026-11-10T07:01:20Z planner-S-0652
S-0652 is planned: three tasks in two layers, 45m, 27.50 USD a week.

## Tasks

| Layer | Task | Touches | After |
|-------|------|---------|-------|
| 1 | T-2841 The dev server's 404 test allows the fake vite time to start, and a start delay reproduces the missing 404 | `flai/internal/components/devserver_test.go` | none |
| 1 | T-2843 The owners' MCP test server sets a slow threshold no call reaches, and `FLAI_SLOW_REQUEST=1ns` reproduces the slow log | `flai/internal/mcpserver/server_test.go` | none |
| 2 | T-2845 Close I-0329 with what fixed it | the I-0329 file, `design/issues/summary.md` | T-2841, T-2843 |

T-2841 and T-2843 share no path, so they run together.

## Causes, from the code

- The dev server test gives `StartDevServer` 600 ms. Under load the fake vite is not serving by then, so the error lacks the 404 the test asserts.
- The MCP test asserts that nothing was logged. `ownersServer` leaves `Options.Slow` zero, so a `thread_reply` over 500 ms logs a slow "request answered" event.

## Figures

- Forecast 45m, raised from `flai forecast`'s 5m: two fixes, each reproduced, and an integration-tier close-out of 10 to 20 minutes under load. Delivery 2026-11-10T08:52Z.
- Cost of delay 27.50 USD a week, as `flai cod` gives it from flai's 11m input.

## Assumptions

- `after: S-0647`. S-0647, in progress, makes the dev server tests check the pid flai started and adds a start delay to the fake vite. T-2841 builds on both.
- Both fixes change only tests. No product code changes.
- The run was started as asked again, but S-0652 had no earlier plan and no thread. So this is a first plan, not a revision.
- I-0329 now has two instances, 15m in all, while flai's input counts one, 11m. I left the input as it is because it is yours. Raise `time_lost_per_cycle` if you want the cost of delay to count both.

Nothing to split, merge, or drop.

### 2026-11-10T07:05:36Z orchestrator
The plan is approved.

- **Tasks.** T-2841 and T-2843 each fix one test and reproduce its failure first, with a start delay and `FLAI_SLOW_REQUEST=1ns`. They share no path. T-2845 waits for both, then closes I-0329. That fits both criteria.
- **Touches.** The two test files, I-0329, and `design/issues/summary.md`. Leaving out `server.go` is right, because `Options.Slow` already exists.
- **After.** `after: S-0647` is right. T-2841 reads the pid check and the fake vite's start delay that S-0647 writes, and the Notes name it.
- **Figures.** The 45m forecast states why it is over flai's 5m. The 27.50 USD a week value is what `flai cod` gives from flai's 11m input. The input stays as flai set it.

I am resolving this thread.

Source: design/conventions/strategic-agents.md § As the orchestrator

### 2026-11-10T07:07:19Z orchestrator
Resolved: Plan approved under plan_backlog_stories: tasks, touches, after, and figures fit S-0652.
