---
id: S-0801
type: story
nature: improvement
title: Publish what is accepted and not yet released
status: in-progress
owner: alex
created: 2026-11-16T01:54:24Z
updated: 2026-11-16T03:51:37Z
transitions:
  - to: ready
    at: 2026-11-16T01:54:26Z
    by: system-flow
  - to: in-progress
    at: 2026-11-16T03:51:37Z
    by: agent-S-0801
tags: []
touches: [wip/kanban/stories/S-0801-publish-what-is-accepted-and-not-yet-released.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 81
  estimated: true
  priced_by: estimate
  turns:
    - day: 2026-11-16
      ceremony: 1
      work: 18
  models:
    - model: claude-opus-5-5
      input: 40
      output: 347
      cache_read: 2299912
      cache_write: 154002
      cost: 0.9305
  strategic:
    - kind: orchestrator
      seconds: 18
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 9
          output: 58
          cache_read: 272937
          cache_write: 62551
          cost: 0.1023
formula: publish
---
# S-0801 Publish what is accepted and not yet released

## Goal

Publish what is accepted and not yet released, by the project's release policy.

## Acceptance criteria

- [ ] Fetch the remote's branch and tags
- [ ] Evaluate the release policy, stopping unless it is met
- [ ] Publish what is accepted and not yet released

## Tasks
- T-3420 Fetch the remote's branch and tags
- T-3421 Evaluate the release policy, stopping unless it is met
- T-3422 Publish what is accepted and not yet released

## Notes

Made by `flai formula run publish` from `design/workflows/formulas/publish.md`, with the inputs `judgement=`.
