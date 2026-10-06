---
id: S-0246
type: story
nature: improvement
title: flai guard refuses a sub-agent's shell command whose heredoc text reads like a flai write
status: backlog
owner: alex
created: 2026-10-03T17:49:40Z
updated: 2026-10-06T18:11:46Z
transitions: []
tags: []
topics: [cli]
touches: [flai/internal/guard/guard.go, flai/internal/guard/guard_test.go, flai/cmd/guard_test.go, design/system/flai-cli.md, design/system/agent-context.md, docs/users/flai-reference.md, docs/users/flai.md, design/issues/I-0058-flai-guard-refuses-a-sub-agent-s-shell-command-whose-heredoc-text-reads-like-a-flai-write.md, design/issues/summary.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
cost_of_delay:
  inputs:
    time_lost_per_cycle: 3m
    by: planner-S-0246
    at: 2026-10-05T05:46:14Z
  value: 7.5
  by: planner-S-0246
  at: 2026-10-05T05:46:19Z
forecast:
  duration: 30m
  delivery: 2026-10-07T04:55:00Z
  basis: "Its own forecast of 30m; 25th in the pull order with an in-progress limit of 3, behind S-0226, S-0284, S-0278, S-0223, S-0224, S-0227, S-0295, S-0229, S-0296, S-0212, S-0213, S-0214, S-0215, S-0216, S-0228, S-0232, S-0233, S-0234, S-0235, S-0236, S-0237, S-0238, S-0239, S-0241 and S-0245."
  by: flai
  at: 2026-10-06T18:11:46Z
---
# S-0246 flai guard refuses a sub-agent's shell command whose heredoc text reads like a flai write

## Goal

This story remediates [I-0058](../../../design/issues/I-0058-flai-guard-refuses-a-sub-agent-s-shell-command-whose-heredoc-text-reads-like-a-flai-write.md), "flai guard refuses a sub-agent's shell command whose heredoc text reads like a flai write". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0058 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0058 is closed with `flai issue close I-0058 --reason` saying what fixed it

## Tasks
- T-0920 flai guard reads a heredoc's text as input, not as commands, unless a shell runs it
- T-0926 flai guard's design and docs say how it reads a heredoc, and I-0058 is closed

## Notes

### Planning

Planned by planner-S-0246 on 2026-10-05. The plan thread, TH-0141, lists the tasks, their layers, and the proposed fix. T-0920 is layer 1, and T-0926 follows it in layer 2.

Topics: `cli` added, because the guard is designed in `flai-cli.md` and `agent-context.md`, whose topic it is.

Touches:

- The story declared none, so `flai touches suggest` started from the guard's package.
- `flai/internal/guard/guard.go`: from the layout. `commands` splits a line at newlines, so a heredoc's lines are read as commands. That is the cause.
- `flai/internal/guard/guard_test.go`: from the layout. This is where the reproducing test of criterion 1 goes.
- `flai/cmd/guard_test.go`: co-change, 38%. It holds one hook-input test of a multi-line command.
- `design/system/flai-cli.md`: co-change 25%, and design. Its `flai guard` row says how the line is split.
- `design/system/agent-context.md`: co-change, 38%. It records the guard's refinements of ADR-0060.
- `docs/users/flai-reference.md` and `docs/users/flai.md`: co-change, 25% each. Both describe what the guard reads.
- I-0058's file and `design/issues/summary.md`: from criterion 2. `flai issue close` writes both.
- Not taken from `flai touches suggest`:
  - `flai/cmd/guard.go` (62%). The hook's decoding does not change.
  - `flai/internal/harness/harness.go` and its test (25%). Planner prompt work that came with the guard's planner rules, which this fix leaves alone.

Forecast: 30m.

- `flai forecast` gave 16m: 86 s per unit of size, times size 11 from 2 criteria and 9 touches.
- Raised for two reasons. Recognising heredocs in the guard's hand-written tokenizer needs about ten edge-case tests: quoted, escaped, and `<<-` delimiters, a body fed to a shell, and command substitutions. The docs-and-issue task runs after it, not beside it. S-0175 built the whole guard in 40m.
- Delivery is flai's, 29th in the pull order. flai recomputes it as the board moves.

Cost of delay: 7.50 USD/week, from `flai cod`, kept as computed.

- The input `time_lost_per_cycle: 3m` is the operator's, given on TH-0138: "3m".
- 3m is the cost recorded on I-0058 over two instances, S-0230's T-0711 and S-0209's T-0791. In each, a task sub-agent was refused and redid its edit with the Edit tool.
