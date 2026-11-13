---
id: S-0770
type: story
nature: improvement
title: Publish what is accepted and not yet released
status: done
owner: alex
created: 2026-11-13T13:16:45Z
updated: 2026-11-13T13:22:23Z
transitions:
  - to: ready
    at: 2026-11-13T13:16:46Z
    by: system-flow
  - to: in-progress
    at: 2026-11-13T13:16:53Z
    by: agent-S-0770
  - to: review
    at: 2026-11-13T13:19:20Z
    by: agent-S-0770
  - to: done
    at: 2026-11-13T13:22:23Z
    by: orchestrator
tags: []
touches: [wip/kanban/stories/S-0770-publish-what-is-accepted-and-not-yet-released.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 158
  priced_by: harness
  turns:
    - day: 2026-11-13
      ceremony: 1
      work: 24
  models:
    - model: claude-opus-5-5
      input: 54
      output: 8513
      cache_read: 3671532
      cache_write: 171406
      cost: 2.276
  strategic:
    - kind: orchestrator
      seconds: 187
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 31
          output: 219
          cache_read: 1120675
          cache_write: 107390
          cost: 0.3565
        - model: claude-sonnet-5-5
          input: 2
          output: 24
          cache_read: 0
          cache_write: 12907
          cost: 0.01
formula: publish
---
# S-0770 Publish what is accepted and not yet released

## Goal

Publish what is accepted and not yet released, by the project's release policy.

## Acceptance criteria

- [x] Fetch the remote's branch and tags
- [x] Evaluate the release policy, stopping unless it is met
- [x] Publish what is accepted and not yet released

## Tasks
- T-3297 Fetch the remote's branch and tags
- T-3298 Evaluate the release policy, stopping unless it is met
- T-3299 Publish what is accepted and not yet released

## Notes

Made by `flai formula run publish` from `design/workflows/formulas/publish.md`, with the inputs `judgement=`.

Every step ran in the main checkout, each `flai` step as the manifest's `flai.command`, `scripts/flai.sh`:

| Step | Task | Exit | Result |
|------|------|------|--------|
| fetch | T-3297 | 0 | `git fetch --tags` |
| evaluate | T-3298 | 0 | threshold met: S-0740 accepted and not yet released, count 1 of 1 |
| publish | T-3299 | 0 | flai 1.90.0 → 1.91.0, tag `flai/v1.91.0`, branch and tag pushed |

Released: `flai/v1.91.0` (minor, S-0740). T-3297's exit is recorded here, not in its own Notes: it was closed before its Notes were written, and a closed task is not edited.

### Accepted by the orchestrator

- Verified: 8fb630e21efe049df57bea90534d94f6f3539466
- At: 2026-11-13T13:22:23Z

Verdict: pass
- 1: T-3297 done
- 2: T-3298 done
- 3: T-3299 done
