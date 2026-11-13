---
id: S-0768
type: story
nature: improvement
title: Publish what is accepted and not yet released
status: cancelled
owner: alex
created: 2026-11-13T11:43:32Z
updated: 2026-11-13T11:45:33Z
transitions:
  - to: ready
    at: 2026-11-13T11:43:34Z
    by: system-flow
  - to: in-progress
    at: 2026-11-13T11:43:41Z
    by: agent-S-0768
  - to: cancelled
    at: 2026-11-13T11:45:33Z
    by: alex
tags: []
touches: [wip/kanban/stories/S-0768-publish-what-is-accepted-and-not-yet-released.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 150
  priced_by: harness
  turns:
    - day: 2026-11-13
      ceremony: 3
      work: 17
  models:
    - model: claude-opus-5-5
      input: 44
      output: 8512
      cache_read: 2797585
      cache_write: 174206
      cost: 2.1236
  strategic:
    - kind: orchestrator
      seconds: 51
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 20
          output: 142
          cache_read: 669284
          cache_write: 36081
          cost: 0.2039
        - model: claude-sonnet-5-5
          input: 0
          output: 5
          cache_read: 0
          cache_write: 2761
          cost: 0.0021
formula: publish
---
# S-0768 Publish what is accepted and not yet released

## Goal

Publish what is accepted and not yet released, by the project's release policy.

## Acceptance criteria

- [ ] Fetch the remote's branch and tags
- [ ] Evaluate the release policy, stopping unless it is met
- [ ] Publish what is accepted and not yet released

## Tasks
- T-3288 Fetch the remote's branch and tags
- T-3289 Evaluate the release policy, stopping unless it is met
- T-3290 Publish what is accepted and not yet released

## Notes

Made by `flai formula run publish` from `design/workflows/formulas/publish.md`, with the inputs `judgement=`.
- 2026-11-13T11:45:33Z: moved to cancelled: REDUNDANT
