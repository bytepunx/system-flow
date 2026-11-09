---
id: S-0635
type: story
nature: improvement
title: Publish what is accepted and not yet released
status: review
owner: alex
created: 2026-11-08T01:14:12Z
updated: 2026-11-09T11:24:53Z
transitions:
  - to: ready
    at: 2026-11-08T01:15:24Z
    by: system-flow
  - to: in-progress
    at: 2026-11-09T11:23:00Z
    by: agent-S-0635
  - to: review
    at: 2026-11-09T11:24:53Z
    by: system-flow
tags: []
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 122
  priced_by: harness
  turns:
    - day: 2026-11-09
      ceremony: 3
      work: 22
  models:
    - model: claude-opus-5-5
      input: 54
      output: 7113
      cache_read: 3489563
      cache_write: 150043
      cost: 2.0407
  strategic:
    - kind: orchestrator
      seconds: 666
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 137
          output: 834
          cache_read: 5309203
          cache_write: 578667
          cost: 1.5314
        - model: claude-sonnet-5-5
          input: 6
          output: 47
          cache_read: 93704
          cache_write: 15803
          cost: 0.0809
formula: publish
---
# S-0635 Publish what is accepted and not yet released

## Goal

Publish what is accepted and not yet released, by the project's release policy.

## Acceptance criteria

- [x] Fetch the remote's branch and tags
- [x] Evaluate the release policy, stopping unless it is met
- [x] Publish what is accepted and not yet released

## Tasks
- T-2694 Fetch the remote's branch and tags
- T-2695 Evaluate the release policy, stopping unless it is met
- T-2696 Publish what is accepted and not yet released

## Notes

Made by `flai formula run publish` from `design/workflows/formulas/publish.md`, with no inputs.
