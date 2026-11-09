---
id: S-0635
type: story
nature: improvement
title: Publish what is accepted and not yet released
status: in-progress
owner: alex
created: 2026-11-08T01:14:12Z
updated: 2026-11-09T11:23:00Z
transitions:
  - to: ready
    at: 2026-11-08T01:15:24Z
    by: system-flow
  - to: in-progress
    at: 2026-11-09T11:23:00Z
    by: agent-S-0635
tags: []
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 55
  estimated: true
  priced_by: estimate
  turns:
    - day: 2026-11-09
      work: 12
  models:
    - model: claude-opus-5-5
      input: 26
      output: 238
      cache_read: 1241925
      cache_write: 138492
      cost: 0.5034
  strategic:
    - kind: orchestrator
      seconds: 666
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 137
          output: 834
          cache_read: 5309203
          cache_write: 578667
          cost: 1.5314
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
