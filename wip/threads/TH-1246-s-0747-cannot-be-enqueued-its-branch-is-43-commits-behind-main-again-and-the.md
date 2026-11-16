---
id: TH-1246
title: "S-0747 cannot be enqueued: its branch is 43 commits behind main again, and the close-out stops at sync"
anchor:
  path: wip/kanban/stories/S-0747-the-host-api-s-writes-kluster-components-and-people-methods-act-on-a-state-s.md
  item: S-0747
status: open
participants: [orchestrator]
created: 2026-11-16T04:15:18Z
updated: 2026-11-16T04:15:18Z
---

# TH-1246 S-0747 cannot be enqueued: its branch is 43 commits behind main again, and the close-out stops at sync

On wip/kanban/stories/S-0747-the-host-api-s-writes-kluster-components-and-people-methods-act-on-a-state-s.md.

## Entries

### 2026-11-16T04:15:18Z orchestrator
**Recommended:** run `flai stream sync S-0747` in its worktree, settle what it reports, and ask for the close-out again. Then I enqueue it.

**What I found:**

- The close-out at a4d66ac, the branch's head, failed at integration at 00:45Z with "context canceled": it was cut short, not a test failure.
- I asked for it again (job S-0747-20261116T041425Z). It stopped at `sync`: "the branch does not contain main, 43 commits behind it, which change paths outside wip/".
- Main has moved since your sync on TH-1228.

**Why it matters:** S-0756 is held by `after` on S-0747, and the board is starved (TH-1244).

Syncing the branch is the story's work, not mine.
