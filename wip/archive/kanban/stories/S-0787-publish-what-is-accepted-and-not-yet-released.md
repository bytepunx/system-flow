---
id: S-0787
type: story
nature: improvement
title: Publish what is accepted and not yet released
status: done
owner: alex
created: 2026-11-14T18:00:18Z
updated: 2026-11-14T18:04:02Z
transitions:
  - to: ready
    at: 2026-11-14T18:00:20Z
    by: system-flow
  - to: in-progress
    at: 2026-11-14T18:00:31Z
    by: agent-S-0787
  - to: review
    at: 2026-11-14T18:02:49Z
    by: system-flow
  - to: done
    at: 2026-11-14T18:04:02Z
    by: orchestrator
tags: []
touches: [wip/kanban/stories/S-0787-publish-what-is-accepted-and-not-yet-released.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 146
  priced_by: harness
  turns:
    - day: 2026-11-14
      ceremony: 2
      hand_edits: 3
      work: 19
  models:
    - model: claude-opus-5-5
      input: 52
      output: 6732
      cache_read: 3552198
      cache_write: 176354
      cost: 2.2561
  strategic:
    - kind: orchestrator
      seconds: 89
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 23
          output: 142
          cache_read: 955697
          cache_write: 140613
          cost: 0.3268
        - model: claude-sonnet-5-5
          input: 8
          output: 838
          cache_read: 93102
          cache_write: 37209
          cost: 0.1024
formula: publish
---
# S-0787 Publish what is accepted and not yet released

## Goal

Publish what is accepted and not yet released, by the project's release policy.

## Acceptance criteria

- [x] Fetch the remote's branch and tags
- [x] Evaluate the release policy, stopping unless it is met
- [x] Publish what is accepted and not yet released

## Tasks
- T-3361 Fetch the remote's branch and tags
- T-3362 Evaluate the release policy, stopping unless it is met
- T-3363 Publish what is accepted and not yet released

## Notes

Made by `flai formula run publish` from `design/workflows/formulas/publish.md`, with the inputs `judgement=`.

### Accepted by the orchestrator

- Verified: 9d2500190057e24c967fad501367e64314f7d1b3
- At: 2026-11-14T18:04:02Z

Verdict: pass, the publish formula's three steps are done and the close-out passed at 9d2500190057
- 1: T-3361 done
- 2: T-3362 done
- 3: T-3363 done
