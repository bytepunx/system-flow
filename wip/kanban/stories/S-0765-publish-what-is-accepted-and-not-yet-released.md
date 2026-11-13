---
id: S-0765
type: story
nature: improvement
title: Publish what is accepted and not yet released
status: in-progress
owner: alex
created: 2026-11-13T09:17:34Z
updated: 2026-11-13T09:17:43Z
transitions:
  - to: ready
    at: 2026-11-13T09:17:36Z
    by: system-flow
  - to: in-progress
    at: 2026-11-13T09:17:43Z
    by: agent-S-0765
tags: []
touches: [wip/kanban/stories/S-0765-publish-what-is-accepted-and-not-yet-released.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 73
  estimated: true
  priced_by: estimate
  turns:
    - day: 2026-11-13
      work: 14
  models:
    - model: claude-opus-5-5
      input: 30
      output: 233
      cache_read: 1694233
      cache_write: 148700
      cost: 0.687
  strategic:
    - kind: orchestrator
      seconds: 24
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 6
          output: 26
          cache_read: 314705
          cache_write: 3897
          cost: 0.0919
formula: publish
---
# S-0765 Publish what is accepted and not yet released

## Goal

Publish what is accepted and not yet released, by the project's release policy.

## Acceptance criteria

- [ ] Fetch the remote's branch and tags
- [ ] Evaluate the release policy, stopping unless it is met
- [ ] Publish what is accepted and not yet released

## Tasks
- T-3281 Fetch the remote's branch and tags
- T-3282 Evaluate the release policy, stopping unless it is met
- T-3283 Publish what is accepted and not yet released

## Notes

Made by `flai formula run publish` from `design/workflows/formulas/publish.md`, with the inputs `judgement=`.
