---
id: TH-1224
title: "The board is starved: nothing in ready can start and no backlog story can be promoted"
anchor:
  path: wip/kanban/board.md
status: resolved
participants: [flai serve]
created: 2026-11-15T21:01:55Z
updated: 2026-11-15T21:11:44Z
---

# TH-1224 The board is starved: nothing in ready can start and no backlog story can be promoted

On wip/kanban/board.md.

## Entries

### 2026-11-15T21:01:55Z flai serve
No story in ready can start, no backlog story can be promoted, and 1 in-progress lane stands idle. flai serve adds an entry here when the list below changes, and resolves this thread once a story can start.

This host has no reminder channel, so without one nothing reminds you of this thread off the dashboard. Set one in a shell on the host with `flai serve notify set --url <webhook>` or `flai serve notify set --command <command>`.

Each backlog story, in the pull order, and the decision it waits on:

- S-0241 Measure planned stories again with forked task sub-agents once a headless session offers forks: a draft; finalize it. It waits on TH-1165 S-0241 stays a draft: it is blocked until a headless session offers forks.
- S-0752 flai stats reports each timer a workflow defines beside the baseline's figures, for every item and over the window, past items included: held by `after`; the root of its chain is S-0750 flai stats computes lead, cycle, queue, and review time from the baseline workflow's timers, and the metrics read a state's kind, with every figure byte for byte as today (in-progress).
- S-0753 The dashboard's board draws the columns, their kinds, and the moves between them from what flai gives, and holds no list of states of its own: held by `after`; the root of its chain is S-0739 workitem's epic following, board view, task plan, holds, order, and promotion read a state's kind instead of its name (review).
- S-0754 The dashboard's charts draw the workflow's states and each timer it defines from what flai gives, and an operator's own timer gets a chart beside lead and cycle time: held by `after`; the roots of its chain are S-0750 flai stats computes lead, cycle, queue, and review time from the baseline workflow's timers, and the metrics read a state's kind, with every figure byte for byte as today (in-progress) and S-0739 workitem's epic following, board view, task plan, holds, order, and promotion read a state's kind instead of its name (review).
- S-0755 flai reads a project's own workflow file in place of the baseline, the template ships the baseline, flai upgrade keeps the project's file, and the guides say how to change the board: held by `after`; the root of its chain is S-0741 flai check validates the workflow and the items against it, and its item rules read a state's kind instead of its name (review).
- S-0756 In a project made from the template, a state added between in-progress and review with a rule and a timer of its own, in the workflow file alone, is followed by the board, the moves, the check, the pull, the acceptance, and the charts: held by `after`; the roots of its chain are S-0739 workitem's epic following, board view, task plan, holds, order, and promotion read a state's kind instead of its name (review), S-0741 flai check validates the workflow and the items against it, and its item rules read a state's kind instead of its name (review), S-0742 flai serve starts, restarts, and stops agents, plans again, and finds its reasons by a state's kind and a transition's effect instead of a state's name (review), S-0747 The host API's writes, kluster, components, and people methods act on a state's kind and who may make a transition instead of a state's name (review), S-0748 A role file's to-review and to-backlog conditions become one condition that names any state the workflow has, and the baseline role files use it (in-progress) and S-0750 flai stats computes lead, cycle, queue, and review time from the baseline workflow's timers, and the metrics read a state's kind, with every figure byte for byte as today (in-progress).
- S-0758 The dashboard's item page, inbox, review, activity, deployments, task plan, and cards read a state's kind from flai instead of comparing its name: held by `after`; the root of its chain is S-0739 workitem's epic following, board view, task plan, holds, order, and promotion read a state's kind instead of its name (review).
- S-0790 flai check finds `board.order` outside the story at close-out: a draft; finalize it.

### 2026-11-15T21:04:40Z flai serve
The list changed. No story in ready can start, no backlog story can be promoted, and 1 in-progress lane stands idle.

This host has no reminder channel, so without one nothing reminds you of this thread off the dashboard. Set one in a shell on the host with `flai serve notify set --url <webhook>` or `flai serve notify set --command <command>`.

Each backlog story, in the pull order, and the decision it waits on:

