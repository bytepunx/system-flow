---
id: S-0731
type: story
nature: improvement
title: Publish what is accepted and not yet released
status: in-progress
owner: alex
created: 2026-11-12T03:44:58Z
updated: 2026-11-12T03:45:07Z
transitions:
  - to: ready
    at: 2026-11-12T03:44:59Z
    by: system-flow
  - to: in-progress
    at: 2026-11-12T03:45:07Z
    by: agent-S-0731
tags: []
touches: [wip/kanban/stories/S-0731-publish-what-is-accepted-and-not-yet-released.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 119
  estimated: true
  priced_by: estimate
  turns:
    - day: 2026-11-12
      hand_edits: 3
      work: 13
  models:
    - model: claude-opus-5-5
      input: 34
      output: 274
      cache_read: 1838507
      cache_write: 149490
      cost: 0.739
  strategic:
    - kind: orchestrator
      seconds: 7
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 3
          output: 22
          cache_read: 75717
          cache_write: 15259
          cost: 0.0258
formula: publish
---
# S-0731 Publish what is accepted and not yet released

## Goal

Publish what is accepted and not yet released, by the project's release policy.

## Acceptance criteria

- [ ] Fetch the remote's branch and tags
- [ ] Evaluate the release policy, stopping unless it is met
- [ ] Publish what is accepted and not yet released

## Tasks
- T-3135 Fetch the remote's branch and tags
- T-3136 Evaluate the release policy, stopping unless it is met
- T-3137 Publish what is accepted and not yet released

## Notes

Made by `flai formula run publish` from `design/workflows/formulas/publish.md`, with the inputs `judgement=`.
