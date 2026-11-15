---
id: S-0791
type: story
nature: improvement
title: Publish what is accepted and not yet released
status: in-progress
owner: alex
created: 2026-11-15T07:18:49Z
updated: 2026-11-15T20:59:23Z
transitions:
  - to: ready
    at: 2026-11-15T07:18:50Z
    by: system-flow
  - to: in-progress
    at: 2026-11-15T20:59:23Z
    by: agent-S-0791
tags: []
touches: [wip/kanban/stories/S-0791-publish-what-is-accepted-and-not-yet-released.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 94
  estimated: true
  priced_by: estimate
  turns:
    - day: 2026-11-15
      ceremony: 1
      work: 14
  models:
    - model: claude-opus-5-5
      input: 32
      output: 310
      cache_read: 1725233
      cache_write: 152456
      cost: 0.71
  strategic:
    - kind: orchestrator
      seconds: 76
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 24
          output: 142
          cache_read: 929169
          cache_write: 90690
          cost: 0.3071
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