- S-0241 Measure planned stories again with forked task sub-agents once a headless session offers forks: a draft; finalize it. It waits on TH-1165 S-0241 stays a draft: it is blocked until a headless session offers forks.
- S-0752 flai stats reports each timer a workflow defines beside the baseline's figures, for every item and over the window, past items included: held by `after`; the root of its chain is S-0750 flai stats computes lead, cycle, queue, and review time from the baseline workflow's timers, and the metrics read a state's kind, with every figure byte for byte as today (in-progress).
- S-0753 The dashboard's board draws the columns, their kinds, and the moves between them from what flai gives, and holds no list of states of its own: held by `after`; the root of its chain is S-0739 workitem's epic following, board view, task plan, holds, order, and promotion read a state's kind instead of its name (review).
- S-0754 The dashboard's charts draw the workflow's states and each timer it defines from what flai gives, and an operator's own timer gets a chart beside lead and cycle time: held by `after`; the roots of its chain are S-0750 flai stats computes lead, cycle, queue, and review time from the baseline workflow's timers, and the metrics read a state's kind, with every figure byte for byte as today (in-progress) and S-0739 workitem's epic following, board view, task plan, holds, order, and promotion read a state's kind instead of its name (review).
- S-0755 flai reads a project's own workflow file in place of the baseline, the template ships the baseline, flai upgrade keeps the project's file, and the guides say how to change the board: held by `after`; the root of its chain is S-0741 flai check validates the workflow and the items against it, and its item rules read a state's kind instead of its name (review).
- S-0756 In a project made from the template, a state added between in-progress and review with a rule and a timer of its own, in the workflow file alone, is followed by the board, the moves, the check, the pull, the acceptance, and the charts: held by `after`; the roots of its chain are S-0739 workitem's epic following, board view, task plan, holds, order, and promotion read a state's kind instead of its name (review), S-0741 flai check validates the workflow and the items against it, and its item rules read a state's kind instead of its name (review), S-0742 flai serve starts, restarts, and stops agents, plans again, and finds its reasons by a state's kind and a transition's effect instead of a state's name (review), S-0747 The host API's writes, kluster, components, and people methods act on a state's kind and who may make a transition instead of a state's name (review), S-0748 A role file's to-review and to-backlog conditions become one condition that names any state the workflow has, and the baseline role files use it (in-progress) and S-0750 flai stats computes lead, cycle, queue, and review time from the baseline workflow's timers, and the metrics read a state's kind, with every figure byte for byte as today (in-progress).
- S-0758 The dashboard's item page, inbox, review, activity, deployments, task plan, and cards read a state's kind from flai instead of comparing its name: held by `after`; the root of its chain is S-0739 workitem's epic following, board view, task plan, holds, order, and promotion read a state's kind instead of its name (review).
- S-0790 flai check finds `board.order` outside the story at close-out: a draft; finalize it. It waits on TH-1225 S-0790 cost of delay inputs: time lost per cycle 2m.

### 2026-11-15T21:05:56Z flai serve
The list changed. No story in ready can start, no backlog story can be promoted, and 1 in-progress lane stands idle.

This host has no reminder channel, so without one nothing reminds you of this thread off the dashboard. Set one in a shell on the host with `flai serve notify set --url <webhook>` or `flai serve notify set --command <command>`.

Each backlog story, in the pull order, and the decision it waits on:

- S-0241 Measure planned stories again with forked task sub-agents once a headless session offers forks: a draft; finalize it. It waits on TH-1165 S-0241 stays a draft: it is blocked until a headless session offers forks.
- S-0752 flai stats reports each timer a workflow defines beside the baseline's figures, for every item and over the window, past items included: held by `after`; the root of its chain is S-0750 flai stats computes lead, cycle, queue, and review time from the baseline workflow's timers, and the metrics read a state's kind, with every figure byte for byte as today (in-progress).
- S-0753 The dashboard's board draws the columns, their kinds, and the moves between them from what flai gives, and holds no list of states of its own: held by `after`; the root of its chain is S-0739 workitem's epic following, board view, task plan, holds, order, and promotion read a state's kind instead of its name (review).
- S-0754 The dashboard's charts draw the workflow's states and each timer it defines from what flai gives, and an operator's own timer gets a chart beside lead and cycle time: held by `after`; the roots of its chain are S-0750 flai stats computes lead, cycle, queue, and review time from the baseline workflow's timers, and the metrics read a state's kind, with every figure byte for byte as today (in-progress) and S-0739 workitem's epic following, board view, task plan, holds, order, and promotion read a state's kind instead of its name (review).
- S-0755 flai reads a project's own workflow file in place of the baseline, the template ships the baseline, flai upgrade keeps the project's file, and the guides say how to change the board: held by `after`; the root of its chain is S-0741 flai check validates the workflow and the items against it, and its item rules read a state's kind instead of its name (review).
- S-0756 In a project made from the template, a state added between in-progress and review with a rule and a timer of its own, in the workflow file alone, is followed by the board, the moves, the check, the pull, the acceptance, and the charts: held by `after`; the roots of its chain are S-0739 workitem's epic following, board view, task plan, holds, order, and promotion read a state's kind instead of its name (review), S-0741 flai check validates the workflow and the items against it, and its item rules read a state's kind instead of its name (review), S-0742 flai serve starts, restarts, and stops agents, plans again, and finds its reasons by a state's kind and a transition's effect instead of a state's name (review), S-0747 The host API's writes, kluster, components, and people methods act on a state's kind and who may make a transition instead of a state's name (review), S-0748 A role file's to-review and to-backlog conditions become one condition that names any state the workflow has, and the baseline role files use it (in-progress) and S-0750 flai stats computes lead, cycle, queue, and review time from the baseline workflow's timers, and the metrics read a state's kind, with every figure byte for byte as today (in-progress).
- S-0758 The dashboard's item page, inbox, review, activity, deployments, task plan, and cards read a state's kind from flai instead of comparing its name: held by `after`; the root of its chain is S-0739 workitem's epic following, board view, task plan, holds, order, and promotion read a state's kind instead of its name (review).
- S-0790 flai check finds `board.order` outside the story at close-out: a draft; finalize it.

### 2026-11-15T21:11:44Z flai serve
Resolved: S-0790 flai check finds `board.order` outside the story at close-out can be promoted from the backlog, so the board is no longer starved.
