---
id: S-0680
type: story
nature: improvement
title: flai guard refuses the orchestrator flai serve status, a read, because its reads do not list serve status
status: backlog
owner: alex
created: 2026-11-10T08:34:57Z
updated: 2026-11-10T08:47:25Z
transitions: []
tags: [flai, guard]
topics: [guard, orchestration]
touches: [flai/internal/guard/guard_tables.go, flai/internal/guard/guard.go, flai/internal/guard/hostconfig.go, flai/internal/guard/guard_test.go, flai/internal/guard/hostconfig_test.go, flai/internal/guard/orchestrate_test.go, design/system/strategic-agents.md, docs/users/flai.md, design/issues/I-0318-flai-guard-refuses-the-orchestrator-flai-serve-status-a-read-because-its-reads.md, design/issues/summary.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
draft: true
cost_of_delay:
  inputs:
    time_lost_per_cycle: 2m
    by: flai
    at: 2026-11-10T08:34:57Z
  value: 5
  by: planner-S-0680
  at: 2026-11-10T08:47:25Z
forecast:
  duration: 45m
  delivery: 2026-11-10T12:19:00Z
  basis: "flai forecast's 22m (107 s per unit of size over 30 done improvement stories, size 12) raised to 45m for the guard's exhaustive read and write table test and the integration and smoke tiers a close-out runs for a story touching flai/; delivery moved by the same 23m."
  by: planner-S-0680
  at: 2026-11-10T08:47:25Z
---
# S-0680 flai guard refuses the orchestrator flai serve status, a read, because its reads do not list serve status

## Goal

This story remediates [I-0318](../../../design/issues/I-0318-flai-guard-refuses-the-orchestrator-flai-serve-status-a-read-because-its-reads.md), "flai guard refuses the orchestrator flai serve status, a read, because its reads do not list serve status". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0318 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0318 is closed with `flai issue close I-0318 --reason` saying what fixed it

## Tasks
- T-2901 The guard counts the read forms of flai serve among the reads
- T-2902 The design and the user guide say the read forms of flai serve are open to every agent
- T-2903 Close I-0318 with what fixed it

## Notes

Cost of delay inputs set by flai from I-0318. time_lost_per_cycle 2m: 2m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-11-09T09:06:05Z, 1 day before this story; under one cycle counts as one).

### Planning

Proposed fix: add `serve` to `flagReads` in `flai/internal/guard/guard_tables.go` with a form function that passes only the read forms ADR-0220 lists (`status`, `actions`, `journal`, and the show and list forms), so the guard agrees with `askForHost` in `flai/internal/guard/hostconfig.go`, which already tells every agent those forms are open. The forms that start, stop, or write stay refused.

Touches, each a file:

| Touch | From |
|-------|------|
| `flai/internal/guard/guard_tables.go` | layout: `cliReads` and `flagReads`, where the `serve` entry goes |
| `flai/internal/guard/guard.go` | co-change (75% with the two tables) and layout: `reads`, which consults them |
| `flai/internal/guard/hostconfig.go` | design: `askForHost`, `hostValues`, and `writesHost`, which the form function reuses |
| `flai/internal/guard/hostconfig_test.go` | co-change (75%), and named by I-0318: its `reads` list |
| `flai/internal/guard/orchestrate_test.go` | named by I-0318: `TestARedirectionNeitherAllowsNorRefusesTheOrchestratorsCommand` pins the refusal |
| `flai/internal/guard/guard_test.go` | layout: the sub-agent read tests |
| `design/system/strategic-agents.md` | design: § Its permissions and the guard says `flai serve status` is still refused (I-0318) |
| `docs/users/flai.md` | design: what a sub-agent may run, per documentation.md |
| `design/issues/I-0318-…-reads.md`, `design/issues/summary.md` | criterion 2: `flai issue close` writes both |

`flai touches suggest` also listed `flai/internal/guard/call.go` (50%). I left it out: it parses the call line, which S-0613 already fixed for redirections. No folder touch. No `after`: the story reads nothing that an open story writes.

Forecast: flai gave 22m from 107 s per unit of size times size 12. I raised it to 45m because T-2901's table test covers every read and write form of `flai serve` for two kinds of agent, and a close-out of a story touching `flai/` runs the integration and smoke tiers. Delivery moves by the same 23m, to 2026-11-10T12:19:00Z.

Cost of delay: 5 USD a week, as `flai cod` worked it out from the operator's input of 2m lost per 168h cycle at 150 USD an hour. Kept as it stands.
