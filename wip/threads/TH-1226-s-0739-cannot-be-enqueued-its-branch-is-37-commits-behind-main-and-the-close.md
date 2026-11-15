---
id: TH-1226
title: "S-0739 cannot be enqueued: its branch is 37 commits behind main, and the close-out stops at sync"
anchor:
  path: wip/kanban/stories/S-0739-workitem-s-epic-following-board-view-task-plan-holds-order-and-promotion-read-a.md
  item: S-0739
status: open
participants: [orchestrator]
created: 2026-11-15T21:09:40Z
updated: 2026-11-15T21:09:40Z
---

# TH-1226 S-0739 cannot be enqueued: its branch is 37 commits behind main, and the close-out stops at sync

On wip/kanban/stories/S-0739-workitem-s-epic-following-board-view-task-plan-holds-order-and-promotion-read-a.md.

## Entries

### 2026-11-15T21:09:40Z orchestrator
**Recommended:** run `flai stream sync S-0739` in its worktree, settle what it reports, and ask for the close-out again (the MCP tool `verify`). Then I enqueue it.

**What I found:**

- The accept dry run at da74987 had one blocker: `verify_behind_main` (ADR-0291). Its verify ran over main at 9fb9d24.
- The close-out I asked for, job S-0739-20261115T210803Z, stopped at `sync`: "the branch does not contain main, 37 commits behind it, which change paths outside wip/".
- `queue_add` now refuses it as `verify_failed` at step sync.
- The verifier matched all 5 criteria to changed files at da74987. Every changed file is inside its touches.

**Why it matters:** review is full (5 of 5), which starves the board (TH-1224).

Rebasing the story branch is the story's work, not mine.
