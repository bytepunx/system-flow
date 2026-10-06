---
id: S-0294
type: story
nature: remediation
title: Claude Code ends a headless agent ten minutes after its turn ends, even while its background sub-agent is still working, and flai serve leaves the story in progress with no agent
status: backlog
owner: alex
created: 2026-10-06T11:44:49Z
updated: 2026-10-06T22:56:13Z
transitions: []
tags: [cli, serve]
topics: [automation]
touches: [flai/internal/serve/agents.go, flai/internal/serve/agents_test.go, flai/internal/serve/restart.go, flai/internal/config/config.go, flai/internal/config/config_test.go, flai/cmd/serve_actions.go, flai/cmd/serve_actions_test.go, design/adrs, design/system/flai-cli.md, design/system/workflow.md, docs/users/flai.md, docs/users/flai-reference.md, docs/operators/index.md, docs/operators/settings.md, design/issues/I-0084-claude-code-ends-a-headless-agent-ten-minutes-after-its-turn-ends-even-while-its-background-sub-agent-is-still-working-and-flai-serve-leaves-the-story-in-progress-with-no-agent.md, design/issues/summary.md]
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
      seconds: 328
      estimated: true
      models:
        - model: claude-haiku-4-5-20251001
          input: 507
          output: 4242
          cache_read: 3640961
          cache_write: 173470
          cost: 0.7745
        - model: claude-opus-5-5
          input: 829
          output: 59561
          cache_read: 40813838
          cache_write: 658779
          cost: 23.5617
cost_of_delay:
  inputs:
    time_lost_per_cycle: 20m
    by: flai
    at: 2026-10-06T11:44:49Z
  value: 50
  by: planner-S-0294
  at: 2026-10-06T22:53:45Z
forecast:
  duration: 40m
  delivery: 2026-10-07T11:05:00Z
  basis: "Its own forecast of 40m; 41st in the pull order with an in-progress limit of 3, behind S-0299, S-0301, S-0300, S-0228, S-0261, S-0269, S-0270, S-0271, S-0212, S-0213, S-0214, S-0215, S-0216, S-0232, S-0233, S-0234, S-0235, S-0236, S-0237, S-0238, S-0239, S-0241, S-0245, S-0246, S-0251, S-0254, S-0265, S-0272, S-0273, S-0274, S-0275, S-0277, S-0279, S-0280, S-0281, S-0286, S-0287, S-0288, S-0289, S-0290, S-0291 and S-0293."
  by: flai
  at: 2026-10-06T22:56:13Z
finalized:
  by: alex
  at: 2026-10-06T22:50:34Z
---
# S-0294 Claude Code ends a headless agent ten minutes after its turn ends, even while its background sub-agent is still working, and flai serve leaves the story in progress with no agent

## Goal

This story remediates [I-0084](../../../design/issues/I-0084-claude-code-ends-a-headless-agent-ten-minutes-after-its-turn-ends-even-while-its-background-sub-agent-is-still-working-and-flai-serve-leaves-the-story-in-progress-with-no-agent.md), "Claude Code ends a headless agent ten minutes after its turn ends, even while its background sub-agent is still working, and flai serve leaves the story in progress with no agent". The issue recommends this solution:

Two separate things went wrong, and each has its own fix.

1. **The agent has no safe way to wait.** The prompt and `delegation.md` name one that keeps the session alive for a sub-agent of any length: a launch with `run_in_background: false`, once tested against a sub-agent that runs longer than ten minutes. This is the first change I-0083 recommends, and one story can make it for both issues.
2. **A story whose agent ended without finishing waits for a person.** flai serve could restart such an agent itself: once, or a small number of times, when the run ended with the story still in progress, no thread awaiting the operator, and no block on the story; and tell the operator when the limit is reached. Whether flai serve restarts on its own is the operator's decision, since ADR-0043 gave the restart to a person.

S-0285 made the first remediation ([ADR-0092](../adrs/0092-a-story-s-agent-waits-for-a-sub-agent-by-launching-it-in-the-foreground-and.md)). It measured the ten minutes again on 2.1.290: a turn ended at 10:34:48Z with a background sub-agent out, and the process exited at 10:44:49Z with the sub-agent cut off. It also measured the remedy: a launch with `run_in_background` false returned an 11-minute sub-agent's result as the tool's result, and three launched in one message ran together. The start prompt and `delegation.md` now name that way and say never to end the turn on a background sub-agent. The second remediation is still open.

## Acceptance criteria
- [ ] The cause I-0084 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0084 is closed with `flai issue close I-0084 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0084. time_lost_per_cycle 20m: 10m per occurrence × 2 occurrences ÷ 1 cycle of 168h (first reported 2026-10-06T07:00:22Z, 0.2 days before this story; under one cycle counts as one).
