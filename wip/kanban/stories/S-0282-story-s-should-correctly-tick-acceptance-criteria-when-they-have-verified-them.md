---
id: S-0282
type: story
nature: remediation
title: Story's should correctly tick acceptance criteria when they have verified them
status: backlog
owner: alex
created: 2026-10-06T02:59:32Z
updated: 2026-10-06T02:59:32Z
transitions: []
tags: [cli]
touches: [flai/cmd]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
cost_of_delay:
  inputs:
    time_lost_per_cycle: 30m
    by: alex
    at: 2026-10-06T02:59:32Z
---
# S-0282 Story's should correctly tick acceptance criteria when they have verified them

## Goal

When an agent's story is validating completed criteria or a task sub-agent is completing criteria, it should correctly tick that acceptance criteria as completed on the story. flai should expose a command in the CLI, HTTP, and MCP for ticking acceptance criteria to prevent the AI's from needing to edit the file or the specific syntax change to make.

## Acceptance criteria
- [ ] criteria are ticked consistently when they are completed, regardless by which agent, through flai

## Tasks

## Notes
