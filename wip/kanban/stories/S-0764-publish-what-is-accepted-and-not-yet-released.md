---
id: S-0764
type: story
nature: improvement
title: Publish what is accepted and not yet released
status: in-progress
owner: alex
created: 2026-11-13T08:55:19Z
updated: 2026-11-13T09:00:03Z
transitions:
  - to: ready
    at: 2026-11-13T08:55:20Z
    by: system-flow
  - to: in-progress
    at: 2026-11-13T08:55:26Z
    by: agent-S-0764
blocked:
  - from: 2026-11-13T08:58:38Z
    until: 2026-11-13T09:00:03Z
    reason: awaiting the operator on TH-1142
tags: []
touches: [wip/kanban/stories/S-0764-publish-what-is-accepted-and-not-yet-released.md, design/issues/I-0359-flai-release-pending-fails-when-another-process-holds-the-main-checkout-s-git.md, design/issues/summary.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 195
  priced_by: harness
  turns:
    - day: 2026-11-13
      work: 28
  models:
    - model: claude-opus-5-5
      input: 60
      output: 8348
      cache_read: 4156102
      cache_write: 169292
      cost: 2.3528
  strategic:
    - kind: orchestrator
      seconds: 111
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 28
          output: 231
          cache_read: 868595
          cache_write: 66012
          cost: 0.2694
formula: publish
---
# S-0764 Publish what is accepted and not yet released

## Goal

Publish what is accepted and not yet released, by the project's release policy.

## Acceptance criteria

- [x] Fetch the remote's branch and tags
- [x] Evaluate the release policy, stopping unless it is met
- [ ] Publish what is accepted and not yet released

## Tasks
- T-3273 Fetch the remote's branch and tags
- T-3274 Evaluate the release policy, stopping unless it is met
- T-3275 Publish what is accepted and not yet released

## Notes

Made by `flai formula run publish` from `design/workflows/formulas/publish.md`, with the inputs `judgement=`.
