---
id: S-0718
type: story
nature: improvement
title: Publish what is accepted and not yet released
status: in-progress
owner: alex
created: 2026-11-11T01:44:41Z
updated: 2026-11-11T22:13:46Z
transitions:
  - to: ready
    at: 2026-11-11T01:44:44Z
    by: system-flow
  - to: in-progress
    at: 2026-11-11T22:13:46Z
    by: agent-S-0718
tags: []
touches: [wip/kanban/stories/S-0718-publish-what-is-accepted-and-not-yet-released.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 92
  estimated: true
  priced_by: estimate
  turns:
    - day: 2026-11-11
      work: 16
  models:
    - model: claude-opus-5-5
      input: 36
      output: 284
      cache_read: 1982134
      cache_write: 152353
      cost: 0.7923
  strategic:
    - kind: orchestrator
      seconds: 253
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 77
          output: 448
          cache_read: 2878292
          cache_write: 399275
          cost: 0.9142
formula: publish
---
# S-0718 Publish what is accepted and not yet released

## Goal

Publish what is accepted and not yet released, by the project's release policy.

## Acceptance criteria

- [ ] Fetch the remote's branch and tags
- [ ] Evaluate the release policy, stopping unless it is met
- [ ] Publish what is accepted and not yet released

## Tasks
- T-3080 Fetch the remote's branch and tags
- T-3081 Evaluate the release policy, stopping unless it is met
- T-3082 Publish what is accepted and not yet released

## Notes

Made by `flai formula run publish` from `design/workflows/formulas/publish.md`, with the inputs `judgement=`.
