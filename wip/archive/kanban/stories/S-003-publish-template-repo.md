---
id: S-003
type: story
nature: feature
title: Publish the template as its own repository
status: done
parent: E-001
owner: agent
created: 2026-09-15T16:09:00Z
updated: 2026-09-17T00:15:07Z
transitions:
  - to: ready
    at: 2026-09-17T00:07:13Z
    by: agent
  - to: in-progress
    at: 2026-09-17T00:07:13Z
    by: agent
  - to: review
    at: 2026-09-17T00:11:38Z
    by: agent
  - to: done
    at: 2026-09-17T00:15:07Z
    by: alex
tags: []
---

# S-003 Publish the template as its own repository

## Goal
Create github.com/bytepunx/system-flow-template from ./template, tag 0.1.0, and switch this repo's manifest and flai default to it.

## Acceptance criteria
- [x] Repository exists with template.yaml at root and a CHANGELOG.md
- [x] flai new against the published repo produces the same output as against ./template
- [x] ./template remains the development copy with a documented sync procedure

## Tasks
- T-063 Create the template repository with gh and push the template history
- T-064 Verify flai new against the published repo matches ./template
- T-065 Document the sync procedure and point defaults at the published repo

## Notes
- Published private at bytepunx/system-flow-template, v1.0.0. I-009 found and fixed on the way.
