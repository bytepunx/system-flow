---
id: S-0646
type: story
nature: improvement
title: Publish what is accepted and not yet released
status: ready
owner: alex
created: 2026-11-09T22:22:55Z
updated: 2026-11-09T22:23:05Z
transitions:
  - to: ready
    at: 2026-11-09T22:23:05Z
    by: system-flow
tags: []
touches: [wip/kanban/stories/S-0646-publish-what-is-accepted-and-not-yet-released.md]
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
      seconds: 68
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 19
          output: 117
          cache_read: 518250
          cache_write: 78648
          cost: 0.1597
formula: publish
---
# S-0646 Publish what is accepted and not yet released

## Goal

Publish what is accepted and not yet released, by the project's release policy.

## Acceptance criteria

- [ ] Fetch the remote's branch and tags
- [ ] Evaluate the release policy, stopping unless it is met
- [ ] Publish what is accepted and not yet released

## Tasks
- T-2822 Fetch the remote's branch and tags
- T-2823 Evaluate the release policy, stopping unless it is met
- T-2824 Publish what is accepted and not yet released

## Notes

Made by `flai formula run publish` from `design/workflows/formulas/publish.md`, with no inputs.
