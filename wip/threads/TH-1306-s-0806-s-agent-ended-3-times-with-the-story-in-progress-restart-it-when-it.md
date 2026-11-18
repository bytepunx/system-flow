---
id: TH-1306
title: "S-0806's agent ended 3 times with the story in progress; restart it when it should go on"
anchor:
  path: wip/kanban/stories/S-0806-the-workflow-file-takes-verify-on-a-transition-runs-on-a-state-and-verify-at.md
  item: S-0806
status: resolved
participants: [flai serve, agent-S-0806]
created: 2026-11-18T04:06:06Z
updated: 2026-11-18T07:19:14Z
---

# TH-1306 S-0806's agent ended 3 times with the story in progress; restart it when it should go on

On wip/kanban/stories/S-0806-the-workflow-file-takes-verify-on-a-transition-runs-on-a-state-and-verify-at.md.

## Entries

### 2026-11-18T04:06:06Z flai serve
S-0806's agent has ended 3 times in a row with S-0806 in progress, and flai serve has restarted it on its own 2 times, as many as agent.auto_restarts allows on this host, and restarts it no more: the story has no agent now. Its last run ended (exit 1) with S-0806 in in-progress.

Its agent events since the first of those runs, oldest first:

- 2026-11-18T04:05:39Z started: `started again in its session, TH-1289 answered`
- 2026-11-18T04:05:42Z failed: `ended (exit 1) with S-0806 in in-progress`
- 2026-11-18T04:05:48Z started: `flai serve started it again on its own, automatic restart 1 of 2: its last agent ended (exit 1) with S-0806 in in-progress`
- 2026-11-18T04:05:50Z failed: `ended (exit 1) with S-0806 in in-progress`
- 2026-11-18T04:05:57Z started: `flai serve started it again on its own, automatic restart 2 of 2: its last agent ended (exit 1) with S-0806 in in-progress`
- 2026-11-18T04:05:59Z failed: `ended (exit 1) with S-0806 in in-progress`

flai serve agent restart S-0806, or Retry on the story's page, starts it again, and its count again from 0. flai serve agent set --auto-restarts <n> changes how many times flai serve restarts a story's agent on its own.

### 2026-11-18T07:19:14Z agent-S-0806
Resolved: The operator restarted S-0806's agent, and it is working the story again.
