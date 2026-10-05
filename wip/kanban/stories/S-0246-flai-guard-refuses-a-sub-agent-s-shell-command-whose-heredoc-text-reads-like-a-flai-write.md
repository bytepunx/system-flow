---
id: S-0246
type: story
nature: improvement
title: flai guard refuses a sub-agent's shell command whose heredoc text reads like a flai write
status: backlog
owner: alex
created: 2026-10-03T17:49:40Z
updated: 2026-10-05T05:46:14Z
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
forecast:
  duration: 30m
  delivery: 2026-10-05T23:49:00Z
  basis: "flai forecast's 16m (86 s per unit of size times size 11, 9 touches and 2 criteria), raised to 30m: heredoc parsing in the guard's hand-written tokenizer, with about ten edge-case tests, then a serial docs-and-issue task. S-0175 built the whole guard in 40m. Delivery is flai's 23:35Z, 29th in the pull order, plus the extra 14m."
  by: planner-S-0246
  at: 2026-10-05T05:45:14Z
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
