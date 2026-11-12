---
id: TH-1092
title: "The board is starved: nothing in ready can start and no backlog story can be promoted"
anchor:
  path: wip/kanban/board.md
status: resolved
participants: [flai serve]
created: 2026-11-12T15:50:09Z
updated: 2026-11-12T15:53:12Z
---

# TH-1092 The board is starved: nothing in ready can start and no backlog story can be promoted

On wip/kanban/board.md.

## Entries

### 2026-11-12T15:50:09Z flai serve
No story in ready can start, no backlog story can be promoted, and 1 in-progress lane stands idle. flai serve adds an entry here when the list below changes, and resolves this thread once a story can start.

This host has no reminder channel, so without one nothing reminds you of this thread off the dashboard. Set one in a shell on the host with `flai serve notify set --url <webhook>` or `flai serve notify set --command <command>`.

Each backlog story, in the pull order, and the decision it waits on:

- S-0241 Measure planned stories again with forked task sub-agents once a headless session offers forks: a draft; finalize it.
- S-0659 An ADR settles how a story in progress is deployed into its kluster namespace at its agent's request: who asks, what the agent may run, when the namespace goes, and how many live at once: no cost of delay value.
- S-0660 A story's agent asks flai serve to deploy its story from its worktree while it is in progress, and the namespace stays until the story is accepted or cancelled: held by `after`; the root of its chain is S-0659 An ADR settles how a story in progress is deployed into its kluster namespace at its agent's request: who asks, what the agent may run, when the namespace goes, and how many live at once (backlog).
- S-0661 A story's agent runs the project's checks against its story's kluster namespace before review, and its own commands with the namespace's variables: held by `after`; the root of its chain is S-0659 An ADR settles how a story in progress is deployed into its kluster namespace at its agent's request: who asks, what the agent may run, when the namespace goes, and how many live at once (backlog).
- S-0662 The move to review, the close-out, and the review page say what was deployed and checked in the story's namespace and at which commit, and agents are told to verify there first: held by `after`; the root of its chain is S-0659 An ADR settles how a story in progress is deployed into its kluster namespace at its agent's request: who asks, what the agent may run, when the namespace goes, and how many live at once (backlog).
- S-0663 The Deployments page shows the namespace of each story in progress, and the operator deploys, checks, and tears down from it: held by `after`; the root of its chain is S-0659 An ADR settles how a story in progress is deployed into its kluster namespace at its agent's request: who asks, what the agent may run, when the namespace goes, and how many live at once (backlog).
- S-0723 flai check finds `item.blocked-closed` outside the story at close-out: held (no-touches): declares no touches, so it may change what S-0683 (in progress) and S-0709 (in progress) change; starts when it declares touches that overlap no story's in progress, or when S-0683 and S-0709 move to review, are cancelled, or are sent back; no forecast duration; no cost of delay value.
- S-0725 The container harness's design and operator docs still say the close-out and the permission_prompt check do not run in the image, and that the orchestrator's verify and accept run inside its container: a draft; finalize it.
- S-0726 A git rebase --continue run by hand after a sync's stop meets later stops on issue files and summary.md that flai's sync would have merged itself: a draft; finalize it.
- S-0730 A story's agent can overwrite a work item file whole with Write, front matter included, and nothing refuses it: a draft; finalize it.
- S-0732 flai's MCP servers listen on every address their clients need, loopback and the container engine's gateway together, set as a list or found by flai for the host's engine: a draft; finalize it.

### 2026-11-12T15:52:45Z flai serve
The list changed. No story in ready can start, no backlog story can be promoted, and 1 in-progress lane stands idle.

This host has no reminder channel, so without one nothing reminds you of this thread off the dashboard. Set one in a shell on the host with `flai serve notify set --url <webhook>` or `flai serve notify set --command <command>`.

Each backlog story, in the pull order, and the decision it waits on:

