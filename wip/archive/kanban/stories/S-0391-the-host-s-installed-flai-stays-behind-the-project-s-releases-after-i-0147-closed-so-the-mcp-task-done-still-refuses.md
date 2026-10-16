---
id: S-0391
type: story
nature: improvement
title: The host's installed flai stays behind the project's releases after I-0147 closed, so the MCP task_done still refuses
status: done
owner: alex
created: 2026-10-15T05:02:18Z
updated: 2026-10-16T05:24:55Z
transitions:
  - to: ready
    at: 2026-10-15T05:03:58Z
    by: alex
  - to: in-progress
    at: 2026-10-15T05:04:08Z
    by: agent-S-0391
  - to: review
    at: 2026-10-15T23:06:44Z
    by: agent-S-0391
  - to: done
    at: 2026-10-16T05:24:55Z
    by: alex
tags: []
topics: [release, host, mcp]
touches: [docs/operators/index.md, docs/operators/runbooks/update.md, design/system/flai-cli.md, design/issues/I-0149-the-host-s-installed-flai-stays-behind-the-project-s-releases-after-i-0147-closed-so-the-mcp-task-done-still-refuses.md, design/issues/summary.md, design/adrs/0161-flai-serve-s-upgrade-thread-names-a-flai-release-only-once-it-is-published-and.md, design/adrs/README.md, flai/cmd/serve.go, flai/cmd/serve_actions.go, flai/internal/serve/serve.go, flai/internal/serve/upgrade.go, flai/internal/serve/upgrade_test.go, docs/operators/runbooks/migrate.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 2437
  priced_by: harness
  turns:
    - day: 2026-10-15
      ceremony: 6
      hand_edits: 1
      work: 53
  models:
    - model: claude-haiku-5-5
      input: 40
      output: 36940
      cache_read: 1495709
      cache_write: 127052
      cost: 0.1501
    - model: claude-opus-5-5
      input: 204
      output: 54170
      cache_read: 13339113
      cache_write: 853410
      cost: 9.857
  strategic:
    - kind: orchestrator
      seconds: 66212
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 178
          output: 3463
          cache_read: 67583624
          cache_write: 47401
          cost: 16.9545
        - model: claude-sonnet-5-5
          input: 10
          output: 75
          cache_read: 123089
          cache_write: 47501
          cost: 0.1396
cost_of_delay:
  inputs:
    time_lost_per_cycle: 2m
    by: flai
    at: 2026-10-15T05:02:18Z
  value: 5
  by: planner-S-0391
  at: 2026-10-15T05:05:54Z
forecast:
  duration: 20m
  delivery: 2026-10-15T05:30:00Z
  basis: "flai forecast's 11m (87 s per unit of size over 14 medium improvement stories, size 7) raised to 20m for verifying the upgraded host through MCP task_done, counted from the 05:04Z start and assuming the operator upgrades the host promptly."
  by: planner-S-0391
  at: 2026-10-15T05:05:54Z
finalized:
  by: alex
  at: 2026-10-15T05:02:58Z
---
# S-0391 The host's installed flai stays behind the project's releases after I-0147 closed, so the MCP task_done still refuses

## Goal

This story remediates [I-0149](../../../design/issues/I-0149-the-host-s-installed-flai-stays-behind-the-project-s-releases-after-i-0147-closed-so-the-mcp-task-done-still-refuses.md), "The host's installed flai stays behind the project's releases after I-0147 closed, so the MCP task_done still refuses". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [x] The cause I-0149 describes no longer occurs, with a test that reproduces it where one fits
- [x] I-0149 is closed with `flai issue close I-0149 --reason` saying what fixed it

## Tasks
- T-1599 flai serve's upgrade thread names a release only once it is published
- T-1600 Operator docs: the one manual upgrade a host older than flai 1.50.1 needs, and which sessions ADR-0160 helps
- T-1601 Verify the host caught up after the operator's upgrade and close I-0149

## Notes

Cost of delay inputs set by flai from I-0149. time_lost_per_cycle 2m: 2m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-10-15T04:41:02Z, 0 days before this story; under one cycle counts as one).

### Planning

Planned by planner-S-0391 on 2026-10-15. The story moved to in-progress at 05:04Z while it was being planned, so the planner drafted no tasks: the story's agent writes them, and the plan's thread proposes three.

The cause, read from I-0149's instance and the host: the installed flai is 1.48.0 (`~/.flai/bin/flai version`), and the host has no `auto-upgrade` action on. ADR-0158's minute check and thread came in flai 1.50.1, and ADR-0160's close with the installed flai in 1.51.0. A host older than 1.50.1 runs neither, so it never notices a release and never asks. Nothing new in the code can reach a binary already installed. The remedy proposed is the operator's one `flai host upgrade`, with `auto-upgrade` on if the operator chooses, and a note in the docs that a host older than 1.50.1 is upgraded by hand once. `flai/internal/serve/upgrade_test.go` already reproduces the check on a newer host, so no new test fits unless the story's agent finds a gap.

Touches, all files, no folder:

| Touch | From |
|-------|------|
| `docs/operators/index.md` | design: § Keeping the host's flai current says the host looks every minute, which holds only from 1.50.1 |
| `docs/operators/runbooks/update.md` | design: steps 5 and 6 predate the bootstrap note and ADR-0160 |
| `design/system/flai-cli.md` | design: § Versions, the host's flai, and the tree describes ADR-0158's check |
| the I-0149 issue file | criterion 2 |
| `design/issues/summary.md` | layout: `flai issue close` regenerates it |

`flai touches suggest` listed `docs/users/flai.md` (33% co-change) and `design/system/agent-narrative.md` (3%). Both were left out: no command or agent rule changes.

Figures:

- Forecast 20m: flai forecast gave 11m (size 7, 87 s per unit). It was raised for verifying the upgraded host through MCP `task_done`. Delivery counts from the 05:04Z start and assumes the operator upgrades promptly.
- Cost of delay 5 USD a week: flai cod's figure from the inputs (2m per cycle at 150 USD an hour), kept. The inputs are the operator's. They likely understate the cost: while the host lags, every story's agent meets the refusal at each task close, as I-0147's five instances in two days showed.
