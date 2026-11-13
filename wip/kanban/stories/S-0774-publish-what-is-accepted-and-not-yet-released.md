---
id: S-0774
type: story
nature: improvement
title: Publish what is accepted and not yet released
status: in-progress
owner: alex
created: 2026-11-13T15:43:49Z
updated: 2026-11-13T15:44:00Z
transitions:
  - to: ready
    at: 2026-11-13T15:43:50Z
    by: system-flow
  - to: in-progress
    at: 2026-11-13T15:44:00Z
    by: agent-S-0774
tags: []
touches: [wip/kanban/stories/S-0774-publish-what-is-accepted-and-not-yet-released.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 74
  estimated: true
  priced_by: estimate
  turns:
    - day: 2026-11-13
      work: 16
  models:
    - model: claude-opus-5-5
      input: 34
      output: 343
      cache_read: 1908677
      cache_write: 149585
      cost: 0.7706
  strategic:
    - kind: orchestrator
      seconds: 11
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 2
          output: 30
          cache_read: 106489
          cache_write: 519
          cost: 0.0312
formula: publish
---
# S-0774 Publish what is accepted and not yet released

## Goal

Publish what is accepted and not yet released, by the project's release policy.

## Acceptance criteria

- [ ] Fetch the remote's branch and tags
- [ ] Evaluate the release policy, stopping unless it is met
- [ ] Publish what is accepted and not yet released

## Tasks
- T-3310 Fetch the remote's branch and tags
- T-3311 Evaluate the release policy, stopping unless it is met
- T-3312 Publish what is accepted and not yet released

## Notes

Made by `flai formula run publish` from `design/workflows/formulas/publish.md`, with the inputs `judgement=`.
