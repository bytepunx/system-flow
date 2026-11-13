---
id: S-0762
type: story
nature: improvement
title: Publish what is accepted and not yet released
status: in-progress
owner: alex
created: 2026-11-13T07:24:38Z
updated: 2026-11-13T07:24:46Z
transitions:
  - to: ready
    at: 2026-11-13T07:24:39Z
    by: system-flow
  - to: in-progress
    at: 2026-11-13T07:24:46Z
    by: agent-S-0762
tags: []
touches: [wip/kanban/stories/S-0762-publish-what-is-accepted-and-not-yet-released.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 99
  estimated: true
  priced_by: estimate
  turns:
    - day: 2026-11-13
      hand_edits: 2
      work: 14
  models:
    - model: claude-opus-5-5
      input: 34
      output: 225
      cache_read: 1942662
      cache_write: 153737
      cost: 0.7799
  strategic:
    - kind: orchestrator
      seconds: 27
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 10
          output: 58
          cache_read: 305652
          cache_write: 15623
          cost: 0.0921
formula: publish
---
# S-0762 Publish what is accepted and not yet released

## Goal

Publish what is accepted and not yet released, by the project's release policy.

## Acceptance criteria

- [ ] Fetch the remote's branch and tags
- [ ] Evaluate the release policy, stopping unless it is met
- [ ] Publish what is accepted and not yet released

## Tasks
- T-3265 Fetch the remote's branch and tags
- T-3266 Evaluate the release policy, stopping unless it is met
- T-3267 Publish what is accepted and not yet released

## Notes

Made by `flai formula run publish` from `design/workflows/formulas/publish.md`, with the inputs `judgement=`.
