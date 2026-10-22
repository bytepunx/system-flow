---
id: TH-0586
title: "S-0483's agent ended 3 times with the story in progress; restart it when it should go on"
anchor:
  path: wip/kanban/stories/S-0483-accepting-a-story-that-adds-a-flaiover-dependency-leaves-the-main-checkout-s-flaiover-node-modules-without-it-so-every-close-out-s-integration-tier-fails-testeverylistedflaiovercomponentcompilesintheharness-until-someone-installs.md
  item: S-0483
status: resolved
participants: [flai serve, agent-S-0483]
created: 2026-10-21T15:03:03Z
updated: 2026-10-21T20:44:52Z
---

# TH-0586 S-0483's agent ended 3 times with the story in progress; restart it when it should go on

On wip/kanban/stories/S-0483-accepting-a-story-that-adds-a-flaiover-dependency-leaves-the-main-checkout-s-flaiover-node-modules-without-it-so-every-close-out-s-integration-tier-fails-testeverylistedflaiovercomponentcompilesintheharness-until-someone-installs.md.

## Entries

### 2026-10-21T15:03:03Z flai serve
S-0483's agent has ended 3 times in a row with S-0483 in progress, and flai serve has restarted it on its own 2 times, as many as agent.auto_restarts allows on this host, and restarts it no more: the story has no agent now. Its last run ended (exit 1) with S-0483 in in-progress.

flai serve agent restart S-0483, or Retry on the story's page, starts it again, and its count again from 0. flai serve agent set --auto-restarts <n> changes how many times flai serve restarts a story's agent on its own.

### 2026-10-21T20:44:52Z agent-S-0483
Resolved: The operator restarted S-0483's agent; it is closing the story out.
