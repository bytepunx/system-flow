---
id: S-0729
type: story
nature: improvement
title: Publish what is accepted and not yet released
status: in-progress
owner: alex
created: 2026-11-12T03:39:40Z
updated: 2026-11-12T03:40:05Z
transitions:
  - to: ready
    at: 2026-11-12T03:39:41Z
    by: system-flow
  - to: in-progress
    at: 2026-11-12T03:40:05Z
    by: agent-S-0729
tags: []
touches: [wip/kanban/stories/S-0729-publish-what-is-accepted-and-not-yet-released.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 91
  estimated: true
  priced_by: estimate
  turns:
    - day: 2026-11-12
      work: 16
  models:
    - model: claude-opus-5-5
      input: 36
      output: 226
      cache_read: 1842394
      cache_write: 151873
      cost: 0.7412
  strategic:
    - kind: orchestrator
      seconds: 29
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 9
          output: 59
          cache_read: 271384
          cache_write: 39938
          cost: 0.0883
formula: publish
---
# S-0729 Publish what is accepted and not yet released

## Goal

Publish what is accepted and not yet released, by the project's release policy.

## Acceptance criteria

- [ ] Fetch the remote's branch and tags
- [ ] Evaluate the release policy, stopping unless it is met
- [ ] Publish what is accepted and not yet released

## Tasks
- T-3132 Fetch the remote's branch and tags
- T-3133 Evaluate the release policy, stopping unless it is met
- T-3134 Publish what is accepted and not yet released

## Notes

Made by `flai formula run publish` from `design/workflows/formulas/publish.md`, with the inputs `judgement=`.
