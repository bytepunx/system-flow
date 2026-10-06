---
id: S-0285
type: story
nature: improvement
title: A story's agent waits for its sub-agents with wait_for_events, which cannot see them finish, so each wait runs to its timeout
status: backlog
owner: alex
created: 2026-10-06T06:03:25Z
updated: 2026-10-06T06:03:25Z
transitions: []
tags: []
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
draft: true
cost_of_delay:
  inputs:
    time_lost_per_cycle: 1h18m
    by: flai
    at: 2026-10-06T06:03:25Z
---
# S-0285 A story's agent waits for its sub-agents with wait_for_events, which cannot see them finish, so each wait runs to its timeout

## Goal

This story remediates [I-0083](../../../design/issues/I-0083-a-story-s-agent-waits-for-its-sub-agents-with-wait-for-events-which-cannot-see-them-finish-so-each-wait-runs-to-its-timeout.md), "A story's agent waits for its sub-agents with wait_for_events, which cannot see them finish, so each wait runs to its timeout". The issue recommends this solution:

Two changes, both needed. The wording T-0869 added did not hold in two stories, so better wording is not enough without a check in flai.

1. **Say how to wait.** The start prompt in `flai/internal/harness`, `delegation.md` here and in the template, and the `wait_for_events` tool description say what the agent does after it launches a sub-agent with nothing else to do: it ends its turn with no tool call, as S-0219's agent did. They say that the session is kept while a sub-agent runs, and that the sub-agent's notice begins the next turn. They say when `wait_for_events` is right: only for a thread awaiting the designer. The sentence no longer depends on the agent having chosen the background. A launch with `run_in_background: false`, which returns the result as the tool's result, is the other candidate; test both against the Claude Code that flai serve runs, 2.1.290 when this was recorded, and name the one that works.
2. **Make flai catch it.** `flai guard` already runs as a `PreToolUse` hook on the story agent's `mcp__flai__` calls, so flai sees the `wait_for_events` call when it is made. Either of these keeps a mistaken wait short:
   - A Claude Code hook on a sub-agent's stop (`SubagentStop`) tells flai, and a `wait_for_events` held by that session returns at once, saying which sub-agent finished. Confirm first what that hook's input carries.
   - `flai guard` refuses a `wait_for_events` call made while a sub-agent of the session is running and no thread on the story awaits the designer, and its refusal says to end the turn.

The measure, on a story flai serve works after the change is installed: no wait on a sub-agent outlasts the sub-agent by more than a few seconds.

## Acceptance criteria
- [ ] The cause I-0083 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0083 is closed with `flai issue close I-0083 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0083. time_lost_per_cycle 1h18m: 39m per occurrence × 2 occurrences ÷ 1 cycle of 168h (first reported 2026-10-06T06:02:33Z, 0 days before this story; under one cycle counts as one).
