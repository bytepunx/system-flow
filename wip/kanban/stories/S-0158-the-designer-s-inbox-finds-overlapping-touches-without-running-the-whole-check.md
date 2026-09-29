---
id: S-0158
type: story
nature: improvement
title: The designer's inbox finds overlapping touches without running the whole check
status: backlog
parent: E-0012
owner: alex
created: 2026-09-29T07:00:29Z
updated: 2026-09-29T07:00:29Z
transitions: []
tags: [cli]
topics: [server-side, back-end]
touches: [flai/internal/hostapi, flai/internal/check]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
---
# S-0158 The designer's inbox finds overlapping touches without running the whole check

## Goal

`inbox.designer`, which the inbox badge and `/api/projects` ask for, runs the whole of `flai check` to find one rule's findings, overlapping touches: 153 ms of its 163 ms. Cause 3 of `design/system/server-performance.md`. Find the overlaps with the rule alone.

## Acceptance criteria
- [ ] `inbox.designer` on this repository takes under 20 ms, and its event has no `check.run` phase.
- [ ] Its overlap entries are the same as the `wip.overlap` findings of `flai check`: the rule keeps one implementation, and a behaviour test compares the two.

## Tasks

## Notes

Measured by S-0152: `inbox.designer` 163 ms, `check.run=152.9`. Proposed: export the overlap rule from `internal/check` as a function over items and call it from the inbox.
