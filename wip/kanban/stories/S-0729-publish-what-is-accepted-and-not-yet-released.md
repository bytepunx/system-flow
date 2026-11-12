---
id: S-0729
type: story
nature: improvement
title: Publish what is accepted and not yet released
status: review
owner: alex
created: 2026-11-12T03:39:40Z
updated: 2026-11-12T03:43:08Z
transitions:
  - to: ready
    at: 2026-11-12T03:39:41Z
    by: system-flow
  - to: in-progress
    at: 2026-11-12T03:40:05Z
    by: agent-S-0729
  - to: review
    at: 2026-11-12T03:43:03Z
    by: agent-S-0729
tags: []
touches: [wip/kanban/stories/S-0729-publish-what-is-accepted-and-not-yet-released.md, design/issues/I-0347-flai-check-finds-item-blocked-closed-outside-the-story-at-close-out.md, design/issues/summary.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 197
  priced_by: harness
  turns:
    - day: 2026-11-12
      ceremony: 2
      work: 28
  models:
    - model: claude-opus-5-5
      input: 66
      output: 7990
      cache_read: 4418357
      cache_write: 164079
      cost: 2.3564
  strategic:
    - kind: orchestrator
      seconds: 33
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 12
          output: 79
          cache_read: 341899
          cache_write: 66269
          cost: 0.1158
formula: publish
---
# S-0729 Publish what is accepted and not yet released

## Goal

Publish what is accepted and not yet released, by the project's release policy.

## Acceptance criteria

- [x] Fetch the remote's branch and tags
- [x] Evaluate the release policy, stopping unless it is met
- [x] Publish what is accepted and not yet released

## Tasks
- T-3132 Fetch the remote's branch and tags
- T-3133 Evaluate the release policy, stopping unless it is met
- T-3134 Publish what is accepted and not yet released

## Notes

Made by `flai formula run publish` from `design/workflows/formulas/publish.md`, with the inputs `judgement=`.
