---
id: S-0511
type: story
nature: remediation
title: flai guard lets a session flai serve started run the flai serve sub-commands that write the host configuration
status: backlog
owner: alex
created: 2026-10-26T16:01:48Z
updated: 2026-10-28T20:51:09Z
transitions: []
tags: [cli, guard]
topics: [cli]
touches: [flai/internal/guard/hostconfig.go, flai/internal/guard/hostconfig_test.go, flai/internal/guard/call.go, flai/internal/guard/guard.go, flai/cmd/guard.go, docs/users/flai-reference.md, docs/users/flai.md, docs/operators/index.md, docs/operators/settings.md, design/system/flai-cli.md, design/system/agent-adapters.md, design/adrs/README.md, design/adrs, design/issues/I-0210-flai-guard-lets-a-session-flai-serve-started-run-the-flai-serve-sub-commands-that-write-the-host-configuration.md, design/issues/summary.md]
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
      seconds: 881
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 29
          output: 397
          cache_read: 2782157
          cache_write: 13945
          cost: 0.614
cost_of_delay:
  inputs:
    penalty_per_week: 500
    by: orchestrator
    at: 2026-10-28T08:21:26Z
  value: 500
  by: planner-S-0511
  at: 2026-10-28T08:22:46Z
forecast:
  duration: 45m
  delivery: 2026-10-29T05:24:00Z
  basis: "Its own forecast of 45m; 24th in the pull order with an in-progress limit of 5, behind S-0416, S-0530, S-0419, S-0423, S-0418, S-0421, S-0422, S-0420, S-0417, S-0241, S-0406, S-0407, S-0408, S-0409, S-0410, S-0411, S-0412, S-0424, S-0425, S-0506, S-0507, S-0508, S-0509 and S-0510."
  by: flai
  at: 2026-10-28T20:51:09Z
finalized:
  by: alex
  at: 2026-10-28T12:56:42Z
---
# S-0511 flai guard lets a session flai serve started run the flai serve sub-commands that write the host configuration

## Goal

This story remediates [I-0210](../../../design/issues/I-0210-flai-guard-lets-a-session-flai-serve-started-run-the-flai-serve-sub-commands-that-write-the-host-configuration.md), "flai guard lets a session flai serve started run the flai serve sub-commands that write the host configuration". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0210 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0210 is closed with `flai issue close I-0210 --reason` saying what fixed it

## Tasks
- T-2184 An ADR records that flai guard refuses every served session and sub-agent the flai serve sub-commands that write the host configuration or turn a host action on or off
- T-2185 flai guard refuses a served session and a sub-agent the flai serve, flai config, and flai template commands that write the host configuration, with a test across every session kind
- T-2186 flai guard's help names the host configuration rule and the reference is regenerated
- T-2187 The user, operator, and design documents say the host configuration is the operator's and which flai serve sub-commands the guard refuses an agent
- T-2188 I-0210 is closed with the reason naming the ADR and the guard rule, and the criteria are ticked

## Notes

### Planning

Written by planner-S-0511 on 2026-10-28; the plan is on TH-0674. The story declared no touches, so `flai touches suggest` had nothing to start from; every touch below comes from the issue's instance, the code layout, and the precedents the guard already has for a rule that binds every served session, `shared` in `flai/internal/guard/shared.go` (S-0295) and `forgePush` in `flai/internal/guard/guard.go` (S-0393).

Touches and where each came from:

- `flai/internal/guard/hostconfig.go`, `flai/internal/guard/hostconfig_test.go` (layout): the new rule and its test, a file of their own as `shared.go` is.
- `flai/internal/guard/call.go` (layout): `Decide` runs the rule for every caller beside `shared` and `forgePush`.
- `flai/internal/guard/guard.go` (layout): the package comment lists each rule.
- `flai/cmd/guard.go`, `docs/users/flai-reference.md` (layout, co-change in S-0369 and S-0295): the command's help names the rule, and the reference is generated from it with `make flai-reference`.
- `docs/users/flai.md`, `docs/operators/index.md`, `docs/operators/settings.md` (design, co-change in S-0369, S-0393, S-0498): the sections that say what the guard refuses, which `flai serve` sub-commands are shell-only, and that `FLAI_CONFIG` is not passed to an agent.
- `design/system/flai-cli.md`, `design/system/agent-adapters.md` (design): the `flai guard` row of the commands table and § The guard.
- `design/adrs/README.md` (design) and the folder `design/adrs` (kept as a folder): the rule is a decision that binds every served session, as ADR-0096 and ADR-0166 record the shared paths and `flai forge push`; `flai adr new` allocates the ADR's number and slug when T-2184 runs, so no task can name the file yet. The folder touch stays on the story alone.
- `design/issues/I-0210-...md`, `design/issues/summary.md` (declared by the criteria): `flai issue close` writes both.

The template changes nothing: its `.claude/settings.json` runs the installed `flai guard` through a generic hook that names no rule, and `safety.md` already forbids changing the operator's global tool configuration.

Layers: T-2184 alone; then T-2185, T-2186, and T-2187, which wait for the ADR's number and share no path; then T-2188.

Forecast: `flai forecast` gave 26m, the median of 91 s per unit of size over 37 done remediation stories on claude-opus-5-5 in the large band times size 17 (2 criteria, 15 touches), and delivery 2026-10-29T21:08Z at 29th in the pull order. Raised to 45m: S-0369, the closest shape (a guard rule across every served session, its test, an ADR, help, and docs, in five tasks), took 33m of agent time, and this rule has more to parse, a third word and flags across eleven `flai serve` sub-commands, with a test over eight session kinds. Delivery moved by the same 19m to 2026-10-29T21:30Z.

Cost of delay: 500 USD a week, as `flai cod` works it out from the one input, `penalty_per_week: 500`, which the orchestrator set on TH-0674 from the planner's recommendation; the story has no epic and no time lost per cycle. The value stands as computed.

Holds: S-0504, in progress, touches `docs/users/flai.md`, `design/adrs/README.md`, the folder `design/adrs`, and `design/issues/summary.md`, so this story is held while S-0504 is in progress once it is ready. The documents are shared by most stories; nothing here can be narrowed to avoid it.
