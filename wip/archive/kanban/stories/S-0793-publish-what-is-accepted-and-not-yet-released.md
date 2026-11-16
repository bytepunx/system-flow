---
id: S-0793
type: story
nature: improvement
title: Publish what is accepted and not yet released
status: done
owner: alex
created: 2026-11-15T21:17:19Z
updated: 2026-11-16T00:36:08Z
transitions:
  - to: ready
    at: 2026-11-15T21:17:21Z
    by: system-flow
  - to: in-progress
    at: 2026-11-15T21:30:08Z
    by: agent-S-0793
  - to: review
    at: 2026-11-15T21:34:05Z
    by: agent-S-0793
  - to: done
    at: 2026-11-16T00:36:08Z
    by: orchestrator
blocked:
  - from: 2026-11-15T21:32:20Z
    until: 2026-11-15T21:32:46Z
    reason: awaiting the operator on TH-1230
tags: []
touches: [wip/kanban/stories/S-0793-publish-what-is-accepted-and-not-yet-released.md, design/issues/I-0359-flai-release-pending-fails-when-another-process-holds-the-main-checkout-s-git.md, design/issues/summary.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 214
  priced_by: harness
  turns:
    - day: 2026-11-15
      ceremony: 3
      hand_edits: 1
      work: 39
  models:
    - model: claude-opus-5-5
      input: 98
      output: 13866
      cache_read: 7102648
      cache_write: 181424
      cost: 3.1496
  strategic:
    - kind: orchestrator
      seconds: 103
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 42
          output: 344
          cache_read: 1261980
          cache_write: 212274
          cost: 0.4462
        - model: claude-sonnet-5-5
          input: 4
          output: 30
          cache_read: 20402
          cache_write: 25113
          cost: 0.036
formula: publish
---
# S-0793 Publish what is accepted and not yet released

## Goal

Publish what is accepted and not yet released, by the project's release policy.

## Acceptance criteria

- [x] Fetch the remote's branch and tags
- [x] Evaluate the release policy, stopping unless it is met
- [x] Publish what is accepted and not yet released

## Tasks
- T-3384 Fetch the remote's branch and tags
- T-3385 Evaluate the release policy, stopping unless it is met
- T-3386 Publish what is accepted and not yet released

## Notes

Made by `flai formula run publish` from `design/workflows/formulas/publish.md`, with the inputs `judgement=`.

Published on the retry TH-1230 asked for (D-0603): flai 1.98.0 → 1.99.0, a minor bump from S-0741, tag `flai/v1.99.0` on the release commit `a299ee4dc chore: publish flai 1.99.0`, pushed with the branch. The first run had exited 1 on a briefly held `.git/index.lock`, recorded in I-0359.

### Accepted by the orchestrator

- Verified: 612141ec7fccf5f3e65c7a72247ce11f28dd4ce4
- Batch: 20261116-002222
- Enqueued verified: 48b3a1532be9af9642b368f5dc68b060a4b968e6
- At: 2026-11-16T00:36:08Z

Verdict: S-0793 meets all three criteria at 48b3a15; the close-out passed every step there.
- 1: T-3384 done
- 2: T-3385 done
- 3: T-3386 done
