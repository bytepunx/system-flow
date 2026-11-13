---
id: S-0778
type: story
nature: improvement
title: Publish what is accepted and not yet released
status: done
owner: alex
created: 2026-11-13T20:23:36Z
updated: 2026-11-13T22:47:29Z
transitions:
  - to: ready
    at: 2026-11-13T20:23:37Z
    by: system-flow
  - to: in-progress
    at: 2026-11-13T20:23:45Z
    by: agent-S-0778
  - to: review
    at: 2026-11-13T20:27:15Z
    by: system-flow
  - to: done
    at: 2026-11-13T22:47:29Z
    by: alex
tags: []
touches: [wip/kanban/stories/S-0778-publish-what-is-accepted-and-not-yet-released.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 236
  priced_by: harness
  turns:
    - day: 2026-11-13
      ceremony: 4
      hand_edits: 3
      work: 24
  models:
    - model: claude-opus-5-5
      input: 66
      output: 7476
      cache_read: 4565226
      cache_write: 159669
      cost: 2.3402
  strategic:
    - kind: orchestrator
      seconds: 95
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 26
          output: 224
          cache_read: 1065657
          cache_write: 100098
          cost: 0.3425
formula: publish
---
# S-0778 Publish what is accepted and not yet released

## Goal

Publish what is accepted and not yet released, by the project's release policy.

## Acceptance criteria

- [x] Fetch the remote's branch and tags
- [x] Evaluate the release policy, stopping unless it is met
- [x] Publish what is accepted and not yet released

## Tasks
- T-3326 Fetch the remote's branch and tags
- T-3327 Evaluate the release policy, stopping unless it is met
- T-3328 Publish what is accepted and not yet released

## Notes

Made by `flai formula run publish` from `design/workflows/formulas/publish.md`, with the inputs `judgement=`.
