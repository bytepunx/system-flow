---
id: T-0555
type: task
nature: feature
title: The thread pane shows the agent working on the operator's reply with an animation
status: done
parent: S-0154
owner: alex
created: 2026-09-29T07:08:18Z
updated: 2026-09-29T07:20:35Z
transitions:
  - to: ready
    at: 2026-09-29T07:08:32Z
    by: agent-S-0154
  - to: in-progress
    at: 2026-09-29T07:18:25Z
    by: agent-S-0154
  - to: done
    at: 2026-09-29T07:20:35Z
    by: agent-S-0154
stream: S-0154
tags: []
usage:
  source: log
  seconds: 130
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 22
      output: 6212
      cache_read: 1585089
      cache_write: 20704
      cost: 0.607
---

# T-0555 The thread pane shows the agent working on the operator's reply with an animation

## Work

When the shown thread's last entry is the operator's and the story's agent is working, or has ended on that thread and is being started again, the thread shows a pulsing line under the reply saying the agent is working on it.

## Done when

A behaviour test shows the line for a working agent after an operator reply, and not after an agent's entry, for a resolved thread, or with no agent.

## Notes
