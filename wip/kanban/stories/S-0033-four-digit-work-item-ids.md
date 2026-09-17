---
id: S-0033
type: story
nature: improvement
title: Four-digit work item IDs
status: review
parent: E-0001
owner: alex
created: 2026-09-17T05:52:19Z
updated: 2026-09-17T06:41:33Z
transitions:
  - to: ready
    at: 2026-09-17T06:33:20Z
    by: alex
  - to: in-progress
    at: 2026-09-17T06:33:20Z
    by: alex
  - to: review
    at: 2026-09-17T06:41:33Z
    by: alex
tags: [cli, template]
---

# S-0033 Four-digit work item IDs

## Goal
Work item IDs are zero-padded to four digits (E-0001, S-0034, T-0101) so a project does not outgrow its ID width within its first months; existing three-digit IDs keep working and the standard says how a project migrates.

## Acceptance criteria
- [x] `flai epic|story|task new` allocate four-digit IDs; the next ID is one more than the highest existing number regardless of width
- [x] `flai check`, the dashboard reader, search, and release accept three- and four-digit IDs side by side; sorting is numeric
- [x] `flai migrate ids` (or a documented one-off) renames existing items, narratives, and archive files to four digits and rewrites every reference in bodies, README indexes, board order, issues, and design docs, with `git mv`; dry run first
- [x] work-hierarchy.md, ADR-0003 (refined, not reversed), the template's CLAUDE.md and item templates, and docs/users say four digits
- [x] This repository migrated; flai check --strict clean; CI green

## Tasks
- T-0105 Four-digit allocation in flai; lookups accept short and padded forms; numeric sort verified
- T-0106 flai migrate ids: rename items and narratives with git mv, rewrite references, dry run, tests
- T-0107 Docs: work-hierarchy, ADR-0017 refining ADR-0003, template CLAUDE.md and templates, docs/users, flai-cli
- T-0108 Migrate this repository; flai check --strict; CI green

## Notes
- Requested 2026-09-17 after task IDs passed T-0100. Migration touches every reference in the repo, so the command and a dry run matter more than the format change.
