---
id: S-0641
type: story
nature: improvement
title: Publish what is accepted and not yet released
status: in-progress
owner: alex
created: 2026-11-09T16:56:49Z
updated: 2026-11-09T16:56:55Z
transitions:
  - to: ready
    at: 2026-11-09T16:56:50Z
    by: system-flow
  - to: in-progress
    at: 2026-11-09T16:56:55Z
    by: agent-S-0641
tags: []
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: sum
  seconds: 14
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 10
      output: 106
      cache_read: 744681
      cache_write: 2494
      cost: 0.2726
  strategic:
    - kind: orchestrator
      seconds: 10
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 5
          output: 34
          cache_read: 89751
          cache_write: 53176
          cost: 0.0377
formula: publish
---
# S-0641 Publish what is accepted and not yet released

## Goal

Publish what is accepted and not yet released, by the project's release policy.

## Acceptance criteria

- [ ] Fetch the remote's branch and tags
- [ ] Evaluate the release policy, stopping unless it is met
- [ ] Publish what is accepted and not yet released

## Tasks
- T-2732 Fetch the remote's branch and tags
- T-2733 Evaluate the release policy, stopping unless it is met
- T-2734 Publish what is accepted and not yet released

## Notes

Made by `flai formula run publish` from `design/workflows/formulas/publish.md`, with no inputs.
