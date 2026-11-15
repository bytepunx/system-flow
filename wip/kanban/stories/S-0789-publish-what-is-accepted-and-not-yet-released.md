---
id: S-0789
type: story
nature: improvement
title: Publish what is accepted and not yet released
status: in-progress
owner: alex
created: 2026-11-14T23:49:29Z
updated: 2026-11-15T01:56:58Z
transitions:
  - to: ready
    at: 2026-11-14T23:49:30Z
    by: system-flow
  - to: in-progress
    at: 2026-11-15T01:56:58Z
    by: agent-S-0789
tags: []
touches: [wip/kanban/stories/S-0789-publish-what-is-accepted-and-not-yet-released.md]
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
      seconds: 29
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 8
          output: 41
          cache_read: 299283
          cache_write: 36332
          cost: 0.1003
formula: publish
---
# S-0789 Publish what is accepted and not yet released

## Goal

Publish what is accepted and not yet released, by the project's release policy.

## Acceptance criteria

- [ ] Fetch the remote's branch and tags
- [ ] Evaluate the release policy, stopping unless it is met
- [ ] Publish what is accepted and not yet released

## Tasks
- T-3366 Fetch the remote's branch and tags
- T-3367 Evaluate the release policy, stopping unless it is met
- T-3368 Publish what is accepted and not yet released

## Notes

Made by `flai formula run publish` from `design/workflows/formulas/publish.md`, with the inputs `judgement=`.
