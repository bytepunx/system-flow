---
id: S-0246
type: story
nature: improvement
title: flai guard refuses a sub-agent's shell command whose heredoc text reads like a flai write
status: backlog
owner: alex
created: 2026-10-03T17:49:40Z
updated: 2026-10-05T05:46:20Z
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
  delivery: 2026-10-05T23:41:00Z
  basis: "Its own forecast of 30m; 29th in the pull order with an in-progress limit of 3, behind S-0276, S-0217, S-0218, S-0219, S-0220, S-0221, S-0222, S-0226, S-0212, S-0213, S-0214, S-0215, S-0216, S-0223, S-0224, S-0227, S-0228, S-0229, S-0232, S-0233, S-0234, S-0235, S-0236, S-0237, S-0238, S-0239, S-0241 and S-0245."
  by: flai
  at: 2026-10-05T05:46:20Z
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
