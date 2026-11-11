---
id: S-0728
type: story
nature: improvement
title: Publish what is accepted and not yet released
status: in-progress
owner: alex
created: 2026-11-11T23:40:03Z
updated: 2026-11-11T23:41:09Z
transitions:
  - to: ready
    at: 2026-11-11T23:40:04Z
    by: system-flow
  - to: in-progress
    at: 2026-11-11T23:40:13Z
    by: agent-S-0728
tags: []
touches: [wip/kanban/stories/S-0728-publish-what-is-accepted-and-not-yet-released.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 94
  estimated: true
  priced_by: estimate
  turns:
    - day: 2026-11-11
      work: 17
  models:
    - model: claude-opus-5-5
      input: 36
      output: 225
      cache_read: 1993996
      cache_write: 147713
      cost: 0.7955
  strategic:
    - kind: orchestrator
      seconds: 17
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 6
          output: 31
          cache_read: 194798
          cache_write: 19067
          cost: 0.0606
formula: publish
---
# S-0728 Publish what is accepted and not yet released

## Goal

Publish what is accepted and not yet released, by the project's release policy.

## Acceptance criteria

- [x] Fetch the remote's branch and tags
- [ ] Evaluate the release policy, stopping unless it is met
- [ ] Publish what is accepted and not yet released

## Tasks
- T-3127 Fetch the remote's branch and tags
- T-3128 Evaluate the release policy, stopping unless it is met
- T-3129 Publish what is accepted and not yet released

## Notes

Made by `flai formula run publish` from `design/workflows/formulas/publish.md`, with the inputs `judgement=`.
