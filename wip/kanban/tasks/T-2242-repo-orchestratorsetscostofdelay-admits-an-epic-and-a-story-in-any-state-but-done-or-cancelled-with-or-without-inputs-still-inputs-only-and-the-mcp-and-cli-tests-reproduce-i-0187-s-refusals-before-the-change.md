---
id: T-2242
type: task
nature: improvement
title: Repo.OrchestratorSetsCostOfDelay admits an epic and a story in any state but done or cancelled, with or without inputs, still inputs only, and the MCP and CLI tests reproduce I-0187's refusals before the change
status: done
parent: S-0523
owner: alex
created: 2026-10-28T14:31:58Z
updated: 2026-10-30T14:03:44Z
transitions:
  - to: ready
    at: 2026-10-30T13:07:20Z
    by: agent-S-0523
  - to: in-progress
    at: 2026-10-30T13:07:22Z
    by: agent-S-0523
  - to: done
    at: 2026-10-30T14:03:44Z
    by: agent-S-0523
stream: S-0523
tags: [flai, orchestrator, cost-of-delay]
touches: [flai/internal/workitem/promotable.go, flai/internal/mcpserver/items_write.go, flai/internal/mcpserver/items_write_test.go, flai/cmd/edit.go, flai/cmd/edit_test.go, flai/cmd/orchestrate_permissions_test.go]
usage:
  source: log
  seconds: 3382
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 115
      output: 23417
      cache_read: 7088305
      cache_write: 376236
      cost: 4.097
---
# T-2242 Repo.OrchestratorSetsCostOfDelay admits an epic and a story in any state but done or cancelled, with or without inputs, still inputs only, and the MCP and CLI tests reproduce I-0187's refusals before the change

## Work

Second layer: waits for T-2241, since the rule and its refusal texts follow the ADR's words. Runs beside the guard task, with which it shares no path.

In `flai/internal/workitem/promotable.go`, change `Repo.OrchestratorSetsCostOfDelay` as the ADR decides. Keep: the `plan_backlog_stories` permission check (`OrchestratorPermits`), the refusal of anything `besides` the inputs (a value, a removal, another field), and the refusal of a task, an archived item, and an item that is `done` or `cancelled`. Drop: the refusal of an epic, of a story not in the backlog, of a story that has inputs already, and of a story whose epic has inputs. Reword each refusal and the comment to say what the orchestrator may set: cost of delay inputs on an epic or a story that is open, done and cancelled left out, and nothing else. Keep the permission's refusal text `the orchestrator sets ... only with orchestration.permissions.plan_backlog_stories, which is off: ask the operator with thread_open on <id>`, since `cmd/orchestrate_permissions_test.go` pins its shape.

Adjust the comments of `setsCostOfDelay` in `flai/internal/mcpserver/items_write.go` and `orchestratorSetsCostOfDelay` in `flai/cmd/edit.go`, which say "a backlog story whose epic has none either"; their code stays, since both only collect `besides` and call the rule.

Tests, reproducing I-0187 first: in `flai/internal/mcpserver/items_write_test.go`, `TestTheOrchestratorGivesABacklogStoryItsFirstCostOfDelayInputs` asserts today that an in-progress story, an epic, a story with inputs, and a story whose epic has inputs are refused; turn those four cases into successes that stamp the inputs with `by` orchestrator (ADR-0080) and keep the value untouched, keep the task, the archived item, the value, the removal, the extra field, and the permission-off cases as refusals, and add a `done` story and a `cancelled` story as refusals. Rename the test to say what it now proves. In `flai/cmd/edit_test.go`, the case `S-0003 is cancelled, not in the backlog` stays a refusal with the new wording; add an epic and an in-progress story that succeed on the command line under `FLAI_ROLE=orchestrate`. In `flai/cmd/orchestrate_permissions_test.go`, `costCall(v, "S-0008")` on a ready story is refused today with `S-0008 is ready, not in the backlog`; make it pass with the permission on and stay refused with it off, and add the same pair for an epic of the fixture board.

## Done when

- `go test ./internal/workitem/... ./internal/mcpserver/... ./cmd/...` passes in `flai/`, and `scripts/flai.sh test` on the changed paths is clean.
- With `plan_backlog_stories` on, `item_edit` with `cost_of_delay` inputs alone succeeds on an epic, on a ready story, on an in-progress story, and on a story that has inputs, and `flai edit --time-lost-per-cycle 4h E-nnnn` succeeds under `FLAI_ROLE=orchestrate`, each stamping the inputs as the orchestrator's.
- A value, a removal, another field, a task, an archived item, a `done` or `cancelled` item, and the permission off are each still refused, naming why.
- The tests named above cover each case, and the refusal texts match the ADR's words.

## Notes

Drafted by planner-S-0523. The guard's half is T-2243; without it the MCP tool's call is refused before flai sees it, so the end-to-end check in `flai/cmd/orchestrate_permissions_test.go`, which runs through the guard and flai, passes only once both are done.
