---
id: S-0821
type: story
nature: improvement
title: Publish what is accepted and not yet released
status: cancelled
owner: alex
created: 2026-11-17T04:29:08Z
updated: 2026-11-17T07:46:43Z
transitions:
  - to: ready
    at: 2026-11-17T04:29:10Z
    by: system-flow
  - to: cancelled
    at: 2026-11-17T07:46:43Z
    by: alex
tags: []
touches: [wip/kanban/stories/S-0821-publish-what-is-accepted-and-not-yet-released.md]
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
      seconds: 64
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 16
          output: 98
          cache_read: 624075
          cache_write: 78791
          cost: 0.2173
formula: publish
---
# S-0821 Publish what is accepted and not yet released

## Goal

Publish what is accepted and not yet released, by the project's release policy.

## Acceptance criteria

- [ ] Fetch the remote's branch and tags
- [ ] Evaluate the release policy, stopping unless it is met
- [ ] Publish what is accepted and not yet released

## Tasks
- T-3510 Fetch the remote's branch and tags
- T-3511 Evaluate the release policy, stopping unless it is met
- T-3512 Publish what is accepted and not yet released

## Notes

Made by `flai formula run publish` from `design/workflows/formulas/publish.md`, with the inputs `judgement=`.
- 2026-11-17T07:46:43Z: moved to cancelled: not needed
