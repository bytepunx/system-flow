---
id: S-0769
type: story
nature: improvement
title: Publish what is accepted and not yet released
status: in-progress
owner: alex
created: 2026-11-13T12:40:07Z
updated: 2026-11-13T12:52:42Z
transitions:
  - to: ready
    at: 2026-11-13T12:40:09Z
    by: system-flow
  - to: in-progress
    at: 2026-11-13T12:52:42Z
    by: agent-S-0769
tags: []
touches: [wip/kanban/stories/S-0769-publish-what-is-accepted-and-not-yet-released.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 79
  estimated: true
  priced_by: estimate
  turns:
    - day: 2026-11-13
      hand_edits: 1
      work: 11
  models:
    - model: claude-opus-5-5
      input: 26
      output: 185
      cache_read: 1317244
      cache_write: 148136
      cost: 0.5472
  strategic:
    - kind: orchestrator
      seconds: 38
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 12
          output: 55
          cache_read: 410793
          cache_write: 21451
          cost: 0.1252
formula: publish
---
# S-0769 Publish what is accepted and not yet released

## Goal

Publish what is accepted and not yet released, by the project's release policy.

## Acceptance criteria

- [ ] Fetch the remote's branch and tags
- [ ] Evaluate the release policy, stopping unless it is met
- [ ] Publish what is accepted and not yet released

## Tasks
- T-3294 Fetch the remote's branch and tags
- T-3295 Evaluate the release policy, stopping unless it is met
- T-3296 Publish what is accepted and not yet released

## Notes

Made by `flai formula run publish` from `design/workflows/formulas/publish.md`, with the inputs `judgement=`.
