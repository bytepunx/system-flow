---
id: S-0794
type: story
nature: improvement
title: Publish what is accepted and not yet released
status: in-progress
owner: alex
created: 2026-11-15T21:47:08Z
updated: 2026-11-16T00:38:45Z
transitions:
  - to: ready
    at: 2026-11-15T21:47:09Z
    by: system-flow
  - to: in-progress
    at: 2026-11-16T00:36:45Z
    by: agent-S-0794
tags: []
touches: [wip/kanban/stories/S-0794-publish-what-is-accepted-and-not-yet-released.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 131
  estimated: true
  priced_by: estimate
  turns:
    - day: 2026-11-16
      ceremony: 1
      work: 17
  models:
    - model: claude-opus-5-5
      input: 38
      output: 313
      cache_read: 2093213
      cache_write: 155137
      cost: 0.8514
  strategic:
    - kind: orchestrator
      seconds: 41
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 9
          output: 47
          cache_read: 245547
          cache_write: 24490
          cost: 0.0818
formula: publish
---
# S-0794 Publish what is accepted and not yet released

## Goal

Publish what is accepted and not yet released, by the project's release policy.

## Acceptance criteria

- [x] Fetch the remote's branch and tags
- [x] Evaluate the release policy, stopping unless it is met
- [ ] Publish what is accepted and not yet released

## Tasks
- T-3387 Fetch the remote's branch and tags
- T-3388 Evaluate the release policy, stopping unless it is met
- T-3389 Publish what is accepted and not yet released

## Notes

Made by `flai formula run publish` from `design/workflows/formulas/publish.md`, with the inputs `judgement=`.
