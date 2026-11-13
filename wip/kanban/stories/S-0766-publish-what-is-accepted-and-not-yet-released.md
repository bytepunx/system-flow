---
id: S-0766
type: story
nature: improvement
title: Publish what is accepted and not yet released
status: in-progress
owner: alex
created: 2026-11-13T11:40:39Z
updated: 2026-11-13T11:40:48Z
transitions:
  - to: ready
    at: 2026-11-13T11:40:40Z
    by: system-flow
  - to: in-progress
    at: 2026-11-13T11:40:48Z
    by: agent-S-0766
tags: []
touches: [wip/kanban/stories/S-0766-publish-what-is-accepted-and-not-yet-released.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 83
  estimated: true
  priced_by: estimate
  turns:
    - day: 2026-11-13
      hand_edits: 1
      work: 15
  models:
    - model: claude-opus-5-5
      input: 34
      output: 248
      cache_read: 1864497
      cache_write: 155470
      cost: 0.7536
  strategic:
    - kind: orchestrator
      seconds: 20
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 8
          output: 56
          cache_read: 188813
          cache_write: 36408
          cost: 0.0651
formula: publish
---
# S-0766 Publish what is accepted and not yet released

## Goal

Publish what is accepted and not yet released, by the project's release policy.

## Acceptance criteria

- [ ] Fetch the remote's branch and tags
- [ ] Evaluate the release policy, stopping unless it is met
- [ ] Publish what is accepted and not yet released

## Tasks
- T-3284 Fetch the remote's branch and tags
- T-3285 Evaluate the release policy, stopping unless it is met
- T-3286 Publish what is accepted and not yet released

## Notes

Made by `flai formula run publish` from `design/workflows/formulas/publish.md`, with the inputs `judgement=`.
