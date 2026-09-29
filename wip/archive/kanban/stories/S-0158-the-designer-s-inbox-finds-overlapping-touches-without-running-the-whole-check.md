---
id: S-0158
type: story
nature: improvement
title: The designer's inbox finds overlapping touches without running the whole check
status: done
parent: E-0012
owner: alex
created: 2026-09-29T07:00:29Z
updated: 2026-09-29T19:41:29Z
transitions:
  - to: ready
    at: 2026-09-29T19:18:28Z
    by: alex
  - to: in-progress
    at: 2026-09-29T19:30:14Z
    by: agent-S-0158
  - to: review
    at: 2026-09-29T19:33:41Z
    by: agent-S-0158
  - to: done
    at: 2026-09-29T19:41:29Z
    by: alex
tags: [cli]
topics: [server-side, back-end]
touches: [flai/internal/hostapi, flai/internal/check, design/system/flai-cli.md, design/system/server-performance.md, docs/users/flai.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 231
  models:
    - model: claude-opus-5-5
      input: 82
      output: 18481
      cache_read: 3570747
      cache_write: 109169
      cost: 1.9574
---
# S-0158 The designer's inbox finds overlapping touches without running the whole check

## Goal

`inbox.designer`, which the inbox badge and `/api/projects` ask for, runs the whole of `flai check` to find one rule's findings, overlapping touches: 153 ms of its 163 ms. Cause 3 of `design/system/server-performance.md`. Find the overlaps with the rule alone.

## Acceptance criteria
- [x] `inbox.designer` on this repository takes under 20 ms, and its event has no `check.run` phase.
- [x] Its overlap entries are the same as the `wip.overlap` findings of `flai check`: the rule keeps one implementation, and a behaviour test compares the two.

## Tasks
- T-0567 check exports the overlap rule and inbox.designer calls it instead of the whole check
- T-0568 Measure inbox.designer on this repository and record it in the design

## Notes

Measured by S-0152: `inbox.designer` 163 ms, `check.run=152.9`. Proposed: export the overlap rule from `internal/check` as a function over items and call it from the inbox.

Done as proposed: `check.Overlaps(repo, items)`, which `check.Run` also uses. Measured on 2026-09-29 with `flai hostapi --timing inbox.designer`, three runs each on this repository: median 178 ms before (`check.run` 167), 12.7 ms after (`threads.read` 5.6, `repo.list` 4.7, `check.overlap` under 0.1), the answer byte for byte the same. `TestInboxDesignerOverlapsAreTheChecksFindings` compares the inbox's overlaps with `check.Run`'s `wip.overlap` findings; `TestReadMethodsMarkWhereTheTimeGoes` fails if `check.run` is a phase of it again.
