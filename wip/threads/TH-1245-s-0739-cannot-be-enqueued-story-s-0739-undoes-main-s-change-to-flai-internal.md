---
id: TH-1245
title: "S-0739 cannot be enqueued: story/S-0739 undoes main's change to flai/internal/workitem/board.go"
anchor:
  path: wip/kanban/stories/S-0739-workitem-s-epic-following-board-view-task-plan-holds-order-and-promotion-read-a.md
  item: S-0739
status: resolved
participants: [orchestrator, alex]
created: 2026-11-16T04:14:49Z
updated: 2026-11-17T02:47:35Z
---

# TH-1245 S-0739 cannot be enqueued: story/S-0739 undoes main's change to flai/internal/workitem/board.go

On wip/kanban/stories/S-0739-workitem-s-epic-following-board-view-task-plan-holds-order-and-promotion-read-a.md.

## Entries

### 2026-11-16T04:14:49Z orchestrator
**Recommended:** agent-S-0739 syncs with main (`flai stream sync S-0739`; the branch is 35 commits behind) and checks `flai/internal/workitem/board.go` against main. It keeps what main added, or restores it with `flai stream restore S-0739 --path flai/internal/workitem/board.go` and applies the story's change again. Then it asks for the close-out again, and I enqueue the story.

**What I found:**

- The close-out passed every step at 7e66d35, the branch's head.
- The merge queue refuses to batch the story with `undone_main`: "story/S-0739 undoes the main branch's change to flai/internal/workitem/board.go" (ADR-0155).
- S-0790's board-lock commits to board.go (97cda524e, ee472f11f) are already in the branch. So the branch's own diff removes lines main has, among them `DefaultWIPLimits`, `LimitedColumns`, and the `SetLimit` column check.

**Alternative:** if those lines are meant to go, because the columns and limits now come from the workflow's kinds, accept the removal yourself: `flai accept S-0739 --keep-removal flai/internal/workitem/board.go --reason "<why>"`. I cannot do that from the queue.

**Why it matters:** S-0753, S-0754, S-0756, and S-0758 are held by `after` on S-0739, and the board is starved (TH-1244).

### 2026-11-17T02:47:35Z alex
Resolved: S-0739 was accepted
