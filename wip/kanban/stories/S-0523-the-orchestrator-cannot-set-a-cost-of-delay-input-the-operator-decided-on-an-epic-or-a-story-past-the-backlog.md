---
id: S-0523
type: story
nature: improvement
title: The orchestrator cannot set a cost of delay input the operator decided on an epic or a story past the backlog
status: ready
owner: alex
created: 2026-10-28T12:53:43Z
updated: 2026-10-30T00:12:55Z
transitions:
  - to: ready
    at: 2026-10-30T00:10:23Z
    by: orchestrator
tags: [flai, orchestrator, cost-of-delay, guard]
topics: [orchestration, planning]
touches: [design/adrs, design/adrs/README.md, design/system/strategic-agents.md, design/system/flai-cli.md, design/system/project-manifest.md, design/conventions/strategic-agents.md, template/root/design/conventions/strategic-agents.md, design/workflows/roles/orchestrator.md, template/root/design/workflows/roles/orchestrator.md, template/CHANGELOG.md, docs/users/flai.md, docs/users/flai-reference.md, docs/users/flaiover.md, docs/operators/settings.md, docs/operators/index.md, flai/internal/manifest/settings.go, flai/internal/guard/guard.go, flai/internal/guard/orchestrate_test.go, flai/internal/guard/guard_test.go, flai/internal/workitem/promotable.go, flai/internal/mcpserver/items_write.go, flai/internal/mcpserver/items_write_test.go, flai/cmd/edit.go, flai/cmd/edit_test.go, flai/cmd/orchestrate_permissions_test.go, flai/internal/harness/harness.go, flai/internal/harness/harness_test.go, design/issues/I-0187-the-orchestrator-cannot-set-a-cost-of-delay-input-the-operator-decided-on-an-epic-or-a-story-past-the-backlog.md, design/issues/summary.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: sum
  seconds: 0
  models: []
  strategic:
    - kind: orchestrator
      seconds: 622
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 114
          output: 1888
          cache_read: 37807309
          cache_write: 30382
          cost: 8.3067
cost_of_delay:
  inputs:
    time_lost_per_cycle: 30m
    by: orchestrator
    at: 2026-10-28T14:34:21Z
  value: 75
  by: planner-S-0523
  at: 2026-10-28T14:35:47Z
forecast:
  duration: 1h15m
  delivery: 2026-10-30T01:54:00Z
  basis: "Its own forecast of 1h15m; 2nd in the pull order with an in-progress limit of 5, behind S-0419, S-0511, S-0527 and S-0424."
  by: flai
  at: 2026-10-30T00:12:55Z
finalized:
  by: alex
  at: 2026-10-28T12:59:11Z
---
# S-0523 The orchestrator cannot set a cost of delay input the operator decided on an epic or a story past the backlog

## Goal

This story remediates [I-0187](../../../design/issues/I-0187-the-orchestrator-cannot-set-a-cost-of-delay-input-the-operator-decided-on-an-epic-or-a-story-past-the-backlog.md), "The orchestrator cannot set a cost of delay input the operator decided on an epic or a story past the backlog". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0187 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0187 is closed with `flai issue close I-0187 --reason` saying what fixed it

## Tasks
- T-2241 An ADR refining ADR-0119 lets the orchestrator add or change cost of delay inputs on an epic or an open story under plan_backlog_stories, inputs only
- T-2242 Repo.OrchestratorSetsCostOfDelay admits an epic and a story in any state but done or cancelled, with or without inputs, still inputs only, and the MCP and CLI tests reproduce I-0187's refusals before the change
- T-2243 flai guard lets the orchestrator's item_edit and flai edit give an epic cost of delay inputs as it does a story, and its refusal texts name what it may set
- T-2244 The design, the convention and its template copy, the orchestrator's prompt, the permission's settings text, and the guides say the orchestrator sets cost of delay inputs the operator decided on an epic or an open story
- T-2245 I-0187 is closed with the reason naming the ADR and the widened rule

## Notes

### Planning

Planned by planner-S-0523 on 2026-10-28; the plan is on TH-0694. The story declared no touches, so `flai touches suggest` had nothing to start from; every touch comes from the issue's instances, the design, and the code layout.

**The solution proposed from I-0187's instances.** Two gates stop the orchestrator today. `flai guard` lets `item_edit` with `cost_of_delay` through for a story ID alone (`setsCostInputs` and `setsCostInputsOnly` test `isStory`), so an epic is refused whatever the permissions. flai itself, `Repo.OrchestratorSetsCostOfDelay` in `flai/internal/workitem/promotable.go`, refuses a story that is not in the backlog and one that has inputs already, or whose epic has. The proposal, recorded in a new ADR refining ADR-0119 (T-2241): under `plan_backlog_stories` the orchestrator may add or change the cost of delay inputs of an epic, or of a story in any state but done or cancelled, archived ones left out, still inputs only, never a value, never a removal, and nothing else in the call; the convention says it does so only for a figure the operator decided on a thread, or the planner's recommendation on a story planner's thread for a backlog story that has none, as today. On TH-0694 the orchestrator recommends instead a permission of its own, `set_cost_of_delay`, off by default; the operator's answer is pending there, and the tasks are written for the planner's option until it comes.

**Where each touch came from.**

- Layout, the guard (T-2243): `flai/internal/guard/guard.go` (`setsCostInputs`, `setsCostInputsOnly`, the `editsMCP` and `editsCLI` refusal texts), `flai/internal/guard/orchestrate_test.go` (refuses `E-0001` with `cost_of_delay` today), `flai/internal/guard/guard_test.go` (the CLI form).
- Layout, flai's own rule (T-2242): `flai/internal/workitem/promotable.go` (`OrchestratorSetsCostOfDelay`), `flai/internal/mcpserver/items_write.go` and `items_write_test.go` (`TestTheOrchestratorGivesABacklogStoryItsFirstCostOfDelayInputs` asserts the in-progress and epic refusals), `flai/cmd/edit.go` (`orchestratorSetsCostOfDelay`), `flai/cmd/edit_test.go`, `flai/cmd/orchestrate_permissions_test.go` (`S-0008 is ready, not in the backlog`).
- Layout, the permission's description (T-2244): `flai/internal/manifest/settings.go` (the risk text), `docs/users/flai-reference.md` (generated from it by `make flai-reference`).
- Design (T-2241, T-2244): `design/system/strategic-agents.md` (What it does with each permission, A story planner's threads, Its permissions and the guard), `design/system/flai-cli.md` (the `flai edit` and `flai mcp` rows), `design/system/project-manifest.md` (the permission's row), `design/adrs/README.md` (the index of the new ADR), `design/adrs` (see the folder touch below).
- Design, the convention and the prompt (T-2244): `design/conventions/strategic-agents.md` and `template/root/design/conventions/strategic-agents.md` (the cost of delay rules under As the orchestrator and Asking the operator), `design/workflows/roles/orchestrator.md` and `template/root/design/workflows/roles/orchestrator.md` (the orchestrator's prompt as S-0414 placed it, rules 3 and 10, added after S-0414's acceptance during this plan), `template/CHANGELOG.md` (the convention and role file change), `flai/internal/harness/harness.go` (the same prompt text until S-0419 reads it from the role file) and `harness_test.go`.
- Design, the guides (T-2244): `docs/users/flai.md` (When the orchestrator plans stories), `docs/users/flaiover.md`, `docs/operators/settings.md`, `docs/operators/index.md` (the permission tables).
- Declared by the criteria (T-2245): `design/issues/I-0187-the-orchestrator-cannot-set-a-cost-of-delay-input-the-operator-decided-on-an-epic-or-a-story-past-the-backlog.md` and `design/issues/summary.md`, which `flai issue close` rewrites.

**Folder touch kept.** `design/adrs`: the story adds one ADR, whose number `flai adr new` allots when the story's agent runs it, so no task can name the file yet. Ready stories S-0416, S-0421, and S-0423 keep the same folder touch for the same reason.

**Forecast.** `flai forecast` gave 44m: the median 89 s per unit of size over 98 done improvement stories on claude-opus-5-5 in the large band, times size 29 (2 criteria, 27 touches). Raised to 1h15m: the five tasks run in four serial layers, T-2242 and T-2243 rewrite pinned tests in three packages, and T-2244 rewrites sixteen documents, the generated reference, the two role files, and the template changelog among them. Delivery is flai's 2026-10-29T04:12Z, 34th in the pull order behind the E-0024, E-0025, and E-0029 stories, moved by the 31m added.

**Cost of delay.** 75.00 USD a week, as `flai cod` works it out from the orchestrator's input on TH-0694, `time_lost_per_cycle: 30m` at 150 USD an hour, one 168h cycle a week. The figure stands: each input decided on an epic or a started story goes back to the operator to run by hand, and the epic's drafts stay incomplete meanwhile.
