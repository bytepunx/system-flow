---
id: S-033
type: story
nature: improvement
title: Four-digit work item IDs
status: backlog
parent: E-001
owner: alex
created: 2026-09-17T05:52:19Z
updated: 2026-09-17T05:52:19Z
transitions: []
tags: [cli, template]
---

# S-033 Four-digit work item IDs

## Goal
Work item IDs are zero-padded to four digits (E-0001, S-0034, T-0101) so a project does not outgrow its ID width within its first months; existing three-digit IDs keep working and the standard says how a project migrates.

## Acceptance criteria
- [ ] `flai epic|story|task new` allocate four-digit IDs; the next ID is one more than the highest existing number regardless of width
- [ ] `flai check`, the dashboard reader, search, and release accept three- and four-digit IDs side by side; sorting is numeric
- [ ] `flai migrate ids` (or a documented one-off) renames existing items, narratives, and archive files to four digits and rewrites every reference in bodies, README indexes, board order, issues, and design docs, with `git mv`; dry run first
- [ ] work-hierarchy.md, ADR-0003 (refined, not reversed), the template's CLAUDE.md and item templates, and docs/users say four digits
- [ ] This repository migrated; flai check --strict clean; CI green

## Tasks

## Notes
- Requested 2026-09-17 after task IDs passed T-100. Migration touches every reference in the repo, so the command and a dry run matter more than the format change.
