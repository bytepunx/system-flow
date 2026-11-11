---
id: S-0718
type: story
nature: improvement
title: Publish what is accepted and not yet released
status: ready
owner: alex
created: 2026-11-11T01:44:41Z
updated: 2026-11-11T01:44:44Z
transitions:
  - to: ready
    at: 2026-11-11T01:44:44Z
    by: system-flow
tags: []
touches: [wip/kanban/stories/S-0718-publish-what-is-accepted-and-not-yet-released.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: sum
  seconds: 0
  models: []
  strategic:
    - kind: orchestrator
      seconds: 26
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 7
          output: 45
          cache_read: 258312
          cache_write: 46901
          cost: 0.0842
formula: publish
---
# S-0718 Publish what is accepted and not yet released

## Goal

Publish what is accepted and not yet released, by the project's release policy.

## Acceptance criteria

- [ ] Fetch the remote's branch and tags
- [ ] Evaluate the release policy, stopping unless it is met
- [ ] Publish what is accepted and not yet released

## Tasks
- T-3080 Fetch the remote's branch and tags
- T-3081 Evaluate the release policy, stopping unless it is met
- T-3082 Publish what is accepted and not yet released

## Notes

Made by `flai formula run publish` from `design/workflows/formulas/publish.md`, with the inputs `judgement=`.
