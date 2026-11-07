---
id: S-0588
type: story
nature: improvement
title: Publish what is accepted and not yet released
status: ready
owner: alex
created: 2026-11-07T09:24:05Z
updated: 2026-11-07T09:24:27Z
transitions:
  - to: ready
    at: 2026-11-07T09:24:27Z
    by: system-flow
tags: []
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: sum
  seconds: 0
  models: []
  strategic:
    - kind: orchestrator
      seconds: 83
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 20
          output: 135
          cache_read: 728549
          cache_write: 84191
          cost: 0.2035
        - model: claude-sonnet-5-5
          input: 3
          output: 20
          cache_read: 24558
          cache_write: 7399
          cost: 0.0231
formula: publish
---
# S-0588 Publish what is accepted and not yet released

## Goal

Publish what is accepted and not yet released, by the project's release policy.

## Acceptance criteria

- [ ] Fetch the remote's branch and tags
- [ ] Evaluate the release policy, stopping unless it is met
- [ ] Publish what is accepted and not yet released

## Tasks
- T-2566 Fetch the remote's branch and tags
- T-2567 Evaluate the release policy, stopping unless it is met
- T-2568 Publish what is accepted and not yet released

## Notes

Made by `flai formula run publish` from `design/workflows/formulas/publish.md`, with no inputs.
