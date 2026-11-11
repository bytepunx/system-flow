---
id: T-3077
type: task
nature: improvement
title: Close I-0287 with what fixed it
status: done
parent: S-0675
owner: alex
created: 2026-11-10T10:37:00Z
updated: 2026-11-10T20:16:17Z
transitions:
  - to: ready
    at: 2026-11-10T20:14:17Z
    by: system-flow
  - to: in-progress
    at: 2026-11-10T20:14:18Z
    by: system-flow
  - to: done
    at: 2026-11-10T20:16:17Z
    by: agent-S-0675
stream: S-0675
tags: [issues]
touches: [design/issues/I-0287-the-merge-queue-builds-one-story-batches-and-runs-three-verifications-per-story.md, design/issues/summary.md]
after: [T-3074, T-3075]
usage:
  source: log
  seconds: 119
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 6
      output: 1584
      cache_read: 608180
      cache_write: 14720
      cost: 0.2459
---
# T-3077 Close I-0287 with what fixed it

## Work

Run `flai issue close I-0287 --reason "<reason>"`. The reason names the ADR T-3073 wrote, the pass-over T-3074 built in `flai/internal/mergequeue/land.go` and the test that reproduces I-0287, and the `dropped` count T-3075 added to `queue.weeks[]`. It says that `queue.weeks[].dropped` and `verifications_per_landed` in `flai stats` measure the fix after the release. It waits for T-3074 and T-3075 because it closes the issue on what they built. It runs beside the documentation task, since the two share no file.

## Done when

- I-0287 is closed with that reason, and `design/issues/summary.md` no longer lists it as open.
- Criterion 2 of S-0675 can be ticked.

## Notes
