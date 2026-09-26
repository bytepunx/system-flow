---
id: S-0131
type: story
nature: feature
title: flai stream sync reports conflicts with other open story branches and changes outside the story's touches
status: review
parent: E-0009
owner: alex
created: 2026-09-26T07:59:22Z
updated: 2026-09-26T18:15:33Z
transitions:
  - to: ready
    at: 2026-09-26T17:50:35Z
    by: alex
  - to: in-progress
    at: 2026-09-26T18:04:23Z
    by: agent-S-0131
  - to: review
    at: 2026-09-26T18:15:33Z
    by: agent-S-0131
tags: [cli]
touches: [flai/cmd, flai/internal/mcpserver, docs, design/system, flai/internal/gitver, design/tech]
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

- [x] `flai stream sync S-nnnn` trial-merges the story branch with every other open story branch using `git merge-tree --write-tree`, writing nothing to any worktree, and prints each branch it conflicts with and the conflicting paths; `--json` carries them.
- [x] A conflict found at sync reaches both stories' agents through the MCP `inbox` and the designer's inbox, naming the two stories and the paths.
- [x] Sync lists the paths the branch changed since main that fall outside the story's claim, and says to widen `touches`.
- [x] Integration tests with real git cover a conflicting pair, a clean pair, and a change outside the claim; `make test`, `make integration`, and lint pass.
- [x] `design/system/workflow.md`, `flai-cli.md`, and `docs/users/flai.md` describe it.

## Tasks
- T-0477 flai stream sync trial-merges the branch with every other open story branch and reports the conflicts
- T-0478 A conflict found at sync is a thread flai writes on the story, naming both stories and the paths
- T-0479 flai stream sync lists the paths the branch changed outside the story's touches
- T-0480 Describe the sync checks for builders and users

## Notes

- `syncStoryBranch` is in `flai/cmd/branch.go`. `git merge-tree --write-tree` needs git 2.38 or newer; say so in `design/tech` if a minimum is not already recorded.
- A conflict reaches the agents and the designer as a thread written by `flai`, one per pair of stories (see the narrative's Decisions). Verified on this project: this branch's sync opened TH-0024 for S-0130 and S-0131 on `design/system/workflow.md`.
- Smoke (`make smoke`) stops at the repository check on `wip.overlap` warnings between S-0130 and S-0131, which ran at once because the host flai predates holds (TH-0023), and its markdown lint fails on archived S-0122/S-0124 files and TH-0012, already red on main (TH-0017, I-0027). The template render and installer steps pass.
