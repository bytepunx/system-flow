---
id: S-0772
type: story
nature: improvement
title: Publish what is accepted and not yet released
status: done
owner: alex
created: 2026-11-13T13:42:42Z
updated: 2026-11-13T15:22:35Z
transitions:
  - to: ready
    at: 2026-11-13T13:42:43Z
    by: system-flow
  - to: in-progress
    at: 2026-11-13T15:18:17Z
    by: agent-S-0772
  - to: review
    at: 2026-11-13T15:21:12Z
    by: agent-S-0772
  - to: done
    at: 2026-11-13T15:22:35Z
    by: orchestrator
tags: []
touches: [wip/kanban/stories/S-0772-publish-what-is-accepted-and-not-yet-released.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 187
  priced_by: harness
  turns:
    - day: 2026-11-13
      ceremony: 3
      work: 29
  models:
    - model: claude-opus-5-5
      input: 70
      output: 8654
      cache_read: 4859053
      cache_write: 171176
      cost: 2.5146
  strategic:
    - kind: orchestrator
      seconds: 19
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 7
          output: 41
          cache_read: 173745
          cache_write: 17026
          cost: 0.0554
formula: publish
---
# S-0772 Publish what is accepted and not yet released

## Goal

Publish what is accepted and not yet released, by the project's release policy.

## Acceptance criteria

- [x] Fetch the remote's branch and tags
- [x] Evaluate the release policy, stopping unless it is met
- [x] Publish what is accepted and not yet released (skipped: step evaluate stopped the run with exit 2)

## Tasks
- T-3303 Fetch the remote's branch and tags
- T-3304 Evaluate the release policy, stopping unless it is met
- T-3305 Publish what is accepted and not yet released

## Notes

Made by `flai formula run publish` from `design/workflows/formulas/publish.md`, with the inputs `judgement=`.

Step evaluate stopped the run with exit 2; skipped T-3305.

- 2026-11-13: nothing was released. S-0771, the duplicate instance TH-1161 names, had already published flai 1.92.0 and template 1.1.37. The evaluation found 0 of 1 accepted stories waiting. Nothing was tagged or pushed.

### Accepted by the orchestrator

- Verified: ce920472bd874c0125d6328740341d21f1b4f889
- At: 2026-11-13T15:22:35Z

Verdict: pass at ce920472bd874c0125d6328740341d21f1b4f889
- 1: T-3303 done
- 2: T-3304 done
- 3: T-3305 skipped
