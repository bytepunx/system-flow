---
id: S-0661
type: story
nature: feature
title: A story's agent runs the project's checks against its story's kluster namespace before review, and its own commands with the namespace's variables
status: ready
parent: E-0035
owner: alex
created: 2026-11-10T08:09:09Z
updated: 2026-11-13T07:30:41Z
transitions:
  - to: ready
    at: 2026-11-13T07:26:16Z
    by: alex
tags: [cli]
topics: [kluster, deployment]
touches: [flai/internal/serve/kluster.go, flai/internal/serve/kluster_test.go, flai/internal/serve/checks.go, flai/internal/serve/checks_test.go, flai/internal/kluster/request.go, flai/internal/kluster/request_test.go, flai/cmd/kluster.go, flai/cmd/kluster_test.go, flai/cmd/checks.go, flai/cmd/checks_test.go, flai/internal/mcpserver/kluster.go, flai/internal/mcpserver/kluster_test.go, flai/internal/mcpserver/folder.go, flai/internal/mcpserver/server_test.go, flai/internal/guard/guard_tables.go, flai/internal/guard/guard_rules.go, flai/internal/guard/kluster_test.go, design/system/kluster.md, design/system/flai-cli.md, docs/users/flai.md, docs/users/flai-reference.md]
after: [S-0660, S-0659]
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
    - kind: planner
      seconds: 20
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 10
          output: 3472
          cache_read: 354925
          cache_write: 26353
          cost: 0.3513
    - kind: orchestrator
      seconds: 257
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 30
          output: 213
          cache_read: 1484011
          cache_write: 120720
          cost: 0.4388
        - model: claude-sonnet-5-5
          input: 8
          output: 851
          cache_read: 77831
          cache_write: 27685
          cost: 0.0795
cost_of_delay:
  value: 986.84
  by: planner-S-0661
  at: 2026-11-10T09:41:33Z
forecast:
  duration: 1h15m
  delivery: 2026-11-13T12:49:00Z
  basis: "Its own forecast of 1h15m; 6th in the pull order with an in-progress limit of 3, behind S-0660, S-0733, S-0735, S-0745, S-0734, S-0751 and S-0740."
  by: flai
  at: 2026-11-13T07:30:41Z
finalized:
  by: orchestrator
  at: 2026-11-10T09:44:44Z
---
# S-0661 A story's agent runs the project's checks against its story's kluster namespace before review, and its own commands with the namespace's variables

## Goal

A story's agent runs the project's checks against its story's namespace before review, and runs its own commands with the namespace's variables, so that it verifies a criterion that needs the whole system and ticks it.

## Acceptance criteria

- [ ] A story's agent runs the project's checks against its own story's namespace, deployed at its branch's head, with one flai command and one MCP tool, and gets each check's result and log
- [ ] The run gets the namespace's variables, `KUBECONFIG` and `FLAI_KLUSTER_NAMESPACE` among them, as a checks run in review does, and is refused when the namespace is not deployed at the head, saying to deploy first
- [ ] The agent gets the same variables for a command of its own, such as a request to a service in its namespace, in the way the ADR decides
- [ ] `flai guard` refuses an agent's command against another story's namespace or the default namespace, as far as the ADR says it can tell
- [ ] Checks runs against the cluster still run one at a time per project, after the job queue
- [ ] A failed run asked for by the agent opens no thread to the operator
- [ ] The run is recorded with its commit, its outcome, and who asked, where the next story reads it
- [ ] `design/system/kluster.md`, `docs/users/flai.md`, and its reference say all of it

## Tasks
- T-2863 flai serve's kluster pass runs the checks a story's agent asks for against its namespace at the head, after the job queue, records who asked, and opens no thread on a failure
- T-2864 flai guard refuses a story's agent a command against another story's namespace or the default namespace
- T-2865 flai kluster env gives a story's agent its own namespace's variables for a command of its own, refusing one not deployed at the head
- T-2866 A story's agent asks flai serve for a checks run against its namespace with one flai command and one MCP tool, and gets each check's result and log
- T-2867 kluster.md, flai-cli.md, flai.md, and the generated reference say how a story's agent checks its namespace and runs its own commands there

## Notes

- `KlusterEnv` in `flai/internal/serve/kluster.go` decides the variables for every checks run today.

### Planning

Touches, by where each came from:

| Touch | From |
|-------|------|
| `flai/internal/serve/kluster.go`, `kluster_test.go` | Declared in the Notes (`KlusterEnv`); the layout: the kluster pass's queue, `checks`, and `p.report` are there |
| `flai/internal/serve/checks.go`, `checks_test.go` | Layout: `ChecksRun` and `RunChecksIn`, where the record gains who asked |
| `flai/internal/kluster/request.go`, `request_test.go` | S-0660's plan: its deploy request lives there, and the checks request is a second kind of it |
| `flai/cmd/kluster.go`, `kluster_test.go` | Layout: the `flai kluster` commands, for `flai kluster env` |
| `flai/cmd/checks.go`, `checks_test.go` | Layout: `flai checks run`, for its `--job` request |
| `flai/internal/mcpserver/kluster.go`, `kluster_test.go`, `folder.go`, `server_test.go` | S-0660's plan: its deploy tool lives in `kluster.go`, registered in `folder.go` and listed in `TestToolsAreAdvertised` |
| `flai/internal/guard/guard_rules.go`, `guard_tables.go`, `kluster_test.go` | Design (kluster.md § Commands, the guard) and layout: `klusterWrites`, the MCP tables, and the guard's kluster tests |
| `design/system/kluster.md`, `docs/users/flai.md`, `docs/users/flai-reference.md` | Criterion 8 |
| `design/system/flai-cli.md` | Co-change: 77% of commits that change the seed paths also change it |

No folder touch.

`after`:

- S-0660 (declared): T-2863 and T-2866 build on its deploy request in `flai/internal/kluster/request.go` and its MCP tool in `flai/internal/mcpserver/kluster.go`: the hand-off to flai serve, its refusals, and its answer. They refuse with the deploy command it builds.
- S-0659: criteria 3 and 4 build "in the way the ADR decides". T-2864 reads the ADR's guard rule, and T-2865 reads how the agent gets the variables.

Forecast: `flai forecast` gives 53m from size 29 (8 criteria, 21 touches). Raised to 1h15m: S-0578 and S-0579, the nearest kluster stories, took 61m and 51m of agent time, and this one adds a guard rule and an MCP tool to a serve pass. The delivery, 2026-11-10T16:49:00Z, is the start flai's play-out of the pull order gives (11th, in-progress limit 5) plus 1h15m.

Cost of delay: 986.84 USD a week, as `flai cod S-0661` works it out: S-0661's share of the operator's `revenue_per_week: 1250` on E-0035 (TH-0980, after the orchestrator's recommendation from D-0479), split by forecast duration between the epic's two open stories without inputs of their own, S-0661 (1h15m) and S-0659. Not adjusted: the split follows the operator's input and flai's rule, and S-0660 and S-0662, which carry their own inputs, take no share of it.
