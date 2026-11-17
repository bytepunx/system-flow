---
id: S-0826
type: story
nature: improvement
title: Publish what is accepted and not yet released
status: in-progress
owner: alex
created: 2026-11-17T08:16:46Z
updated: 2026-11-17T08:17:31Z
transitions:
  - to: ready
    at: 2026-11-17T08:16:48Z
    by: system-flow
  - to: in-progress
    at: 2026-11-17T08:17:31Z
    by: agent-S-0826
tags: []
touches: [wip/kanban/stories/S-0826-publish-what-is-accepted-and-not-yet-released.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 118
  estimated: true
  priced_by: estimate
  turns:
    - day: 2026-11-17
      ceremony: 1
      hand_edits: 1
      work: 14
  models:
    - model: claude-opus-5-5
      input: 34
      output: 290
      cache_read: 1899338
      cache_write: 156613
      cost: 0.7834
  strategic:
    - kind: orchestrator
      seconds: 22
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 3
          output: 15
          cache_read: 203350
          cache_write: 13983
          cost: 0.0673
formula: publish
---
# S-0826 Publish what is accepted and not yet released

## Goal

Publish what is accepted and not yet released, by the project's release policy.

## Acceptance criteria

- [ ] Fetch the remote's branch and tags
- [ ] Evaluate the release policy, stopping unless it is met
- [ ] Publish what is accepted and not yet released

## Tasks
- T-3517 Fetch the remote's branch and tags
- T-3518 Evaluate the release policy, stopping unless it is met
- T-3519 Publish what is accepted and not yet released

## Notes

Made by `flai formula run publish` from `design/workflows/formulas/publish.md`, with the inputs `judgement=`.
