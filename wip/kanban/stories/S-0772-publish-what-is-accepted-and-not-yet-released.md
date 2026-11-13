---
id: S-0772
type: story
nature: improvement
title: Publish what is accepted and not yet released
status: in-progress
owner: alex
created: 2026-11-13T13:42:42Z
updated: 2026-11-13T15:19:40Z
transitions:
  - to: ready
    at: 2026-11-13T13:42:43Z
    by: system-flow
  - to: in-progress
    at: 2026-11-13T15:18:17Z
    by: agent-S-0772
tags: []
touches: [wip/kanban/stories/S-0772-publish-what-is-accepted-and-not-yet-released.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 86
  estimated: true
  priced_by: estimate
  turns:
    - day: 2026-11-13
      ceremony: 1
      work: 14
  models:
    - model: claude-opus-5-5
      input: 32
      output: 314
      cache_read: 1662575
      cache_write: 151766
      cost: 0.679
  strategic:
    - kind: orchestrator
      seconds: 19
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 7
          output: 41
          cache_read: 173745
          cache_write: 17026
          cost: 0.0554
formula: publish
---
# S-0772 Publish what is accepted and not yet released

## Goal

Publish what is accepted and not yet released, by the project's release policy.

## Acceptance criteria

- [ ] Fetch the remote's branch and tags
- [x] Evaluate the release policy, stopping unless it is met
- [x] Publish what is accepted and not yet released (skipped: step evaluate stopped the run with exit 2)

## Tasks
- T-3303 Fetch the remote's branch and tags
- T-3304 Evaluate the release policy, stopping unless it is met
- T-3305 Publish what is accepted and not yet released

## Notes

Made by `flai formula run publish` from `design/workflows/formulas/publish.md`, with the inputs `judgement=`.

Step evaluate stopped the run with exit 2; skipped T-3305.
