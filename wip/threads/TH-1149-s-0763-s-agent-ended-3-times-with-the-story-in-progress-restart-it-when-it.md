---
id: TH-1149
title: "S-0763's agent ended 3 times with the story in progress; restart it when it should go on"
anchor:
  path: wip/kanban/stories/S-0763-each-deploy-gives-a-story-s-namespace-a-kubeconfig-held-to-it-and-a-story-s.md
  item: S-0763
status: open
participants: [flai serve]
created: 2026-11-13T09:48:02Z
updated: 2026-11-13T09:48:02Z
---

# TH-1149 S-0763's agent ended 3 times with the story in progress; restart it when it should go on

On wip/kanban/stories/S-0763-each-deploy-gives-a-story-s-namespace-a-kubeconfig-held-to-it-and-a-story-s.md.

## Entries

### 2026-11-13T09:48:02Z flai serve
S-0763's agent has ended 3 times in a row with S-0763 in progress, and flai serve has restarted it on its own 2 times, as many as agent.auto_restarts allows on this host, and restarts it no more: the story has no agent now. Its last run ended (exit 1) with S-0763 in in-progress.

Its agent events since the first of those runs, oldest first:

- 2026-11-13T09:06:47Z started: `it entered ready`
- 2026-11-13T09:47:09Z failed: `ended (exit 1) with S-0763 in in-progress`
- 2026-11-13T09:47:14Z started: `flai serve started it again on its own, automatic restart 1 of 2: its last agent ended (exit 1) with S-0763 in in-progress`
- 2026-11-13T09:47:44Z failed: `ended (exit 1) with S-0763 in in-progress`
- 2026-11-13T09:47:50Z started: `flai serve started it again on its own, automatic restart 2 of 2: its last agent ended (exit 1) with S-0763 in in-progress`
- 2026-11-13T09:47:57Z failed: `ended (exit 1) with S-0763 in in-progress`

flai serve agent restart S-0763, or Retry on the story's page, starts it again, and its count again from 0. flai serve agent set --auto-restarts <n> changes how many times flai serve restarts a story's agent on its own.
