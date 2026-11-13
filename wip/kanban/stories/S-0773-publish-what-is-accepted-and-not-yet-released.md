---
id: S-0773
type: story
nature: improvement
title: Publish what is accepted and not yet released
status: in-progress
owner: alex
created: 2026-11-13T15:20:30Z
updated: 2026-11-13T15:21:24Z
transitions:
  - to: ready
    at: 2026-11-13T15:20:32Z
    by: system-flow
  - to: in-progress
    at: 2026-11-13T15:21:24Z
    by: agent-S-0773
tags: []
touches: [wip/kanban/stories/S-0773-publish-what-is-accepted-and-not-yet-released.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 63
  estimated: true
  priced_by: estimate
  turns:
    - day: 2026-11-13
      work: 11
  models:
    - model: claude-opus-5-5
      input: 24
      output: 229
      cache_read: 955719
      cache_write: 145331
      cost: 0.4121
  strategic:
    - kind: orchestrator
      seconds: 37
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 6
          output: 38
          cache_read: 286426
          cache_write: 8982
          cost: 0.0859
formula: publish
---
# S-0773 Publish what is accepted and not yet released

## Goal

Publish what is accepted and not yet released, by the project's release policy.

## Acceptance criteria

- [ ] Fetch the remote's branch and tags
- [ ] Evaluate the release policy, stopping unless it is met
- [ ] Publish what is accepted and not yet released

## Tasks
- T-3307 Fetch the remote's branch and tags
- T-3308 Evaluate the release policy, stopping unless it is met
- T-3309 Publish what is accepted and not yet released

## Notes

Made by `flai formula run publish` from `design/workflows/formulas/publish.md`, with the inputs `judgement=`.
