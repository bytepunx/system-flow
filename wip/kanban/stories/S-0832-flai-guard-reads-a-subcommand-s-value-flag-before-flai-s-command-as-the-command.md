---
id: S-0832
type: story
nature: remediation
title: flai guard reads a subcommand's value flag before flai's command as the command, so a rule can miss the command flai runs
status: in-progress
owner: alex
created: 2026-11-17T08:17:26Z
updated: 2026-11-17T10:19:25Z
transitions:
  - to: ready
    at: 2026-11-17T10:04:07Z
    by: orchestrator
  - to: in-progress
    at: 2026-11-17T10:19:25Z
    by: agent-S-0832
tags: [flai, guard]
touches: [flai/internal/guard/guard_refusals.go, flai/internal/guard/guard_rules.go, flai/internal/guard/authors.go, flai/internal/guard/guard_tables.go, flai/internal/guard/guard.go, flai/internal/guard/roles.go, flai/internal/guard/guard_test.go, flai/cmd/guard.go, flai/cmd/guard_test.go, design/system/flai-cli.md, design/issues/I-0358-flai-guard-reads-a-subcommand-s-value-flag-before-flai-s-command-as-the-command.md, design/issues/summary.md]
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
      seconds: 5
      estimated: true
      models:
        - model: claude-haiku-5-5
          input: 2
          output: 10
          cache_read: 186612
          cache_write: 2504
          cost: 0.0129
        - model: claude-opus-5-5
          input: 11
          output: 1051
          cache_read: 706887
          cache_write: 9979
          cost: 0.2437
    - kind: orchestrator
      seconds: 21
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 5
          output: 63
          cache_read: 242900
          cache_write: 44342
          cost: 0.0894
        - model: claude-sonnet-5-5
          input: 2
          output: 268
          cache_read: 14873
          cache_write: 15254
          cost: 0.0244
cost_of_delay:
  inputs:
    time_lost_per_cycle: 10m
    by: flai
    at: 2026-11-17T08:17:26Z
  value: 25
  by: planner-S-0832
  at: 2026-11-17T08:30:33Z
forecast:
  duration: 40m
  delivery: 2026-11-17T11:07:00Z
  basis: "Its own forecast of 40m; 2nd in the pull order with an in-progress limit of 3, behind S-0754, S-0806, S-0823 and S-0811."
  by: flai
  at: 2026-11-17T10:14:13Z
finalized:
  by: orchestrator
  at: 2026-11-17T08:33:54Z
---
# S-0832 flai guard reads a subcommand's value flag before flai's command as the command, so a rule can miss the command flai runs

## Goal

This story remediates [I-0358](../../../design/issues/I-0358-flai-guard-reads-a-subcommand-s-value-flag-before-flai-s-command-as-the-command.md), "flai guard reads a subcommand's value flag before flai's command as the command, so a rule can miss the command flai runs". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0358 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0358 is closed with `flai issue close I-0358 --reason` saying what fixed it

## Tasks
- T-3530 flai guard reads the command and subcommand after flai as cobra finds them, so a value flag before the command takes its value with it
- T-3532 flai guard takes the flags that take no value from flai's command tree, and a test holds its reading of the command to cobra's Find
- T-3537 flai-cli.md says how flai guard reads flai's command before its flags, and I-0358 is closed

## Notes

Cost of delay inputs set by flai from I-0358. time_lost_per_cycle 10m: 10m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-11-13T08:42:15Z, 4 days before this story; under one cycle counts as one).

### Planning

Proposed fix, from I-0358's one instance: the guard stops listing the flags that take a value before flai's command (`flaiValues`, `authorValues`) and reads the command as cobra's `Find` does, a flag before it taking the next word unless flai's command tree defines it with no value. T-3530 builds the reading in `flai/internal/guard`, T-3532 fills its table from flai's command tree and tests it against cobra itself, and T-3537 documents it and closes I-0358.

Touches, all files, no folder:

| Touch | Source |
|-------|--------|
| `flai/internal/guard/guard_refusals.go` | design: I-0358 names `subcommands()` there |
| `flai/internal/guard/guard_rules.go` | layout: `flaiValues` is defined there |
| `flai/internal/guard/authors.go` | layout: `authorValues` and `names` read the command with the same gap, and lack `--wait` |
| `flai/internal/guard/guard_tables.go` | layout: `readsRoles` and `hostAPIMethod` read the command with `flaiValues` |
| `flai/internal/guard/guard.go` | co-change (66% with the two guard files) and layout: the `Guard` struct and the planner's moves |
| `flai/internal/guard/roles.go` | layout: a role's plans read the command by position with `flaiValues` |
| `flai/internal/guard/guard_test.go` | co-change (43%): the tests that reproduce I-0358 |
| `flai/cmd/guard.go` | layout: builds the `Guard` from the command tree |
| `flai/cmd/guard_test.go` | co-change (36%): the test against cobra's `Find` |
| `design/system/flai-cli.md` | co-change (23%): the `flai guard` row says how the guard reads |
| I-0358's file and `design/issues/summary.md` | criterion 2: `flai issue close` writes both |

The story declared no touches. `docs/users/flai-reference.md` (52% co-change) is left out: it is generated from the command tree, and no flag or command changes.

No `after`: the story reads nothing another open story writes.

Forecast: `flai forecast` gave 28m from 14 units of size (2 criteria, 12 touches); set to 40m because the fix is proposed first and T-3532's test walks flai's whole command tree, work the size does not count. Delivery moves by the same 12m, to 2026-11-17T16:29:00Z.

Cost of delay: 25 USD a week, as `flai cod` works it out from the operator's input of 10m lost per 168h cycle at 150 USD an hour. It stands: the one instance cost 10m and no other is recorded.
