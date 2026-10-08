---
id: S-0287
type: story
nature: remediation
title: flai guard lets a task sub-agent run flai adr new but refuses flai adr topics
status: ready
owner: alex
created: 2026-10-06T09:56:50Z
updated: 2026-10-08T04:25:41Z
transitions:
  - to: ready
    at: 2026-10-08T04:24:25Z
    by: orchestrator
tags: [cli, flai, guard]
topics: [cli, conventions, template]
touches: [flai/internal/guard/guard.go, flai/internal/guard/guard_test.go, flai/cmd/guard_test.go, design/system/agent-context.md, design/conventions/delegation.md, template/root/design/conventions/delegation.md, docs/users/flai.md, design/issues/I-0062-flai-guard-lets-a-task-sub-agent-run-flai-adr-new-but-refuses-flai-adr-topics.md, design/issues/summary.md]
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
      seconds: 212
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 62
          output: 939
          cache_read: 19442604
          cache_write: 26355
          cost: 4.7971
cost_of_delay:
  inputs:
    time_lost_per_cycle: 12m
    by: flai
    at: 2026-10-06T09:56:50Z
  value: 30
  by: planner-S-0287
  at: 2026-10-07T23:18:54Z
forecast:
  duration: 45m
  delivery: 2026-10-08T06:15:00Z
  basis: "Its own forecast of 45m; 7th in the pull order with an in-progress limit of 3, behind S-0232, S-0316, S-0324, S-0318, S-0320, S-0319, S-0309, S-0312 and S-0326."
  by: flai
  at: 2026-10-08T04:25:41Z
finalized:
  by: orchestrator
  at: 2026-10-08T04:24:21Z
---
# S-0287 flai guard lets a task sub-agent run flai adr new but refuses flai adr topics

## Goal

This story remediates [I-0062](../../../design/issues/I-0062-flai-guard-lets-a-task-sub-agent-run-flai-adr-new-but-refuses-flai-adr-topics.md), "flai guard lets a task sub-agent run flai adr new but refuses flai adr topics". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0062 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0062 is closed with `flai issue close I-0062 --reason` saying what fixed it

## Tasks
- T-1270 flai guard lets a sub-agent run flai adr new, topics, and accept, and adr_new, without a commit
- T-1271 An ADR refining ADR-0060 records that a sub-agent writes an ADR but never commits it, and agent-context.md says so
- T-1272 delegation.md, its template copy, and flai.md say a task sub-agent records an ADR with flai adr new and leaves the commit to the story's agent
- T-1273 Close I-0062 saying the guard lets a sub-agent write an ADR and keeps the commit the story's agent's

## Notes

Cost of delay inputs set by flai from I-0062. time_lost_per_cycle 12m: 4m per occurrence × 3 occurrences ÷ 1 cycle of 168h (first reported 2026-10-03T07:27:30Z, 3.1 days before this story; under one cycle counts as one).

### Planning

Proposed solution, from the three instances: an ADR is a document, not a work item, a thread, or history, which is all [ADR-0060](../../../design/adrs/0060-a-claude-code-pretooluse-hook-flai-guard-refuses-any-sub-agent-s-call-that.md) guards. Today `flai guard` refuses a sub-agent every `flai adr` command and the MCP `adr_new`, so task sub-agents copied `internal/adr` by hand (S-0249) or guessed the number (S-0220). The fix lets a sub-agent run `flai adr new` (with `--print-body` too), `flai adr topics`, and `flai adr accept`, and call `adr_new`, so long as none commits: `--commit`, `--autocommit`, and `adr_new`'s `commit` stay the story's agent's, since they write history and the story's touches. A new ADR refines ADR-0060 with this. The alternative, refusing all of them and saying so in `delegation.md`, keeps the hand-copying I-0062 records; it is named on the plan's thread.

Touches, each a file; no folder touch was kept:

| Touch | From |
|-------|------|
| `flai/internal/guard/guard.go` | layout: the sub-agent rules (`subAgent`, `cliReads`, `flagReads`, `MCPReads`) live here |
| `flai/internal/guard/guard_test.go` | co-change, 92% with `guard.go` |
| `flai/cmd/guard_test.go` | co-change, 46%: the test through the binary that reproduces I-0062 |
| `design/system/agent-context.md` | design: § Sub-agents says what the guard refuses a sub-agent |
| `design/conventions/delegation.md`, `template/root/design/conventions/delegation.md` | design: § As a task sub-agent says what a task sub-agent may not do; the template's copy changes in the same story |
| `docs/users/flai.md` | design: its sub-agent paragraph lists what a sub-agent may run |
| `design/issues/I-0062-…md`, `design/issues/summary.md` | criterion 2: `flai issue close` writes both |

The new ADR and its row in `design/adrs/README.md` are not declared: the story's agent records it with `flai adr new --commit`, which adds both to the touches.

Forecast: flai gave 5m (134 s per unit of size times size 2, from two criteria and no touches). It is raised to 45m: S-0246, the last guard remediation with a rule, tests, design, and docs, took 1845 s (31m), and this story adds an ADR and a template convention change. Delivery is flai's 2026-10-08T04:25Z, tenth in the pull order, pushed back by the 40m difference.

Cost of delay: 30 USD a week, as `flai cod` worked it out from the operator's input (12m lost per 168h cycle at 150 USD an hour). It stands: the three instances in four days match the input's rate.
