---
id: S-0816
type: story
nature: improvement
title: Publish what is accepted and not yet released
status: in-progress
owner: alex
created: 2026-11-17T01:33:37Z
updated: 2026-11-17T01:33:58Z
transitions:
  - to: ready
    at: 2026-11-17T01:33:39Z
    by: system-flow
  - to: in-progress
    at: 2026-11-17T01:33:58Z
    by: agent-S-0816
tags: []
touches: [wip/kanban/stories/S-0816-publish-what-is-accepted-and-not-yet-released.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 134
  estimated: true
  priced_by: estimate
  turns:
    - day: 2026-11-17
      ceremony: 1
      work: 15
  models:
    - model: claude-opus-5-5
      input: 34
      output: 340
      cache_read: 1619468
      cache_write: 149679
      cost: 0.6724
  strategic:
    - kind: orchestrator
      seconds: 56
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 9
          output: 47
          cache_read: 413495
          cache_write: 40106
          cost: 0.1392
formula: publish
---
# S-0816 Publish what is accepted and not yet released

## Goal

Publish what is accepted and not yet released, by the project's release policy.

## Acceptance criteria

- [ ] Fetch the remote's branch and tags
- [ ] Evaluate the release policy, stopping unless it is met
- [ ] Publish what is accepted and not yet released

## Tasks
- T-3488 Fetch the remote's branch and tags
- T-3489 Evaluate the release policy, stopping unless it is met
- T-3490 Publish what is accepted and not yet released

## Notes

Made by `flai formula run publish` from `design/workflows/formulas/publish.md`, with the inputs `judgement=`.
