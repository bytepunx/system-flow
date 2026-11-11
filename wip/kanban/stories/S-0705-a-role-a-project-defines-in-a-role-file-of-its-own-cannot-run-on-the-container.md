---
id: S-0705
type: story
nature: improvement
title: A role a project defines in a role file of its own cannot run on the container harness, since the start request does not carry its file to pilot
status: backlog
owner: alex
created: 2026-11-10T09:11:33Z
updated: 2026-11-11T01:45:02Z
transitions: []
tags: [flai]
topics: [harness, docker, code]
touches: [flai/internal/harness/harness.go, flai/internal/harness/container.go, flai/internal/harness/container_test.go, flai/internal/pilot/request.go, flai/internal/pilot/request_test.go, flai/internal/pilot/run.go, flai/internal/pilot/run_test.go, flai/internal/serve/container.go, flai/internal/serve/container_test.go, design/system/strategic-agents.md, design/system/flai-cli.md, design/system/agent-adapters.md, docs/users/flai.md, docs/operators/index.md, design/issues/I-0258-a-role-a-project-defines-in-a-role-file-of-its-own-cannot-run-on-the-container.md, design/issues/summary.md]
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
        - model: claude-haiku-5-5
          input: 72
          output: 323
          cache_read: 4891802
          cache_write: 140025
          cost: 0.339
        - model: claude-opus-5-5
          input: 68
          output: 8452
          cache_read: 5066414
          cache_write: 92957
          cost: 1.9533
    - kind: orchestrator
      seconds: 6
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 1
          output: 2
          cache_read: 102139
          cache_write: 491
          cost: 0.0281
cost_of_delay:
  inputs:
    time_lost_per_cycle: 5m
    by: flai
    at: 2026-11-10T09:11:33Z
  value: 12.5
  by: planner-S-0705
  at: 2026-11-10T10:29:53Z
forecast:
  duration: 2h
  delivery: 2026-11-11T19:05:00Z
  basis: "Its own forecast of 2h; 27th in the pull order with an in-progress limit of 3, behind S-0674, S-0676, S-0688, S-0690, S-0691, S-0678, S-0686, S-0681, S-0718, S-0241, S-0659, S-0660, S-0661, S-0662, S-0663, S-0683, S-0689, S-0692, S-0693, S-0694, S-0695, S-0696, S-0697, S-0698, S-0699, S-0700, S-0701, S-0702 and S-0703."
  by: flai
  at: 2026-11-11T01:45:02Z
finalized:
  by: orchestrator
  at: 2026-11-10T11:04:44Z
---
# S-0705 A role a project defines in a role file of its own cannot run on the container harness, since the start request does not carry its file to pilot

## Goal

This story remediates [I-0258](../../../design/issues/I-0258-a-role-a-project-defines-in-a-role-file-of-its-own-cannot-run-on-the-container.md), "A role a project defines in a role file of its own cannot run on the container harness, since the start request does not carry its file to pilot". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0258 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0258 is closed with `flai issue close I-0258 --reason` saying what fixed it

## Tasks
- T-2973 The container adapter carries a defined role's name and inputs to pilot and names its container and cache by the role's name
- T-2975 pilot loads a defined role's file by its name from the mounted checkout and starts the inner harness as that role
- T-2976 flai serve readies, names, and stops a defined role's container by the role's name
- T-2979 The design and docs say a defined role runs on the container harness, and I-0258 is closed

## Notes

Cost of delay inputs set by flai from I-0258. time_lost_per_cycle 5m: 5m per occurrence × 1 occurrence ÷ 1.1 cycles of 168h (first reported 2026-11-02T19:06:43Z, 7.6 days before this story).

### Planning

Proposed solution, from the issue's one instance (S-0509) and ADR-0247 §3, which leaves this to "a story that carries its file to pilot": the start request carries the defined role's name and its inputs, never its file. pilot loads the file by that name from the mounted checkout, as it loads a built-in role's, so the inner harness starts with the definition, `FLAI_ROLE`, `FLAI_ITEM`, and start prompt it gets on the host. The container and its cache folder are named for the role's name, `flai-agent-<name>-<session>` and `.flai-cache/agents/<name>`, so two defined roles that share a verb do not share a cache. The mounts stay ADR-0247's. No new ADR is planned, since ADR-0247 §3 foresees this change.

Touches, all files, no folder:

| Touch | From |
|-------|------|
| `flai/internal/harness/harness.go`, `container.go`, `container_test.go` | layout: `Request`'s `json:"-"` fields and `containerSubject`, the refusal the issue names |
| `flai/internal/pilot/request.go`, `request_test.go`, `run.go`, `run_test.go` | layout: `builtin` in `run.go`, which loads a built-in role's file, and `Check` |
| `flai/internal/serve/container.go`, `container_test.go` | layout: `containerReady` and `AgentRun.containerSubject` name the cache and the stop by the run's verb |
| `design/system/strategic-agents.md`, `design/system/flai-cli.md`, `design/system/agent-adapters.md` | design: each says the defined role is refused on `container` |
| `docs/users/flai.md`, `docs/operators/index.md` | design: each says the same to users and operators |
| the I-0258 issue file, `design/issues/summary.md` | criteria: `flai issue close` writes both |

`flai touches suggest` named no co-change above 10%. `flai/internal/harness/harness_test.go` (6%) is left out: `task done` widens the touches if the request's change reaches it.

Forecast: flai gave 31m, 103 s a unit of size over 155 large improvement stories, times size 18. Raised to 2h: four tasks across three packages, tests with docker faked, and a close-out that runs the integration and smoke tiers for a change under `flai/`. Today's close-outs of such stories took 8 to 21 minutes alone (S-0672, S-0682, S-0669). Delivery is flai's 15:42Z, set from the pull order, moved by the added time.

Cost of delay: 12.50 USD a week, `flai cod`'s figure from the operator's input, 5m lost a 168h cycle at 150 USD an hour. It stands.

No `after`: the story reads no output of another open story.
