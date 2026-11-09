---
id: T-2764
type: task
nature: remediation
title: Test that the orchestrator sets an in-progress story's cost of delay input only while set_cost_of_delay is on
status: done
parent: S-0622
owner: alex
created: 2026-11-09T18:05:20Z
updated: 2026-11-09T22:10:29Z
transitions:
  - to: ready
    at: 2026-11-09T21:55:52Z
    by: agent-S-0622
  - to: in-progress
    at: 2026-11-09T21:55:53Z
    by: agent-S-0622
  - to: done
    at: 2026-11-09T22:10:29Z
    by: agent-S-0622
stream: S-0622
tags: [cli]
touches: [flai/cmd/orchestrate_permissions_test.go]
usage:
  source: sum
  seconds: 0
  models: []
  strategic:
    - kind: orchestrator
      seconds: 3
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 0
          output: 8
          cache_read: 26246
          cache_write: 880
          cost: 0.0072
---
# T-2764 Test that the orchestrator sets an in-progress story's cost of delay input only while set_cost_of_delay is on

## Work

Write the regression test for I-0284's second half, the case it records. S-0549 was in progress when the operator confirmed `time_lost_per_cycle` 20m on TH-0775. The orchestrator could not write it because `set_cost_of_delay` was off.

`testTheOrchestratorSetsCostOfDelayInputsOnlyWithSetCostOfDelay` in `flai/cmd/orchestrate_permissions_test.go` covers the backlog story S-0005, the ready story S-0008, and the epic E-0003, but no story in progress. The fixture already has one, S-0007. Extend that test, or add a case beside it under `TestAsTheOrchestrator`, for S-0007. Cover both `flai edit --time-lost-per-cycle` and the MCP tool `item_edit`:

- With `set_cost_of_delay` on, the input is set, `by` is `orchestrator`, and no value is set.
- With it off and `plan_backlog_stories` on, flai refuses the edit with the refusal that names `orchestration.permissions.set_cost_of_delay`, and the board is unchanged.

The task waits for nothing. It shares no path with the other test task, so the two run together.

## Done when

- The test sets S-0007's input with the permission on and is refused with it off, by Bash and by MCP.
- `flai test flai/cmd/orchestrate_permissions_test.go` passes.

## Notes
