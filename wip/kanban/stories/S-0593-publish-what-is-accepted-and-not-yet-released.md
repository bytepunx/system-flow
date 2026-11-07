---
id: S-0593
type: story
nature: improvement
title: Publish what is accepted and not yet released
status: in-progress
owner: alex
created: 2026-11-07T15:17:08Z
updated: 2026-11-07T15:17:37Z
transitions:
  - to: ready
    at: 2026-11-07T15:17:32Z
    by: system-flow
  - to: in-progress
    at: 2026-11-07T15:17:37Z
    by: agent-S-0593
tags: []
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 78
  estimated: true
  priced_by: estimate
  turns:
    - day: 2026-11-07
      work: 15
  models:
    - model: claude-opus-5-5
      input: 32
      output: 208
      cache_read: 1685001
      cache_write: 140190
      cost: 0.6542
  strategic:
    - kind: orchestrator
      seconds: 21
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 4
          output: 27
          cache_read: 79654
          cache_write: 35995
          cost: 0.0292
formula: publish
---
# S-0593 Publish what is accepted and not yet released

## Goal

Publish what is accepted and not yet released, by the project's release policy.

## Acceptance criteria

- [ ] Fetch the remote's branch and tags
- [ ] Evaluate the release policy, stopping unless it is met
- [ ] Publish what is accepted and not yet released

## Tasks
- T-2586 Fetch the remote's branch and tags
- T-2587 Evaluate the release policy, stopping unless it is met
- T-2588 Publish what is accepted and not yet released

## Notes

Made by `flai formula run publish` from `design/workflows/formulas/publish.md`, with no inputs.
