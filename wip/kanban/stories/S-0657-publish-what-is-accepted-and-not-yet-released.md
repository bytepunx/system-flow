---
id: S-0657
type: story
nature: improvement
title: Publish what is accepted and not yet released
status: review
owner: alex
created: 2026-11-10T07:11:22Z
updated: 2026-11-10T07:56:34Z
transitions:
  - to: ready
    at: 2026-11-10T07:11:28Z
    by: system-flow
  - to: in-progress
    at: 2026-11-10T07:53:51Z
    by: agent-S-0657
  - to: review
    at: 2026-11-10T07:56:34Z
    by: system-flow
tags: []
touches: [wip/kanban/stories/S-0657-publish-what-is-accepted-and-not-yet-released.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 173
  priced_by: harness
  turns:
    - day: 2026-11-10
      ceremony: 2
      hand_edits: 1
      work: 22
  models:
    - model: claude-opus-5-5
      input: 54
      output: 7008
      cache_read: 3518591
      cache_write: 150844
      cost: 2.0508
  strategic:
    - kind: orchestrator
      seconds: 66
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 16
          output: 84
          cache_read: 602248
          cache_write: 97976
          cost: 0.1886
formula: publish
---
# S-0657 Publish what is accepted and not yet released

## Goal

Publish what is accepted and not yet released, by the project's release policy.

## Acceptance criteria

- [x] Fetch the remote's branch and tags
- [x] Evaluate the release policy, stopping unless it is met
- [x] Publish what is accepted and not yet released (skipped: step evaluate stopped the run with exit 2)

## Tasks
- T-2853 Fetch the remote's branch and tags
- T-2854 Evaluate the release policy, stopping unless it is met
- T-2855 Publish what is accepted and not yet released

## Notes

Made by `flai formula run publish` from `design/workflows/formulas/publish.md`, with the inputs `judgement=`.

Step evaluate stopped the run with exit 2; skipped T-2855.
