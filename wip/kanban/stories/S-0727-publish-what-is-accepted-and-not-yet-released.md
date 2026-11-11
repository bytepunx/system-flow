---
id: S-0727
type: story
nature: improvement
title: Publish what is accepted and not yet released
status: in-progress
owner: alex
created: 2026-11-11T23:22:51Z
updated: 2026-11-11T23:23:01Z
transitions:
  - to: ready
    at: 2026-11-11T23:22:52Z
    by: system-flow
  - to: in-progress
    at: 2026-11-11T23:23:01Z
    by: agent-S-0727
tags: []
touches: [wip/kanban/stories/S-0727-publish-what-is-accepted-and-not-yet-released.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 93
  estimated: true
  priced_by: estimate
  turns:
    - day: 2026-11-11
      ceremony: 1
      work: 15
  models:
    - model: claude-opus-5-5
      input: 36
      output: 275
      cache_read: 1828233
      cache_write: 146642
      cost: 0.7334
  strategic:
    - kind: orchestrator
      seconds: 15
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 7
          output: 67
          cache_read: 233265
          cache_write: 25734
          cost: 0.0733
formula: publish
---
# S-0727 Publish what is accepted and not yet released

## Goal

Publish what is accepted and not yet released, by the project's release policy.

## Acceptance criteria

- [ ] Fetch the remote's branch and tags
- [ ] Evaluate the release policy, stopping unless it is met
- [ ] Publish what is accepted and not yet released

## Tasks
- T-3124 Fetch the remote's branch and tags
- T-3125 Evaluate the release policy, stopping unless it is met
- T-3126 Publish what is accepted and not yet released

## Notes

Made by `flai formula run publish` from `design/workflows/formulas/publish.md`, with the inputs `judgement=`.
