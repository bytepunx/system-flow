---
id: TH-0670
title: "S-0503's agent ended 3 times with the story in progress; restart it when it should go on"
anchor:
  path: wip/kanban/stories/S-0503-no-test-tier-selects-scripts-so-a-story-that-changes-only-shell-scripts-and-the-script-tests-it-adds-pass-flai-verify-and-the-merge-queue-without-any-of-them-running.md
  item: S-0503
status: resolved
participants: [flai serve, alex]
created: 2026-10-27T20:35:20Z
updated: 2026-10-28T06:19:20Z
---

# TH-0670 S-0503's agent ended 3 times with the story in progress; restart it when it should go on

On wip/kanban/stories/S-0503-no-test-tier-selects-scripts-so-a-story-that-changes-only-shell-scripts-and-the-script-tests-it-adds-pass-flai-verify-and-the-merge-queue-without-any-of-them-running.md.

## Entries

### 2026-10-27T20:35:20Z flai serve
S-0503's agent has ended 3 times in a row with S-0503 in progress, and flai serve has restarted it on its own 2 times, as many as agent.auto_restarts allows on this host, and restarts it no more: the story has no agent now. Its last run ended (exit 1) with S-0503 in in-progress.

flai serve agent restart S-0503, or Retry on the story's page, starts it again, and its count again from 0. flai serve agent set --auto-restarts <n> changes how many times flai serve restarts a story's agent on its own.

### 2026-10-28T06:19:20Z alex
Resolved.
