---
id: S-0771
type: story
nature: improvement
title: Publish what is accepted and not yet released
status: in-progress
owner: alex
created: 2026-11-13T13:41:40Z
updated: 2026-11-13T14:35:00Z
transitions:
  - to: ready
    at: 2026-11-13T13:41:41Z
    by: system-flow
  - to: in-progress
    at: 2026-11-13T14:35:00Z
    by: agent-S-0771
tags: []
touches: [wip/kanban/stories/S-0771-publish-what-is-accepted-and-not-yet-released.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 111
  estimated: true
  priced_by: estimate
  turns:
    - day: 2026-11-13
      ceremony: 1
      work: 15
  models:
    - model: claude-opus-5-5
      input: 36
      output: 264
      cache_read: 1819386
      cache_write: 155756
      cost: 0.7385
  strategic:
    - kind: orchestrator
      seconds: 40
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 10
          output: 67
          cache_read: 364419
          cache_write: 2820
          cost: 0.1066
formula: publish
---
# S-0771 Publish what is accepted and not yet released

## Goal

Publish what is accepted and not yet released, by the project's release policy.

## Acceptance criteria

- [ ] Fetch the remote's branch and tags
- [ ] Evaluate the release policy, stopping unless it is met
- [ ] Publish what is accepted and not yet released

## Tasks
- T-3300 Fetch the remote's branch and tags
- T-3301 Evaluate the release policy, stopping unless it is met
- T-3302 Publish what is accepted and not yet released

## Notes

Made by `flai formula run publish` from `design/workflows/formulas/publish.md`, with the inputs `judgement=`.