- S-0241 Measure planned stories again with forked task sub-agents once a headless session offers forks: a draft; finalize it.
- S-0659 An ADR settles how a story in progress is deployed into its kluster namespace at its agent's request: who asks, what the agent may run, when the namespace goes, and how many live at once: no cost of delay value.
- S-0660 A story's agent asks flai serve to deploy its story from its worktree while it is in progress, and the namespace stays until the story is accepted or cancelled: held by `after`; the root of its chain is S-0659 An ADR settles how a story in progress is deployed into its kluster namespace at its agent's request: who asks, what the agent may run, when the namespace goes, and how many live at once (backlog).
- S-0661 A story's agent runs the project's checks against its story's kluster namespace before review, and its own commands with the namespace's variables: held by `after`; the root of its chain is S-0659 An ADR settles how a story in progress is deployed into its kluster namespace at its agent's request: who asks, what the agent may run, when the namespace goes, and how many live at once (backlog).
- S-0662 The move to review, the close-out, and the review page say what was deployed and checked in the story's namespace and at which commit, and agents are told to verify there first: held by `after`; the root of its chain is S-0659 An ADR settles how a story in progress is deployed into its kluster namespace at its agent's request: who asks, what the agent may run, when the namespace goes, and how many live at once (backlog).
- S-0663 The Deployments page shows the namespace of each story in progress, and the operator deploys, checks, and tears down from it: held by `after`; the root of its chain is S-0659 An ADR settles how a story in progress is deployed into its kluster namespace at its agent's request: who asks, what the agent may run, when the namespace goes, and how many live at once (backlog).
- S-0723 flai check finds `item.blocked-closed` outside the story at close-out: no forecast duration; no cost of delay value.
- S-0725 The container harness's design and operator docs still say the close-out and the permission_prompt check do not run in the image, and that the orchestrator's verify and accept run inside its container: a draft; finalize it.
- S-0726 A git rebase --continue run by hand after a sync's stop meets later stops on issue files and summary.md that flai's sync would have merged itself: a draft; finalize it.
- S-0730 A story's agent can overwrite a work item file whole with Write, front matter included, and nothing refuses it: a draft; finalize it.
- S-0732 flai's MCP servers listen on every address their clients need, loopback and the container engine's gateway together, set as a list or found by flai for the host's engine: a draft; finalize it.

### 2026-11-12T15:53:01Z flai serve
The list changed. No story in ready can start, no backlog story can be promoted, and 1 in-progress lane stands idle.

This host has no reminder channel, so without one nothing reminds you of this thread off the dashboard. Set one in a shell on the host with `flai serve notify set --url <webhook>` or `flai serve notify set --command <command>`.

Each backlog story, in the pull order, and the decision it waits on:

- S-0241 Measure planned stories again with forked task sub-agents once a headless session offers forks: a draft; finalize it.
- S-0659 An ADR settles how a story in progress is deployed into its kluster namespace at its agent's request: who asks, what the agent may run, when the namespace goes, and how many live at once: no cost of delay value.
- S-0660 A story's agent asks flai serve to deploy its story from its worktree while it is in progress, and the namespace stays until the story is accepted or cancelled: held by `after`; the root of its chain is S-0659 An ADR settles how a story in progress is deployed into its kluster namespace at its agent's request: who asks, what the agent may run, when the namespace goes, and how many live at once (backlog).
- S-0661 A story's agent runs the project's checks against its story's kluster namespace before review, and its own commands with the namespace's variables: held by `after`; the root of its chain is S-0659 An ADR settles how a story in progress is deployed into its kluster namespace at its agent's request: who asks, what the agent may run, when the namespace goes, and how many live at once (backlog).
- S-0662 The move to review, the close-out, and the review page say what was deployed and checked in the story's namespace and at which commit, and agents are told to verify there first: held by `after`; the root of its chain is S-0659 An ADR settles how a story in progress is deployed into its kluster namespace at its agent's request: who asks, what the agent may run, when the namespace goes, and how many live at once (backlog).
- S-0663 The Deployments page shows the namespace of each story in progress, and the operator deploys, checks, and tears down from it: held by `after`; the root of its chain is S-0659 An ADR settles how a story in progress is deployed into its kluster namespace at its agent's request: who asks, what the agent may run, when the namespace goes, and how many live at once (backlog).
- S-0723 flai check finds `item.blocked-closed` outside the story at close-out: no forecast duration; no cost of delay value.
- S-0725 The container harness's design and operator docs still say the close-out and the permission_prompt check do not run in the image, and that the orchestrator's verify and accept run inside its container: a draft; finalize it. It waits on TH-1093 Plan for S-0725: correct the two stale container passages and close I-0350.
- S-0726 A git rebase --continue run by hand after a sync's stop meets later stops on issue files and summary.md that flai's sync would have merged itself: a draft; finalize it.
- S-0730 A story's agent can overwrite a work item file whole with Write, front matter included, and nothing refuses it: a draft; finalize it.
- S-0732 flai's MCP servers listen on every address their clients need, loopback and the container engine's gateway together, set as a list or found by flai for the host's engine: a draft; finalize it.

### 2026-11-12T15:53:12Z flai serve
Resolved: S-0725 The container harness's design and operator docs still say the close-out and the permission_prompt check do not run in the image, and that the orchestrator's verify and accept run inside its container can be promoted from the backlog, so the board is no longer starved.
