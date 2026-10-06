---
id: T-1006
type: task
nature: improvement
title: flai serve's start prompt says how a story's agent waits for its sub-agents
status: done
parent: S-0285
owner: alex
created: 2026-10-06T10:38:07Z
updated: 2026-10-06T10:49:47Z
transitions:
  - to: ready
    at: 2026-10-06T10:46:39Z
    by: agent-S-0285
  - to: in-progress
    at: 2026-10-06T10:46:39Z
    by: agent-S-0285
  - to: done
    at: 2026-10-06T10:49:47Z
    by: agent-S-0285
stream: S-0285
tags: []
touches: [flai/internal/harness]
usage:
  source: log
  seconds: 188
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 14
      output: 5586
      cache_read: 757861
      cache_write: 22236
      cost: 0.4025
---
# T-1006 flai serve's start prompt says how a story's agent waits for its sub-agents

## Work

The start prompt's sentence on waiting (`delegation` in `flai/internal/harness/harness.go`) says how to wait in the way the narrative's experiment found to work for a sub-agent of any length, that ending the turn with a sub-agent out loses the session after ten minutes, and that `wait_for_events` is only for a thread awaiting the designer. It no longer depends on the agent having chosen the background. `harness_test.go` holds the new sentence. Waits for nothing: first layer.

## Done when

- [ ] the prompt says how to wait and when `wait_for_events` is right, with a test holding it

## Notes
