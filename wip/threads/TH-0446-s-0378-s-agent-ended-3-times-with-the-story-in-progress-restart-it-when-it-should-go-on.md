---
id: TH-0446
title: "S-0378's agent ended 3 times with the story in progress; restart it when it should go on"
anchor:
  path: wip/kanban/stories/S-0378-the-dashboard-opens-a-markdown-editor-from-a-file-link-in-a-report-a-refusal-or-a-page-so-the-operator-can-fix-markdown-that-blocks-an-acceptance.md
  item: S-0378
status: open
participants: [flai serve]
created: 2026-10-13T20:27:17Z
updated: 2026-10-13T20:27:17Z
---

# TH-0446 S-0378's agent ended 3 times with the story in progress; restart it when it should go on

On wip/kanban/stories/S-0378-the-dashboard-opens-a-markdown-editor-from-a-file-link-in-a-report-a-refusal-or-a-page-so-the-operator-can-fix-markdown-that-blocks-an-acceptance.md.

## Entries

### 2026-10-13T20:27:17Z flai serve
S-0378's agent has ended 3 times in a row with S-0378 in progress, and flai serve has restarted it on its own 2 times, as many as agent.auto_restarts allows on this host, and restarts it no more: the story has no agent now. Its last run ended (exit 0) with S-0378 in in-progress.

flai serve agent restart S-0378, or Retry on the story's page, starts it again, and its count again from 0. flai serve agent set --auto-restarts <n> changes how many times flai serve restarts a story's agent on its own.
