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
  seconds: 58
  estimated: true
  priced_by: estimate
  turns:
    - day: 2026-11-10
      ceremony: 1
      hand_edits: 1
      work: 11
  models:
    - model: claude-opus-5-5
      input: 30
      output: 172
      cache_read: 1561205
      cache_write: 146322
      cost: 0.6281
  strategic:
    - kind: orchestrator
      seconds: 10
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 2
          output: 4
          cache_read: 76674
          cache_write: 1019
          cost: 0.021
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
