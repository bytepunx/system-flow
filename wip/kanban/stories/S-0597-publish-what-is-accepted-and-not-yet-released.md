---
id: S-0597
type: story
nature: improvement
title: Publish what is accepted and not yet released
status: in-progress
owner: alex
created: 2026-11-07T22:31:23Z
updated: 2026-11-07T23:04:49Z
transitions:
  - to: ready
    at: 2026-11-07T22:31:51Z
    by: system-flow
  - to: in-progress
    at: 2026-11-07T23:04:49Z
    by: agent-S-0597
tags: []
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 103
  estimated: true
  priced_by: estimate
  turns:
    - day: 2026-11-07
      work: 15
  models:
    - model: claude-opus-5-5
      input: 32
      output: 170
      cache_read: 1746374
      cache_write: 140915
      cost: 0.6776
  strategic:
    - kind: orchestrator
      seconds: 42
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 8
          output: 56
          cache_read: 295809
          cache_write: 18532
          cost: 0.0797
formula: publish
---
# S-0597 Publish what is accepted and not yet released

## Goal

Publish what is accepted and not yet released, by the project's release policy.

## Acceptance criteria

- [ ] Fetch the remote's branch and tags
- [ ] Evaluate the release policy, stopping unless it is met
- [ ] Publish what is accepted and not yet released

## Tasks
- T-2612 Fetch the remote's branch and tags
- T-2613 Evaluate the release policy, stopping unless it is met
- T-2614 Publish what is accepted and not yet released

## Notes

Made by `flai formula run publish` from `design/workflows/formulas/publish.md`, with no inputs.
