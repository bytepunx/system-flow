---
id: TH-0566
title: "S-0460's agent ended 3 times with the story in progress; restart it when it should go on"
anchor:
  path: wip/kanban/stories/S-0460-a-page-may-carry-a-left-sidebar-of-its-sections-and-the-item-page-shows-one-section-at-a-time-from-it-instead-of-one-long-scroll.md
  item: S-0460
status: resolved
participants: [flai serve, agent-S-0460, orchestrator, alex]
created: 2026-10-21T09:40:47Z
updated: 2026-10-21T11:57:35Z
---

# TH-0566 S-0460's agent ended 3 times with the story in progress; restart it when it should go on

On wip/kanban/stories/S-0460-a-page-may-carry-a-left-sidebar-of-its-sections-and-the-item-page-shows-one-section-at-a-time-from-it-instead-of-one-long-scroll.md.

## Entries

### 2026-10-21T09:40:47Z flai serve
S-0460's agent has ended 3 times in a row with S-0460 in progress, and flai serve has restarted it on its own 2 times, as many as agent.auto_restarts allows on this host, and restarts it no more: the story has no agent now. Its last run ended (exit 1) with S-0460 in in-progress.

flai serve agent restart S-0460, or Retry on the story's page, starts it again, and its count again from 0. flai serve agent set --auto-restarts <n> changes how many times flai serve restarts a story's agent on its own.

### 2026-10-21T11:36:39Z agent-S-0460
S-0460's agent is running again after the restart. T-1963 was in the worktree, uncommitted. Its tests failed because a test helper named `history` shadowed the global, and prettier had split a line of text. Both are fixed, and T-1963 is closed (ff7e417f). Next: T-1964 (docs), then the close-out and review. Why the earlier runs exited 1 is not in the worktree or the narrative.

### 2026-10-21T11:45:54Z orchestrator (recommendation)
Recommended: resolve this thread. Nothing is left to restart, and the thread is now the only thing keeping S-0460 out of the merge queue.

- The agent was restarted, and it finished T-1963 and T-1964 and moved S-0460 to review at 11:44Z.
- flai verify passed every step at the branch head, 03aa9014. The verifier matched all five criteria to changed files, and every changed file is within the story's touches.
- `queue_add` refuses S-0460 for one reason only: "thread TH-0566 is not resolved on S-0460 or its tasks".
- I resolve only threads I opened, so this one is yours.

Once it is resolved, I enqueue S-0460 at 03aa9014 with the evidence I have. S-0461 waits on S-0460 by `after`, and is next in ready.

Source: design/conventions/strategic-agents.md § As the orchestrator

### 2026-10-21T11:57:35Z alex
Resolved: S-0460 was accepted
