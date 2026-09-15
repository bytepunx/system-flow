---
id: S-001
type: story
nature: feature
title: Define the standard
status: done
parent: E-001
owner: agent
created: 2026-09-15T16:10:00Z
updated: 2026-09-15T18:05:00Z
transitions:
  - to: ready
    at: 2026-09-15T16:11:00Z
    by: agent
  - to: in-progress
    at: 2026-09-15T16:12:00Z
    by: agent
  - to: review
    at: 2026-09-15T16:45:00Z
    by: agent
  - to: done
    at: 2026-09-15T18:05:00Z
    by: alex
tags: []
---

# S-001 Define the standard

## Goal
Write the living design for the repository layout, work hierarchy, workflow, agent narrative, metrics, manifest, template, CLI, and dashboard, and capture each decision as an ADR, so that template and tooling work can start from a settled base.

## Acceptance criteria
- [x] design/system has a document for each topic listed in its README
- [x] Every decision that shapes tooling has an ADR in design/adrs
- [x] design/tech records every chosen technology with version and rationale
- [x] This repo's wip/ folder follows the standard it defines

## Tasks
- T-001 Write design/system documents
- T-002 Write ADRs 0001 to 0012
- T-003 Write design/tech documents
- T-004 Bootstrap wip/ for this repo

## Notes
