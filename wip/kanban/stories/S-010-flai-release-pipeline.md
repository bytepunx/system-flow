---
id: S-010
type: story
nature: feature
title: Release pipeline for flai
status: review
parent: E-002
owner: agent
created: 2026-09-15T16:09:00Z
updated: 2026-09-16T04:28:57Z
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
  - to: in-progress
    at: 2026-09-15T23:01:11Z
    by: alex
  - to: review
    at: 2026-09-16T04:28:57Z
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
- T-046 gh repository creation and release bump rules in git.md

## Notes
- Operator standards added during review: gh for remotes with org and visibility asked first, release bumps by delivery type. Repository created private at bytepunx/system-flow.
- Dry run exposed a git-ignored test fixture and an empty kanban/tasks folder in the committed tree; both fixed here.
- Dependabot config and Makefile release targets added while here; design/tech/ci.md already promised Dependabot.
- 2026-09-15T23:01:11Z: moved to in-progress: operator added standards: create remotes with gh (ask org and visibility), release bumps by delivery type
