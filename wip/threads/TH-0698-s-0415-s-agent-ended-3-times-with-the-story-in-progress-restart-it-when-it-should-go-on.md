---
id: TH-0698
title: "S-0415's agent ended 3 times with the story in progress; restart it when it should go on"
anchor:
  path: wip/kanban/stories/S-0415-flai-prime-role-gives-a-role-a-project-defines-its-pack-within-the-budget-from-the-topics-and-the-item-kind-its-role-file-names-and-the-three-built-in-roles-prime-from-their-files-unchanged.md
  item: S-0415
status: resolved
participants: [flai serve, agent-S-0415, alex]
created: 2026-10-28T15:14:38Z
updated: 2026-10-28T20:49:37Z
---

# TH-0698 S-0415's agent ended 3 times with the story in progress; restart it when it should go on

On wip/kanban/stories/S-0415-flai-prime-role-gives-a-role-a-project-defines-its-pack-within-the-budget-from-the-topics-and-the-item-kind-its-role-file-names-and-the-three-built-in-roles-prime-from-their-files-unchanged.md.

## Entries

### 2026-10-28T15:14:38Z flai serve
S-0415's agent has ended 3 times in a row with S-0415 in progress, and flai serve has restarted it on its own 2 times, as many as agent.auto_restarts allows on this host, and restarts it no more: the story has no agent now. Its last run ended (exit 1) with S-0415 in in-progress.

flai serve agent restart S-0415, or Retry on the story's page, starts it again, and its count again from 0. flai serve agent set --auto-restarts <n> changes how many times flai serve restarts a story's agent on its own.

### 2026-10-28T20:38:03Z agent-S-0415
Restarted, S-0415's agent took up where the earlier runs stopped: T-1726 to T-1728 are done, the close-out passed every step, every criterion is ticked, and S-0415 is moving to review now. Nothing more is needed here.

### 2026-10-28T20:49:37Z alex
Resolved: S-0415 was accepted
