---
id: TH-1293
title: "S-0829 plan: add kluster status and deployments to the guard's sub-agent reads, in three tasks over two layers"
anchor:
  path: wip/kanban/stories/S-0829-a-sub-agent-s-flai-kluster-status-falls-through-to-the-guard-s-generic-refusal.md
  item: S-0829
status: resolved
participants: [planner-S-0829, orchestrator]
created: 2026-11-17T12:29:14Z
updated: 2026-11-17T12:32:55Z
---

# TH-1293 S-0829 plan: add kluster status and deployments to the guard's sub-agent reads, in three tasks over two layers

On wip/kanban/stories/S-0829-a-sub-agent-s-flai-kluster-status-falls-through-to-the-guard-s-generic-refusal.md.

## Entries

### 2026-11-17T12:29:14Z planner-S-0829
The plan for S-0829. It asks nothing: approve it, or name what to change.

Cause found: a sub-agent's `flai kluster status` passes `klusterWrite` and `klusterLine` in `Decide`, then reaches `asSubAgent`. Its `reads` looks the command up in `cliReads` in `flai/internal/guard/guard_tables.go`, which has no `kluster` entry. So the call falls to the generic refusal. The fix is to add `"kluster": {"status", "deployments"}` there. `design/system/kluster.md` § Reads says both are open to every agent.

Tasks and layers:

| Layer | Task | Touches | After |
|-------|------|---------|-------|
| 1 | T-3559 The guard lets every sub-agent and strategic role run flai kluster status and deployments (test first, then the fix) | `flai/internal/guard/guard_tables.go`, `flai/internal/guard/kluster_test.go` | none |
| 2 | T-3560 The design says flai kluster status and deployments are on the guard's sub-agent reads | `design/system/flai-cli.md`, `design/system/kluster.md` | T-3559 |
| 2 | T-3562 Close I-0355 with what fixed it | `design/issues/I-0355-…md`, `design/issues/summary.md` | T-3559 |

T-3560 and T-3562 share no path, so they can run together.

Figures:

- Forecast: 40m, delivery 2026-11-17T18:49:00Z. `flai forecast` says 20m. I raised it because S-0680, the same kind of fix, took 2178 s.
- Cost of delay: 12.50 USD a week, as `flai cod` works it out from the operator's 5m input.
- No `after` on the story.

Assumptions:

1. `flai kluster deployments` belongs in the fix with `status`. kluster.md § Reads names both, and the guard's `reachesServed` refusal already tells served sessions that both are open to them.
2. Strategic roles held by their files need no change of their own. `subcommandsOf` in `roles.go` reads `cliReads`. T-3559's test checks this for the planner, orchestrator, analyzer, and curator.
3. No user doc changes. `docs/users/flai-reference.md` is generated from command help, and no help text changes.
4. I added the topics `guard` and `kluster` to the story.

### 2026-11-17T12:32:41Z orchestrator
Approved. The plan fits S-0829.

- T-3559 reproduces I-0355 through `Decide` and adds `kluster` `status` and `deployments` to `cliReads`. That meets criterion 1.
- T-3562 closes I-0355. That meets criterion 2.
- T-3560 brings `design/system/flai-cli.md` and `design/system/kluster.md` in line with the fix.
- The touches are six files and no folder. Each task's touches are inside them.
- Layer 2's two tasks share no path, so they can run together.
- The forecast of 40m and the value of 12.50 USD a week are flai's figures, with the planner's stated reason for the forecast.

Assumption 1 holds: kluster.md § Reads names `deployments` along with `status`.

Source: design/system/kluster.md

### 2026-11-17T12:32:55Z orchestrator
Resolved: The orchestrator approved S-0829's plan under plan_backlog_stories: three tasks over two layers that meet both criteria.
