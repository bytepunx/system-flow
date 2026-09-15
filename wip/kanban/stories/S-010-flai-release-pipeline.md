---
id: S-010
type: story
nature: feature
title: Release pipeline for flai
status: review
parent: E-002
owner: agent
created: 2026-09-15T16:09:00Z
updated: 2026-09-15T18:18:23Z
transitions:
  - to: ready
    at: 2026-09-15T18:11:50Z
    by: agent
  - to: in-progress
    at: 2026-09-15T18:11:50Z
    by: agent
  - to: review
    at: 2026-09-15T18:18:23Z
    by: agent
tags: []
---

# S-010 Release pipeline for flai

## Goal
Tagged releases publish binaries for linux, darwin, and windows.

## Acceptance criteria
- [x] GoReleaser config builds all targets
- [ ] Tag flai/v0.1.0 produces a GitHub release (needs the repo pushed to GitHub and the tag pushed, a human step)
- [x] README documents install via go install and release download

## Tasks
- T-030 GoReleaser config with version injection
- T-031 Release workflow on flai/v* tags
- T-032 Snapshot build verification and install docs

## Notes
- Dry run exposed a git-ignored test fixture and an empty kanban/tasks folder in the committed tree; both fixed here.
- Dependabot config and Makefile release targets added while here; design/tech/ci.md already promised Dependabot.
