---
id: S-0294
type: story
nature: remediation
title: Claude Code ends a headless agent ten minutes after its turn ends, even while its background sub-agent is still working, and flai serve leaves the story in progress with no agent
status: backlog
owner: alex
created: 2026-10-06T11:44:49Z
updated: 2026-10-06T22:50:34Z
transitions: []
tags: []
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
cost_of_delay:
  inputs:
    time_lost_per_cycle: 20m
    by: flai
    at: 2026-10-06T11:44:49Z
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
