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
      seconds: 20
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 5
          output: 46
          cache_read: 136836
          cache_write: 14920
          cost: 0.0378
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
