---
id: S-029
type: story
nature: feature
title: flai release computes the bump, tags, and pushes on acceptance
status: backlog
parent: E-002
owner: alex
created: 2026-09-16T15:50:25Z
updated: 2026-09-16T15:50:25Z
transitions: []
tags: [cli, release]
---

# S-029 flai release computes the bump, tags, and pushes on acceptance

## Goal
Acceptance is the release trigger without an agent remembering: `flai release` reads the accepted item, works out which sub-projects it touched and the bump its delivery type implies, tags, pushes, and reports, so every merge to main has a matching release.

## Acceptance criteria
- [ ] `flai release <id>` computes the bump from the item: epic done major; feature story minor; remediation, improvement, docs-only, or dependency update patch; research and experiment refuse with an explanation
- [ ] Scope comes from `system-flow.yaml` projects and the files the story's commits touched: the component the story delivers to (from the epic or the story tags) gets the delivery-type bump, every other touched component gets a patch; nothing touched means no tag and a clear message
- [ ] The tag is annotated with the item ID and title, created on the acceptance commit, and pushed; `--dry-run` prints what would happen
- [ ] `flai move <id> done` calls release unless `--no-release`; the template changelog gets an entry when the template is bumped
- [ ] Pre-release suffixes for research builds are designed as a follow-up note, not implemented

## Tasks

## Notes
- From the operator on 2026-09-16: releases are automatic on acceptance regardless of visibility; research and experiment stay on branches; pre-release semver suffixes are a later discussion.
