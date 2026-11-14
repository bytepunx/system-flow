---
id: TH-1202
title: "S-0788's agent ended 3 times with the story in progress; restart it when it should go on"
anchor:
  path: wip/kanban/stories/S-0788-stories-in-review-that-cannot-be-merged-need-an-action.md
  item: S-0788
status: resolved
participants: [flai serve, agent-S-0788]
created: 2026-11-14T18:58:32Z
updated: 2026-11-14T23:13:51Z
---

# TH-1202 S-0788's agent ended 3 times with the story in progress; restart it when it should go on

On wip/kanban/stories/S-0788-stories-in-review-that-cannot-be-merged-need-an-action.md.

## Entries

### 2026-11-14T18:58:32Z flai serve
S-0788's agent has ended 3 times in a row with S-0788 in progress, and flai serve has restarted it on its own 2 times, as many as agent.auto_restarts allows on this host, and restarts it no more: the story has no agent now. Its last run ended (exit 1) with S-0788 in in-progress.

Its agent events since the first of those runs, oldest first:

- 2026-11-14T18:02:49Z started: `it entered ready`
- 2026-11-14T18:58:15Z failed: `ended (an exit code nobody saw) with S-0788 in in-progress`
- 2026-11-14T18:58:19Z started: `flai serve started it again on its own, automatic restart 1 of 2: its last agent ended (an exit code nobody saw) with S-0788 in in-progress`
- 2026-11-14T18:58:21Z failed: `ended (exit 1) with S-0788 in in-progress`
- 2026-11-14T18:58:26Z started: `flai serve started it again on its own, automatic restart 2 of 2: its last agent ended (exit 1) with S-0788 in in-progress`
- 2026-11-14T18:58:28Z failed: `ended (exit 1) with S-0788 in in-progress`

flai serve agent restart S-0788, or Retry on the story's page, starts it again, and its count again from 0. flai serve agent set --auto-restarts <n> changes how many times flai serve restarts a story's agent on its own.

### 2026-11-14T23:13:51Z agent-S-0788
Resolved: The operator restarted S-0788's agent, and it is running again. The close-out passed at a2a165287, the criteria are verified, and the story is moving to review.
