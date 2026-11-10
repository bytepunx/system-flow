---
id: S-0664
type: story
nature: improvement
title: Publish what is accepted and not yet released
status: in-progress
owner: alex
created: 2026-11-10T08:33:20Z
updated: 2026-11-10T08:33:28Z
transitions:
  - to: ready
    at: 2026-11-10T08:33:22Z
    by: system-flow
  - to: in-progress
    at: 2026-11-10T08:33:28Z
    by: agent-S-0664
tags: []
touches: [wip/kanban/stories/S-0664-publish-what-is-accepted-and-not-yet-released.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 90
  estimated: true
  priced_by: estimate
  turns:
    - day: 2026-11-10
      ceremony: 1
      hand_edits: 1
      work: 14
  models:
    - model: claude-opus-5-5
      input: 36
      output: 243
      cache_read: 2048807
      cache_write: 149496
      cost: 0.8087
  strategic:
    - kind: orchestrator
      seconds: 17
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 7
          output: 31
          cache_read: 219462
          cache_write: 33426
          cost: 0.0684
formula: publish
---
# S-0664 Publish what is accepted and not yet released

## Goal

Publish what is accepted and not yet released, by the project's release policy.

## Acceptance criteria

- [ ] Fetch the remote's branch and tags
- [ ] Evaluate the release policy, stopping unless it is met
- [ ] Publish what is accepted and not yet released

## Tasks
- T-2890 Fetch the remote's branch and tags
- T-2891 Evaluate the release policy, stopping unless it is met
- T-2892 Publish what is accepted and not yet released

## Notes

Made by `flai formula run publish` from `design/workflows/formulas/publish.md`, with the inputs `judgement=`.
