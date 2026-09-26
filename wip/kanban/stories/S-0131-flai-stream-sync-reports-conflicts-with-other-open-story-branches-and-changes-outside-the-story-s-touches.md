---
id: S-0131
type: story
nature: feature
title: flai stream sync reports conflicts with other open story branches and changes outside the story's touches
status: ready
parent: E-0009
owner: alex
created: 2026-09-26T07:59:22Z
updated: 2026-09-26T17:50:35Z
transitions:
  - to: ready
    at: 2026-09-26T17:50:35Z
    by: alex
tags: [cli]
touches: [flai/cmd, flai/internal/mcpserver, docs, design/system]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
---
# S-0131 flai stream sync reports conflicts with other open story branches and changes outside the story's touches

## Goal

Build the safety net of [ADR-0046](../../../design/adrs/0046-a-ready-story-whose-claim-overlaps-an-open-story-s-is-held-yellow-and-with-its.md) at sync: catch what the declared claims missed while both stories are still open, not after one is accepted.

## Acceptance criteria

- [ ] `flai stream sync S-nnnn` trial-merges the story branch with every other open story branch using `git merge-tree --write-tree`, writing nothing to any worktree, and prints each branch it conflicts with and the conflicting paths; `--json` carries them.
- [ ] A conflict found at sync reaches both stories' agents through the MCP `inbox` and the designer's inbox, naming the two stories and the paths.
- [ ] Sync lists the paths the branch changed since main that fall outside the story's claim, and says to widen `touches`.
- [ ] Integration tests with real git cover a conflicting pair, a clean pair, and a change outside the claim; `make test`, `make integration`, and lint pass.
- [ ] `design/system/workflow.md`, `flai-cli.md`, and `docs/users/flai.md` describe it.

## Tasks

## Notes

- `syncStoryBranch` is in `flai/cmd/branch.go`. `git merge-tree --write-tree` needs git 2.38 or newer; say so in `design/tech` if a minimum is not already recorded.
