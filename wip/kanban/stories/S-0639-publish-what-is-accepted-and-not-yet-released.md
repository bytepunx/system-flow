---
id: S-0639
type: story
nature: improvement
title: Publish what is accepted and not yet released
status: in-progress
owner: alex
created: 2026-11-09T11:24:54Z
updated: 2026-11-09T11:25:00Z
transitions:
  - to: ready
    at: 2026-11-09T11:24:55Z
    by: system-flow
  - to: in-progress
    at: 2026-11-09T11:25:00Z
    by: agent-S-0639
tags: []
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 59
  estimated: true
  priced_by: estimate
  turns:
    - day: 2026-11-09
      ceremony: 1
      work: 12
  models:
    - model: claude-opus-5-5
      input: 28
      output: 226
      cache_read: 1380338
      cache_write: 138655
      cost: 0.5539
  strategic:
    - kind: orchestrator
      seconds: 28
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 6
          output: 33
          cache_read: 282982
          cache_write: 2156
          cost: 0.075
formula: publish
---
# S-0639 Publish what is accepted and not yet released

## Goal

Publish what is accepted and not yet released, by the project's release policy.

## Acceptance criteria

- [ ] Fetch the remote's branch and tags
- [ ] Evaluate the release policy, stopping unless it is met
- [ ] Publish what is accepted and not yet released

## Tasks
- T-2729 Fetch the remote's branch and tags
- T-2730 Evaluate the release policy, stopping unless it is met
- T-2731 Publish what is accepted and not yet released

## Notes

Made by `flai formula run publish` from `design/workflows/formulas/publish.md`, with no inputs.
