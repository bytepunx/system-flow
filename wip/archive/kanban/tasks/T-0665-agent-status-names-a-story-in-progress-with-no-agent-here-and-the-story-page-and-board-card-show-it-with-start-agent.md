---
id: T-0665
type: task
nature: feature
title: agent.status names a story in progress with no agent here, and the story page and board card show it with Start agent
status: done
parent: S-0177
owner: alex
created: 2026-10-01T10:07:00Z
updated: 2026-10-01T10:19:45Z
transitions:
  - to: ready
    at: 2026-10-01T10:07:15Z
    by: agent-S-0177
  - to: in-progress
    at: 2026-10-01T10:15:12Z
    by: agent-S-0177
  - to: done
    at: 2026-10-01T10:19:45Z
    by: agent-S-0177
stream: S-0177
tags: []
touches: [flai/internal/serve, flaiover/src]
usage:
  source: log
  seconds: 273
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 78
      output: 27257
      cache_read: 5633137
      cache_write: 118141
      cost: 2.4194
---
# T-0665 agent.status names a story in progress with no agent here, and the story page and board card show it with Start agent

## Work

- `serve.Activity` gives a story in progress with no run here an entry: `waiting`, why "begun by <agent> on <host|another host> at <time>; no agent here", and `elsewhere: {agent, at, host}`, with a stand-in run as a held story has.
- flaiover: `StoryAgent` shows the line and offers Start agent, which posts `restart`. `BoardCard` shows the line. `storyActivity` keeps these entries while the action is off, as it keeps the held ones.
- Go and vitest tests. `design/system/flaiover-dashboard.md`, `design/system/flai-cli.md`, and `docs/users/flaiover.md` say so.

## Done when

The tests and lint pass and the change is committed.

## Notes
