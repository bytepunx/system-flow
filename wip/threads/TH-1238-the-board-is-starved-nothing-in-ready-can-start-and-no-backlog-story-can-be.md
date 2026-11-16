---
id: TH-1238
title: "The board is starved: nothing in ready can start and no backlog story can be promoted"
anchor:
  path: wip/kanban/board.md
status: resolved
participants: [flai serve]
created: 2026-11-16T01:36:49Z
updated: 2026-11-16T01:44:45Z
---

# TH-1238 The board is starved: nothing in ready can start and no backlog story can be promoted

On wip/kanban/board.md.

## Entries

### 2026-11-16T01:36:49Z flai serve
No story in ready can start, no backlog story can be promoted, and 1 in-progress lane stands idle. flai serve adds an entry here when the list below changes, and resolves this thread once a story can start.

This host has no reminder channel, so without one nothing reminds you of this thread off the dashboard. Set one in a shell on the host with `flai serve notify set --url <webhook>` or `flai serve notify set --command <command>`.

Each backlog story, in the pull order, and the decision it waits on:

- S-0241 Measure planned stories again with forked task sub-agents once a headless session offers forks: a draft; finalize it. It waits on TH-1165 S-0241 stays a draft: it is blocked until a headless session offers forks.
- S-0752 flai stats reports each timer a workflow defines beside the baseline's figures, for every item and over the window, past items included: held by `after`; the root of its chain is S-0750 flai stats computes lead, cycle, queue, and review time from the baseline workflow's timers, and the metrics read a state's kind, with every figure byte for byte as today (in-progress).
- S-0753 The dashboard's board draws the columns, their kinds, and the moves between them from what flai gives, and holds no list of states of its own: held by `after`; the root of its chain is S-0739 workitem's epic following, board view, task plan, holds, order, and promotion read a state's kind instead of its name (review).
- S-0754 The dashboard's charts draw the workflow's states and each timer it defines from what flai gives, and an operator's own timer gets a chart beside lead and cycle time: held by `after`; the roots of its chain are S-0750 flai stats computes lead, cycle, queue, and review time from the baseline workflow's timers, and the metrics read a state's kind, with every figure byte for byte as today (in-progress) and S-0739 workitem's epic following, board view, task plan, holds, order, and promotion read a state's kind instead of its name (review).
- S-0756 In a project made from the template, a state added between in-progress and review with a rule and a timer of its own, in the workflow file alone, is followed by the board, the moves, the check, the pull, the acceptance, and the charts: held by `after`; the roots of its chain are S-0739 workitem's epic following, board view, task plan, holds, order, and promotion read a state's kind instead of its name (review), S-0742 flai serve starts, restarts, and stops agents, plans again, and finds its reasons by a state's kind and a transition's effect instead of a state's name (review), S-0747 The host API's writes, kluster, components, and people methods act on a state's kind and who may make a transition instead of a state's name (review) and S-0750 flai stats computes lead, cycle, queue, and review time from the baseline workflow's timers, and the metrics read a state's kind, with every figure byte for byte as today (in-progress).
- S-0758 The dashboard's item page, inbox, review, activity, deployments, task plan, and cards read a state's kind from flai instead of comparing its name: held by `after`; the root of its chain is S-0739 workitem's epic following, board view, task plan, holds, order, and promotion read a state's kind instead of its name (review).
- S-0797 The merge queue's preflight holds the queue on main's files when the installed flai is older than the flai main's files are written for, and the hold cannot clear until a release is installed: a draft; finalize it.
- S-0798 A close-out asked for a story whose branch is behind main stops at sync without syncing, and the failed report makes a story the merge queue could batch not eligible: a draft; finalize it.
- S-0799 After a sync moves a story's head, the review page, flai queue, and the acceptance preview report the verify failure from before the sync while a new close-out runs at the new head: a draft; finalize it.

### 2026-11-16T01:39:44Z flai serve
The list changed. No story in ready can start, no backlog story can be promoted, and 1 in-progress lane stands idle.

This host has no reminder channel, so without one nothing reminds you of this thread off the dashboard. Set one in a shell on the host with `flai serve notify set --url <webhook>` or `flai serve notify set --command <command>`.

Each backlog story, in the pull order, and the decision it waits on:

- S-0241 Measure planned stories again with forked task sub-agents once a headless session offers forks: a draft; finalize it. It waits on TH-1165 S-0241 stays a draft: it is blocked until a headless session offers forks.
- S-0752 flai stats reports each timer a workflow defines beside the baseline's figures, for every item and over the window, past items included: held by `after`; the root of its chain is S-0750 flai stats computes lead, cycle, queue, and review time from the baseline workflow's timers, and the metrics read a state's kind, with every figure byte for byte as today (in-progress).
- S-0753 The dashboard's board draws the columns, their kinds, and the moves between them from what flai gives, and holds no list of states of its own: held by `after`; the root of its chain is S-0739 workitem's epic following, board view, task plan, holds, order, and promotion read a state's kind instead of its name (review).
- S-0754 The dashboard's charts draw the workflow's states and each timer it defines from what flai gives, and an operator's own timer gets a chart beside lead and cycle time: held by `after`; the roots of its chain are S-0750 flai stats computes lead, cycle, queue, and review time from the baseline workflow's timers, and the metrics read a state's kind, with every figure byte for byte as today (in-progress) and S-0739 workitem's epic following, board view, task plan, holds, order, and promotion read a state's kind instead of its name (review).
- S-0756 In a project made from the template, a state added between in-progress and review with a rule and a timer of its own, in the workflow file alone, is followed by the board, the moves, the check, the pull, the acceptance, and the charts: held by `after`; the roots of its chain are S-0739 workitem's epic following, board view, task plan, holds, order, and promotion read a state's kind instead of its name (review), S-0742 flai serve starts, restarts, and stops agents, plans again, and finds its reasons by a state's kind and a transition's effect instead of a state's name (review), S-0747 The host API's writes, kluster, components, and people methods act on a state's kind and who may make a transition instead of a state's name (review) and S-0750 flai stats computes lead, cycle, queue, and review time from the baseline workflow's timers, and the metrics read a state's kind, with every figure byte for byte as today (in-progress).
- S-0758 The dashboard's item page, inbox, review, activity, deployments, task plan, and cards read a state's kind from flai instead of comparing its name: held by `after`; the root of its chain is S-0739 workitem's epic following, board view, task plan, holds, order, and promotion read a state's kind instead of its name (review).
- S-0797 The merge queue's preflight holds the queue on main's files when the installed flai is older than the flai main's files are written for, and the hold cannot clear until a release is installed: a draft; finalize it.
- S-0798 A close-out asked for a story whose branch is behind main stops at sync without syncing, and the failed report makes a story the merge queue could batch not eligible: a draft; finalize it.
- S-0799 After a sync moves a story's head, the review page, flai queue, and the acceptance preview report the verify failure from before the sync while a new close-out runs at the new head: a draft; finalize it.
- S-0800 flai stream sync reports a rebased branch undoes main's line when the branch removes a closing brace elsewhere in a Go file main added a closing brace to, and the merge queue then never batches the story: a draft; finalize it.

### 2026-11-16T01:44:45Z flai serve
Resolved: S-0800 flai stream sync reports a rebased branch undoes main's line when the branch removes a closing brace elsewhere in a Go file main added a closing brace to, and the merge queue then never batches the story can be promoted from the backlog, so the board is no longer starved.
