---
id: S-0802
type: story
nature: improvement
title: Publish what is accepted and not yet released
status: in-progress
owner: alex
created: 2026-11-16T04:12:37Z
updated: 2026-11-16T04:12:46Z
transitions:
  - to: ready
    at: 2026-11-16T04:12:39Z
    by: system-flow
  - to: in-progress
    at: 2026-11-16T04:12:46Z
    by: agent-S-0802
tags: []
touches: [wip/kanban/stories/S-0802-publish-what-is-accepted-and-not-yet-released.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 77
  estimated: true
  priced_by: estimate
  turns:
    - day: 2026-11-16
      work: 14
  models:
    - model: claude-opus-5-5
      input: 32
      output: 241
      cache_read: 1687446
      cache_write: 153694
      cost: 0.6984
  strategic:
    - kind: orchestrator
      seconds: 31
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 10
          output: 82
          cache_read: 342050
          cache_write: 28310
          cost: 0.1132
formula: publish
---
# S-0802 Publish what is accepted and not yet released

## Goal

Publish what is accepted and not yet released, by the project's release policy.

## Acceptance criteria

- [ ] Fetch the remote's branch and tags
- [ ] Evaluate the release policy, stopping unless it is met
- [ ] Publish what is accepted and not yet released

## Tasks
- T-3423 Fetch the remote's branch and tags
- T-3424 Evaluate the release policy, stopping unless it is met
- T-3425 Publish what is accepted and not yet released

## Notes

Made by `flai formula run publish` from `design/workflows/formulas/publish.md`, with the inputs `judgement=`.
