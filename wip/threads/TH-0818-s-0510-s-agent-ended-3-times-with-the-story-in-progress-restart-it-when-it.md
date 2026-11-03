---
id: TH-0818
title: "S-0510's agent ended 3 times with the story in progress; restart it when it should go on"
anchor:
  path: wip/kanban/stories/S-0510-a-project-s-manifest-names-an-image-layer-for-its-test-tiers-tools-flai-serve-runs-it-flai-check-warns-when-it-is-not-built-and-this-repository-s-layer-runs-the-close-out-inside-the-container.md
  item: S-0510
status: resolved
participants: [flai serve, agent-S-0510]
created: 2026-11-03T01:09:46Z
updated: 2026-11-03T01:19:04Z
---

# TH-0818 S-0510's agent ended 3 times with the story in progress; restart it when it should go on

On wip/kanban/stories/S-0510-a-project-s-manifest-names-an-image-layer-for-its-test-tiers-tools-flai-serve-runs-it-flai-check-warns-when-it-is-not-built-and-this-repository-s-layer-runs-the-close-out-inside-the-container.md.

## Entries

### 2026-11-03T01:09:46Z flai serve
S-0510's agent has ended 3 times in a row with S-0510 in progress, and flai serve has restarted it on its own 2 times, as many as agent.auto_restarts allows on this host, and restarts it no more: the story has no agent now. Its last run ended (exit 2) with S-0510 in in-progress.

Its agent events since the first of those runs, oldest first:

- 2026-11-03T01:07:02Z started: `the operator restarted it: its last agent could not be started: rmhclod: executable file not found in $PATH`
- 2026-11-03T01:08:55Z failed: `ended (an exit code nobody saw) with S-0510 in in-progress`
- 2026-11-03T01:08:59Z started: `flai serve started it again on its own, automatic restart 1 of 2: its last agent ended (an exit code nobody saw) with S-0510 in in-progress`
- 2026-11-03T01:09:21Z failed: `ended (exit 2) with S-0510 in in-progress`
- 2026-11-03T01:09:25Z started: `flai serve started it again on its own, automatic restart 2 of 2: its last agent ended (exit 2) with S-0510 in in-progress`
- 2026-11-03T01:09:43Z failed: `ended (exit 2) with S-0510 in in-progress`

flai serve agent restart S-0510, or Retry on the story's page, starts it again, and its count again from 0. flai serve agent set --auto-restarts <n> changes how many times flai serve restarts a story's agent on its own.

### 2026-11-03T01:19:02Z agent-S-0510
S-0510's agent is running again (this session) and goes on from where it was: TH-0814 is answered and applied, every criterion is ticked, and what is left is the sync, the close-out, and the move to review. Nothing is owed here; resolving.

### 2026-11-03T01:19:04Z agent-S-0510
Resolved: The agent runs again and carries S-0510 to review.
