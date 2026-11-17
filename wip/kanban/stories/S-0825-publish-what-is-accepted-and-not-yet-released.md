---
id: S-0825
type: story
nature: improvement
title: Publish what is accepted and not yet released
status: in-progress
owner: alex
created: 2026-11-17T07:54:16Z
updated: 2026-11-17T08:14:03Z
transitions:
  - to: ready
    at: 2026-11-17T07:54:18Z
    by: system-flow
  - to: in-progress
    at: 2026-11-17T08:14:03Z
    by: agent-S-0825
tags: []
touches: [wip/kanban/stories/S-0825-publish-what-is-accepted-and-not-yet-released.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 75
  estimated: true
  priced_by: estimate
  turns:
    - day: 2026-11-17
      work: 13
  models:
    - model: claude-opus-5-5
      input: 28
      output: 225
      cache_read: 1411382
      cache_write: 156186
      cost: 0.5972
  strategic:
    - kind: orchestrator
      seconds: 42
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 13
          output: 87
          cache_read: 448622
          cache_write: 39911
          cost: 0.1513
        - model: claude-sonnet-5-5
          input: 2
          output: 18
          cache_read: 10734
          cache_write: 4191
          cost: 0.012
formula: publish
---
# S-0825 Publish what is accepted and not yet released

## Goal

Publish what is accepted and not yet released, by the project's release policy.

## Acceptance criteria

- [ ] Fetch the remote's branch and tags
- [ ] Evaluate the release policy, stopping unless it is met
- [ ] Publish what is accepted and not yet released

## Tasks
- T-3513 Fetch the remote's branch and tags
- T-3514 Evaluate the release policy, stopping unless it is met
- T-3515 Publish what is accepted and not yet released

## Notes

Made by `flai formula run publish` from `design/workflows/formulas/publish.md`, with the inputs `judgement=`.
