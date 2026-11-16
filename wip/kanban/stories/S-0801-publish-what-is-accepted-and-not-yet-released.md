---
id: S-0801
type: story
nature: improvement
title: Publish what is accepted and not yet released
status: review
owner: alex
created: 2026-11-16T01:54:24Z
updated: 2026-11-16T03:53:50Z
transitions:
  - to: ready
    at: 2026-11-16T01:54:26Z
    by: system-flow
  - to: in-progress
    at: 2026-11-16T03:51:37Z
    by: agent-S-0801
  - to: review
    at: 2026-11-16T03:53:50Z
    by: agent-S-0801
tags: []
touches: [wip/kanban/stories/S-0801-publish-what-is-accepted-and-not-yet-released.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 146
  priced_by: harness
  turns:
    - day: 2026-11-16
      ceremony: 2
      work: 29
  models:
    - model: claude-opus-5-5
      input: 66
      output: 7686
      cache_read: 4555352
      cache_write: 162556
      cost: 2.3655
  strategic:
    - kind: orchestrator
      seconds: 52
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 23
          output: 192
          cache_read: 738184
          cache_write: 134997
          cost: 0.2664
        - model: claude-sonnet-5-5
          input: 2
          output: 24
          cache_read: 7329
          cache_write: 6164
          cost: 0.0108
formula: publish
---
# S-0801 Publish what is accepted and not yet released

## Goal

Publish what is accepted and not yet released, by the project's release policy.

## Acceptance criteria

- [x] Fetch the remote's branch and tags
- [x] Evaluate the release policy, stopping unless it is met
- [x] Publish what is accepted and not yet released

## Tasks
- T-3420 Fetch the remote's branch and tags
- T-3421 Evaluate the release policy, stopping unless it is met
- T-3422 Publish what is accepted and not yet released

## Notes

Made by `flai formula run publish` from `design/workflows/formulas/publish.md`, with the inputs `judgement=`.
