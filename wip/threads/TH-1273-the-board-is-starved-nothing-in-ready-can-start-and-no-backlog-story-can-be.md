---
id: TH-1273
title: "The board is starved: nothing in ready can start and no backlog story can be promoted"
anchor:
  path: wip/kanban/board.md
status: open
participants: [flai serve]
created: 2026-11-17T08:14:03Z
updated: 2026-11-17T08:14:03Z
---

# TH-1273 The board is starved: nothing in ready can start and no backlog story can be promoted

On wip/kanban/board.md.

## Entries

### 2026-11-17T08:14:03Z flai serve
No story in ready can start, no backlog story can be promoted, and 1 in-progress lane stands idle. flai serve adds an entry here when the list below changes, and resolves this thread once a story can start.

This host has no reminder channel, so without one nothing reminds you of this thread off the dashboard. Set one in a shell on the host with `flai serve notify set --url <webhook>` or `flai serve notify set --command <command>`.

Each backlog story, in the pull order, and the decision it waits on:

- S-0241 Measure planned stories again with forked task sub-agents once a headless session offers forks: a draft; finalize it. It waits on TH-1165 S-0241 stays a draft: it is blocked until a headless session offers forks.
- S-0809 The merge queue's batch runs the suites verify_at.batch lists over the merge of its stories: held by `after`; the root of its chain is S-0805 A tier in the manifest belongs to a suite, flai test --suite runs one, and the template's tiers are labelled static, behavior, integration, and smoke (review).
- S-0810 flai serve runs a state's suites and the cluster's deploy and checks as a job when a story enters it, skips the cluster's steps where none is configured, and sends the story back when a step fails: held by `after`; the root of its chain is S-0805 A tier in the manifest belongs to a suite, flai test --suite runs one, and the template's tiers are labelled static, behavior, integration, and smoke (review).
- S-0811 The release's publish runs the suites verify_at.publish lists at the main branch's head before it tags, and publishes nothing when one fails: held by `after`; the root of its chain is S-0805 A tier in the manifest belongs to a suite, flai test --suite runs one, and the template's tiers are labelled static, behavior, integration, and smoke (review).
- S-0812 flai test --at ci runs the suites verify_at.ci lists, and the template's and this repository's CI call it: held by `after`; the root of its chain is S-0805 A tier in the manifest belongs to a suite, flai test --suite runs one, and the template's tiers are labelled static, behavior, integration, and smoke (review).
- S-0813 A tier's queue and when are read as the suites of the batch and the move to review, flai upgrade writes them into a project's file, and the two keys are retired: held by `after`; the root of its chain is S-0805 A tier in the manifest belongs to a suite, flai test --suite runs one, and the template's tiers are labelled static, behavior, integration, and smoke (review).
- S-0814 The review page and the board show the suites each step needs and ran, and follow a state's job in the activity stream: held by `after`; the root of its chain is S-0805 A tier in the manifest belongs to a suite, flai test --suite runs one, and the template's tiers are labelled static, behavior, integration, and smoke (review).
- S-0815 In a project made from the template, review needs only the fast suites and an integration state between review and done deploys where a cluster is configured and runs the integration and smoke suites: held by `after`; the roots of its chain are S-0805 A tier in the manifest belongs to a suite, flai test --suite runs one, and the template's tiers are labelled static, behavior, integration, and smoke (review), S-0750 flai stats computes lead, cycle, queue, and review time from the baseline workflow's timers, and the metrics read a state's kind, with every figure byte for byte as today (review) and S-0758 The dashboard's item page, inbox, review, activity, deployments, task plan, and cards read a state's kind from flai instead of comparing its name (in-progress).
- S-0822 A rebase and resolve job that passes records stayed in review without reading the story, and a move out of review neither stops the job nor is noticed by it: a draft; finalize it.
- S-0823 A story the operator moves from review to in-progress by hand gets no agent and no thread, so it waits in in-progress with nothing to bring it back: a draft; finalize it.
- S-0824 A rebase and resolve job runs the whole close-out again each time main moves under it, so one took three rounds and 33 minutes while other stories were accepted: a draft; finalize it.
