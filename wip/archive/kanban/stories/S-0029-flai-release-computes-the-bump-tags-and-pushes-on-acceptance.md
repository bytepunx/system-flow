---
id: S-0029
type: story
nature: feature
title: flai release computes the bump, tags, and pushes on acceptance
status: done
parent: E-0002
owner: alex
created: 2026-09-16T15:50:25Z
updated: 2026-09-17T03:15:26Z
transitions:
  - to: ready
    at: 2026-09-17T02:04:20Z
    by: agent
  - to: in-progress
    at: 2026-09-17T02:04:20Z
    by: agent
  - to: review
    at: 2026-09-17T02:09:10Z
    by: agent
  - to: done
    at: 2026-09-17T03:15:26Z
    by: alex
tags: [cli, release]
---

# S-0029 flai release computes the bump, tags, and pushes on acceptance

## Goal
Acceptance is the release trigger without an agent remembering: `flai release` reads the accepted item, works out which sub-projects it touched and the bump its delivery type implies, tags, pushes, and reports, so every merge to main has a matching release.

## Acceptance criteria
- [x] `flai release <id>` computes the bump from the item: epic done major; feature story minor; remediation, improvement, docs-only, or dependency update patch; research and experiment refuse with an explanation
- [x] Scope comes from `system-flow.yaml` projects and the files the story's commits touched: the component the story delivers to (from the epic or the story tags) gets the delivery-type bump, every other touched component gets a patch; nothing touched means no tag and a clear message
- [x] The tag is annotated with the item ID and title, created on the acceptance commit, and pushed; `--dry-run` prints what would happen
- [x] Acceptance is one command, `flai accept`, that moves to done, archives, releases (template changelog entry when the template is bumped), commits, tags on that commit, and pushes; `--no-release` skips the release and `flai move ... done` logs a hint
- [x] Pre-release suffixes for research builds are designed as a follow-up note, not implemented

## Tasks
- T-0069 release package: components, versions from tags, touched files from story commits, bump plan
- T-0070 flai release and flai accept commands
- T-0071 Tests on a synthetic git repo, docs, conventions and manifest updates

## Notes
- From the operator on 2026-09-16: releases are automatic on acceptance regardless of visibility; research and experiment stay on branches; pre-release semver suffixes are a later discussion.
