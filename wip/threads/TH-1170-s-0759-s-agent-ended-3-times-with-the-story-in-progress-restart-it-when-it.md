---
id: TH-1170
title: "S-0759's agent ended 3 times with the story in progress; restart it when it should go on"
anchor:
  path: wip/kanban/stories/S-0759-an-operator-names-a-build-command-of-their-own-that-flai-runs-in-place-of-its.md
  item: S-0759
status: resolved
participants: [flai serve, agent-S-0759, alex]
created: 2026-11-13T18:10:30Z
updated: 2026-11-13T19:45:12Z
---

# TH-1170 S-0759's agent ended 3 times with the story in progress; restart it when it should go on

On wip/kanban/stories/S-0759-an-operator-names-a-build-command-of-their-own-that-flai-runs-in-place-of-its.md.

## Entries

### 2026-11-13T18:10:30Z flai serve
S-0759's agent has ended 3 times in a row with S-0759 in progress, and flai serve has restarted it on its own 2 times, as many as agent.auto_restarts allows on this host, and restarts it no more: the story has no agent now. Its last run ended (exit 1) with S-0759 in in-progress.

Its agent events since the first of those runs, oldest first:

- 2026-11-13T17:54:48Z started: `started again in its session, TH-1110 answered`
- 2026-11-13T17:55:50Z failed: `ended (exit 1) with S-0759 in in-progress`
- 2026-11-13T17:55:55Z started: `flai serve started it again on its own, automatic restart 1 of 2: its last agent ended (exit 1) with S-0759 in in-progress`
- 2026-11-13T18:04:01Z failed: `ended (exit 1) with S-0759 in in-progress`
- 2026-11-13T18:04:05Z started: `flai serve started it again on its own, automatic restart 2 of 2: its last agent ended (exit 1) with S-0759 in in-progress`
- 2026-11-13T18:10:25Z failed: `ended (exit 1) with S-0759 in in-progress`

flai serve agent restart S-0759, or Retry on the story's page, starts it again, and its count again from 0. flai serve agent set --auto-restarts <n> changes how many times flai serve restarts a story's agent on its own.

### 2026-11-13T19:44:49Z agent-S-0759
Restarted, and it went on. The close-out had already passed at 3b7a681. The diff review found criteria 1 to 5 met, so S-0759 is going to review now. This thread can be resolved.

### 2026-11-13T19:45:12Z alex
Resolved: S-0759 was accepted
