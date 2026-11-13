---
id: S-0780
type: story
nature: improvement
title: Publish what is accepted and not yet released
status: in-progress
owner: alex
created: 2026-11-13T20:51:49Z
updated: 2026-11-13T20:53:59Z
transitions:
  - to: ready
    at: 2026-11-13T20:51:50Z
    by: system-flow
  - to: in-progress
    at: 2026-11-13T20:52:08Z
    by: agent-S-0780
tags: []
touches: [wip/kanban/stories/S-0780-publish-what-is-accepted-and-not-yet-released.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 118
  estimated: true
  priced_by: estimate
  turns:
    - day: 2026-11-13
      ceremony: 1
      work: 17
  models:
    - model: claude-opus-5-5
      input: 38
      output: 290
      cache_read: 2395011
      cache_write: 156724
      cost: 0.958
  strategic:
    - kind: orchestrator
      seconds: 35
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 9
          output: 64
          cache_read: 390727
          cache_write: 43963
          cost: 0.128
formula: publish
---
# S-0780 Publish what is accepted and not yet released

## Goal

Publish what is accepted and not yet released, by the project's release policy.

## Acceptance criteria

- [x] Fetch the remote's branch and tags
- [x] Evaluate the release policy, stopping unless it is met
- [ ] Publish what is accepted and not yet released

## Tasks
- T-3334 Fetch the remote's branch and tags
- T-3335 Evaluate the release policy, stopping unless it is met
- T-3336 Publish what is accepted and not yet released

## Notes

Made by `flai formula run publish` from `design/workflows/formulas/publish.md`, with the inputs `judgement=`.
