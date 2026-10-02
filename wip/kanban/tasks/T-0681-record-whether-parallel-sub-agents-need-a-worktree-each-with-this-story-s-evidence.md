---
id: T-0681
type: task
nature: experiment
title: Record whether parallel sub-agents need a worktree each, with this story's evidence
status: done
parent: S-0176
owner: arobson
created: 2026-10-01T11:42:11Z
updated: 2026-10-01T12:08:28Z
transitions:
  - to: ready
    at: 2026-10-01T12:07:10Z
    by: agent-S-0176
  - to: in-progress
    at: 2026-10-01T12:07:10Z
    by: agent-S-0176
  - to: done
    at: 2026-10-01T12:08:28Z
    by: agent-S-0176
stream: S-0176
tags: []
touches: [design/system/agent-context.md]
usage:
  source: log
  seconds: 78
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 38
      output: 12107
      cache_read: 2675597
      cache_write: 38179
      cost: 1.0138
---
# T-0681 Record whether parallel sub-agents need a worktree each, with this story's evidence

## Work

- Record what happened when this story's layers ran: the first with its sub-agents in the story's one worktree, and the second with each task in its own worktree from the story's branch, merged back by the story's agent. Note every interference (a build or test broken by another sub-agent's half-done edit, a file written by two, a merge conflict) and what each mode cost in time and tokens.
- Say which mode the conventions should ask for, and why, in `design/system/agent-context.md`.

## Done when

- agent-context.md records both modes with the evidence from this story's logs, and the conventions written in this story ask for the mode it recommends.

## Notes
