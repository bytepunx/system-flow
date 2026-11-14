---
id: S-0786
type: story
nature: improvement
title: Publish what is accepted and not yet released
status: in-progress
owner: alex
created: 2026-11-14T17:32:28Z
updated: 2026-11-14T17:32:37Z
transitions:
  - to: ready
    at: 2026-11-14T17:32:29Z
    by: system-flow
  - to: in-progress
    at: 2026-11-14T17:32:37Z
    by: agent-S-0786
tags: []
touches: [wip/kanban/stories/S-0786-publish-what-is-accepted-and-not-yet-released.md]
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
    - day: 2026-11-14
      ceremony: 1
      work: 16
  models:
    - model: claude-opus-5-5
      input: 38
      output: 333
      cache_read: 2036265
      cache_write: 155456
      cost: 0.8244
  strategic:
    - kind: orchestrator
      seconds: 9
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 5
          output: 30
          cache_read: 130005
          cache_write: 36925
          cost: 0.0495
formula: publish
---
# S-0786 Publish what is accepted and not yet released

## Goal

Publish what is accepted and not yet released, by the project's release policy.

## Acceptance criteria

- [ ] Fetch the remote's branch and tags
- [ ] Evaluate the release policy, stopping unless it is met
- [ ] Publish what is accepted and not yet released

## Tasks
- T-3358 Fetch the remote's branch and tags
- T-3359 Evaluate the release policy, stopping unless it is met
- T-3360 Publish what is accepted and not yet released

## Notes

Made by `flai formula run publish` from `design/workflows/formulas/publish.md`, with the inputs `judgement=`.
