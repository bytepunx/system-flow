---
id: S-0657
type: story
nature: improvement
title: Publish what is accepted and not yet released
status: ready
owner: alex
created: 2026-11-10T07:11:22Z
updated: 2026-11-10T07:11:28Z
transitions:
  - to: ready
    at: 2026-11-10T07:11:28Z
    by: system-flow
tags: []
touches: [wip/kanban/stories/S-0657-publish-what-is-accepted-and-not-yet-released.md]
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
      seconds: 38
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 5
          output: 30
          cache_read: 246489
          cache_write: 22130
          cost: 0.0723
formula: publish
---
# S-0657 Publish what is accepted and not yet released

## Goal

Publish what is accepted and not yet released, by the project's release policy.

## Acceptance criteria

- [ ] Fetch the remote's branch and tags
- [ ] Evaluate the release policy, stopping unless it is met
- [ ] Publish what is accepted and not yet released

## Tasks
- T-2853 Fetch the remote's branch and tags
- T-2854 Evaluate the release policy, stopping unless it is met
- T-2855 Publish what is accepted and not yet released

## Notes

Made by `flai formula run publish` from `design/workflows/formulas/publish.md`, with the inputs `judgement=`.
