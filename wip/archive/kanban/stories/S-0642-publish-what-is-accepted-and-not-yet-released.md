---
id: S-0642
type: story
nature: improvement
title: Publish what is accepted and not yet released
status: done
owner: alex
created: 2026-11-09T17:28:50Z
updated: 2026-11-09T21:55:02Z
transitions:
  - to: ready
    at: 2026-11-09T17:28:54Z
    by: system-flow
  - to: in-progress
    at: 2026-11-09T21:53:03Z
    by: agent-S-0642
  - to: review
    at: 2026-11-09T21:54:50Z
    by: system-flow
  - to: done
    at: 2026-11-09T21:55:02Z
    by: alex
tags: []
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 89
  estimated: true
  priced_by: estimate
  turns:
    - day: 2026-11-09
      ceremony: 1
      work: 17
  models:
    - model: claude-opus-5-5
      input: 38
      output: 322
      cache_read: 2301139
      cache_write: 148731
      cost: 0.8954
  strategic:
    - kind: orchestrator
      seconds: 116
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 40
          output: 262
          cache_read: 1194915
          cache_write: 250282
          cost: 0.3832
formula: publish
---
# S-0642 Publish what is accepted and not yet released

## Goal

Publish what is accepted and not yet released, by the project's release policy.

## Acceptance criteria

- [x] Fetch the remote's branch and tags
- [x] Evaluate the release policy, stopping unless it is met
- [x] Publish what is accepted and not yet released

## Tasks
- T-2760 Fetch the remote's branch and tags
- T-2761 Evaluate the release policy, stopping unless it is met
- T-2762 Publish what is accepted and not yet released

## Notes

Made by `flai formula run publish` from `design/workflows/formulas/publish.md`, with no inputs.
