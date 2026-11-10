---
id: S-0657
type: story
nature: improvement
title: Publish what is accepted and not yet released
status: in-progress
owner: alex
created: 2026-11-10T07:11:22Z
updated: 2026-11-10T07:55:30Z
transitions:
  - to: ready
    at: 2026-11-10T07:11:28Z
    by: system-flow
  - to: in-progress
    at: 2026-11-10T07:53:51Z
    by: agent-S-0657
tags: []
touches: [wip/kanban/stories/S-0657-publish-what-is-accepted-and-not-yet-released.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 76
  estimated: true
  priced_by: estimate
  turns:
    - day: 2026-11-10
      ceremony: 1
      hand_edits: 1
      work: 11
  models:
    - model: claude-opus-5-5
      input: 28
      output: 189
      cache_read: 1436769
      cache_write: 139976
      cost: 0.5801
  strategic:
    - kind: orchestrator
      seconds: 55
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 14
          output: 75
          cache_read: 538682
          cache_write: 83862
          cost: 0.1677
formula: publish
---
# S-0657 Publish what is accepted and not yet released

## Goal

Publish what is accepted and not yet released, by the project's release policy.

## Acceptance criteria

- [ ] Fetch the remote's branch and tags
- [x] Evaluate the release policy, stopping unless it is met
- [x] Publish what is accepted and not yet released (skipped: step evaluate stopped the run with exit 2)

## Tasks
- T-2853 Fetch the remote's branch and tags
- T-2854 Evaluate the release policy, stopping unless it is met
- T-2855 Publish what is accepted and not yet released

## Notes

Made by `flai formula run publish` from `design/workflows/formulas/publish.md`, with the inputs `judgement=`.

Step evaluate stopped the run with exit 2; skipped T-2855.
