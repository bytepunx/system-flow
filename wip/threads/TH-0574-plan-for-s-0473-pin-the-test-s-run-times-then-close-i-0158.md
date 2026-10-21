---
id: TH-0574
title: "Plan for S-0473: pin the test's run times, then close I-0158"
anchor:
  path: wip/kanban/stories/S-0473-testasettledreplystartstheagentwaitingonitonce-in-flai-internal-serve-fails-now-and-then-in-the-close-out-s-integration-tier-on-a-branch-that-does-not-touch-serve.md
  item: S-0473
status: resolved
participants: [planner-S-0473, orchestrator, agent-S-0473, alex]
created: 2026-10-21T12:08:46Z
updated: 2026-10-21T13:05:46Z
---

# TH-0574 Plan for S-0473: pin the test's run times, then close I-0158

On wip/kanban/stories/S-0473-testasettledreplystartstheagentwaitingonitonce-in-flai-internal-serve-fails-now-and-then-in-the-close-out-s-integration-tier-on-a-branch-that-does-not-touch-serve.md.

## Entries

### 2026-10-21T12:08:46Z planner-S-0473
The plan is written: the story's touches, a 45m forecast, a cost of delay value of 112.50 USD a week, and two tasks.

The operator finalized S-0473 and agent-S-0473 started it while the planner was planning. That agent then wrote T-1992, which covers the same fix and the issue close.

## Cause

The test's two `putEnded` calls stamp the runs with the wall-clock second. When they straddle a second boundary, the message from `id` to `other` falls strictly after `other`'s run (`messaged` in `flai/internal/serve/agents.go`, through `messages.ToSince`). The first look then starts `other`'s agent on MS-0001, which fails at `started before the reply`. The fix belongs in the test, not in `messaged`.

## Tasks and layers

| Layer | Task | Touches | After |
| --- | --- | --- | --- |
| 1 | T-1986: reproduce with a one-second shift, then pin `other`'s run to the message's second | `flai/internal/serve/agents_test.go` | none |
| 2 | T-1989: `flai issue close I-0158 --commit` with the reason | the I-0158 file, `design/issues/summary.md` | T-1986 |

## Proposal

T-1992, written by the story's agent, duplicates T-1986 and T-1989. It sits in layer 1 beside T-1986 and touches the same file. Recommended: the story's agent keeps one plan and cancels the other's tasks. Either T-1992 alone, or T-1986 then T-1989. The planner cancels nothing.

## Assumptions

- No production change is needed. T-1986 says to plan again and widen its touches if the reproduction points to `agents.go`.
- The other tests in `agents_test.go` that use `putEnded` make one run each and set their times with `endedAgo`, so they do not share the flake.
- The cost of delay value stands on the operator's inputs, 45m per 168h cycle. The instances suggest these inputs understate: 5 occurrences in 4.3 days is about 8 a week, or about 73m per cycle, and the merge-queue instance also cost a story its place in the queue. Recommended: the operator raises `time_lost_per_cycle` to 73m, which would give about 182.50 USD a week. Planning does not wait on this.

### 2026-10-21T12:10:12Z orchestrator
Approved under plan_backlog_stories: the cause, a run-time stamp straddling a second boundary, and the fix in the test alone. Forecast 45m.

On the duplicate: S-0473 is in progress, so its agent keeps one plan. It reviews the planner's tasks against its own and cancels or changes what it would plan differently, saying why in `## Decisions`. That is either T-1992 alone, or T-1986 then T-1989. The planner cancels nothing, as you say.

The cost of delay input follows as a recommendation for the operator, so this thread stays open for it.

Source: design/conventions/work-management.md § Rules

### 2026-10-21T12:10:16Z orchestrator (recommendation)
Recommended for the operator: raise S-0473's `time_lost_per_cycle` from 45m to 73m, which is about 182.50 USD a week. The planner reckons I-0158 at 5 occurrences in 4.3 days, about 8 a week, and one of them sent S-0463 back out of the merge queue.

It changes no order today, since the policy is throughput and S-0473 is already in progress. It does record what the flake cost for `flai stats`. The story's inputs are yours, and my permission does not reach a story in progress, so set it with `flai edit S-0473` and its cost of delay flags. Otherwise leave 45m as it is.

Source: design/system/strategic-agents.md § Cost of delay

### 2026-10-21T12:11:15Z agent-S-0473
I kept T-1992 and cancelled T-1986 and T-1989 as duplicates. T-1992 is done: the test puts the other run a second earlier, which reproduced "started before the reply" on every run before the fix. The fix puts the waiting story's run, and so its message, a minute before the other run started, and the test now passes. I-0158 is closed on the branch. The cost-of-delay input is the operator's to set, and I leave it to them.

### 2026-10-21T12:43:36Z alex
Set COD to 190

### 2026-10-21T12:43:55Z orchestrator
I cannot set it: S-0473 is in progress, and flai lets the orchestrator set a cost of delay on a backlog story only. It refused me the same on S-0464 earlier today. To set the value to 190 USD a week, run:

```sh
flai edit S-0473 --cost-of-delay-value 190
```

Or, to set the input the planner suggested so that `flai cod` works the value out, `flai edit S-0473 --time-lost-per-cycle 73m`, which gives about 182.50.

Source: design/adrs/0119-with-plan-backlog-stories-on-the-orchestrator-asks-the-planner-for-each-backlog.md

### 2026-10-21T13:05:46Z orchestrator
Resolved: The plan was settled: S-0473's agent kept T-1992 and the story is in review. The operator's "Set COD to 190" stays theirs to run, since flai gives the orchestrator a cost of delay on backlog stories only: flai edit S-0473 --cost-of-delay-value 190.
