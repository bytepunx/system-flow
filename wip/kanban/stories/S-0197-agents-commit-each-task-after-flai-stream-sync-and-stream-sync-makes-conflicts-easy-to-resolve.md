---
id: S-0197
type: story
nature: improvement
title: Agents commit each task after flai stream sync, and stream sync makes conflicts easy to resolve
status: ready
owner: arobson
created: 2026-10-02T09:46:13Z
updated: 2026-10-02T23:25:04Z
transitions:
  - to: ready
    at: 2026-10-02T23:25:04Z
    by: alex
tags: [flai, template]
touches: [flai/internal/storygit, flai/cmd/stream.go, flai/internal/harness, scripts/close-out.sh, template/]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
---
# S-0197 Agents commit each task after flai stream sync, and stream sync makes conflicts easy to resolve

## Goal

[ADR-0069](../../../design/adrs/0069-agents-commit-each-task-on-the-story-branch-after-flai-stream-sync-and-never.md) (2026-10-02) refines ADR-0019: agents use `flai stream sync` for a story branch's git operations, never rebase or merge it by hand, and commit each task on `story/S-nnnn` after syncing and resolving conflicts; they sync again before review. The git convention says so as of 2026-10-02. The rest of the system has to agree and make it easy.

## Acceptance criteria
- [ ] `flai stream sync` refuses to sync over uncommitted changes in the worktree, saying what to commit or stash, and on a conflict lists each conflicting path and how to continue or abort
- [ ] `harness.Prompt` for `claude-code`, the work-management convention, and `scripts/close-out.sh` (here and in the template) say to sync, resolve, test, and commit at each task, and to sync again before review
- [ ] `design/system/flai-cli.md`, `design/system/workflow.md`, and the user guide describe the per-task cycle
- [ ] Tests cover a sync refused over uncommitted changes and a sync that reports a conflict

## Tasks

## Notes
