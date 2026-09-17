---
id: S-0008
type: story
nature: feature
title: flai check and flai stats
status: done
parent: E-0002
owner: agent
created: 2026-09-15T16:09:00Z
updated: 2026-09-15T18:11:49Z
transitions:
  - to: ready
    at: 2026-09-15T18:00:59Z
    by: agent
  - to: in-progress
    at: 2026-09-15T18:00:59Z
    by: agent
  - to: review
    at: 2026-09-15T18:09:36Z
    by: agent
  - to: done
    at: 2026-09-15T18:11:49Z
    by: alex
tags: []
---

# S-0008 flai check and flai stats

## Goal
The reference implementation of validation and metrics.

## Acceptance criteria
- [x] check reports every rule in work-hierarchy.md and workflow.md with file and line
- [x] stats implements every aggregate in metrics.md with --json output
- [x] A fixture repo in testdata has known expected numbers

## Tasks
- T-0024 Check: item, structure, and relationship rules
- T-0025 Check: narratives, board, and documentation front matter
- T-0026 Stats: per-item metrics and aggregates
- T-0027 Stats: burn-up and cumulative flow series
- T-0028 Fixture repo with expected numbers
- T-0029 Docs and design updates for check and stats

## Notes
- First real run of check on this repo found an unchecked criterion in archived S-0002 and backlog stories in the board order; both resolved (see narrative decisions).
- Root CI workflow system-flow-check.yml added; it builds flai from source so CI never depends on a published release.
