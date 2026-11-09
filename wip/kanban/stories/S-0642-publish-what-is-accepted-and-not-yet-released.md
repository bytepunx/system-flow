---
id: S-0642
type: story
nature: improvement
title: Publish what is accepted and not yet released
status: in-progress
owner: alex
created: 2026-11-09T17:28:50Z
updated: 2026-11-09T21:53:03Z
transitions:
  - to: ready
    at: 2026-11-09T17:28:54Z
    by: system-flow
  - to: in-progress
    at: 2026-11-09T21:53:03Z
    by: agent-S-0642
tags: []
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 58
  estimated: true
  priced_by: estimate
  turns:
    - day: 2026-11-09
      ceremony: 1
      work: 12
  models:
    - model: claude-opus-5-5
      input: 28
      output: 185
      cache_read: 1505034
      cache_write: 142856
      cost: 0.6022
  strategic:
    - kind: orchestrator
      seconds: 102
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 32
          output: 198
          cache_read: 1020838
          cache_write: 181123
          cost: 0.3186
formula: publish
---
# S-0642 Publish what is accepted and not yet released

## Goal

Publish what is accepted and not yet released, by the project's release policy.

## Acceptance criteria

- [ ] Fetch the remote's branch and tags
- [ ] Evaluate the release policy, stopping unless it is met
- [ ] Publish what is accepted and not yet released

## Tasks
- T-2760 Fetch the remote's branch and tags
- T-2761 Evaluate the release policy, stopping unless it is met
- T-2762 Publish what is accepted and not yet released

## Notes

Made by `flai formula run publish` from `design/workflows/formulas/publish.md`, with no inputs.
