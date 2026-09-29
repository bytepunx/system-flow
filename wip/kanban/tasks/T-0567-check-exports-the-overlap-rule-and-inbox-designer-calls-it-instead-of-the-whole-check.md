---
id: T-0567
type: task
nature: feature
title: check exports the overlap rule and inbox.designer calls it instead of the whole check
status: done
parent: S-0158
owner: alex
created: 2026-09-29T19:30:30Z
updated: 2026-09-29T19:33:11Z
transitions:
  - to: ready
    at: 2026-09-29T19:30:43Z
    by: agent-S-0158
  - to: in-progress
    at: 2026-09-29T19:30:44Z
    by: agent-S-0158
  - to: done
    at: 2026-09-29T19:33:11Z
    by: agent-S-0158
stream: S-0158
tags: []
usage:
  source: log
  seconds: 147
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 46
      output: 10468
      cache_read: 2022512
      cache_write: 61835
      cost: 1.1087
---

# T-0567 check exports the overlap rule and inbox.designer calls it instead of the whole check

## Work

Export the `wip.overlap` rule from `flai/internal/check` as `check.Overlaps(repo, items)`, the one implementation that `check.Run` also uses. Make `inbox.designer` in `flai/internal/hostapi/people.go` call it over the items it has already listed, and drop its `check.run` phase.

## Done when

- `check.Run` and `inbox.designer` both find overlaps through `check.Overlaps`.
- A behaviour test in `internal/hostapi` builds a repository with overlapping in-progress touches and asserts that `inbox.designer`'s overlap entries equal the `wip.overlap` findings of `check.Run`.
- `scripts/flai-test.sh` passes.

## Notes
