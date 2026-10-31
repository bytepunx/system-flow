---
id: D-0439
title: "S-0551 needs a cost of delay input: recommend time lost 2h per cycle: Recommendation: set `time_lost_per_cycle: 2h` on S-0551, as the planner recommends"
kind: answer
sources: [TH-0772]
topics: [git, orchestration, planning, conventions, dashboard, testing, flai, cli, go, flaiover, sveltekit, template, code, all]
paths: [flai/internal/serve, design/workflows/roles, design/system/workflow.md, design/system/strategic-agents.md, design/conventions/work-management.md, docs/users/flai.md, design/adrs, design/adrs/README.md, design/system/metrics.md, design/system/flai-cli.md, design/system/flaiover-dashboard.md, flai/internal/storygit/trial.go, flai/internal/storygit/trial_test.go, flai/internal/workitem/hold.go, flai/internal/workitem/hold_test.go, flai/internal/check/check.go, flai/internal/check/check_test.go, flai/internal/serve/agents.go, flai/internal/serve/agents_test.go, flai/internal/serve/hold_test.go, flai/internal/serve/restart.go, flai/internal/serve/start.go, flai/internal/mcpserver/work.go, flai/internal/mcpserver/work_test.go, flai/internal/mcpserver/orchestrate.go, flai/internal/workitem/promote.go, flai/internal/workitem/boardview.go, flai/internal/storystart/start.go, flai/internal/storystart/start_test.go, flai/internal/inbox/inbox.go, design/workflows/roles/planner.md, template/root/design/workflows/roles/planner.md, flai/internal/roles/baseline/planner.md, flai/internal/harness/harness.go, flai/internal/harness/harness_test.go, template/root/.claude/agents/planner.md, ".claude/agents/planner.md", design/conventions/strategic-agents.md, template/root/design/conventions/strategic-agents.md, template/root/design/conventions/work-management.md, template/CHANGELOG.md, flai/internal/metrics/concurrency.go, flai/internal/metrics/concurrency_test.go, flai/internal/metrics/metrics.go, flai/cmd/stats.go, flai/cmd/stats_test.go, flaiover/src/lib/viz/catalog.ts, flaiover/src/lib/viz/catalog.test.ts, flaiover/src/lib/viz/charts.ts, flaiover/src/lib/viz/charts.test.ts, flaiover/src/lib/viz/explain.ts, flaiover/src/lib/sitemenu.ts, "flaiover/src/routes/charts/[kind]/charts.svelte.test.ts", docs/users/flaiover.md, flai/internal/workitem/repo.go, flai/internal/workitem/share.go, flai/internal/check/after_reason_test.go, flai/internal/check/scope.go, flai/internal/check/scope_test.go]
date: 2026-10-31
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0439 S-0551 needs a cost of delay input: recommend time lost 2h per cycle: Recommendation: set `time_lost_per_cycle: 2h` on S-0551, as the planner recommends

The operator, alex, answered TH-0772 at 2026-10-31T09:09:42Z.

What was asked:

```text
**Recommendation:** set `time_lost_per_cycle: 2h` on S-0551: `flai edit S-0551 --time-lost-per-cycle 2h`.
```

The answer:

````text
**Recommendation:** set `time_lost_per_cycle: 2h` on S-0551, as the planner recommends:

```bash
flai edit S-0551 --time-lost-per-cycle 2h
```

## The figure, and not 1h or 4h

- **The goal is time lost.** The story removes the hold on stories that would have merged cleanly. That is waiting time, not revenue or a penalty.
- **The story points to no alternative.** It does not say the in-progress limit causes most of the wait, which would argue for 1h. It does not count the host's idle time as lost, which would argue for 4h.
- **It matches D-0316.** You set 2h on S-0427, also a flow-overhead story.

## Why this needs you

I tried to set it, and flai refused: "S-0551 is ready, not in the backlog: the orchestrator gives cost of delay inputs to a story in the backlog alone". You moved S-0551 to ready at 2026-10-31T07:24:00Z, so its input is yours to set. A confirmation of this entry does not write the input. Running the command does.

For the operator alone: flai serve does not confirm this after orchestration.permissions.confirm_after.

Source: design/adrs/0119-with-plan-backlog-stories-on-the-orchestrator-asks-the-planner-for-each-backlog.md § Decision
````

Source: TH-0772.
