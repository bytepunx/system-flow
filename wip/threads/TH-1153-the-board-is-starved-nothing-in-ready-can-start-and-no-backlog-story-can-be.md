---
id: TH-1153
title: "The board is starved: nothing in ready can start and no backlog story can be promoted"
anchor:
  path: wip/kanban/board.md
status: resolved
participants: [flai serve]
created: 2026-11-13T11:45:34Z
updated: 2026-11-13T11:48:34Z
---

# TH-1153 The board is starved: nothing in ready can start and no backlog story can be promoted

On wip/kanban/board.md.

## Entries

### 2026-11-13T11:45:34Z flai serve
No story in ready can start, no backlog story can be promoted, and 1 in-progress lane stands idle. flai serve adds an entry here when the list below changes, and resolves this thread once a story can start.

This host has no reminder channel, so without one nothing reminds you of this thread off the dashboard. Set one in a shell on the host with `flai serve notify set --url <webhook>` or `flai serve notify set --command <command>`.

Each backlog story, in the pull order, and the decision it waits on:

- S-0241 Measure planned stories again with forked task sub-agents once a headless session offers forks: a draft; finalize it.
- S-0730 A story's agent can overwrite a work item file whole with Write, front matter included, and nothing refuses it: a draft; finalize it.
- S-0736 An ADR fixes the workflow file's format and place under design/workflows, its state kinds, transition effects, rules, and timers, and how a workflow changes with items on the board: a draft; finalize it. It waits on TH-1113 S-0736 cannot be finalized yet: no forecast delivery and no cost of delay value until E-0034 has inputs.
- S-0737 flai reads the baseline workflow it carries, and the board's states and transition map come from it with nothing else changed: a draft; finalize it. It waits on TH-1114 S-0737 cannot be finalized yet: no forecast delivery and no cost of delay value until E-0034 has inputs.
- S-0738 The rules Repo.Move enforces are named rules from a vocabulary flai implements, listed on the baseline's transitions, each refusing or warning: a draft; finalize it. It waits on TH-1115 S-0738 cannot be finalized yet: no forecast delivery and no cost of delay value until E-0034 has inputs.
- S-0739 workitem's epic following, board view, task plan, holds, order, and promotion read a state's kind instead of its name: a draft; finalize it. It waits on TH-1116 S-0739 cannot be finalized yet: no forecast delivery and no cost of delay value until E-0034 has inputs.
- S-0741 flai check validates the workflow and the items against it, and its item rules read a state's kind instead of its name: a draft; finalize it. It waits on TH-1117 S-0741 cannot be finalized yet: no forecast delivery and no cost of delay value until E-0034 has inputs.
- S-0742 flai serve starts, restarts, and stops agents, plans again, and finds its reasons by a state's kind and a transition's effect instead of a state's name: a draft; finalize it. It waits on TH-1118 S-0742 cannot be finalized yet: no forecast delivery and no cost of delay value until E-0034 has inputs.
- S-0743 flai serve's kluster deploys, forge pass, and merge queue act on a state's kind and a transition's effect instead of review and done: a draft; finalize it. It waits on TH-1119 S-0743 cannot be finalized yet: no forecast delivery and no cost of delay value until E-0034 has inputs.
- S-0744 flai move, flai accept, the formula commands, and the acceptance preview act on a transition's effect and who may make it instead of a state's name: a draft; finalize it. It waits on TH-1120 S-0744 cannot be finalized yet: no forecast delivery and no cost of delay value until E-0034 has inputs.
- S-0746 The MCP tools, story start, task done, the sync, item edits, and the inbox act on a state's kind and who may make a transition instead of a state's name: a draft; finalize it. It waits on TH-1121 S-0746 cannot be finalized yet: no forecast delivery and no cost of delay value until E-0034 has inputs.
- S-0747 The host API's writes, kluster, components, and people methods act on a state's kind and who may make a transition instead of a state's name: a draft; finalize it. It waits on TH-1122 S-0747 cannot be finalized yet: no forecast delivery and no cost of delay value until E-0034 has inputs.
- S-0748 A role file's to-review and to-backlog conditions become one condition that names any state the workflow has, and the baseline role files use it: a draft; finalize it. It waits on TH-1123 S-0748 cannot be finalized yet: no forecast delivery and no cost of delay value until E-0034 has inputs.
- S-0749 The guard, messages between stories, forecasts, prime packs, the release evaluation, the digest, and turn counts read a state's kind instead of its name: a draft; finalize it. It waits on TH-1124 S-0749 cannot be finalized yet: no forecast delivery and no cost of delay value until E-0034 has inputs.
- S-0750 flai stats computes lead, cycle, queue, and review time from the baseline workflow's timers, and the metrics read a state's kind, with every figure byte for byte as today: a draft; finalize it. It waits on TH-1125 S-0750 cannot be finalized yet: no forecast delivery and no cost of delay value until E-0034 has inputs.
- S-0752 flai stats reports each timer a workflow defines beside the baseline's figures, for every item and over the window, past items included: a draft; finalize it. It waits on TH-1126 S-0752 cannot be finalized yet: no forecast delivery and no cost of delay value until E-0034 has inputs.
- S-0753 The dashboard's board draws the columns, their kinds, and the moves between them from what flai gives, and holds no list of states of its own: a draft; finalize it. It waits on TH-1127 S-0753 cannot be finalized yet: no forecast delivery and no cost of delay value until E-0034 has inputs.
- S-0754 The dashboard's charts draw the workflow's states and each timer it defines from what flai gives, and an operator's own timer gets a chart beside lead and cycle time: a draft; finalize it. It waits on TH-1128 S-0754 cannot be finalized yet: no forecast delivery and no cost of delay value until E-0034 has inputs.
- S-0755 flai reads a project's own workflow file in place of the baseline, the template ships the baseline, flai upgrade keeps the project's file, and the guides say how to change the board: a draft; finalize it. It waits on TH-1129 S-0755 cannot be finalized yet: no forecast delivery and no cost of delay value until E-0034 has inputs.
- S-0756 In a project made from the template, a state added between in-progress and review with a rule and a timer of its own, in the workflow file alone, is followed by the board, the moves, the check, the pull, the acceptance, and the charts: a draft; finalize it. It waits on TH-1130 S-0756 cannot be finalized yet: no forecast delivery and no cost of delay value until E-0034 has inputs.
- S-0758 The dashboard's item page, inbox, review, activity, deployments, task plan, and cards read a state's kind from flai instead of comparing its name: a draft; finalize it. It waits on TH-1131 S-0758 cannot be finalized yet: no forecast delivery and no cost of delay value until E-0034 has inputs.
- S-0767 flai stream sync reports a rebased branch undoes main's line when both changed one long markdown line and the branch keeps every word of main's: a draft; finalize it.

### 2026-11-13T11:48:34Z flai serve
Resolved: S-0767 flai stream sync reports a rebased branch undoes main's line when both changed one long markdown line and the branch keeps every word of main's can be promoted from the backlog, so the board is no longer starved.
