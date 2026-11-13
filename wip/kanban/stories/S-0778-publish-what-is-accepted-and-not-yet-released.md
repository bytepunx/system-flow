---
id: S-0778
type: story
nature: improvement
title: Publish what is accepted and not yet released
status: in-progress
owner: alex
created: 2026-11-13T20:23:36Z
updated: 2026-11-13T20:25:52Z
transitions:
  - to: ready
    at: 2026-11-13T20:23:37Z
    by: system-flow
  - to: in-progress
    at: 2026-11-13T20:23:45Z
    by: agent-S-0778
tags: []
touches: [wip/kanban/stories/S-0778-publish-what-is-accepted-and-not-yet-released.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 121
  estimated: true
  priced_by: estimate
  turns:
    - day: 2026-11-13
      ceremony: 1
      hand_edits: 2
      work: 15
  models:
    - model: claude-opus-5-5
      input: 38
      output: 373
      cache_read: 2172562
      cache_write: 150140
      cost: 0.8719
  strategic:
    - kind: orchestrator
      seconds: 23
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 10
          output: 78
          cache_read: 225848
          cache_write: 91586
          cost: 0.0933
formula: publish
---
# S-0778 Publish what is accepted and not yet released

## Goal

Publish what is accepted and not yet released, by the project's release policy.

## Acceptance criteria

- [x] Fetch the remote's branch and tags
- [x] Evaluate the release policy, stopping unless it is met
- [ ] Publish what is accepted and not yet released

## Tasks
- T-3326 Fetch the remote's branch and tags
- T-3327 Evaluate the release policy, stopping unless it is met
- T-3328 Publish what is accepted and not yet released

## Notes

Made by `flai formula run publish` from `design/workflows/formulas/publish.md`, with the inputs `judgement=`.
