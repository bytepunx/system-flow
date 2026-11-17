---
id: S-0820
type: story
nature: improvement
title: Publish what is accepted and not yet released
status: in-progress
owner: alex
created: 2026-11-17T02:48:43Z
updated: 2026-11-17T03:23:30Z
transitions:
  - to: ready
    at: 2026-11-17T02:48:46Z
    by: system-flow
  - to: in-progress
    at: 2026-11-17T03:23:30Z
    by: agent-S-0820
tags: []
touches: [wip/kanban/stories/S-0820-publish-what-is-accepted-and-not-yet-released.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 160
  estimated: true
  priced_by: estimate
  turns:
    - day: 2026-11-17
      work: 17
  models:
    - model: claude-opus-5-5
      input: 36
      output: 285
      cache_read: 2087829
      cache_write: 154930
      cost: 0.8535
  strategic:
    - kind: orchestrator
      seconds: 39
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 10
          output: 74
          cache_read: 297286
          cache_write: 48245
          cost: 0.1064
formula: publish
---
# S-0820 Publish what is accepted and not yet released

## Goal

Publish what is accepted and not yet released, by the project's release policy.

## Acceptance criteria

- [ ] Fetch the remote's branch and tags
- [ ] Evaluate the release policy, stopping unless it is met
- [ ] Publish what is accepted and not yet released

## Tasks
- T-3500 Fetch the remote's branch and tags
- T-3501 Evaluate the release policy, stopping unless it is met
- T-3502 Publish what is accepted and not yet released

## Notes

Made by `flai formula run publish` from `design/workflows/formulas/publish.md`, with the inputs `judgement=`.
