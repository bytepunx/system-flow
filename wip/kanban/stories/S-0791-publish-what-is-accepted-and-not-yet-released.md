---
id: S-0791
type: story
nature: improvement
title: Publish what is accepted and not yet released
status: ready
owner: alex
created: 2026-11-15T07:18:49Z
updated: 2026-11-15T07:18:50Z
transitions:
  - to: ready
    at: 2026-11-15T07:18:50Z
    by: system-flow
tags: []
touches: [wip/kanban/stories/S-0791-publish-what-is-accepted-and-not-yet-released.md]
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
      seconds: 48
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 14
          output: 93
          cache_read: 558601
          cache_write: 51147
          cost: 0.1834
formula: publish
---
# S-0791 Publish what is accepted and not yet released

## Goal

Publish what is accepted and not yet released, by the project's release policy.

## Acceptance criteria

- [ ] Fetch the remote's branch and tags
- [ ] Evaluate the release policy, stopping unless it is met
- [ ] Publish what is accepted and not yet released

## Tasks
- T-3373 Fetch the remote's branch and tags
- T-3374 Evaluate the release policy, stopping unless it is met
- T-3375 Publish what is accepted and not yet released

## Notes

Made by `flai formula run publish` from `design/workflows/formulas/publish.md`, with the inputs `judgement=`.
