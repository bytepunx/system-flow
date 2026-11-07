---
id: S-0594
type: story
nature: improvement
title: Publish what is accepted and not yet released
status: in-progress
owner: alex
created: 2026-11-07T21:47:12Z
updated: 2026-11-07T21:53:56Z
transitions:
  - to: ready
    at: 2026-11-07T21:47:37Z
    by: system-flow
  - to: in-progress
    at: 2026-11-07T21:53:56Z
    by: agent-S-0594
tags: []
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 137
  estimated: true
  priced_by: estimate
  turns:
    - day: 2026-11-07
      ceremony: 1
      work: 13
  models:
    - model: claude-opus-5-5
      input: 30
      output: 207
      cache_read: 1518446
      cache_write: 140500
      cost: 0.5949
  strategic:
    - kind: orchestrator
      seconds: 15
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 5
          output: 26
          cache_read: 117285
          cache_write: 24867
          cost: 0.0359
formula: publish
---
# S-0594 Publish what is accepted and not yet released

## Goal

Publish what is accepted and not yet released, by the project's release policy.

## Acceptance criteria

- [ ] Fetch the remote's branch and tags
- [ ] Evaluate the release policy, stopping unless it is met
- [ ] Publish what is accepted and not yet released

## Tasks
- T-2589 Fetch the remote's branch and tags
- T-2590 Evaluate the release policy, stopping unless it is met
- T-2591 Publish what is accepted and not yet released

## Notes

Made by `flai formula run publish` from `design/workflows/formulas/publish.md`, with no inputs.
