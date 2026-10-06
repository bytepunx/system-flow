---
id: T-1004
type: task
nature: improvement
title: flai guard refuses a story agent's wait_for_events while a sub-agent of its session runs and no thread on the story is open
status: done
parent: S-0285
owner: alex
created: 2026-10-06T10:38:05Z
updated: 2026-10-06T10:46:33Z
transitions:
  - to: ready
    at: 2026-10-06T10:38:35Z
    by: agent-S-0285
  - to: in-progress
    at: 2026-10-06T10:38:36Z
    by: agent-S-0285
  - to: done
    at: 2026-10-06T10:46:33Z
    by: agent-S-0285
stream: S-0285
tags: []
touches: [flai/internal/guard, flai/cmd/guard.go, flai/cmd/guard_test.go, docs/users/flai-reference.md]
usage:
  source: log
  seconds: 477
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 83
      output: 32313
      cache_read: 4383982
      cache_write: 128630
      cost: 2.3283
---
# T-1004 flai guard refuses a story agent's wait_for_events while a sub-agent of its session runs and no thread on the story is open

## Work

`flai guard` learns which sub-agents of a Claude Code session are running from two more hook events, `SubagentStart` and `SubagentStop`, whose input carries `session_id`, `agent_id`, and `agent_type` (tested on Claude Code 2.1.290; see the narrative). It records them per session under the main checkout's `.flai-cache`, and refuses the story agent's own `mcp__flai__wait_for_events` call (no `agent_id`, `FLAI_STORY` set, no `FLAI_ROLE`) while one runs and no unresolved thread on the story or its tasks is open, saying which sub-agents run and how to wait for them instead. A wait with such a thread open passes, so it still returns on the designer's answer. The command's help says so, and `docs/users/flai-reference.md` is regenerated from it. Waits for nothing: first layer.

## Done when

- [ ] `SubagentStart` and `SubagentStop` inputs keep a session's running sub-agents, in the main checkout's `.flai-cache`, safely when several start at once
- [ ] the story agent's `wait_for_events` is refused, exit 2, naming the running sub-agents and how to wait, while one runs and no thread on the story is open
- [ ] it passes with no sub-agent running, with a thread on the story or one of its tasks open, for a sub-agent's call (refused as before), and outside a story agent's session
- [ ] tests cover each case, through `flai guard` on standard input

## Notes
