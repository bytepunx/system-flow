---
id: T-1843
type: task
nature: feature
title: flai/internal/metrics computes digest.weeks and digest.entries from the threads, the completed stories, and the entries
status: done
parent: S-0434
owner: alex
created: 2026-10-16T06:08:21Z
updated: 2026-10-24T21:52:16Z
transitions:
  - to: ready
    at: 2026-10-24T21:41:46Z
    by: system-flow
  - to: in-progress
    at: 2026-10-24T21:41:46Z
    by: system-flow
  - to: done
    at: 2026-10-24T21:52:16Z
    by: agent-S-0434
stream: S-0434
tags: [cli]
touches: [flai/internal/metrics/digest.go, flai/internal/metrics/digest_test.go, flai/internal/metrics/metrics.go, flai/internal/statsread/statsread.go, flai/internal/statsread/statsread_test.go]
after: [T-1842]
usage:
  source: log
  seconds: 630
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 31
      output: 10019
      cache_read: 1677868
      cache_write: 83371
      cost: 1.0084
---
# T-1843 flai/internal/metrics computes digest.weeks and digest.entries from the threads, the completed stories, and the entries

## Work

- `flai/internal/metrics/digest.go`: from the threads, the items, the owner, and the entries, compute `digest.weeks[]` (`week`, `start`, `threads_to_operator`, `stories_completed`, `per_story`, `asked_again`) over the window's ISO weeks and `digest.entries` (`active`, `compacted`, `absorbed`), as `metrics.md` defines them.
- `metrics.go`: an `Entries` input and a `Digest` field in the report, filled by `digest.go`.
- `flai/internal/statsread/statsread.go`: load the entries with `digest.List` and pass the owner.

Waits for T-1842, whose definitions it follows.

## Done when

`digest_test.go` covers a week with repeats, a week with threads and no completed story, a thread by the operator and one by flai left out, and a project with no digest (the entries all 0, the weeks still counted); `flai test flai/internal/metrics flai/internal/statsread` passes.

## Notes

The thread counts do not depend on the digest, so the weeks before it existed are reported too, which is what S-0435 compares against.
