---
id: S-0583
type: story
nature: remediation
title: Claude Code drops any call to flai's HTTP MCP server that answers after 60 seconds, so wait_for_events, wait_for_work, verify, and a slow task_done fail for an agent in a container
status: backlog
parent: E-0029
owner: alex
created: 2026-11-07T04:46:50Z
updated: 2026-11-07T04:46:50Z
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
    time_lost_per_cycle: 7m
    by: flai
    at: 2026-11-07T04:46:50Z
---
# S-0583 Claude Code drops any call to flai's HTTP MCP server that answers after 60 seconds, so wait_for_events, wait_for_work, verify, and a slow task_done fail for an agent in a container

## Goal

This story remediates [I-0228](../../../design/issues/I-0228-claude-code-drops-any-call-to-flai-s-http-mcp-server-that-answers-after-60-seconds-so-wait-for-events-wait-for-work-verify-and-a-slow-task-done-fail-for-an-agent-in-a-container.md), "Claude Code drops any call to flai's HTTP MCP server that answers after 60 seconds, so wait_for_events, wait_for_work, verify, and a slow task_done fail for an agent in a container". The issue recommends this solution:

S-0506 fixes `wait_for_events` and `wait_for_work`: over HTTP the server holds one call at most 45 seconds, the bound `permission_prompt` holds there, and then answers `timed_out` for the agent to call again ([ADR-0244](../adrs/0244-over-http-one-wait-for-events-or-wait-for-work-call-holds-at-most-the-45.md)). `verify` has answered at once since S-0546 made the close-out a `flai serve` job, whose `verified` event the agent waits for within that bound ([ADR-0229](../adrs/0229-the-close-out-runs-as-a-flai-serve-job-the-story-s-agent-asks-for-with-the.md)). Still open: a `task_done` or `story_start` slower than 60 seconds over HTTP still meets Claude Code's limit while the server finishes it. A per-server `timeout` in the MCP file S-0508's container adapter writes may lift that limit; it does not replace the bound.

## Acceptance criteria
- [ ] The cause I-0228 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0228 is closed with `flai issue close I-0228 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0228. time_lost_per_cycle 7m: 3m per occurrence × 3 occurrences ÷ 1.3 cycles of 168h (first reported 2026-10-29T00:56:36Z, 9.2 days before this story).
