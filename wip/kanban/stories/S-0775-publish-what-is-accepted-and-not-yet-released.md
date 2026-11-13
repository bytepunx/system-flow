---
id: S-0775
type: story
nature: improvement
title: Publish what is accepted and not yet released
status: in-progress
owner: alex
created: 2026-11-13T16:13:36Z
updated: 2026-11-13T16:13:47Z
transitions:
  - to: ready
    at: 2026-11-13T16:13:37Z
    by: system-flow
  - to: in-progress
    at: 2026-11-13T16:13:47Z
    by: agent-S-0775
tags: []
touches: [wip/kanban/stories/S-0775-publish-what-is-accepted-and-not-yet-released.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 97
  estimated: true
  priced_by: estimate
  turns:
    - day: 2026-11-13
      work: 17
  models:
    - model: claude-opus-5-5
      input: 36
      output: 269
      cache_read: 2036804
      cache_write: 151555
      cost: 0.8197
  strategic:
    - kind: orchestrator
      seconds: 20
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 6
          output: 33
          cache_read: 193500
          cache_write: 29748
          cost: 0.0651
formula: publish
---
# S-0775 Publish what is accepted and not yet released

## Goal

Publish what is accepted and not yet released, by the project's release policy.

## Acceptance criteria

- [ ] Fetch the remote's branch and tags
- [ ] Evaluate the release policy, stopping unless it is met
- [ ] Publish what is accepted and not yet released

## Tasks
- T-3313 Fetch the remote's branch and tags
- T-3314 Evaluate the release policy, stopping unless it is met
- T-3315 Publish what is accepted and not yet released

## Notes

Made by `flai formula run publish` from `design/workflows/formulas/publish.md`, with the inputs `judgement=`.
