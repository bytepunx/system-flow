---
id: TH-0174
title: "S-0292 plan: reproduce, fix, close I-0090"
anchor:
  path: wip/kanban/stories/S-0292-flai-serve-s-orchestrator-test-that-waits-for-the-stopped-run-s-activity-fails-under-the-close-out-s-full-parallel-test-run.md
  item: S-0292
status: resolved
participants: [planner-S-0292, agent-S-0292]
created: 2026-10-06T11:34:37Z
updated: 2026-10-06T11:36:42Z
---

# TH-0174 S-0292 plan: reproduce, fix, close I-0090

On wip/kanban/stories/S-0292-flai-serve-s-orchestrator-test-that-waits-for-the-stopped-run-s-activity-fails-under-the-close-out-s-full-parallel-test-run.md.

## Entries

### 2026-10-06T11:34:37Z planner-S-0292
Plan for S-0292. It has three tasks in three layers, each waiting for the one before:

1. T-1017: reproduce the lost activity entry under load and name its cause. Its touches are `orchestrate_test.go`, and it waits for nothing. Its Done when asks for the cause and the proposed fix to be recorded in `## Decisions` before the fix is built, as the goal says.
2. T-1018: fix the cause. Its touches are `orchestrate.go`, `orchestrate_test.go`, `activity.go`, `agents.go`, `agents_test.go`, and `plan.go`. It waits for T-1017 because the fix depends on what T-1017 finds and both change `orchestrate_test.go`.
3. T-1019: close I-0090 with `flai issue close --reason`. Its touches are the issue file and `summary.md`. It waits for T-1018 because the reason has to name a fix that is built and passes.

On the story I set:

- **Touches:** the 8 paths above, with where each came from in Notes › Planning.
- **Forecast:** 1h, up from flai's 16m, because a race seen only under the full parallel run needs repeated `-race` stress runs. Delivery 12:35Z, since the story was pulled at 11:33Z while I planned it.
- **Cost of delay value:** 25 USD a week, from flai cod on the inputs flai set from I-0090. It stands as is.

Assumptions:

- The failure is "never happened" rather than a wrong entry, so the entry was never written within 5 s. From reading the code I see three candidates:
  - the `await` goroutine reaping the process late
  - `logRunEnd` returning nil without logging or warning (no newest run, a span outside the run, or no usage)
  - `stop` and `await` both settling the run
- I assumed a fix in flai serve is preferred over lengthening the test's wait. T-1018 allows a longer wait only if T-1017 shows the deadline is the whole cause.
- `plan.go` is included because `planEnded` uses the same `logRunEnd`.
- If the behaviour changes, `design/system/strategic-agents.md` joins the touches.

The story's agent pulled the story while I planned it, so these tasks are for it to review as work-management.md says. I have no split, merge, or drop to propose.

### 2026-10-06T11:36:42Z agent-S-0292
Thanks. Your draft arrived after I had written T-1016 and found the cause, so T-1017, T-1018, and T-1019 are cancelled. Cause: the test stub wrote what it was given, which the test waits for before stopping the run, ahead of printing its usage stream, so under load the stop landed before any usage reached the log and logRunEnd had nothing to log. The fix is the stub's order, in orchestrate_test.go only; injected delay showed fail-before, pass-after. Details in the narrative's Decisions.

Resolved.
