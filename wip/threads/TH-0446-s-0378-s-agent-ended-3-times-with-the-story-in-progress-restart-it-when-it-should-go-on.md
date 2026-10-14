---
id: TH-0446
title: "S-0378's agent ended 3 times with the story in progress; restart it when it should go on"
anchor:
  path: wip/kanban/stories/S-0378-the-dashboard-opens-a-markdown-editor-from-a-file-link-in-a-report-a-refusal-or-a-page-so-the-operator-can-fix-markdown-that-blocks-an-acceptance.md
  item: S-0378
status: resolved
participants: [flai serve, agent-S-0378, alex]
created: 2026-10-13T20:27:17Z
updated: 2026-10-14T18:25:51Z
---

# TH-0446 S-0378's agent ended 3 times with the story in progress; restart it when it should go on

On wip/kanban/stories/S-0378-the-dashboard-opens-a-markdown-editor-from-a-file-link-in-a-report-a-refusal-or-a-page-so-the-operator-can-fix-markdown-that-blocks-an-acceptance.md.

## Entries

### 2026-10-13T20:27:17Z flai serve
S-0378's agent has ended 3 times in a row with S-0378 in progress, and flai serve has restarted it on its own 2 times, as many as agent.auto_restarts allows on this host, and restarts it no more: the story has no agent now. Its last run ended (exit 0) with S-0378 in in-progress.

flai serve agent restart S-0378, or Retry on the story's page, starts it again, and its count again from 0. flai serve agent set --auto-restarts <n> changes how many times flai serve restarts a story's agent on its own.

### 2026-10-14T03:52:43Z agent-S-0378
Restarted at 03:40Z. The branch is rebased onto main and the close-out passed again at 9f46596b. Every task is done and criteria 1 to 4 are ticked. The move to review is refused for one reason: **TH-0436** still awaits your confirmation of the orchestrator's recommendation (ADR-0145 saves on a story's branch, and `flai accept --json` exits 4 with the findings). Confirm it in the dashboard's inbox or with `flai thread confirm TH-0436`. Then the agent syncs and moves S-0378 to review. A restart before that ends the same way.

### 2026-10-14T18:25:51Z alex
Resolved.
