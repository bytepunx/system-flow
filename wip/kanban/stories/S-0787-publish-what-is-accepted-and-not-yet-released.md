---
id: S-0787
type: story
nature: improvement
title: Publish what is accepted and not yet released
status: in-progress
owner: alex
created: 2026-11-14T18:00:18Z
updated: 2026-11-14T18:00:31Z
transitions:
  - to: ready
    at: 2026-11-14T18:00:20Z
    by: system-flow
  - to: in-progress
    at: 2026-11-14T18:00:31Z
    by: agent-S-0787
tags: []
touches: [wip/kanban/stories/S-0787-publish-what-is-accepted-and-not-yet-released.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 86
  estimated: true
  priced_by: estimate
  turns:
    - day: 2026-11-14
      ceremony: 1
      hand_edits: 2
      work: 11
  models:
    - model: claude-opus-5-5
      input: 30
      output: 225
      cache_read: 1559368
      cache_write: 154205
      cost: 0.6451
  strategic:
    - kind: orchestrator
      seconds: 24
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 9
          output: 42
          cache_read: 407127
          cache_write: 48419
          cost: 0.1357
formula: publish
---
# S-0787 Publish what is accepted and not yet released

## Goal

Publish what is accepted and not yet released, by the project's release policy.

## Acceptance criteria

- [ ] Fetch the remote's branch and tags
- [ ] Evaluate the release policy, stopping unless it is met
- [ ] Publish what is accepted and not yet released

## Tasks
- T-3361 Fetch the remote's branch and tags
- T-3362 Evaluate the release policy, stopping unless it is met
- T-3363 Publish what is accepted and not yet released

## Notes

Made by `flai formula run publish` from `design/workflows/formulas/publish.md`, with the inputs `judgement=`.
