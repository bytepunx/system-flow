---
id: S-0784
type: story
nature: improvement
title: Publish what is accepted and not yet released
status: in-progress
owner: alex
created: 2026-11-14T07:05:35Z
updated: 2026-11-14T07:05:55Z
transitions:
  - to: ready
    at: 2026-11-14T07:05:37Z
    by: system-flow
  - to: in-progress
    at: 2026-11-14T07:05:55Z
    by: agent-S-0784
tags: []
touches: [wip/kanban/stories/S-0784-publish-what-is-accepted-and-not-yet-released.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 117
  estimated: true
  priced_by: estimate
  turns:
    - day: 2026-11-14
      ceremony: 1
      hand_edits: 1
      work: 15
  models:
    - model: claude-opus-5-5
      input: 38
      output: 298
      cache_read: 2193358
      cache_write: 156270
      cost: 0.8832
  strategic:
    - kind: orchestrator
      seconds: 28
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 10
          output: 82
          cache_read: 308067
          cache_write: 53713
          cost: 0.107
formula: publish
---
# S-0784 Publish what is accepted and not yet released

## Goal

Publish what is accepted and not yet released, by the project's release policy.

## Acceptance criteria

- [ ] Fetch the remote's branch and tags
- [ ] Evaluate the release policy, stopping unless it is met
- [ ] Publish what is accepted and not yet released

## Tasks
- T-3351 Fetch the remote's branch and tags
- T-3352 Evaluate the release policy, stopping unless it is met
- T-3353 Publish what is accepted and not yet released

## Notes

Made by `flai formula run publish` from `design/workflows/formulas/publish.md`, with the inputs `judgement=`.
