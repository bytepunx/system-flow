---
id: S-0768
type: story
nature: improvement
title: Publish what is accepted and not yet released
status: in-progress
owner: alex
created: 2026-11-13T11:43:32Z
updated: 2026-11-13T11:43:41Z
transitions:
  - to: ready
    at: 2026-11-13T11:43:34Z
    by: system-flow
  - to: in-progress
    at: 2026-11-13T11:43:41Z
    by: agent-S-0768
tags: []
touches: [wip/kanban/stories/S-0768-publish-what-is-accepted-and-not-yet-released.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 82
  estimated: true
  priced_by: estimate
  turns:
    - day: 2026-11-13
      ceremony: 1
      work: 13
  models:
    - model: claude-opus-5-5
      input: 30
      output: 244
      cache_read: 1536851
      cache_write: 154950
      cost: 0.6313
  strategic:
    - kind: orchestrator
      seconds: 31
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 14
          output: 86
          cache_read: 388237
          cache_write: 19229
          cost: 0.1178
formula: publish
---
# S-0768 Publish what is accepted and not yet released

## Goal

Publish what is accepted and not yet released, by the project's release policy.

## Acceptance criteria

- [ ] Fetch the remote's branch and tags
- [ ] Evaluate the release policy, stopping unless it is met
- [ ] Publish what is accepted and not yet released

## Tasks
- T-3288 Fetch the remote's branch and tags
- T-3289 Evaluate the release policy, stopping unless it is met
- T-3290 Publish what is accepted and not yet released

## Notes

Made by `flai formula run publish` from `design/workflows/formulas/publish.md`, with the inputs `judgement=`.
