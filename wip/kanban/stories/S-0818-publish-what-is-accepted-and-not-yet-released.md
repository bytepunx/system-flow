---
id: S-0818
type: story
nature: improvement
title: Publish what is accepted and not yet released
status: in-progress
owner: alex
created: 2026-11-17T02:27:29Z
updated: 2026-11-17T02:30:25Z
transitions:
  - to: ready
    at: 2026-11-17T02:27:30Z
    by: system-flow
  - to: in-progress
    at: 2026-11-17T02:27:53Z
    by: agent-S-0818
tags: []
touches: [wip/kanban/stories/S-0818-publish-what-is-accepted-and-not-yet-released.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 162
  estimated: true
  priced_by: estimate
  turns:
    - day: 2026-11-17
      ceremony: 1
      work: 18
  models:
    - model: claude-opus-5-5
      input: 42
      output: 417
      cache_read: 2398543
      cache_write: 155414
      cost: 0.9712
  strategic:
    - kind: orchestrator
      seconds: 55
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 13
          output: 83
          cache_read: 505966
          cache_write: 85765
          cost: 0.1818
formula: publish
---
# S-0818 Publish what is accepted and not yet released

## Goal

Publish what is accepted and not yet released, by the project's release policy.

## Acceptance criteria

- [x] Fetch the remote's branch and tags
- [x] Evaluate the release policy, stopping unless it is met
- [ ] Publish what is accepted and not yet released

## Tasks
- T-3492 Fetch the remote's branch and tags
- T-3493 Evaluate the release policy, stopping unless it is met
- T-3494 Publish what is accepted and not yet released

## Notes

Made by `flai formula run publish` from `design/workflows/formulas/publish.md`, with the inputs `judgement=`.
