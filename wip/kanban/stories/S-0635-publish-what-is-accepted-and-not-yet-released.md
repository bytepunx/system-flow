---
id: S-0635
type: story
nature: improvement
title: Publish what is accepted and not yet released
status: ready
owner: alex
created: 2026-11-08T01:14:12Z
updated: 2026-11-08T01:15:24Z
transitions:
  - to: ready
    at: 2026-11-08T01:15:24Z
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
      seconds: 650
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 129
          output: 802
          cache_read: 5058822
          cache_write: 544166
          cost: 1.4565
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

- [ ] Fetch the remote's branch and tags
- [ ] Evaluate the release policy, stopping unless it is met
- [ ] Publish what is accepted and not yet released

## Tasks
- T-2694 Fetch the remote's branch and tags
- T-2695 Evaluate the release policy, stopping unless it is met
- T-2696 Publish what is accepted and not yet released

## Notes

Made by `flai formula run publish` from `design/workflows/formulas/publish.md`, with no inputs.
