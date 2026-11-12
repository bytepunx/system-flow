---
id: TH-1085
title: "The board is starved: nothing in ready can start and no backlog story can be promoted"
anchor:
  path: wip/kanban/board.md
status: resolved
participants: [flai serve]
created: 2026-11-12T00:41:16Z
updated: 2026-11-12T03:39:42Z
---

# TH-1085 The board is starved: nothing in ready can start and no backlog story can be promoted

On wip/kanban/board.md.

## Entries

### 2026-11-12T00:41:16Z flai serve
No story in ready can start, no backlog story can be promoted, and 2 in-progress lanes stand idle. flai serve adds an entry here when the list below changes, and resolves this thread once a story can start.

This host has no reminder channel, so without one nothing reminds you of this thread off the dashboard. Set one in a shell on the host with `flai serve notify set --url <webhook>` or `flai serve notify set --command <command>`.

Each backlog story, in the pull order, and the decision it waits on:

- S-0241 Measure planned stories again with forked task sub-agents once a headless session offers forks: a draft; finalize it.
- S-0659 An ADR settles how a story in progress is deployed into its kluster namespace at its agent's request: who asks, what the agent may run, when the namespace goes, and how many live at once: a draft; finalize it.
- S-0660 A story's agent asks flai serve to deploy its story from its worktree while it is in progress, and the namespace stays until the story is accepted or cancelled: held by `after`; the root of its chain is S-0659 An ADR settles how a story in progress is deployed into its kluster namespace at its agent's request: who asks, what the agent may run, when the namespace goes, and how many live at once (backlog).
- S-0661 A story's agent runs the project's checks against its story's kluster namespace before review, and its own commands with the namespace's variables: held by `after`; the root of its chain is S-0659 An ADR settles how a story in progress is deployed into its kluster namespace at its agent's request: who asks, what the agent may run, when the namespace goes, and how many live at once (backlog).
- S-0662 The move to review, the close-out, and the review page say what was deployed and checked in the story's namespace and at which commit, and agents are told to verify there first: held by `after`; the root of its chain is S-0659 An ADR settles how a story in progress is deployed into its kluster namespace at its agent's request: who asks, what the agent may run, when the namespace goes, and how many live at once (backlog).
- S-0663 The Deployments page shows the namespace of each story in progress, and the operator deploys, checks, and tears down from it: held by `after`; the root of its chain is S-0659 An ADR settles how a story in progress is deployed into its kluster namespace at its agent's request: who asks, what the agent may run, when the namespace goes, and how many live at once (backlog).
- S-0683 Two stories accepted apart leave main red: S-0535 reworded to-backlog's condition, and S-0416's TestADefinedRoleIsHeldToItsFile still expects the old words: a draft; finalize it. It waits on TH-1044 Plan for S-0683: a hand acceptance waits for a verify over main's new paths; six tasks in four layers.
- S-0709 flai stream sync reports a path as undone by the branch when main gains a commit to it while the sync runs: a draft; finalize it. It waits on TH-1014 Plan for S-0709: pin main's commit across the sync's rebase.
- S-0723 flai check finds `item.blocked-closed` outside the story at close-out: a draft; finalize it.
- S-0725 The container harness's design and operator docs still say the close-out and the permission_prompt check do not run in the image, and that the orchestrator's verify and accept run inside its container: a draft; finalize it.
- S-0726 A git rebase --continue run by hand after a sync's stop meets later stops on issue files and summary.md that flai's sync would have merged itself: a draft; finalize it.

### 2026-11-12T03:39:41Z flai serve
The list changed. No story in ready can start, no backlog story can be promoted, and 3 in-progress lanes stand idle.

This host has no reminder channel, so without one nothing reminds you of this thread off the dashboard. Set one in a shell on the host with `flai serve notify set --url <webhook>` or `flai serve notify set --command <command>`.

Each backlog story, in the pull order, and the decision it waits on:

- S-0241 Measure planned stories again with forked task sub-agents once a headless session offers forks: a draft; finalize it.
- S-0659 An ADR settles how a story in progress is deployed into its kluster namespace at its agent's request: who asks, what the agent may run, when the namespace goes, and how many live at once: a draft; finalize it.
- S-0660 A story's agent asks flai serve to deploy its story from its worktree while it is in progress, and the namespace stays until the story is accepted or cancelled: held by `after`; the root of its chain is S-0659 An ADR settles how a story in progress is deployed into its kluster namespace at its agent's request: who asks, what the agent may run, when the namespace goes, and how many live at once (backlog).
- S-0661 A story's agent runs the project's checks against its story's kluster namespace before review, and its own commands with the namespace's variables: held by `after`; the root of its chain is S-0659 An ADR settles how a story in progress is deployed into its kluster namespace at its agent's request: who asks, what the agent may run, when the namespace goes, and how many live at once (backlog).
- S-0662 The move to review, the close-out, and the review page say what was deployed and checked in the story's namespace and at which commit, and agents are told to verify there first: held by `after`; the root of its chain is S-0659 An ADR settles how a story in progress is deployed into its kluster namespace at its agent's request: who asks, what the agent may run, when the namespace goes, and how many live at once (backlog).
- S-0663 The Deployments page shows the namespace of each story in progress, and the operator deploys, checks, and tears down from it: held by `after`; the root of its chain is S-0659 An ADR settles how a story in progress is deployed into its kluster namespace at its agent's request: who asks, what the agent may run, when the namespace goes, and how many live at once (backlog).
- S-0683 Two stories accepted apart leave main red: S-0535 reworded to-backlog's condition, and S-0416's TestADefinedRoleIsHeldToItsFile still expects the old words: a draft; finalize it. It waits on TH-1044 Plan for S-0683: a hand acceptance waits for a verify over main's new paths; six tasks in four layers.
- S-0709 flai stream sync reports a path as undone by the branch when main gains a commit to it while the sync runs: a draft; finalize it. It waits on TH-1014 Plan for S-0709: pin main's commit across the sync's rebase.
- S-0723 flai check finds `item.blocked-closed` outside the story at close-out: a draft; finalize it.
- S-0725 The container harness's design and operator docs still say the close-out and the permission_prompt check do not run in the image, and that the orchestrator's verify and accept run inside its container: a draft; finalize it.
- S-0726 A git rebase --continue run by hand after a sync's stop meets later stops on issue files and summary.md that flai's sync would have merged itself: a draft; finalize it.
- S-0729 Publish what is accepted and not yet released: no forecast duration; no cost of delay value.

### 2026-11-12T03:39:42Z flai serve
Resolved: S-0729 Publish what is accepted and not yet released in ready can start, so the board is no longer starved.
