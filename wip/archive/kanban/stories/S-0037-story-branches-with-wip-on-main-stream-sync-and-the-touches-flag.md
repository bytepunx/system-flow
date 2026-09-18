---
id: S-0037
type: story
nature: feature
title: Story branches with wip on main, stream sync, and the touches flag
status: done
parent: E-0006
owner: alex
created: 2026-09-17T19:46:28Z
updated: 2026-09-18T04:55:23Z
transitions:
  - to: ready
    at: 2026-09-17T20:20:10Z
    by: alex
  - to: in-progress
    at: 2026-09-17T20:20:11Z
    by: alex
  - to: review
    at: 2026-09-17T20:27:32Z
    by: alex
  - to: done
    at: 2026-09-18T04:55:23Z
    by: alex
tags: [dashboard, cli]
touches: [flai/internal/workitem, flai/cmd, flai/internal/check, flaiover/src/routes/docs, docs/users/flai.md]
---

# S-0037 Story branches with wip on main, stream sync, and the touches flag

## Goal
Agents work each story on its own branch in a git worktree while `wip/` stays in the main checkout so the board is always live; `flai stream sync` rebases the branch onto main at task transitions; `flai accept` rebases, merges, tags, and pushes; a `touches` list in front matter tells everyone what a story or task is working on.

## Acceptance criteria
- [x] `flai stream open S-nnnn` creates branch `story/S-nnnn` from main and a worktree under `.flai-cache/worktrees/S-nnnn` (path printed); wip files written by flai from inside a worktree land in the main checkout
- [x] `flai stream sync S-nnnn` rebases the story branch onto main, reports conflicts file by file, and refuses to continue past them; the agent runs it at every task transition (work-management project addition)
- [x] `flai accept` rebases and merges the story branch into main (fast-forward when possible), removes the worktree and branch, then tags and pushes as today; refuses when the rebase has conflicts
- [x] `touches: [paths or components]` is optional front matter on stories and tasks; `flai check` warns `wip.overlap` when two in-progress items touch the same path; `flai board` shows touches; the dashboard shows a badge on documents being worked on
- [x] ADR-0019 accepted; git.md project addition replaces the 'work happens on main' line; work-hierarchy.md documents `touches`; docs/users updated; behavior and integration tests cover open, sync, accept, and overlap

## Tasks
- T-0120 flai stream open creates story/S-nnnn and a worktree under .flai-cache/worktrees; wip writes from a worktree land in the main checkout
- T-0121 flai stream sync rebases the branch onto main and reports conflicts; flai accept rebases, merges, removes worktree and branch before tagging
- T-0122 touches front matter on stories and tasks; flai check wip.overlap; flai board shows touches; dashboard badge on touched documents
- T-0123 Behavior and integration tests with real git; docs/users, work-hierarchy, git convention, CLI design

## Notes
Operator's call 2026-09-17: agents branch, rebase and merge at accept. Board state must stay visible on main, hence worktrees with wip in the main checkout rather than wip on branches. Human edits in the dashboard land on main and surface as small rebase conflicts.
