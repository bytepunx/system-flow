---
id: TH-0566
title: "S-0460's agent ended 3 times with the story in progress; restart it when it should go on"
anchor:
  path: wip/kanban/stories/S-0460-a-page-may-carry-a-left-sidebar-of-its-sections-and-the-item-page-shows-one-section-at-a-time-from-it-instead-of-one-long-scroll.md
  item: S-0460
status: open
participants: [flai serve]
created: 2026-10-21T09:40:47Z
updated: 2026-10-21T09:40:47Z
---

# TH-0566 S-0460's agent ended 3 times with the story in progress; restart it when it should go on

On wip/kanban/stories/S-0460-a-page-may-carry-a-left-sidebar-of-its-sections-and-the-item-page-shows-one-section-at-a-time-from-it-instead-of-one-long-scroll.md.

## Entries

### 2026-10-21T09:40:47Z flai serve
S-0460's agent has ended 3 times in a row with S-0460 in progress, and flai serve has restarted it on its own 2 times, as many as agent.auto_restarts allows on this host, and restarts it no more: the story has no agent now. Its last run ended (exit 1) with S-0460 in in-progress.

flai serve agent restart S-0460, or Retry on the story's page, starts it again, and its count again from 0. flai serve agent set --auto-restarts <n> changes how many times flai serve restarts a story's agent on its own.
