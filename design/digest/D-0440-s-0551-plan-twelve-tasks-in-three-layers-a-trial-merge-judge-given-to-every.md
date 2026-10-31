---
id: D-0440
title: "S-0551 plan: twelve tasks in three layers, a trial-merge judge given to every caller of the holds: Recommendation: take the plan as written"
kind: answer
sources: [TH-0776]
topics: [git, orchestration, planning, conventions, dashboard, testing, flai, cli, go, flaiover, sveltekit, template, code, all]
paths: [flai/internal/serve, design/workflows/roles, design/system/workflow.md, design/system/strategic-agents.md, design/conventions/work-management.md, docs/users/flai.md, design/adrs, design/adrs/README.md, design/system/metrics.md, design/system/flai-cli.md, design/system/flaiover-dashboard.md, flai/internal/storygit/trial.go, flai/internal/storygit/trial_test.go, flai/internal/workitem/hold.go, flai/internal/workitem/hold_test.go, flai/internal/check/check.go, flai/internal/check/check_test.go, flai/internal/serve/agents.go, flai/internal/serve/agents_test.go, flai/internal/serve/hold_test.go, flai/internal/serve/restart.go, flai/internal/serve/start.go, flai/internal/mcpserver/work.go, flai/internal/mcpserver/work_test.go, flai/internal/mcpserver/orchestrate.go, flai/internal/workitem/promote.go, flai/internal/workitem/boardview.go, flai/internal/storystart/start.go, flai/internal/storystart/start_test.go, flai/internal/inbox/inbox.go, design/workflows/roles/planner.md, template/root/design/workflows/roles/planner.md, flai/internal/roles/baseline/planner.md, flai/internal/harness/harness.go, flai/internal/harness/harness_test.go, template/root/.claude/agents/planner.md, ".claude/agents/planner.md", design/conventions/strategic-agents.md, template/root/design/conventions/strategic-agents.md, template/root/design/conventions/work-management.md, template/CHANGELOG.md, flai/internal/metrics/concurrency.go, flai/internal/metrics/concurrency_test.go, flai/internal/metrics/metrics.go, flai/cmd/stats.go, flai/cmd/stats_test.go, flaiover/src/lib/viz/catalog.ts, flaiover/src/lib/viz/catalog.test.ts, flaiover/src/lib/viz/charts.ts, flaiover/src/lib/viz/charts.test.ts, flaiover/src/lib/viz/explain.ts, flaiover/src/lib/sitemenu.ts, "flaiover/src/routes/charts/[kind]/charts.svelte.test.ts", docs/users/flaiover.md, flai/internal/workitem/repo.go, flai/internal/workitem/share.go, flai/internal/check/after_reason_test.go, flai/internal/check/scope.go, flai/internal/check/scope_test.go]
date: 2026-10-31
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0440 S-0551 plan: twelve tasks in three layers, a trial-merge judge given to every caller of the holds: Recommendation: take the plan as written

The operator, alex, answered TH-0776 at 2026-10-31T09:10:54Z.

What was asked:

```text
S-0551 now has twelve tasks in three layers, predicted touches, and a 2h forecast. The cost of delay value waits on TH-0772.
```

The answer:

```text
**Recommendation:** take the plan as written. Do not split it. Leave the declared touches to the story's agent.

## The plan: approved

- **The tasks fit the criteria.**
  - Criterion 1: T-2389, T-2391, T-2395, and T-2404.
  - Criterion 2: T-2393, T-2406, T-2408, and T-2412.
  - Criterion 3: T-2410 and T-2414.
  - Criterion 4: T-2388.
  - T-2416 carries the docs.
- **The layers follow the real dependencies.** The judge is passed in because `workitem` cannot import `storygit`. Holding as today when the trial merge cannot run keeps the change safe.
- **The figures fit.** The 2h forecast over flai's 1h30m states its reason. The cost of delay input is on TH-0772.

## Proposal 1, drop the four folder touches: no change is needed now

- ADR-0096 §3 already narrows the claim. A declared folder that holds a task touch is replaced by those task touches.
  - `flai/internal/serve` already claims only the five serve files of T-2395.
  - `design/workflows/roles` already claims only `planner.md`.
  - Dropping them frees none of S-0546, S-0549, or S-0550. The holds on the board come from S-0548 and S-0552.
- `flai/internal/sync` and `flai/internal/claims` name no package. The story's agent can remove them with `flai touches`. S-0550 declares the same two folders, and their only overlap is with each other.

## Proposal 2, split off the stats and the dashboard: no

- Criterion 3 is one of the four you wrote, so splitting it off changes the story's scope.
- T-2410 and T-2414 sit in layers 2 and 3 and share no file with the hold work. They do not delay criterion 1 inside the story.
- The 2h forecast fits one session. Splitting the story would add a second acceptance and a second metric ADR hand-off, and gain no time.

## What needs you

- On TH-0772: run `flai edit S-0551 --time-lost-per-cycle 2h`.
- Confirm this entry, or name the split you want instead.

Source: design/adrs/0096-a-story-in-review-holds-nothing-an-overlap-inside-the-manifest-s-shared-paths.md § 3. Tasks narrow a folder claim
```

Source: TH-0776.
