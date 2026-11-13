---
id: S-0770
type: story
nature: improvement
title: Publish what is accepted and not yet released
status: in-progress
owner: alex
created: 2026-11-13T13:16:45Z
updated: 2026-11-13T13:16:53Z
transitions:
  - to: ready
    at: 2026-11-13T13:16:46Z
    by: system-flow
  - to: in-progress
    at: 2026-11-13T13:16:53Z
    by: agent-S-0770
tags: []
touches: [wip/kanban/stories/S-0770-publish-what-is-accepted-and-not-yet-released.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 100
  estimated: true
  priced_by: estimate
  turns:
    - day: 2026-11-13
      work: 15
  models:
    - model: claude-opus-5-5
      input: 32
      output: 267
      cache_read: 1685254
      cache_write: 153683
      cost: 0.6871
  strategic:
    - kind: orchestrator
      seconds: 26
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 9
          output: 53
          cache_read: 325012
          cache_write: 33767
          cost: 0.1041
formula: publish
---
# S-0770 Publish what is accepted and not yet released

## Goal

Publish what is accepted and not yet released, by the project's release policy.

## Acceptance criteria

- [ ] Fetch the remote's branch and tags
- [ ] Evaluate the release policy, stopping unless it is met
- [ ] Publish what is accepted and not yet released

## Tasks
- T-3297 Fetch the remote's branch and tags
- T-3298 Evaluate the release policy, stopping unless it is met
- T-3299 Publish what is accepted and not yet released

## Notes

Made by `flai formula run publish` from `design/workflows/formulas/publish.md`, with the inputs `judgement=`.
