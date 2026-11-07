---
id: S-0594
type: story
nature: improvement
title: Publish what is accepted and not yet released
status: done
owner: alex
created: 2026-11-07T21:47:12Z
updated: 2026-11-07T22:02:20Z
transitions:
  - to: ready
    at: 2026-11-07T21:47:37Z
    by: system-flow
  - to: in-progress
    at: 2026-11-07T21:53:56Z
    by: agent-S-0594
  - to: review
    at: 2026-11-07T22:00:53Z
    by: system-flow
  - to: done
    at: 2026-11-07T22:02:20Z
    by: alex
tags: []
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 430
  priced_by: harness
  turns:
    - day: 2026-11-07
      ceremony: 2
      work: 29
  models:
    - model: claude-opus-5-5
      input: 66
      output: 8776
      cache_read: 4469568
      cache_write: 155909
      cost: 2.317
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

- [x] Fetch the remote's branch and tags
- [x] Evaluate the release policy, stopping unless it is met
- [x] Publish what is accepted and not yet released

## Tasks
- T-2589 Fetch the remote's branch and tags
- T-2590 Evaluate the release policy, stopping unless it is met
- T-2591 Publish what is accepted and not yet released

## Notes

Made by `flai formula run publish` from `design/workflows/formulas/publish.md`, with no inputs.

Release published 2026-11-07 by `flai release --pending --json` in the main checkout, pushed:

| Component | From | To | Tag | Items |
|-----------|------|----|-----|-------|
| flai | 1.85.1 | 1.86.0 | `flai/v1.86.0` | S-0407 (minor) |
| template | 1.1.26 | 1.1.27 | v1.1.27 at `system-flow-template` (78fc638) | S-0407 (patch) |
