---
id: S-0727
type: story
nature: improvement
title: Publish what is accepted and not yet released
status: done
owner: alex
created: 2026-11-11T23:22:51Z
updated: 2026-11-11T23:27:00Z
transitions:
  - to: ready
    at: 2026-11-11T23:22:52Z
    by: system-flow
  - to: in-progress
    at: 2026-11-11T23:23:01Z
    by: agent-S-0727
  - to: review
    at: 2026-11-11T23:25:33Z
    by: agent-S-0727
  - to: done
    at: 2026-11-11T23:27:00Z
    by: orchestrator
tags: []
touches: [wip/kanban/stories/S-0727-publish-what-is-accepted-and-not-yet-released.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 167
  priced_by: harness
  turns:
    - day: 2026-11-11
      ceremony: 3
      work: 25
  models:
    - model: claude-opus-5-5
      input: 62
      output: 7195
      cache_read: 3999108
      cache_write: 155387
      cost: 2.1871
  strategic:
    - kind: orchestrator
      seconds: 37
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 17
          output: 133
          cache_read: 506502
          cache_write: 89581
          cost: 0.1687
formula: publish
---
# S-0727 Publish what is accepted and not yet released

## Goal

Publish what is accepted and not yet released, by the project's release policy.

## Acceptance criteria

- [x] Fetch the remote's branch and tags
- [x] Evaluate the release policy, stopping unless it is met
- [x] Publish what is accepted and not yet released

## Tasks
- T-3124 Fetch the remote's branch and tags
- T-3125 Evaluate the release policy, stopping unless it is met
- T-3126 Publish what is accepted and not yet released

## Notes

Made by `flai formula run publish` from `design/workflows/formulas/publish.md`, with the inputs `judgement=`.

### Accepted by the orchestrator

- Verified: b4984d8e7ff22076cee5b75b6d608be916ef3fae
- At: 2026-11-11T23:27:00Z

Verdict: all three publish steps done; close-out passed at b4984d8e7ff22076cee5b75b6d608be916ef3fae
- 1: T-3124 done
- 2: T-3125 done
- 3: T-3126 done
