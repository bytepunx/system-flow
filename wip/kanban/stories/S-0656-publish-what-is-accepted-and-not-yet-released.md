---
id: S-0656
type: story
nature: improvement
title: Publish what is accepted and not yet released
status: in-progress
owner: alex
created: 2026-11-10T06:45:10Z
updated: 2026-11-10T06:58:07Z
transitions:
  - to: ready
    at: 2026-11-10T06:45:12Z
    by: system-flow
  - to: in-progress
    at: 2026-11-10T06:58:07Z
    by: agent-S-0656
tags: []
touches: [wip/kanban/stories/S-0656-publish-what-is-accepted-and-not-yet-released.md]
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
      seconds: 16
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 5
          output: 33
          cache_read: 127642
          cache_write: 16888
          cost: 0.0388
formula: publish
---
# S-0656 Publish what is accepted and not yet released

## Goal

Publish what is accepted and not yet released, by the project's release policy.

## Acceptance criteria

- [ ] Fetch the remote's branch and tags
- [ ] Evaluate the release policy, stopping unless it is met
- [ ] Publish what is accepted and not yet released

## Tasks
- T-2831 Fetch the remote's branch and tags
- T-2832 Evaluate the release policy, stopping unless it is met
- T-2833 Publish what is accepted and not yet released

## Notes

Made by `flai formula run publish` from `design/workflows/formulas/publish.md`, with the inputs `judgement=`.
