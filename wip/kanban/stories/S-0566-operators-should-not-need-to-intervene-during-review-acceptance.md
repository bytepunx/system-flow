---
id: S-0566
type: story
nature: improvement
title: Operators should not need to intervene during review acceptance
status: backlog
owner: alex
created: 2026-11-02T18:21:16Z
updated: 2026-11-02T18:21:16Z
transitions: []
tags: [cli]
touches: [flai/cmd]
agent:
  harness: claude-code
  model: claude-fable-5-1
  config:
    effort: high
cost_of_delay:
  inputs:
    penalty_per_week: 600
    by: alex
    at: 2026-11-02T18:21:16Z
---
# S-0566 Operators should not need to intervene during review acceptance

## Goal

Right now, nearly every story that lands in review fails to merge because it is no longer in sync with main. Rebase fails due to conflicts that flai is not resolving, forcing the operator to try and resolve merge conflicts in files they did not author.

flai exists to automate the development process. It should automatically resolve the rebase on main including any merge conflicts that arise from this operation.

## Acceptance criteria
- [ ] An operator is never instructed to rebase the story’s worktree onto main since flai should do this
- [ ] An operator is never told to resolve merge conflicts during rebase, the story agent should do this
- [ ] files in shared folders like wip and design should always get checked in on review acceptance - there doesn’t need to be a checkbox/permission request since flai can do this directly

## Tasks

## Notes
