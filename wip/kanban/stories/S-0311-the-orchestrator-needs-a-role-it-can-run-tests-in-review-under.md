---
id: S-0311
type: story
nature: feature
title: The Orchestrator needs a role it can run tests in review under
status: ready
owner: alex
created: 2026-10-07T14:24:35Z
updated: 2026-10-07T14:24:36Z
transitions:
  - to: ready
    at: 2026-10-07T14:24:36Z
    by: alex
tags: [cli]
topics: [orchestrator]
touches: [flai/cmd]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: sum
  seconds: 0
  models: []
  strategic:
    - kind: orchestrator
      seconds: 295
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 6
          output: 100
          cache_read: 1004938
          cache_write: 5135
          cost: 0.2636
---
# S-0311 The Orchestrator needs a role it can run tests in review under

## Goal

The orchestrator needs the ability to execute test runs without hitting flai guard's rules about orchestrator permissions in test scenarios. This likely means using whatever FLAI_ROLE identity will avoid tripping the guard's check.

## Acceptance criteria
- [ ] Orchestrator is able to complete test runs as a role that does not break flai's guard checks against the `orchestrator` role

## Tasks

## Notes
