---
id: T-0662
type: task
nature: feature
title: An ADR records how a story begun on another host gets an agent on this one
status: done
parent: S-0177
owner: alex
created: 2026-10-01T10:06:49Z
updated: 2026-10-01T10:08:14Z
transitions:
  - to: ready
    at: 2026-10-01T10:07:15Z
    by: agent-S-0177
  - to: in-progress
    at: 2026-10-01T10:07:15Z
    by: agent-S-0177
  - to: done
    at: 2026-10-01T10:08:14Z
    by: agent-S-0177
stream: S-0177
tags: []
touches: [design/adrs/]
usage:
  source: log
  seconds: 59
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 16
      output: 96
      cache_read: 923804
      cache_write: 12879
      cost: 0.374
---
# T-0662 An ADR records how a story begun on another host gets an agent on this one

## Work

Write an ADR, refining ADR-0043, that decides: restart starts a story in ready or in progress that this host has no run for; its agent is told who began the story, when, and where, and which threads were answered since; the narrative records the host that opened it; `flai stream open` reopens a stream whose worktree is gone, from the local branch, the remote's, or main; `agent.status` names a story in progress with no agent here; the MCP server can start or restart an agent, gated by the `agent` action.

## Done when

The ADR is accepted in `design/adrs/` with its index row, and the narrative's Decisions name it.
