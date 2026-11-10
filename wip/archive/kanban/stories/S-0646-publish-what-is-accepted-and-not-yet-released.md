---
id: S-0646
type: story
nature: improvement
title: Publish what is accepted and not yet released
status: done
owner: alex
created: 2026-11-09T22:22:55Z
updated: 2026-11-10T06:44:52Z
transitions:
  - to: ready
    at: 2026-11-09T22:23:05Z
    by: system-flow
  - to: in-progress
    at: 2026-11-10T06:17:16Z
    by: agent-S-0646
  - to: review
    at: 2026-11-10T06:18:53Z
    by: agent-S-0646
  - to: done
    at: 2026-11-10T06:44:52Z
    by: alex
tags: []
touches: [wip/kanban/stories/S-0646-publish-what-is-accepted-and-not-yet-released.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 105
  priced_by: harness
  turns:
    - day: 2026-11-10
      ceremony: 1
      work: 22
  models:
    - model: claude-opus-5-5
      input: 52
      output: 5736
      cache_read: 3140209
      cache_write: 148971
      cost: 1.9347
  strategic:
    - kind: orchestrator
      seconds: 85
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 25
          output: 140
          cache_read: 676316
          cache_write: 101735
          cost: 0.2082
formula: publish
---
# S-0646 Publish what is accepted and not yet released

## Goal

Publish what is accepted and not yet released, by the project's release policy.

## Acceptance criteria

- [x] Fetch the remote's branch and tags
- [x] Evaluate the release policy, stopping unless it is met
- [x] Publish what is accepted and not yet released (skipped: step evaluate stopped the run with exit 2)

## Tasks
- T-2822 Fetch the remote's branch and tags
- T-2823 Evaluate the release policy, stopping unless it is met
- T-2824 Publish what is accepted and not yet released

## Notes

Made by `flai formula run publish` from `design/workflows/formulas/publish.md`, with no inputs.

Step evaluate stopped the run with exit 2; skipped T-2824.
