---
id: S-0793
type: story
nature: improvement
title: Publish what is accepted and not yet released
status: in-progress
owner: alex
created: 2026-11-15T21:17:19Z
updated: 2026-11-15T21:32:46Z
transitions:
  - to: ready
    at: 2026-11-15T21:17:21Z
    by: system-flow
  - to: in-progress
    at: 2026-11-15T21:30:08Z
    by: agent-S-0793
blocked:
  - from: 2026-11-15T21:32:20Z
    until: 2026-11-15T21:32:46Z
    reason: awaiting the operator on TH-1230
tags: []
touches: [wip/kanban/stories/S-0793-publish-what-is-accepted-and-not-yet-released.md, design/issues/I-0359-flai-release-pending-fails-when-another-process-holds-the-main-checkout-s-git.md, design/issues/summary.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 134
  priced_by: harness
  turns:
    - day: 2026-11-15
      ceremony: 1
      work: 26
  models:
    - model: claude-opus-5-5
      input: 58
      output: 9084
      cache_read: 3904607
      cache_write: 162779
      cost: 2.2651
  strategic:
    - kind: orchestrator
      seconds: 49
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 20
          output: 138
          cache_read: 637950
          cache_write: 89797
          cost: 0.2202
formula: publish
---
# S-0793 Publish what is accepted and not yet released

## Goal

Publish what is accepted and not yet released, by the project's release policy.

## Acceptance criteria

- [x] Fetch the remote's branch and tags
- [x] Evaluate the release policy, stopping unless it is met
- [ ] Publish what is accepted and not yet released

## Tasks
- T-3384 Fetch the remote's branch and tags
- T-3385 Evaluate the release policy, stopping unless it is met
- T-3386 Publish what is accepted and not yet released

## Notes

Made by `flai formula run publish` from `design/workflows/formulas/publish.md`, with the inputs `judgement=`.
