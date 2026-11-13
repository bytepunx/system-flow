---
id: S-0776
type: story
nature: improvement
title: Publish what is accepted and not yet released
status: in-progress
owner: alex
created: 2026-11-13T17:20:42Z
updated: 2026-11-13T17:20:53Z
transitions:
  - to: ready
    at: 2026-11-13T17:20:43Z
    by: system-flow
  - to: in-progress
    at: 2026-11-13T17:20:53Z
    by: agent-S-0776
tags: []
touches: [wip/kanban/stories/S-0776-publish-what-is-accepted-and-not-yet-released.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 84
  estimated: true
  priced_by: estimate
  turns:
    - day: 2026-11-13
      work: 15
  models:
    - model: claude-opus-5-5
      input: 32
      output: 213
      cache_read: 1810106
      cache_write: 148837
      cost: 0.7341
  strategic:
    - kind: orchestrator
      seconds: 19
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 6
          output: 59
          cache_read: 317510
          cache_write: 3112
          cost: 0.0936
formula: publish
---
# S-0776 Publish what is accepted and not yet released

## Goal

Publish what is accepted and not yet released, by the project's release policy.

## Acceptance criteria

- [ ] Fetch the remote's branch and tags
- [ ] Evaluate the release policy, stopping unless it is met
- [ ] Publish what is accepted and not yet released

## Tasks
- T-3316 Fetch the remote's branch and tags
- T-3317 Evaluate the release policy, stopping unless it is met
- T-3318 Publish what is accepted and not yet released

## Notes

Made by `flai formula run publish` from `design/workflows/formulas/publish.md`, with the inputs `judgement=`.
