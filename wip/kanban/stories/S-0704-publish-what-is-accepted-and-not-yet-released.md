---
id: S-0704
type: story
nature: improvement
title: Publish what is accepted and not yet released
status: in-progress
owner: alex
created: 2026-11-10T08:56:47Z
updated: 2026-11-11T01:41:18Z
transitions:
  - to: ready
    at: 2026-11-10T08:56:57Z
    by: system-flow
  - to: in-progress
    at: 2026-11-11T01:41:18Z
    by: agent-S-0704
tags: []
touches: [wip/kanban/stories/S-0704-publish-what-is-accepted-and-not-yet-released.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 84
  estimated: true
  priced_by: estimate
  turns:
    - day: 2026-11-11
      work: 15
  models:
    - model: claude-opus-5-5
      input: 34
      output: 288
      cache_read: 1776225
      cache_write: 154097
      cost: 0.7116
  strategic:
    - kind: orchestrator
      seconds: 226
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 44
          output: 234
          cache_read: 1674458
          cache_write: 142397
          cost: 0.4961
        - model: claude-sonnet-5-5
          input: 7
          output: 57
          cache_read: 75831
          cache_write: 23201
          cost: 0.0743
formula: publish
---
# S-0704 Publish what is accepted and not yet released

## Goal

Publish what is accepted and not yet released, by the project's release policy.

## Acceptance criteria

- [ ] Fetch the remote's branch and tags
- [ ] Evaluate the release policy, stopping unless it is met
- [ ] Publish what is accepted and not yet released

## Tasks
- T-2952 Fetch the remote's branch and tags
- T-2953 Evaluate the release policy, stopping unless it is met
- T-2954 Publish what is accepted and not yet released

## Notes

Made by `flai formula run publish` from `design/workflows/formulas/publish.md`, with the inputs `judgement=`.
