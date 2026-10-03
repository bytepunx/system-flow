---
id: S-0197
type: story
nature: improvement
title: Agents commit each task after flai stream sync, and stream sync makes conflicts easy to resolve
status: done
owner: arobson
created: 2026-10-02T09:46:13Z
updated: 2026-10-03T01:22:17Z
transitions:
  - to: ready
    at: 2026-10-02T23:25:04Z
    by: alex
  - to: in-progress
    at: 2026-10-03T00:38:46Z
    by: agent-S-0197
  - to: review
    at: 2026-10-03T01:09:05Z
    by: agent-S-0197
  - to: done
    at: 2026-10-03T01:22:17Z
    by: alex
tags: [flai, template]
touches: [flai/internal/storygit, flai/cmd/stream.go, flai/cmd/branch.go, flai/cmd/stream_sync_test.go, flai/internal/harness, scripts/close-out.sh, template, design/conventions/git.md, design/conventions/work-management.md, design/conventions/delegation.md, design/system/flai-cli.md, design/system/workflow.md, docs/users/flai.md, docs/users/flai-reference.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 1846
  models:
    - model: claude-opus-5-5
      input: 264
      output: 102017
      cache_read: 8934919
      cache_write: 354851
      cost: 6.0149
    - model: claude-sonnet-5-5
      input: 48
      output: 9573
      cache_read: 746116
      cache_write: 77925
      cost: 0.4399
---
# S-0197 Agents commit each task after flai stream sync, and stream sync makes conflicts easy to resolve

## Goal

[ADR-0069](../../../design/adrs/0069-agents-commit-each-task-on-the-story-branch-after-flai-stream-sync-and-never.md) (2026-10-02) refines ADR-0019: agents use `flai stream sync` for a story branch's git operations, never rebase or merge it by hand, and commit each task on `story/S-nnnn` after syncing and resolving conflicts; they sync again before review. The git convention says so as of 2026-10-02. The rest of the system has to agree and make it easy.

## Acceptance criteria
- [x] `flai stream sync` refuses to sync over uncommitted changes in the worktree, saying what to commit or stash, and on a conflict lists each conflicting path and how to continue or abort
- [x] `harness.Prompt` for `claude-code`, the work-management convention, and `scripts/close-out.sh` (here and in the template) say to sync, resolve, test, and commit at each task, and to sync again before review
- [x] `design/system/flai-cli.md`, `design/system/workflow.md`, and the user guide describe the per-task cycle
- [x] Tests cover a sync refused over uncommitted changes and a sync that reports a conflict

## Tasks
- T-0723 flai stream sync refuses uncommitted changes and lists each conflict with how to continue or abort
- T-0724 The claude-code prompt says to sync, resolve, test, and commit at each task, and to sync again before review
- T-0725 The work-management convention and close-out.sh say the per-task cycle, here and in the template
- T-0726 flai-cli.md, workflow.md, and the user guide describe the per-task cycle and how sync refuses and reports

## Notes

- The designer decided on TH-0072 that, because sync refuses uncommitted changes, each task is committed before it is synced: commit, `flai stream sync`, resolve, test, commit fixes. `git.md` says so in both copies.
