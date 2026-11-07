---
id: S-0568
type: story
nature: improvement
title: flai guard refuses a sub-agent every flai role subcommand, the reads list, show, and definition among them
status: done
owner: alex
created: 2026-11-02T20:31:49Z
updated: 2026-11-07T08:52:49Z
transitions:
  - to: ready
    at: 2026-11-07T08:02:15Z
    by: alex
  - to: in-progress
    at: 2026-11-07T08:19:52Z
    by: agent-S-0568
  - to: review
    at: 2026-11-07T08:50:44Z
    by: agent-S-0568
  - to: done
    at: 2026-11-07T08:52:49Z
    by: orchestrator
tags: []
touches: [flai/internal/guard/guard.go, design/system/flai-cli.md, design/issues/I-0245-flai-guard-refuses-a-sub-agent-every-flai-role-subcommand-the-reads-list-show.md, design/issues/summary.md, wip/kanban/stories/S-0568-flai-guard-refuses-a-sub-agent-every-flai-role-subcommand-the-reads-list-show.md, flai/cmd/guard_test.go, flai/internal/guard/roles_test.go, design/system/strategic-agents.md, docs/users/flai-reference.md, docs/users/flai.md, flai/cmd/guard.go, design/issues/I-0301-flai-serve-starts-a-story-s-agent-for-a-ready-story-the-board-holds-no-touches.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 1927
  priced_by: harness
  turns:
    - day: 2026-11-07
      ceremony: 6
      hand_edits: 1
      work: 103
  models:
    - model: claude-opus-5-5
      input: 282
      output: 31728
      cache_read: 32810094
      cache_write: 352316
      cost: 9.8776
  strategic:
    - kind: orchestrator
      seconds: 258
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 49
          output: 356
          cache_read: 1840617
          cache_write: 217155
          cost: 0.51
        - model: claude-sonnet-5-5
          input: 6
          output: 57
          cache_read: 40544
          cache_write: 29311
          cost: 0.0504
cost_of_delay:
  inputs:
    time_lost_per_cycle: 5m
    by: flai
    at: 2026-11-02T20:31:49Z
  value: 12.5
  by: planner-S-0568
  at: 2026-11-07T08:20:56Z
forecast:
  duration: 12m
  delivery: 2026-11-07T08:36:00Z
  basis: "flai forecast: median 88 s per unit of size over 23 done improvement stories on claude-opus-5-5 in the medium band, times size 8 (2 criteria, 6 touches), counted from the story's start; not adjusted."
  by: planner-S-0568
  at: 2026-11-07T08:20:56Z
finalized:
  by: alex
  at: 2026-11-07T08:01:38Z
---
# S-0568 flai guard refuses a sub-agent every flai role subcommand, the reads list, show, and definition among them

## Goal

This story remediates [I-0245](../../../design/issues/I-0245-flai-guard-refuses-a-sub-agent-every-flai-role-subcommand-the-reads-list-show.md), "flai guard refuses a sub-agent every flai role subcommand, the reads list, show, and definition among them". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [x] The cause I-0245 describes no longer occurs, with a test that reproduces it where one fits
- [x] I-0245 is closed with `flai issue close I-0245 --reason` saying what fixed it

## Tasks
- T-2551 flai guard lets every agent read the roles with flai role list, show, definition, and role alone, and keeps run and --write refused
- T-2553 The guard's help, the user guide, and the design say flai role list, show, and definition are reads
- T-2554 Close I-0245 saying what fixed it

## Notes

Cost of delay inputs set by flai from I-0245. time_lost_per_cycle 5m: 5m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-11-02T00:53:54Z, 0.8 days before this story; under one cycle counts as one).

### Planning

Where each predicted touch came from. The story declared none, so `flai touches suggest` ran from `flai/internal/guard/guard.go`, the file I-0245 names. The story's agent started at 2026-11-07T08:19:52Z, once these were written, and replaced the story's touches with its own, so the story's list and this table differ. Its list stands. The tasks still carry the planner's touches, for the agent to review.

| Predicted touch | Source | Why | On the story now |
|-----------------|--------|-----|------------------|
| `flai/internal/guard/guard.go` | design | I-0245 names its `cliReads` table as the cause | Yes |
| `flai/internal/guard/roles_test.go` | layout | The test that reproduces the refusal; the sibling reads have test files by topic | No: the agent names `flai/internal/guard/role_test.go` |
| `flai/cmd/guard_test.go` | co-change, 31% | The refusal through the hook's input, as I-0245 met it | No |
| `flai/cmd/guard.go` | co-change, 53% | Its help lists the reads | No |
| `docs/users/flai-reference.md` | co-change, 35% | Generated from that help | No |
| `docs/users/flai.md` | co-change, 8% | `## Sub-agents` says what a sub-agent may run | No |
| `design/system/flai-cli.md` | layout | The `flai guard` row lists the reads | Yes |
| `design/system/strategic-agents.md` | layout | `## Role files (S-0414)` holds the role commands | No |
| `design/issues/I-0245-flai-guard-refuses-a-sub-agent-every-flai-role-subcommand-the-reads-list-show.md` | design | Criterion 2 closes it | Yes |
| `design/issues/summary.md` | layout | Closing an issue regenerates it | Yes |

The agent also added this story's own file to the touches. No folder touch and no `after`: nothing here reads another open story's output.

Figures. The forecast stands at 12m from `flai forecast`: median 88 s per unit of size over 23 done improvement stories in the medium band, times size 8, the story's 2 criteria and the 6 touches it holds now. The same run over the 10 predicted touches gave 18m, so `flai task done` widening the touches back would make 12m low. Delivery, 2026-11-07T08:36:00Z, counts from the story's start. Neither figure was adjusted, because no measurement gives an amount to adjust by. The cost of delay value stands at 12.50 USD a week from `flai cod`: 5m lost per 168h cycle at 150 USD an hour. I-0245's one occurrence gives no reason to change it.

### Accepted by the orchestrator

- Verified: 7aee328fb4c94ab032c3f2c34d2aa644bf77fd67
- At: 2026-11-07T08:52:49Z

Verdict: pass at 7aee328fb4c94ab032c3f2c34d2aa644bf77fd67
- 1: flai/internal/guard/guard.go, flai/internal/guard/roles_test.go, flai/cmd/guard.go, flai/cmd/guard_test.go, design/system/flai-cli.md, design/system/strategic-agents.md, docs/users/flai-reference.md, docs/users/flai.md
- 2: design/issues/I-0245-flai-guard-refuses-a-sub-agent-every-flai-role-subcommand-the-reads-list-show.md, design/issues/summary.md
