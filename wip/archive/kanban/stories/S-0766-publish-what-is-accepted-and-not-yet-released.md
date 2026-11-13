---
id: S-0766
type: story
nature: improvement
title: Publish what is accepted and not yet released
status: done
owner: alex
created: 2026-11-13T11:40:39Z
updated: 2026-11-13T11:45:13Z
transitions:
  - to: ready
    at: 2026-11-13T11:40:40Z
    by: system-flow
  - to: in-progress
    at: 2026-11-13T11:40:48Z
    by: agent-S-0766
  - to: review
    at: 2026-11-13T11:43:19Z
    by: agent-S-0766
  - to: done
    at: 2026-11-13T11:45:13Z
    by: alex
tags: []
touches: [wip/kanban/stories/S-0766-publish-what-is-accepted-and-not-yet-released.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 159
  priced_by: harness
  turns:
    - day: 2026-11-13
      ceremony: 1
      hand_edits: 1
      work: 25
  models:
    - model: claude-opus-5-5
      input: 58
      output: 8890
      cache_read: 4119617
      cache_write: 185786
      cost: 2.4882
  strategic:
    - kind: orchestrator
      seconds: 35
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 18
          output: 137
          cache_read: 432162
          cache_write: 77698
          cost: 0.1474
formula: publish
---
# S-0766 Publish what is accepted and not yet released

## Goal

Publish what is accepted and not yet released, by the project's release policy.

## Acceptance criteria

- [x] Fetch the remote's branch and tags
- [x] Evaluate the release policy, stopping unless it is met
- [x] Publish what is accepted and not yet released

## Tasks
- T-3284 Fetch the remote's branch and tags
- T-3285 Evaluate the release policy, stopping unless it is met
- T-3286 Publish what is accepted and not yet released

## Notes

Made by `flai formula run publish` from `design/workflows/formulas/publish.md`, with the inputs `judgement=`.

- 2026-11-13: published flai 1.88.0 → 1.89.0 (minor), tag `flai/v1.89.0`, for S-0734. No other component was released and the template was not published. Branch and tag pushed to `origin`.
- S-0761 was accepted at 11:42:31Z, after step 3 ran: it is not in this release and waits for the next publish.
