---
id: T-3392
type: task
nature: remediation
title: Replace an image-held run with its start in one write, so it is never kept among the past runs
status: done
parent: S-0795
owner: alex
created: 2026-11-16T00:40:58Z
updated: 2026-11-16T00:52:10Z
transitions:
  - to: ready
    at: 2026-11-16T00:45:07Z
    by: agent-S-0795
  - to: in-progress
    at: 2026-11-16T00:45:14Z
    by: agent-S-0795
  - to: done
    at: 2026-11-16T00:52:10Z
    by: agent-S-0795
stream: S-0795
tags: [flai, serve]
touches: [flai/internal/serve/agents_state.go, flai/internal/serve/agents_start.go, flai/internal/serve/agents_test.go]
usage:
  source: log
  seconds: 416
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 15
      output: 2786
      cache_read: 1544888
      cache_write: 29782
      cost: 0.603
---
# T-3392 Replace an image-held run with its start in one write, so it is never kept among the past runs

## Work

Confirm the cause, then fix it. As read at planning: when a start held for its image job goes on, `goOn` in `flai/internal/serve/agents_start.go` calls `h.again`, whose start records the new run with `AgentState.put`. `put` uses `keep` in `flai/internal/serve/agents_state.go`, and `keep` pushes the held run (PID 0) into `PastPlans` because `same` (in `flai/internal/serve/agents.go`) finds the PIDs differ. Only a second `updateAgent` then calls `dropPast` to take it out. Between the two writes, a reader sees the new run live and the held run among the past runs. `TestContainerStrategicStartUnderALayerWaitsForTheHarnessImageAlone` waits for the new run's PID, then reads `PastPlans`. Under load it lands in that window and fails at `container_test.go:1618`. The same holds for the orchestrator's, the analyzer's, the curator's, and a defined role's past runs.

Make the replacement one write. Recommended: `keep` leaves out a `was` that is `imageHeld()`, since such a run never ran. The stop and failure paths clear `ImageHeld` before they write, so they still keep their runs. Keep or remove `dropPast` as the change needs, and say why in the narrative's Decisions. A different fix that closes the window in the product code is fine; a fix to the test alone is not, since the dashboard and `agent.status` read the same state.

Add a unit test in `flai/internal/serve/agents_test.go` beside `TestEveryRolesEarlierRunsAreKept`. It puts an image-held run for each role kind, then puts its start (same item and `Started`, a PID), and finds no past run after that one `put`. It must fail on the code as it was.

Waits for no task.

## Done when

- The new unit test fails without the fix and passes with it.
- No path records an image-held run among a role's past runs, in any one write.
- `flai test` on the changed paths passes.

## Notes
