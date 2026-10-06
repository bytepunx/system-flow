---
id: T-1005
type: task
nature: improvement
title: The wait_for_events description says it is for a thread awaiting the designer, not a sub-agent
status: done
parent: S-0285
owner: alex
created: 2026-10-06T10:38:06Z
updated: 2026-10-06T10:49:46Z
transitions:
  - to: ready
    at: 2026-10-06T10:46:38Z
    by: agent-S-0285
  - to: in-progress
    at: 2026-10-06T10:46:38Z
    by: agent-S-0285
  - to: done
    at: 2026-10-06T10:49:46Z
    by: agent-S-0285
stream: S-0285
tags: []
touches: [flai/internal/mcpserver]
usage:
  source: log
  seconds: 188
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 15
      output: 5710
      cache_read: 774695
      cache_write: 22730
      cost: 0.4114
---
# T-1005 The wait_for_events description says it is for a thread awaiting the designer, not a sub-agent

## Work

The MCP tool description of `wait_for_events` no longer says "Hold this when idle"; it says the tool is for a thread awaiting the designer, that it does not see sub-agents finish, and how a story's agent waits for a sub-agent instead. A test pins the sentence. A test shows a held `wait_for_events` returns on a reply to a thread on the story (add it if none exists). Waits for nothing: first layer.

## Done when

- [ ] the description says what it is for and how to wait for a sub-agent, with a test holding it
- [ ] a test shows a held wait returns on a thread reply

## Notes
