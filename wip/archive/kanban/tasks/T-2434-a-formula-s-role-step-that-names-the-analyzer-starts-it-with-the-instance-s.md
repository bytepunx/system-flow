---
id: T-2434
type: task
nature: improvement
title: A formula's role step that names the analyzer starts it with the instance's focus on the instance, and its end is the agent event on the instance
status: done
parent: S-0424
owner: alex
created: 2026-11-02T06:20:56Z
updated: 2026-11-02T06:34:18Z
transitions:
  - to: ready
    at: 2026-11-02T06:21:17Z
    by: agent-S-0424
  - to: in-progress
    at: 2026-11-02T06:21:18Z
    by: agent-S-0424
  - to: done
    at: 2026-11-02T06:34:18Z
    by: agent-S-0424
stream: S-0424
tags: []
touches: [flai/cmd/role.go, flai/cmd/role_test.go, flai/internal/serve/role.go, flai/internal/serve/role_test.go, flai/internal/serve/analyze.go, flai/internal/serve/analyze_test.go, flai/internal/serve/agents.go, flai/internal/serve/agent_events.go, flai/internal/formulas/formulas.go, flai/internal/formulas/formulas_test.go, design/adrs/README.md, design/adrs/0245-a-formula-s-role-step-that-names-the-built-in-analyzer-starts-it-through-the.md]
after: [T-1786]
usage:
  source: log
  seconds: 780
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 117
      output: 34218
      cache_read: 8709089
      cache_write: 257956
      cost: 4.0106
---
# T-2434 A formula's role step that names the analyzer starts it with the instance's focus on the instance, and its end is the agent event on the instance

## Work

The `analysis-cadence` formula's first step is `role: analyzer`, and `flai check` accepts it, because `roleNames` counts the built-in role files. But the step cannot be worked. `serve.RunRole` refuses a built-in role (`definedRole` in `flai/internal/serve/role.go`), and `flai analyze` and the MCP `analyze` tool refuse an agent `flai serve` started. TH-0588 assumed S-0422's `role_run` would start the analyzer for such a step.

- `runTrigger` (`flai/cmd/role.go`) and `serve.RunRole`: when the agent of a formula's instance asks, on that instance, for the built-in `analyzer` that a role step of its formula names, start it through the analyzer's own start (`serve.Analyze`, with its checks and the `analyze` host action). The focus is the instance's `focus` input, empty for all three, and the trigger is `FormulaTrigger`. The operator's own `flai role run analyzer` stays refused, naming `flai analyze`.
- Give that analyzer run the instance as its item, so that its end is the `agent` event on the instance, `failed` or `ended`, as a defined role's is. The instance's agent waits for that event (ADR-0239).
- `flai check` (`formulas.file`, `flai/internal/formulas/formulas.go`) reports a role step that names a built-in role a role step cannot start: the planner, the orchestrator, or the curator.
- Record the decision with `flai adr new`, refining ADR-0239, without `--commit`.

Third layer: waits for T-1786, whose formula names the analyzer. Shares no path with T-1789.

## Done when

`role_test.go` in `cmd` and `serve`, `analyze_test.go`, and `formulas_test.go` pass, their existing cases unchanged but for the refusal's words, with cases for:

- the instance's agent starting the analyzer with the focus, on the instance;
- the run's end event on the instance;
- the operator refused;
- a role step naming the curator reported.

`flai check --strict` reports nothing.

## Notes

Added by agent-S-0424 on 2026-11-02, after T-1786's review found the gap.
