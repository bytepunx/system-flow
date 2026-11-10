---
id: S-0658
type: story
nature: improvement
title: Publish what is accepted and not yet released
status: done
owner: alex
created: 2026-11-10T07:56:53Z
updated: 2026-11-10T08:37:02Z
transitions:
  - to: ready
    at: 2026-11-10T07:56:58Z
    by: system-flow
  - to: in-progress
    at: 2026-11-10T08:26:48Z
    by: agent-S-0658
  - to: review
    at: 2026-11-10T08:29:45Z
    by: system-flow
  - to: done
    at: 2026-11-10T08:37:02Z
    by: alex
tags: []
touches: [wip/kanban/stories/S-0658-publish-what-is-accepted-and-not-yet-released.md, design/issues/I-0332-a-formula-step-s-bare-flai-command-resolves-to-the-installed-flai-on-path-which.md, design/issues/summary.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 185
  priced_by: harness
  turns:
    - day: 2026-11-10
      ceremony: 1
      hand_edits: 1
      work: 34
  models:
    - model: claude-opus-5-5
      input: 80
      output: 9201
      cache_read: 5416458
      cache_write: 161694
      cost: 2.5612
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

- [x] Fetch the remote's branch and tags
- [x] Evaluate the release policy, stopping unless it is met
- [x] Publish what is accepted and not yet released

## Tasks
- T-2856 Fetch the remote's branch and tags
- T-2857 Evaluate the release policy, stopping unless it is met
- T-2858 Publish what is accepted and not yet released

## Notes

Made by `flai formula run publish` from `design/workflows/formulas/publish.md`, with the inputs `judgement=`.
