---
id: S-0783
type: story
nature: improvement
title: Publish what is accepted and not yet released
status: done
owner: alex
created: 2026-11-13T22:47:23Z
updated: 2026-11-13T22:52:06Z
transitions:
  - to: ready
    at: 2026-11-13T22:47:24Z
    by: system-flow
  - to: in-progress
    at: 2026-11-13T22:47:32Z
    by: agent-S-0783
  - to: review
    at: 2026-11-13T22:49:58Z
    by: system-flow
  - to: done
    at: 2026-11-13T22:52:06Z
    by: orchestrator
tags: []
touches: [wip/kanban/stories/S-0783-publish-what-is-accepted-and-not-yet-released.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 171
  priced_by: harness
  turns:
    - day: 2026-11-13
      ceremony: 1
      hand_edits: 1
      work: 21
  models:
    - model: claude-opus-5-5
      input: 52
      output: 6603
      cache_read: 3198307
      cache_write: 159547
      cost: 2.0483
  strategic:
    - kind: orchestrator
      seconds: 135
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 31
          output: 232
          cache_read: 1289211
          cache_write: 112637
          cost: 0.4134
        - model: claude-sonnet-5-5
          input: 4
          output: 27
          cache_read: 13755
          cache_write: 35493
          cost: 0.0383
formula: publish
---
# S-0783 Publish what is accepted and not yet released

## Goal

Publish what is accepted and not yet released, by the project's release policy.

## Acceptance criteria

- [x] Fetch the remote's branch and tags
- [x] Evaluate the release policy, stopping unless it is met
- [x] Publish what is accepted and not yet released (skipped: step evaluate stopped the run with exit 2)

## Tasks
- T-3340 Fetch the remote's branch and tags
- T-3341 Evaluate the release policy, stopping unless it is met
- T-3342 Publish what is accepted and not yet released

## Notes

Made by `flai formula run publish` from `design/workflows/formulas/publish.md`, with the inputs `judgement=`.

Step evaluate stopped the run with exit 2; skipped T-3342.

### Accepted by the orchestrator

- Verified: c4e3d198e45f0f3b89144402971b5e5b0ec8ecfb
- At: 2026-11-13T22:52:06Z

Verdict: S-0783 meets all three criteria; the publish formula's evaluate step stopped the run cleanly with exit 2, a listed stop, so the publish step was skipped.
- 1: T-3340 done
- 2: T-3341 done
- 3: T-3342 skipped
