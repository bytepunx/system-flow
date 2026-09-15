---
id: S-026
type: story
nature: feature
title: Adopt the conventions in this repository
status: done
parent: E-005
owner: alex
created: 2026-09-15T18:40:50Z
updated: 2026-09-15T23:00:55Z
transitions:
  - to: ready
    at: 2026-09-15T22:37:53Z
    by: agent
  - to: in-progress
    at: 2026-09-15T22:37:53Z
    by: agent
  - to: review
    at: 2026-09-15T22:42:26Z
    by: agent
  - to: in-progress
    at: 2026-09-15T22:49:30Z
    by: alex
  - to: review
    at: 2026-09-15T22:52:23Z
    by: agent
  - to: done
    at: 2026-09-15T23:00:55Z
    by: alex
tags: [conventions]
---

# S-026 Adopt the conventions in this repository

## Goal
Copy the template's design/conventions into this repository after the operator has edited them, add this project's additions below the markers, and change how work is done here to match: CLAUDE.md links to them, the narratives and stories follow them, and anything the conventions contradict in current practice is fixed.

## Acceptance criteria
- [x] design/conventions/ in this repository matches the template above every marker
- [x] Project additions below the markers capture norms specific to this monorepo (template development copy, sub-project layout, dogfooding flai)
- [x] Norms that lived only in CLAUDE.md now live in a convention file and CLAUDE.md links to them
- [x] A pass over current practice: narratives, board, commits, and docs are checked against every convention and deviations are fixed or recorded as project additions
- [x] flai check --strict passes

## Tasks
- T-040 Review template conventions and apply consistency fixes
- T-041 ADR-0014 design/issues and continuous-improvement design
- T-042 Copy conventions into design/conventions with project additions
- T-043 Seed design/issues and add scripts/ with Makefile delegation
- T-044 Align CLAUDE.md, layout docs, and current practice
- T-045 Apply operator answers: test tiers, commit timing, kluster

## Notes
- Operator edits to the template added a tenth topic (continuous improvement), design/issues (ADR-0014), and scripts/; all adopted here.
- Six open questions recorded in the narrative for the operator.
- Blocked until the operator has finished editing the template conventions; pull only when asked. If S-023 is done and the operator has not asked, skip to S-024.
- 2026-09-15T22:49:30Z: moved to in-progress: operator answered open questions 1 to 3: no task-boundary commits, kluster is public, three test tiers with separate commands and folders
