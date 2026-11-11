---
id: T-3095
type: task
nature: remediation
title: Keep a budget state under the harness's resolved program and login, and still read the old shape
status: done
parent: S-0720
owner: alex
created: 2026-11-11T08:54:35Z
updated: 2026-11-11T22:18:12Z
transitions:
  - to: ready
    at: 2026-11-11T21:57:35Z
    by: agent-S-0720
  - to: in-progress
    at: 2026-11-11T21:57:46Z
    by: agent-S-0720
  - to: done
    at: 2026-11-11T22:18:12Z
    by: agent-S-0720
stream: S-0720
tags: [flai, serve, budget]
touches: [flai/internal/serve/budget.go, flai/internal/serve/budget_test.go, flai/internal/serve/agents.go, flai/internal/serve/agents_start.go, flai/internal/serve/orchestrate.go, flai/internal/serve/agents_test.go]
usage:
  source: log
  seconds: 1226
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 68
      output: 19547
      cache_read: 5354424
      cache_write: 90532
      cost: 2.0112
---
# T-3095 Keep a budget state under the harness's resolved program and login, and still read the old shape

## Work

- Widen `budgetKey` in `flai/internal/serve/budget.go`: the provider when the run names one, else the harness with its resolved program and login, so `claude-code` on `rmhclod` and `claude-code` on `meclod` are two states. Keep the key readable, since `flai serve status`, the inbox, and the notice show it.
- Record the resolved program and login on `AgentRun` (`flai/internal/serve/agents.go`) where `startOf` resolves them (`flai/internal/serve/agents_start.go`), so `budgetEnded` keys the end of a run as its start was keyed, without threading the configuration through each end caller.
- Pass the same identity to `budgetStart` at the gate in `agents_start.go`, where the configuration is in scope.
- Read a `budget.json` in the old shape, entries keyed by the bare harness name: it is not refused or lost. Decide whether an old entry holds every program on its harness until it clears, and say which in the code's comment and the story's narrative.
- Waits for nothing: it is the first layer.

## Done when

- A run of one program ending on the limit holds starts on that program and not on another program of the same harness, in a test in `flai/internal/serve/budget_test.go`.
- A run that names a provider keys the state as before.
- A `budget.json` written in the old shape is read, in a test.
- `flai test` on the changed paths passes.

## Notes
