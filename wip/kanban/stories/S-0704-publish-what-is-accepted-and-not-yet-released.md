---
id: S-0704
type: story
nature: improvement
title: Publish what is accepted and not yet released
status: ready
owner: alex
created: 2026-11-10T08:56:47Z
updated: 2026-11-10T08:56:57Z
transitions:
  - to: ready
    at: 2026-11-10T08:56:57Z
    by: system-flow
tags: []
touches: [wip/kanban/stories/S-0704-publish-what-is-accepted-and-not-yet-released.md]
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
      seconds: 107
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 13
          output: 77
          cache_read: 802966
          cache_write: 34172
          cost: 0.2275
        - model: claude-sonnet-5-5
          input: 1
          output: 11
          cache_read: 4883
          cache_write: 5260
          cost: 0.0076
formula: publish
---
# S-0704 Publish what is accepted and not yet released

## Goal

Publish what is accepted and not yet released, by the project's release policy.

## Acceptance criteria

- [ ] Fetch the remote's branch and tags
- [ ] Evaluate the release policy, stopping unless it is met
- [ ] Publish what is accepted and not yet released

## Tasks
- T-2952 Fetch the remote's branch and tags
- T-2953 Evaluate the release policy, stopping unless it is met
- T-2954 Publish what is accepted and not yet released

## Notes

Made by `flai formula run publish` from `design/workflows/formulas/publish.md`, with the inputs `judgement=`.
