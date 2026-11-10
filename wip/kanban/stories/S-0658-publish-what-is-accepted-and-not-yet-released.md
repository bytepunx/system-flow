---
id: S-0658
type: story
nature: improvement
title: Publish what is accepted and not yet released
status: in-progress
owner: alex
created: 2026-11-10T07:56:53Z
updated: 2026-11-10T08:28:53Z
transitions:
  - to: ready
    at: 2026-11-10T07:56:58Z
    by: system-flow
  - to: in-progress
    at: 2026-11-10T08:26:48Z
    by: agent-S-0658
tags: []
touches: [wip/kanban/stories/S-0658-publish-what-is-accepted-and-not-yet-released.md, design/issues/I-0332-a-formula-step-s-bare-flai-command-resolves-to-the-installed-flai-on-path-which.md, design/issues/summary.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 129
  estimated: true
  priced_by: estimate
  turns:
    - day: 2026-11-10
      hand_edits: 1
      work: 24
  models:
    - model: claude-opus-5-5
      input: 56
      output: 531
      cache_read: 3345193
      cache_write: 151504
      cost: 1.2865
  strategic:
    - kind: orchestrator
      seconds: 31
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 9
          output: 53
          cache_read: 270825
          cache_write: 41200
          cost: 0.0841
formula: publish
---
# S-0658 Publish what is accepted and not yet released

## Goal

Publish what is accepted and not yet released, by the project's release policy.

## Acceptance criteria

- [ ] Fetch the remote's branch and tags
- [ ] Evaluate the release policy, stopping unless it is met
- [ ] Publish what is accepted and not yet released

## Tasks
- T-2856 Fetch the remote's branch and tags
- T-2857 Evaluate the release policy, stopping unless it is met
- T-2858 Publish what is accepted and not yet released

## Notes

Made by `flai formula run publish` from `design/workflows/formulas/publish.md`, with the inputs `judgement=`.
