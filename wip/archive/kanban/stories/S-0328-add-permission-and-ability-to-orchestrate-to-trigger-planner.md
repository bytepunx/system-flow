---
id: S-0328
type: story
nature: feature
title: Add Permission and Ability to Orchestrate to Trigger Planner
status: done
owner: alex
created: 2026-10-07T19:35:32Z
updated: 2026-10-07T20:28:28Z
transitions:
  - to: ready
    at: 2026-10-07T19:35:33Z
    by: alex
  - to: in-progress
    at: 2026-10-07T19:35:39Z
    by: agent-S-0328
  - to: review
    at: 2026-10-07T20:24:15Z
    by: agent-S-0328
  - to: done
    at: 2026-10-07T20:28:28Z
    by: alex
tags: [flai, orchestrator, planner]
topics: [planning, orchestration]
touches: [design/adrs, design/adrs/README.md, flai/internal/manifest/manifest.go, flai/internal/manifest/manifest_test.go, flai/internal/manifest/settings.go, flai/internal/manifest/settings_test.go, flai/internal/workitem/plancandidates.go, flai/internal/workitem/plancandidates_test.go, flai/internal/workitem/promotable.go, flai/cmd/plan.go, flai/cmd/plan_candidates_test.go, flai/cmd/plan_orchestrate_test.go, flai/cmd/orchestrate_permissions_test.go, flai/cmd/guard.go, flai/cmd/edit.go, flai/cmd/edit_test.go, flai/internal/guard/guard.go, flai/internal/guard/guard_test.go, flai/internal/guard/orchestrate_test.go, flai/internal/mcpserver/plan.go, flai/internal/mcpserver/plan_test.go, flai/internal/mcpserver/items_write.go, flai/internal/mcpserver/items_write_test.go, flai/internal/harness/harness.go, flai/internal/harness/harness_test.go, ".claude/agents/orchestrator.md", template/root/.claude/agents/orchestrator.md, design/conventions/strategic-agents.md, template/root/design/conventions/strategic-agents.md, template/CHANGELOG.md, design/system/strategic-agents.md, design/system/project-manifest.md, design/system/flai-cli.md, docs/users/flai.md, docs/users/flai-reference.md, docs/users/flaiover.md, docs/operators/settings.md, docs/operators/index.md, flai/internal/guard/thread_test.go, flaiover/src/routes/workflow/orchestrator/orchestrator.svelte.test.ts, design/issues/I-0104-flai-task-done-commits-everything-in-the-worktree-so-two-tasks-of-one-layer-cannot-be-closed-apart.md, design/issues/summary.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 2932
  turns:
    - day: 2026-10-07
      ceremony: 8
      hand_edits: 2
      work: 63
  models:
    - model: claude-haiku-4-5-20251001
      input: 202
      output: 9640
      cache_read: 1612417
      cache_write: 105416
      cost: 0.3414
    - model: claude-opus-5-5
      input: 494
      output: 215557
      cache_read: 36565958
      cache_write: 919181
      cost: 17.3729
  strategic:
    - kind: planner
      seconds: 404
      estimated: true
      models:
        - model: claude-haiku-4-5-20251001
          input: 178
          output: 7345
          cache_read: 1111419
          cache_write: 71743
          cost: 0.2377
        - model: claude-opus-5-5
          input: 206
          output: 38863
          cache_read: 15660150
          cache_write: 172886
          cost: 8.9316
    - kind: orchestrator
      seconds: 75
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 12
          output: 82
          cache_read: 1304552
          cache_write: 12570
          cost: 0.3245
        - model: claude-sonnet-5-5
          input: 10
          output: 62
          cache_read: 113317
          cache_write: 56836
          cost: 0.1467
cost_of_delay:
  inputs:
    penalty_per_week: 800
    by: alex
    at: 2026-10-07T19:35:32Z
  value: 800
  by: planner-S-0328
  at: 2026-10-07T19:41:11Z
forecast:
  duration: 1h45m
  delivery: 2026-10-07T22:00:00Z
  basis: flai forecast's 1h19m (114 s per unit over 29 large-band feature stories, size 41) raised toward S-0219 (40m) plus S-0220 (1h10m), the two stories this one extends, for its serial layers and a new kind of guarded edit.
  by: planner-S-0328
  at: 2026-10-07T19:41:11Z
---
# S-0328 Add Permission and Ability to Orchestrate to Trigger Planner

## Goal

The Orchestrator needs the ability to dispatch the planner on stories in the backlog that are missing a plan. The orchestrator settings should introduce a permission called `plan_backlog_stories` to enable this behavior.

## Acceptance criteria
- [x] The orchestrator settings add a new permission called `plan_backlog_stories` that, when enabled, gives it the permission to delegate planning to the planner agent.
- [x] The orchestrator will identify stories in the backlog that are missing a plan (touches, CoD, and tasks) and, when enabled, will delegate planning of the story to the planning agent.
- [x] The orchestrator will monitor threads from the planner and approve them, answer outstanding questions, and choose a CoD or accept the planner's recommendation.

## Tasks
- T-1168 An ADR records plan_backlog_stories, story candidates, and the orchestrator's handling of the planner's threads
- T-1169 The manifest gains orchestration.permissions.plan_backlog_stories, off by default, in the settings catalog and the operators' index
- T-1170 flai plan --candidates lists the backlog stories that lack a plan: touches, a forecast, a cost of delay value, or tasks
- T-1171 The orchestrator starts the planner for a candidate story only while plan_backlog_stories is on, in flai and in the guard
- T-1172 The orchestrator sets a backlog story's missing cost of delay inputs through item_edit and flai edit, as the ADR allows
- T-1173 The orchestrator's prompt, agent definition, and convention tell it to plan candidate stories and to settle the planner's threads
- T-1174 An ADR records plan_backlog_stories: what it lets the orchestrator plan, answer, and set
- T-1175 The design and the users' and operators' guides describe plan_backlog_stories and the orchestrator's handling of the planner's threads
- T-1176 plan_backlog_stories is a permission, flai plan --candidates lists unplanned backlog stories, and flai holds the orchestrator's story plans and cost of delay edits to it
- T-1177 flai guard lets the orchestrator plan a story, answer and resolve a story planner's threads, and set a backlog story's cost of delay under plan_backlog_stories
- T-1178 The orchestrator's prompt, definition, and convention say what it does with plan_backlog_stories
- T-1179 The design, the guides, the dashboard's settings test, and the template's changelog describe plan_backlog_stories

## Notes

### Planning

Two plans were written at once. The planner drafted T-1168 to T-1173 and T-1175. Meanwhile the story's agent, which had already pulled the story, wrote T-1174 and T-1176 to T-1179 and is working them. The two plans cover the same work, so the planner's seven tasks are duplicates. The plan thread on S-0328 proposes cancelling them and names the files they cover that the agent's tasks do not.

Touches. `flai touches suggest S-0328` gave nothing, since the story declared no touches. Each touch comes from where `plan_backlog_epics` (S-0219) is built and documented today. A sibling permission changes the same places:

| Touches | Source |
|---------|--------|
| `flai/internal/manifest/manifest.go`, `settings.go`, and their tests; `docs/operators/settings.md` | layout: the permission's declaration, its catalog, and the settings index its test requires |
| `flai/internal/workitem/plancandidates.go` and its test; `flai/cmd/plan_candidates_test.go` | layout: `flai plan --candidates`, which lists only epics today |
| `flai/cmd/plan.go`, `plan_orchestrate_test.go`, `orchestrate_permissions_test.go`, `guard.go`; `flai/internal/guard/guard.go` and its tests; `flai/internal/mcpserver/plan.go` and its test | layout: `orchestratorPlans`, the guard's `plansEpics` refusal, and the MCP tool `plan` |
| `flai/internal/workitem/promotable.go`, `flai/internal/mcpserver/items_write.go`, `flai/cmd/edit.go`, and their tests | layout: `Repo.OrchestratorPermits` and the edits it holds, for criterion 3's cost of delay |
| `flai/internal/harness/harness.go` and its test; both `orchestrator.md` agent files; both `strategic-agents.md` conventions; `template/CHANGELOG.md` | layout: `orchestratePrompt` and what the orchestrator is told |
| `design/system/strategic-agents.md`, `project-manifest.md`, `flai-cli.md`; `docs/users/flai.md`, `flai-reference.md`, `flaiover.md`; `docs/operators/index.md` | design: every document that names `plan_backlog_epics` |
| `design/adrs`, `design/adrs/README.md` | design: the new ADR |

The tasks' claim adds `flai/internal/guard/thread_test.go` (T-1177) and `flaiover/src/routes/workflow/orchestrator/orchestrator.svelte.test.ts` (T-1179), both declared by the story's agent.

Folder touch kept: `design/adrs`. The ADR's file name follows its title, which no task could name when it was planned.

Not touched: `system-flow.yaml`. Turning the permission on is the operator's.

Forecast: 1h45m, with delivery at 2026-10-07T22:00:00Z. `flai forecast` gave 1h19m from 29 large-band feature stories at size 41. I raised it because the story extends both S-0219 (40m of agent time) and S-0220 (1h10m), runs in serial layers, and adds a new kind of guarded edit. The delivery scales flai's 21:25Z by the same ratio from the 19:36Z start. The duplicate tasks add no work once one plan is kept.

Cost of delay: 800 USD a week. `flai cod` gives this from the operator's `penalty_per_week: 800`, and it is kept as is.
