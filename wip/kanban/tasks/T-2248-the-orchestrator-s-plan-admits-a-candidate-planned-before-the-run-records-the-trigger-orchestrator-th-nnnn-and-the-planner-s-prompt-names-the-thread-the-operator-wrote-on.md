---
id: T-2248
type: task
nature: improvement
title: "The orchestrator's plan admits a candidate planned before, the run records the trigger orchestrator; TH-nnnn, and the planner's prompt names the thread the operator wrote on"
status: done
parent: S-0524
owner: alex
created: 2026-10-28T14:42:09Z
updated: 2026-10-30T14:20:11Z
transitions:
  - to: ready
    at: 2026-10-30T13:33:00Z
    by: agent-S-0524
  - to: in-progress
    at: 2026-10-30T13:33:01Z
    by: agent-S-0524
  - to: done
    at: 2026-10-30T14:20:11Z
    by: agent-S-0524
stream: S-0524
tags: [flai, orchestrator, planner, mcp]
touches: [flai/cmd/plan.go, flai/cmd/plan_orchestrate_test.go, flai/cmd/orchestrate_permissions_test.go, flai/internal/mcpserver/plan.go, flai/internal/mcpserver/plan_test.go, flai/internal/serve/plan.go, flai/internal/serve/plan_test.go, flai/internal/harness/harness.go, flai/internal/harness/harness_test.go, docs/users/flai-reference.md]
after: [T-2247]
usage:
  source: log
  seconds: 2830
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 122
      output: 580
      cache_read: 7326160
      cache_write: 650247
      cost: 2.7613
---
# T-2248 The orchestrator's plan admits a candidate planned before, the run records the trigger orchestrator; TH-nnnn, and the planner's prompt names the thread the operator wrote on

## Work

Second layer: waits for T-2247, which shares `flai/cmd/plan.go` with it and defines the candidate's `thread` this task carries to the run.

Carry the third kind of candidate through the orchestrator's asking to the planner's run:

- `orchestratorPlans` in `flai/cmd/plan.go` admits an item the candidates list under the new kind as it admits the others: an epic under `plan_backlog_epics`, a story under `plan_backlog_stories`. Its refusal text for an item not listed stays. From the admitted candidate, take its `thread` and pass it to the start.
- `serve.PlanForOrchestrator` in `flai/internal/serve/plan.go` takes the thread, or none, and the run's trigger becomes `orchestrator; TH-nnnn` when one is given, `orchestrator` otherwise, through `planFor`: `serve/agents.json` and the `- Trigger:` line of the activity entry in `wip/agents/planner.md` show it, as the replanner's `accepted S-0210; reordered` does. Keep `TriggerOrchestrator` as the constant the two readers compare against: a run whose trigger starts with it is the orchestrator's.
- The planner's prompt (`planPrompt` in `flai/internal/harness/harness.go`) says, when the run's trigger names a thread, that the item was planned before and is planned again because the operator wrote on that thread since: name the thread and tell the planner to read it with `thread_get` and take in what the operator asked for, drafting the additions and revisiting the rest as the convention says. Today the prompt says the operator asked for the item to be planned; keep that for the other triggers. The prompt's other text is unchanged.
- The MCP tool `plan` in `flai/internal/mcpserver/plan.go` changes no argument: `mcpPlan` in `flai/cmd/plan.go` reads the thread from the candidate it admitted. Its answer names the trigger it recorded.
- `flai guard` is unchanged: it holds the orchestrator's `plan` to the permission by the ID's type and leaves the list to flai (`plansCandidates` in `flai/internal/guard/guard.go`).

Tests:

- `flai/cmd/plan_orchestrate_test.go` and `flai/cmd/orchestrate_permissions_test.go`: under `FLAI_ROLE=orchestrate` with `plan_backlog_epics` on, `flai plan E-nnnn` for an epic listed as planned before and written on starts the run with the trigger `orchestrator; TH-nnnn`; the same epic with the permission off is refused; a story listed the same way under `plan_backlog_stories` likewise.
- `flai/internal/mcpserver/plan_test.go`: the tool's answer for such an item names the trigger with the thread.
- `flai/internal/serve/plan_test.go`: `PlanForOrchestrator` with a thread records `orchestrator; TH-nnnn`, without one `orchestrator`.
- `flai/internal/harness/harness_test.go`: the planner prompt for a run triggered by `orchestrator; TH-0579` names the thread and says the item was planned before; the prompt for `asked` is unchanged.

## Done when

- The orchestrator's `flai plan` and MCP `plan` start the planner on a candidate of the third kind, the run's trigger is `orchestrator; TH-nnnn`, and the planner's prompt names the thread.
- The tests above pass with `go test ./cmd/... ./internal/mcpserver/... ./internal/serve/... ./internal/harness/...` in `flai/`, and `flai test` on the changed paths passes.

## Notes

The trigger's form, `orchestrator; TH-nnnn`, follows the replanner's semicolon-joined triggers in design/system/strategic-agents.md § What a run records, so `wip/agents/planner.md` needs no new line.
