---
id: S-0789
type: story
nature: improvement
title: Publish what is accepted and not yet released
status: done
owner: alex
created: 2026-11-14T23:49:29Z
updated: 2026-11-15T20:01:53Z
transitions:
  - to: ready
    at: 2026-11-14T23:49:30Z
    by: system-flow
  - to: in-progress
    at: 2026-11-15T01:56:58Z
    by: agent-S-0789
  - to: review
    at: 2026-11-15T01:59:04Z
    by: agent-S-0789
  - to: done
    at: 2026-11-15T20:01:53Z
    by: alex
tags: []
touches: [wip/kanban/stories/S-0789-publish-what-is-accepted-and-not-yet-released.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 147
  priced_by: harness
  turns:
    - day: 2026-11-15
      ceremony: 2
      work: 27
  models:
    - model: claude-opus-5-5
      input: 64
      output: 7394
      cache_read: 4118301
      cache_write: 163773
      cost: 2.282
  strategic:
    - kind: orchestrator
      seconds: 98
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 43
          output: 292
          cache_read: 1788393
          cache_write: 146954
          cost: 0.5792
formula: publish
---
# S-0789 Publish what is accepted and not yet released

## Goal

Publish what is accepted and not yet released, by the project's release policy.

## Acceptance criteria

- [x] Fetch the remote's branch and tags
- [x] Evaluate the release policy, stopping unless it is met
- [x] Publish what is accepted and not yet released (skipped: step evaluate stopped the run with exit 2)

## Tasks
- T-3366 Fetch the remote's branch and tags
- T-3367 Evaluate the release policy, stopping unless it is met
- T-3368 Publish what is accepted and not yet released

## Notes

Made by `flai formula run publish` from `design/workflows/formulas/publish.md`, with the inputs `judgement=`.

Step evaluate stopped the run with exit 2; skipped T-3368.
