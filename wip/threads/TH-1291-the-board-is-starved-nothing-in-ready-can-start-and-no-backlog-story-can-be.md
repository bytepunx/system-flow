---
id: TH-1291
title: "The board is starved: nothing in ready can start and no backlog story can be promoted"
anchor:
  path: wip/kanban/board.md
status: resolved
participants: [flai serve]
created: 2026-11-17T12:23:51Z
updated: 2026-11-17T12:35:32Z
---

# TH-1291 The board is starved: nothing in ready can start and no backlog story can be promoted

On wip/kanban/board.md.

## Entries

### 2026-11-17T12:23:51Z flai serve
No story in ready can start, no backlog story can be promoted, and 1 in-progress lane stands idle. flai serve adds an entry here when the list below changes, and resolves this thread once a story can start.

This host has no reminder channel, so without one nothing reminds you of this thread off the dashboard. Set one in a shell on the host with `flai serve notify set --url <webhook>` or `flai serve notify set --command <command>`.

Each backlog story, in the pull order, and the decision it waits on:

- S-0241 Measure planned stories again with forked task sub-agents once a headless session offers forks: a draft; finalize it. It waits on TH-1165 S-0241 stays a draft: it is blocked until a headless session offers forks.
- S-0813 A tier's queue and when are read as the suites of the batch and the move to review, flai upgrade writes them into a project's file, and the two keys are retired: held by `after`; the root of its chain is S-0806 The workflow file takes verify on a transition, runs on a state, and verify_at for the batch, the publish, and CI, and the baseline and the template's file give today's behaviour (in-progress).
- S-0814 The review page and the board show the suites each step needs and ran, and follow a state's job in the activity stream: held by `after`; the root of its chain is S-0806 The workflow file takes verify on a transition, runs on a state, and verify_at for the batch, the publish, and CI, and the baseline and the template's file give today's behaviour (in-progress).
- S-0815 In a project made from the template, review needs only the fast suites and an integration state between review and done deploys where a cluster is configured and runs the integration and smoke suites: held by `after`; the root of its chain is S-0806 The workflow file takes verify on a transition, runs on a state, and verify_at for the batch, the publish, and CI, and the baseline and the template's file give today's behaviour (in-progress).
- S-0822 A rebase and resolve job that passes records stayed in review without reading the story, and a move out of review neither stops the job nor is noticed by it: a draft; finalize it. It waits on TH-1274 Should a move out of review stop a running rebase and resolve job? and TH-1275 S-0822 stays a draft: no touches, forecast, or cost of delay value yet.
- S-0824 A rebase and resolve job runs the whole close-out again each time main moves under it, so one took three rounds and 33 minutes while other stories were accepted: a draft; finalize it. It waits on TH-1287 Plan for S-0824: a resolve job's last check passes over a main that sends back no tier, six tasks in four layers.
- S-0829 A sub-agent's flai kluster status falls through to the guard's generic refusal, against kluster.md's rule that status is a read every agent may run: a draft; finalize it.
- S-0835 flai serve's review pass deploys again a story in review whose deploy passed at its head and only its checks have not, where ADR-0292 §3 runs the checks alone: held by `after`; the root of its chain is S-0806 The workflow file takes verify on a transition, runs on a state, and verify_at for the batch, the publish, and CI, and the baseline and the template's file give today's behaviour (in-progress).
- S-0837 A publish the orchestrator starts is a release_publish call and makes no story, and the baseline publish formula is retired: a draft; finalize it.
- S-0838 flaiover still compares state names in the epic filters and the move route, so an operator's own accepted or cancelled lane is not seen as closed there: a draft; finalize it.
- S-0839 No flai command writes agent.image.args, so a project's layer build arguments can only be written by hand: a draft; finalize it.
- S-0840 flai kluster up cannot refresh a stale kubeconfig: kluster up refuses an existing cluster before flai writes it: a draft; finalize it.
- S-0841 scripts/flai-test.sh's repository check fails on this host while story/S-0001, S-0002, S-0003, and S-0005 are never merged: a draft; finalize it.
- S-0842 workitem's Item.Closed and PlanArchive read the baseline's done and cancelled names, so on a board that renames them flai accept moves a story to the accepted state and then cannot archive it: a draft; finalize it.
- S-0843 Callers of the workflow-aware guard, forecast replay, and share hold still pass the baseline, so a project's own workflow file will not reach them: a draft; finalize it.
- S-0844 ADR-0301 rules for a project's own workflow file that no code enforces yet: operator-only acceptance of a change to it, the dashboard's refusal to save it on main, board.wip-state, and WIP limits on a state the project adds: a draft; finalize it.
- S-0845 A story in review whose close-out fails after main moved waits for the operator to send it back, while its after chain holds the backlog and the board starves: a draft; finalize it.
- S-0846 flai serve's auto-upgrade restart interrupts the close-out jobs it was running, which start again from scratch: a draft; finalize it.
- S-0847 flai serve opens a new board-starvation thread at each flip of the board, 48 in five days, each a thread to the operator: a draft; finalize it.

### 2026-11-17T12:27:21Z flai serve
The list changed. No story in ready can start, no backlog story can be promoted, and 1 in-progress lane stands idle.

This host has no reminder channel, so without one nothing reminds you of this thread off the dashboard. Set one in a shell on the host with `flai serve notify set --url <webhook>` or `flai serve notify set --command <command>`.

Each backlog story, in the pull order, and the decision it waits on:

- S-0241 Measure planned stories again with forked task sub-agents once a headless session offers forks: a draft; finalize it. It waits on TH-1165 S-0241 stays a draft: it is blocked until a headless session offers forks.
- S-0813 A tier's queue and when are read as the suites of the batch and the move to review, flai upgrade writes them into a project's file, and the two keys are retired: held by `after`; the root of its chain is S-0806 The workflow file takes verify on a transition, runs on a state, and verify_at for the batch, the publish, and CI, and the baseline and the template's file give today's behaviour (in-progress).
- S-0814 The review page and the board show the suites each step needs and ran, and follow a state's job in the activity stream: held by `after`; the root of its chain is S-0806 The workflow file takes verify on a transition, runs on a state, and verify_at for the batch, the publish, and CI, and the baseline and the template's file give today's behaviour (in-progress).
- S-0815 In a project made from the template, review needs only the fast suites and an integration state between review and done deploys where a cluster is configured and runs the integration and smoke suites: held by `after`; the root of its chain is S-0806 The workflow file takes verify on a transition, runs on a state, and verify_at for the batch, the publish, and CI, and the baseline and the template's file give today's behaviour (in-progress).
- S-0822 A rebase and resolve job that passes records stayed in review without reading the story, and a move out of review neither stops the job nor is noticed by it: a draft; finalize it. It waits on TH-1274 Should a move out of review stop a running rebase and resolve job? and TH-1275 S-0822 stays a draft: no touches, forecast, or cost of delay value yet.
- S-0824 A rebase and resolve job runs the whole close-out again each time main moves under it, so one took three rounds and 33 minutes while other stories were accepted: a draft; finalize it. It waits on TH-1287 Plan for S-0824: a resolve job's last check passes over a main that sends back no tier, six tasks in four layers.
- S-0829 A sub-agent's flai kluster status falls through to the guard's generic refusal, against kluster.md's rule that status is a read every agent may run: a draft; finalize it.
- S-0835 flai serve's review pass deploys again a story in review whose deploy passed at its head and only its checks have not, where ADR-0292 §3 runs the checks alone: held by `after`; the root of its chain is S-0806 The workflow file takes verify on a transition, runs on a state, and verify_at for the batch, the publish, and CI, and the baseline and the template's file give today's behaviour (in-progress).
- S-0837 A publish the orchestrator starts is a release_publish call and makes no story, and the baseline publish formula is retired: a draft; finalize it. It waits on TH-1292 S-0837: cost of delay inputs — recommend penalty_per_week 90 USD.
- S-0838 flaiover still compares state names in the epic filters and the move route, so an operator's own accepted or cancelled lane is not seen as closed there: a draft; finalize it.
- S-0839 No flai command writes agent.image.args, so a project's layer build arguments can only be written by hand: a draft; finalize it.
- S-0840 flai kluster up cannot refresh a stale kubeconfig: kluster up refuses an existing cluster before flai writes it: a draft; finalize it.
- S-0841 scripts/flai-test.sh's repository check fails on this host while story/S-0001, S-0002, S-0003, and S-0005 are never merged: a draft; finalize it.
- S-0842 workitem's Item.Closed and PlanArchive read the baseline's done and cancelled names, so on a board that renames them flai accept moves a story to the accepted state and then cannot archive it: a draft; finalize it.
- S-0843 Callers of the workflow-aware guard, forecast replay, and share hold still pass the baseline, so a project's own workflow file will not reach them: a draft; finalize it.
- S-0844 ADR-0301 rules for a project's own workflow file that no code enforces yet: operator-only acceptance of a change to it, the dashboard's refusal to save it on main, board.wip-state, and WIP limits on a state the project adds: a draft; finalize it.
- S-0845 A story in review whose close-out fails after main moved waits for the operator to send it back, while its after chain holds the backlog and the board starves: a draft; finalize it.
- S-0846 flai serve's auto-upgrade restart interrupts the close-out jobs it was running, which start again from scratch: a draft; finalize it.
- S-0847 flai serve opens a new board-starvation thread at each flip of the board, 48 in five days, each a thread to the operator: a draft; finalize it.

### 2026-11-17T12:28:57Z flai serve
The list changed. No story in ready can start, no backlog story can be promoted, and 1 in-progress lane stands idle.

This host has no reminder channel, so without one nothing reminds you of this thread off the dashboard. Set one in a shell on the host with `flai serve notify set --url <webhook>` or `flai serve notify set --command <command>`.

Each backlog story, in the pull order, and the decision it waits on:

- S-0241 Measure planned stories again with forked task sub-agents once a headless session offers forks: a draft; finalize it. It waits on TH-1165 S-0241 stays a draft: it is blocked until a headless session offers forks.
- S-0813 A tier's queue and when are read as the suites of the batch and the move to review, flai upgrade writes them into a project's file, and the two keys are retired: held by `after`; the root of its chain is S-0806 The workflow file takes verify on a transition, runs on a state, and verify_at for the batch, the publish, and CI, and the baseline and the template's file give today's behaviour (in-progress).
- S-0814 The review page and the board show the suites each step needs and ran, and follow a state's job in the activity stream: held by `after`; the root of its chain is S-0806 The workflow file takes verify on a transition, runs on a state, and verify_at for the batch, the publish, and CI, and the baseline and the template's file give today's behaviour (in-progress).
- S-0815 In a project made from the template, review needs only the fast suites and an integration state between review and done deploys where a cluster is configured and runs the integration and smoke suites: held by `after`; the root of its chain is S-0806 The workflow file takes verify on a transition, runs on a state, and verify_at for the batch, the publish, and CI, and the baseline and the template's file give today's behaviour (in-progress).
- S-0822 A rebase and resolve job that passes records stayed in review without reading the story, and a move out of review neither stops the job nor is noticed by it: a draft; finalize it. It waits on TH-1274 Should a move out of review stop a running rebase and resolve job? and TH-1275 S-0822 stays a draft: no touches, forecast, or cost of delay value yet.
- S-0824 A rebase and resolve job runs the whole close-out again each time main moves under it, so one took three rounds and 33 minutes while other stories were accepted: a draft; finalize it. It waits on TH-1287 Plan for S-0824: a resolve job's last check passes over a main that sends back no tier, six tasks in four layers.
- S-0829 A sub-agent's flai kluster status falls through to the guard's generic refusal, against kluster.md's rule that status is a read every agent may run: a draft; finalize it.
- S-0835 flai serve's review pass deploys again a story in review whose deploy passed at its head and only its checks have not, where ADR-0292 §3 runs the checks alone: held by `after`; the root of its chain is S-0806 The workflow file takes verify on a transition, runs on a state, and verify_at for the batch, the publish, and CI, and the baseline and the template's file give today's behaviour (in-progress).
- S-0837 A publish the orchestrator starts is a release_publish call and makes no story, and the baseline publish formula is retired: a draft; finalize it.
- S-0838 flaiover still compares state names in the epic filters and the move route, so an operator's own accepted or cancelled lane is not seen as closed there: a draft; finalize it.
- S-0839 No flai command writes agent.image.args, so a project's layer build arguments can only be written by hand: a draft; finalize it.
- S-0840 flai kluster up cannot refresh a stale kubeconfig: kluster up refuses an existing cluster before flai writes it: a draft; finalize it.
- S-0841 scripts/flai-test.sh's repository check fails on this host while story/S-0001, S-0002, S-0003, and S-0005 are never merged: a draft; finalize it.
- S-0842 workitem's Item.Closed and PlanArchive read the baseline's done and cancelled names, so on a board that renames them flai accept moves a story to the accepted state and then cannot archive it: a draft; finalize it.
- S-0843 Callers of the workflow-aware guard, forecast replay, and share hold still pass the baseline, so a project's own workflow file will not reach them: a draft; finalize it.
- S-0844 ADR-0301 rules for a project's own workflow file that no code enforces yet: operator-only acceptance of a change to it, the dashboard's refusal to save it on main, board.wip-state, and WIP limits on a state the project adds: a draft; finalize it.
- S-0845 A story in review whose close-out fails after main moved waits for the operator to send it back, while its after chain holds the backlog and the board starves: a draft; finalize it.
- S-0846 flai serve's auto-upgrade restart interrupts the close-out jobs it was running, which start again from scratch: a draft; finalize it.
- S-0847 flai serve opens a new board-starvation thread at each flip of the board, 48 in five days, each a thread to the operator: a draft; finalize it.

### 2026-11-17T12:29:15Z flai serve
The list changed. No story in ready can start, no backlog story can be promoted, and 1 in-progress lane stands idle.

This host has no reminder channel, so without one nothing reminds you of this thread off the dashboard. Set one in a shell on the host with `flai serve notify set --url <webhook>` or `flai serve notify set --command <command>`.

Each backlog story, in the pull order, and the decision it waits on:

- S-0241 Measure planned stories again with forked task sub-agents once a headless session offers forks: a draft; finalize it. It waits on TH-1165 S-0241 stays a draft: it is blocked until a headless session offers forks.
- S-0813 A tier's queue and when are read as the suites of the batch and the move to review, flai upgrade writes them into a project's file, and the two keys are retired: held by `after`; the root of its chain is S-0806 The workflow file takes verify on a transition, runs on a state, and verify_at for the batch, the publish, and CI, and the baseline and the template's file give today's behaviour (in-progress).
- S-0814 The review page and the board show the suites each step needs and ran, and follow a state's job in the activity stream: held by `after`; the root of its chain is S-0806 The workflow file takes verify on a transition, runs on a state, and verify_at for the batch, the publish, and CI, and the baseline and the template's file give today's behaviour (in-progress).
- S-0815 In a project made from the template, review needs only the fast suites and an integration state between review and done deploys where a cluster is configured and runs the integration and smoke suites: held by `after`; the root of its chain is S-0806 The workflow file takes verify on a transition, runs on a state, and verify_at for the batch, the publish, and CI, and the baseline and the template's file give today's behaviour (in-progress).
- S-0822 A rebase and resolve job that passes records stayed in review without reading the story, and a move out of review neither stops the job nor is noticed by it: a draft; finalize it. It waits on TH-1274 Should a move out of review stop a running rebase and resolve job? and TH-1275 S-0822 stays a draft: no touches, forecast, or cost of delay value yet.
- S-0824 A rebase and resolve job runs the whole close-out again each time main moves under it, so one took three rounds and 33 minutes while other stories were accepted: a draft; finalize it. It waits on TH-1287 Plan for S-0824: a resolve job's last check passes over a main that sends back no tier, six tasks in four layers.
- S-0829 A sub-agent's flai kluster status falls through to the guard's generic refusal, against kluster.md's rule that status is a read every agent may run: a draft; finalize it. It waits on TH-1293 S-0829 plan: add kluster status and deployments to the guard's sub-agent reads, in three tasks over two layers.
- S-0835 flai serve's review pass deploys again a story in review whose deploy passed at its head and only its checks have not, where ADR-0292 §3 runs the checks alone: held by `after`; the root of its chain is S-0806 The workflow file takes verify on a transition, runs on a state, and verify_at for the batch, the publish, and CI, and the baseline and the template's file give today's behaviour (in-progress).
- S-0837 A publish the orchestrator starts is a release_publish call and makes no story, and the baseline publish formula is retired: a draft; finalize it.
- S-0838 flaiover still compares state names in the epic filters and the move route, so an operator's own accepted or cancelled lane is not seen as closed there: a draft; finalize it.
- S-0839 No flai command writes agent.image.args, so a project's layer build arguments can only be written by hand: a draft; finalize it.
- S-0840 flai kluster up cannot refresh a stale kubeconfig: kluster up refuses an existing cluster before flai writes it: a draft; finalize it.
- S-0841 scripts/flai-test.sh's repository check fails on this host while story/S-0001, S-0002, S-0003, and S-0005 are never merged: a draft; finalize it.
- S-0842 workitem's Item.Closed and PlanArchive read the baseline's done and cancelled names, so on a board that renames them flai accept moves a story to the accepted state and then cannot archive it: a draft; finalize it.
- S-0843 Callers of the workflow-aware guard, forecast replay, and share hold still pass the baseline, so a project's own workflow file will not reach them: a draft; finalize it.
- S-0844 ADR-0301 rules for a project's own workflow file that no code enforces yet: operator-only acceptance of a change to it, the dashboard's refusal to save it on main, board.wip-state, and WIP limits on a state the project adds: a draft; finalize it.
- S-0845 A story in review whose close-out fails after main moved waits for the operator to send it back, while its after chain holds the backlog and the board starves: a draft; finalize it.
- S-0846 flai serve's auto-upgrade restart interrupts the close-out jobs it was running, which start again from scratch: a draft; finalize it.
- S-0847 flai serve opens a new board-starvation thread at each flip of the board, 48 in five days, each a thread to the operator: a draft; finalize it.

### 2026-11-17T12:29:47Z flai serve
The list changed. No story in ready can start, no backlog story can be promoted, and 1 in-progress lane stands idle.

This host has no reminder channel, so without one nothing reminds you of this thread off the dashboard. Set one in a shell on the host with `flai serve notify set --url <webhook>` or `flai serve notify set --command <command>`.

Each backlog story, in the pull order, and the decision it waits on:

- S-0241 Measure planned stories again with forked task sub-agents once a headless session offers forks: a draft; finalize it. It waits on TH-1165 S-0241 stays a draft: it is blocked until a headless session offers forks.
- S-0813 A tier's queue and when are read as the suites of the batch and the move to review, flai upgrade writes them into a project's file, and the two keys are retired: held by `after`; the root of its chain is S-0806 The workflow file takes verify on a transition, runs on a state, and verify_at for the batch, the publish, and CI, and the baseline and the template's file give today's behaviour (in-progress).
- S-0814 The review page and the board show the suites each step needs and ran, and follow a state's job in the activity stream: held by `after`; the root of its chain is S-0806 The workflow file takes verify on a transition, runs on a state, and verify_at for the batch, the publish, and CI, and the baseline and the template's file give today's behaviour (in-progress).
- S-0815 In a project made from the template, review needs only the fast suites and an integration state between review and done deploys where a cluster is configured and runs the integration and smoke suites: held by `after`; the root of its chain is S-0806 The workflow file takes verify on a transition, runs on a state, and verify_at for the batch, the publish, and CI, and the baseline and the template's file give today's behaviour (in-progress).
- S-0822 A rebase and resolve job that passes records stayed in review without reading the story, and a move out of review neither stops the job nor is noticed by it: a draft; finalize it. It waits on TH-1274 Should a move out of review stop a running rebase and resolve job? and TH-1275 S-0822 stays a draft: no touches, forecast, or cost of delay value yet.
- S-0824 A rebase and resolve job runs the whole close-out again each time main moves under it, so one took three rounds and 33 minutes while other stories were accepted: a draft; finalize it. It waits on TH-1287 Plan for S-0824: a resolve job's last check passes over a main that sends back no tier, six tasks in four layers.
- S-0829 A sub-agent's flai kluster status falls through to the guard's generic refusal, against kluster.md's rule that status is a read every agent may run: a draft; finalize it. It waits on TH-1293 S-0829 plan: add kluster status and deployments to the guard's sub-agent reads, in three tasks over two layers.
- S-0835 flai serve's review pass deploys again a story in review whose deploy passed at its head and only its checks have not, where ADR-0292 §3 runs the checks alone: held by `after`; the root of its chain is S-0806 The workflow file takes verify on a transition, runs on a state, and verify_at for the batch, the publish, and CI, and the baseline and the template's file give today's behaviour (in-progress).
- S-0837 A publish the orchestrator starts is a release_publish call and makes no story, and the baseline publish formula is retired: a draft; finalize it.
- S-0838 flaiover still compares state names in the epic filters and the move route, so an operator's own accepted or cancelled lane is not seen as closed there: a draft; finalize it.
- S-0839 No flai command writes agent.image.args, so a project's layer build arguments can only be written by hand: a draft; finalize it.
- S-0840 flai kluster up cannot refresh a stale kubeconfig: kluster up refuses an existing cluster before flai writes it: a draft; finalize it. It waits on TH-1294 S-0840 plan: three tasks in three layers to make flai kluster up write a running cluster's kubeconfig again.
- S-0841 scripts/flai-test.sh's repository check fails on this host while story/S-0001, S-0002, S-0003, and S-0005 are never merged: a draft; finalize it.
- S-0842 workitem's Item.Closed and PlanArchive read the baseline's done and cancelled names, so on a board that renames them flai accept moves a story to the accepted state and then cannot archive it: a draft; finalize it.
- S-0843 Callers of the workflow-aware guard, forecast replay, and share hold still pass the baseline, so a project's own workflow file will not reach them: a draft; finalize it.
- S-0844 ADR-0301 rules for a project's own workflow file that no code enforces yet: operator-only acceptance of a change to it, the dashboard's refusal to save it on main, board.wip-state, and WIP limits on a state the project adds: a draft; finalize it.
- S-0845 A story in review whose close-out fails after main moved waits for the operator to send it back, while its after chain holds the backlog and the board starves: a draft; finalize it.
- S-0846 flai serve's auto-upgrade restart interrupts the close-out jobs it was running, which start again from scratch: a draft; finalize it.
- S-0847 flai serve opens a new board-starvation thread at each flip of the board, 48 in five days, each a thread to the operator: a draft; finalize it.

### 2026-11-17T12:30:03Z flai serve
The list changed. No story in ready can start, no backlog story can be promoted, and 1 in-progress lane stands idle.

This host has no reminder channel, so without one nothing reminds you of this thread off the dashboard. Set one in a shell on the host with `flai serve notify set --url <webhook>` or `flai serve notify set --command <command>`.

Each backlog story, in the pull order, and the decision it waits on:

- S-0241 Measure planned stories again with forked task sub-agents once a headless session offers forks: a draft; finalize it. It waits on TH-1165 S-0241 stays a draft: it is blocked until a headless session offers forks.
- S-0813 A tier's queue and when are read as the suites of the batch and the move to review, flai upgrade writes them into a project's file, and the two keys are retired: held by `after`; the root of its chain is S-0806 The workflow file takes verify on a transition, runs on a state, and verify_at for the batch, the publish, and CI, and the baseline and the template's file give today's behaviour (in-progress).
- S-0814 The review page and the board show the suites each step needs and ran, and follow a state's job in the activity stream: held by `after`; the root of its chain is S-0806 The workflow file takes verify on a transition, runs on a state, and verify_at for the batch, the publish, and CI, and the baseline and the template's file give today's behaviour (in-progress).
- S-0815 In a project made from the template, review needs only the fast suites and an integration state between review and done deploys where a cluster is configured and runs the integration and smoke suites: held by `after`; the root of its chain is S-0806 The workflow file takes verify on a transition, runs on a state, and verify_at for the batch, the publish, and CI, and the baseline and the template's file give today's behaviour (in-progress).
- S-0822 A rebase and resolve job that passes records stayed in review without reading the story, and a move out of review neither stops the job nor is noticed by it: a draft; finalize it. It waits on TH-1274 Should a move out of review stop a running rebase and resolve job? and TH-1275 S-0822 stays a draft: no touches, forecast, or cost of delay value yet.
- S-0824 A rebase and resolve job runs the whole close-out again each time main moves under it, so one took three rounds and 33 minutes while other stories were accepted: a draft; finalize it. It waits on TH-1287 Plan for S-0824: a resolve job's last check passes over a main that sends back no tier, six tasks in four layers.
- S-0829 A sub-agent's flai kluster status falls through to the guard's generic refusal, against kluster.md's rule that status is a read every agent may run: a draft; finalize it. It waits on TH-1293 S-0829 plan: add kluster status and deployments to the guard's sub-agent reads, in three tasks over two layers.
- S-0835 flai serve's review pass deploys again a story in review whose deploy passed at its head and only its checks have not, where ADR-0292 §3 runs the checks alone: held by `after`; the root of its chain is S-0806 The workflow file takes verify on a transition, runs on a state, and verify_at for the batch, the publish, and CI, and the baseline and the template's file give today's behaviour (in-progress).
- S-0837 A publish the orchestrator starts is a release_publish call and makes no story, and the baseline publish formula is retired: a draft; finalize it.
- S-0838 flaiover still compares state names in the epic filters and the move route, so an operator's own accepted or cancelled lane is not seen as closed there: a draft; finalize it.
- S-0839 No flai command writes agent.image.args, so a project's layer build arguments can only be written by hand: a draft; finalize it.
- S-0840 flai kluster up cannot refresh a stale kubeconfig: kluster up refuses an existing cluster before flai writes it: a draft; finalize it. It waits on TH-1294 S-0840 plan: three tasks in three layers to make flai kluster up write a running cluster's kubeconfig again.
- S-0841 scripts/flai-test.sh's repository check fails on this host while story/S-0001, S-0002, S-0003, and S-0005 are never merged: a draft; finalize it.
- S-0842 workitem's Item.Closed and PlanArchive read the baseline's done and cancelled names, so on a board that renames them flai accept moves a story to the accepted state and then cannot archive it: a draft; finalize it.
- S-0843 Callers of the workflow-aware guard, forecast replay, and share hold still pass the baseline, so a project's own workflow file will not reach them: a draft; finalize it. It waits on TH-1295 Plan for S-0843: four tasks in three layers pass the project's workflow to the baseline callers left.
- S-0844 ADR-0301 rules for a project's own workflow file that no code enforces yet: operator-only acceptance of a change to it, the dashboard's refusal to save it on main, board.wip-state, and WIP limits on a state the project adds: a draft; finalize it.
- S-0845 A story in review whose close-out fails after main moved waits for the operator to send it back, while its after chain holds the backlog and the board starves: a draft; finalize it.
- S-0846 flai serve's auto-upgrade restart interrupts the close-out jobs it was running, which start again from scratch: a draft; finalize it.
- S-0847 flai serve opens a new board-starvation thread at each flip of the board, 48 in five days, each a thread to the operator: a draft; finalize it.

### 2026-11-17T12:32:43Z flai serve
The list changed. No story in ready can start, no backlog story can be promoted, and 1 in-progress lane stands idle.

This host has no reminder channel, so without one nothing reminds you of this thread off the dashboard. Set one in a shell on the host with `flai serve notify set --url <webhook>` or `flai serve notify set --command <command>`.

Each backlog story, in the pull order, and the decision it waits on:

- S-0241 Measure planned stories again with forked task sub-agents once a headless session offers forks: a draft; finalize it. It waits on TH-1165 S-0241 stays a draft: it is blocked until a headless session offers forks.
- S-0813 A tier's queue and when are read as the suites of the batch and the move to review, flai upgrade writes them into a project's file, and the two keys are retired: held by `after`; the root of its chain is S-0806 The workflow file takes verify on a transition, runs on a state, and verify_at for the batch, the publish, and CI, and the baseline and the template's file give today's behaviour (in-progress).
- S-0814 The review page and the board show the suites each step needs and ran, and follow a state's job in the activity stream: held by `after`; the root of its chain is S-0806 The workflow file takes verify on a transition, runs on a state, and verify_at for the batch, the publish, and CI, and the baseline and the template's file give today's behaviour (in-progress).
- S-0815 In a project made from the template, review needs only the fast suites and an integration state between review and done deploys where a cluster is configured and runs the integration and smoke suites: held by `after`; the root of its chain is S-0806 The workflow file takes verify on a transition, runs on a state, and verify_at for the batch, the publish, and CI, and the baseline and the template's file give today's behaviour (in-progress).
- S-0822 A rebase and resolve job that passes records stayed in review without reading the story, and a move out of review neither stops the job nor is noticed by it: a draft; finalize it. It waits on TH-1274 Should a move out of review stop a running rebase and resolve job? and TH-1275 S-0822 stays a draft: no touches, forecast, or cost of delay value yet.
- S-0824 A rebase and resolve job runs the whole close-out again each time main moves under it, so one took three rounds and 33 minutes while other stories were accepted: a draft; finalize it. It waits on TH-1287 Plan for S-0824: a resolve job's last check passes over a main that sends back no tier, six tasks in four layers.
- S-0829 A sub-agent's flai kluster status falls through to the guard's generic refusal, against kluster.md's rule that status is a read every agent may run: a draft; finalize it. It waits on TH-1293 S-0829 plan: add kluster status and deployments to the guard's sub-agent reads, in three tasks over two layers.
- S-0835 flai serve's review pass deploys again a story in review whose deploy passed at its head and only its checks have not, where ADR-0292 §3 runs the checks alone: held by `after`; the root of its chain is S-0806 The workflow file takes verify on a transition, runs on a state, and verify_at for the batch, the publish, and CI, and the baseline and the template's file give today's behaviour (in-progress).
- S-0837 A publish the orchestrator starts is a release_publish call and makes no story, and the baseline publish formula is retired: a draft; finalize it.
- S-0838 flaiover still compares state names in the epic filters and the move route, so an operator's own accepted or cancelled lane is not seen as closed there: a draft; finalize it.
- S-0839 No flai command writes agent.image.args, so a project's layer build arguments can only be written by hand: a draft; finalize it.
- S-0840 flai kluster up cannot refresh a stale kubeconfig: kluster up refuses an existing cluster before flai writes it: a draft; finalize it. It waits on TH-1294 S-0840 plan: three tasks in three layers to make flai kluster up write a running cluster's kubeconfig again.
- S-0841 scripts/flai-test.sh's repository check fails on this host while story/S-0001, S-0002, S-0003, and S-0005 are never merged: a draft; finalize it.
- S-0842 workitem's Item.Closed and PlanArchive read the baseline's done and cancelled names, so on a board that renames them flai accept moves a story to the accepted state and then cannot archive it: a draft; finalize it. It waits on TH-1296 Plan for S-0842: Closed and PlanArchive read the workflow's closed kinds.
- S-0843 Callers of the workflow-aware guard, forecast replay, and share hold still pass the baseline, so a project's own workflow file will not reach them: a draft; finalize it. It waits on TH-1295 Plan for S-0843: four tasks in three layers pass the project's workflow to the baseline callers left.
- S-0844 ADR-0301 rules for a project's own workflow file that no code enforces yet: operator-only acceptance of a change to it, the dashboard's refusal to save it on main, board.wip-state, and WIP limits on a state the project adds: a draft; finalize it.
- S-0845 A story in review whose close-out fails after main moved waits for the operator to send it back, while its after chain holds the backlog and the board starves: a draft; finalize it.
- S-0846 flai serve's auto-upgrade restart interrupts the close-out jobs it was running, which start again from scratch: a draft; finalize it.
- S-0847 flai serve opens a new board-starvation thread at each flip of the board, 48 in five days, each a thread to the operator: a draft; finalize it.

### 2026-11-17T12:32:45Z flai serve
The list changed. No story in ready can start, no backlog story can be promoted, and 1 in-progress lane stands idle.

This host has no reminder channel, so without one nothing reminds you of this thread off the dashboard. Set one in a shell on the host with `flai serve notify set --url <webhook>` or `flai serve notify set --command <command>`.

Each backlog story, in the pull order, and the decision it waits on:

- S-0241 Measure planned stories again with forked task sub-agents once a headless session offers forks: a draft; finalize it. It waits on TH-1165 S-0241 stays a draft: it is blocked until a headless session offers forks.
- S-0813 A tier's queue and when are read as the suites of the batch and the move to review, flai upgrade writes them into a project's file, and the two keys are retired: held by `after`; the root of its chain is S-0806 The workflow file takes verify on a transition, runs on a state, and verify_at for the batch, the publish, and CI, and the baseline and the template's file give today's behaviour (in-progress).
- S-0814 The review page and the board show the suites each step needs and ran, and follow a state's job in the activity stream: held by `after`; the root of its chain is S-0806 The workflow file takes verify on a transition, runs on a state, and verify_at for the batch, the publish, and CI, and the baseline and the template's file give today's behaviour (in-progress).
- S-0815 In a project made from the template, review needs only the fast suites and an integration state between review and done deploys where a cluster is configured and runs the integration and smoke suites: held by `after`; the root of its chain is S-0806 The workflow file takes verify on a transition, runs on a state, and verify_at for the batch, the publish, and CI, and the baseline and the template's file give today's behaviour (in-progress).
- S-0822 A rebase and resolve job that passes records stayed in review without reading the story, and a move out of review neither stops the job nor is noticed by it: a draft; finalize it. It waits on TH-1274 Should a move out of review stop a running rebase and resolve job? and TH-1275 S-0822 stays a draft: no touches, forecast, or cost of delay value yet.
- S-0824 A rebase and resolve job runs the whole close-out again each time main moves under it, so one took three rounds and 33 minutes while other stories were accepted: a draft; finalize it. It waits on TH-1287 Plan for S-0824: a resolve job's last check passes over a main that sends back no tier, six tasks in four layers.
- S-0829 A sub-agent's flai kluster status falls through to the guard's generic refusal, against kluster.md's rule that status is a read every agent may run: a draft; finalize it. It waits on TH-1293 S-0829 plan: add kluster status and deployments to the guard's sub-agent reads, in three tasks over two layers.
- S-0835 flai serve's review pass deploys again a story in review whose deploy passed at its head and only its checks have not, where ADR-0292 §3 runs the checks alone: held by `after`; the root of its chain is S-0806 The workflow file takes verify on a transition, runs on a state, and verify_at for the batch, the publish, and CI, and the baseline and the template's file give today's behaviour (in-progress).
- S-0837 A publish the orchestrator starts is a release_publish call and makes no story, and the baseline publish formula is retired: a draft; finalize it.
- S-0838 flaiover still compares state names in the epic filters and the move route, so an operator's own accepted or cancelled lane is not seen as closed there: a draft; finalize it. It waits on TH-1297 Plan for S-0838: read each lane's kind from the board's flow in the epic filters, the move route, and ProposeChange.
- S-0839 No flai command writes agent.image.args, so a project's layer build arguments can only be written by hand: a draft; finalize it.
- S-0840 flai kluster up cannot refresh a stale kubeconfig: kluster up refuses an existing cluster before flai writes it: a draft; finalize it. It waits on TH-1294 S-0840 plan: three tasks in three layers to make flai kluster up write a running cluster's kubeconfig again.
- S-0841 scripts/flai-test.sh's repository check fails on this host while story/S-0001, S-0002, S-0003, and S-0005 are never merged: a draft; finalize it.
- S-0842 workitem's Item.Closed and PlanArchive read the baseline's done and cancelled names, so on a board that renames them flai accept moves a story to the accepted state and then cannot archive it: a draft; finalize it. It waits on TH-1296 Plan for S-0842: Closed and PlanArchive read the workflow's closed kinds.
- S-0843 Callers of the workflow-aware guard, forecast replay, and share hold still pass the baseline, so a project's own workflow file will not reach them: a draft; finalize it. It waits on TH-1295 Plan for S-0843: four tasks in three layers pass the project's workflow to the baseline callers left.
- S-0844 ADR-0301 rules for a project's own workflow file that no code enforces yet: operator-only acceptance of a change to it, the dashboard's refusal to save it on main, board.wip-state, and WIP limits on a state the project adds: a draft; finalize it.
- S-0845 A story in review whose close-out fails after main moved waits for the operator to send it back, while its after chain holds the backlog and the board starves: a draft; finalize it.
- S-0846 flai serve's auto-upgrade restart interrupts the close-out jobs it was running, which start again from scratch: a draft; finalize it.
- S-0847 flai serve opens a new board-starvation thread at each flip of the board, 48 in five days, each a thread to the operator: a draft; finalize it.

### 2026-11-17T12:32:48Z flai serve
The list changed. No story in ready can start, no backlog story can be promoted, and 1 in-progress lane stands idle.

This host has no reminder channel, so without one nothing reminds you of this thread off the dashboard. Set one in a shell on the host with `flai serve notify set --url <webhook>` or `flai serve notify set --command <command>`.

Each backlog story, in the pull order, and the decision it waits on:

- S-0241 Measure planned stories again with forked task sub-agents once a headless session offers forks: a draft; finalize it. It waits on TH-1165 S-0241 stays a draft: it is blocked until a headless session offers forks.
- S-0813 A tier's queue and when are read as the suites of the batch and the move to review, flai upgrade writes them into a project's file, and the two keys are retired: held by `after`; the root of its chain is S-0806 The workflow file takes verify on a transition, runs on a state, and verify_at for the batch, the publish, and CI, and the baseline and the template's file give today's behaviour (in-progress).
- S-0814 The review page and the board show the suites each step needs and ran, and follow a state's job in the activity stream: held by `after`; the root of its chain is S-0806 The workflow file takes verify on a transition, runs on a state, and verify_at for the batch, the publish, and CI, and the baseline and the template's file give today's behaviour (in-progress).
- S-0815 In a project made from the template, review needs only the fast suites and an integration state between review and done deploys where a cluster is configured and runs the integration and smoke suites: held by `after`; the root of its chain is S-0806 The workflow file takes verify on a transition, runs on a state, and verify_at for the batch, the publish, and CI, and the baseline and the template's file give today's behaviour (in-progress).
- S-0822 A rebase and resolve job that passes records stayed in review without reading the story, and a move out of review neither stops the job nor is noticed by it: a draft; finalize it. It waits on TH-1274 Should a move out of review stop a running rebase and resolve job? and TH-1275 S-0822 stays a draft: no touches, forecast, or cost of delay value yet.
- S-0824 A rebase and resolve job runs the whole close-out again each time main moves under it, so one took three rounds and 33 minutes while other stories were accepted: a draft; finalize it. It waits on TH-1287 Plan for S-0824: a resolve job's last check passes over a main that sends back no tier, six tasks in four layers.
- S-0829 A sub-agent's flai kluster status falls through to the guard's generic refusal, against kluster.md's rule that status is a read every agent may run: a draft; finalize it. It waits on TH-1293 S-0829 plan: add kluster status and deployments to the guard's sub-agent reads, in three tasks over two layers.
- S-0835 flai serve's review pass deploys again a story in review whose deploy passed at its head and only its checks have not, where ADR-0292 §3 runs the checks alone: held by `after`; the root of its chain is S-0806 The workflow file takes verify on a transition, runs on a state, and verify_at for the batch, the publish, and CI, and the baseline and the template's file give today's behaviour (in-progress).
- S-0837 A publish the orchestrator starts is a release_publish call and makes no story, and the baseline publish formula is retired: a draft; finalize it.
- S-0838 flaiover still compares state names in the epic filters and the move route, so an operator's own accepted or cancelled lane is not seen as closed there: a draft; finalize it. It waits on TH-1297 Plan for S-0838: read each lane's kind from the board's flow in the epic filters, the move route, and ProposeChange.
- S-0839 No flai command writes agent.image.args, so a project's layer build arguments can only be written by hand: a draft; finalize it. It waits on TH-1298 S-0839 plan: flai manifest set writes agent.image.layer and agent.image.args.<NAME>.
- S-0840 flai kluster up cannot refresh a stale kubeconfig: kluster up refuses an existing cluster before flai writes it: a draft; finalize it. It waits on TH-1294 S-0840 plan: three tasks in three layers to make flai kluster up write a running cluster's kubeconfig again.
- S-0841 scripts/flai-test.sh's repository check fails on this host while story/S-0001, S-0002, S-0003, and S-0005 are never merged: a draft; finalize it.
- S-0842 workitem's Item.Closed and PlanArchive read the baseline's done and cancelled names, so on a board that renames them flai accept moves a story to the accepted state and then cannot archive it: a draft; finalize it. It waits on TH-1296 Plan for S-0842: Closed and PlanArchive read the workflow's closed kinds.
- S-0843 Callers of the workflow-aware guard, forecast replay, and share hold still pass the baseline, so a project's own workflow file will not reach them: a draft; finalize it. It waits on TH-1295 Plan for S-0843: four tasks in three layers pass the project's workflow to the baseline callers left.
- S-0844 ADR-0301 rules for a project's own workflow file that no code enforces yet: operator-only acceptance of a change to it, the dashboard's refusal to save it on main, board.wip-state, and WIP limits on a state the project adds: a draft; finalize it.
- S-0845 A story in review whose close-out fails after main moved waits for the operator to send it back, while its after chain holds the backlog and the board starves: a draft; finalize it.
- S-0846 flai serve's auto-upgrade restart interrupts the close-out jobs it was running, which start again from scratch: a draft; finalize it.
- S-0847 flai serve opens a new board-starvation thread at each flip of the board, 48 in five days, each a thread to the operator: a draft; finalize it.

### 2026-11-17T12:33:23Z flai serve
The list changed. No story in ready can start, no backlog story can be promoted, and 1 in-progress lane stands idle.

This host has no reminder channel, so without one nothing reminds you of this thread off the dashboard. Set one in a shell on the host with `flai serve notify set --url <webhook>` or `flai serve notify set --command <command>`.

Each backlog story, in the pull order, and the decision it waits on:

- S-0241 Measure planned stories again with forked task sub-agents once a headless session offers forks: a draft; finalize it. It waits on TH-1165 S-0241 stays a draft: it is blocked until a headless session offers forks.
- S-0813 A tier's queue and when are read as the suites of the batch and the move to review, flai upgrade writes them into a project's file, and the two keys are retired: held by `after`; the root of its chain is S-0806 The workflow file takes verify on a transition, runs on a state, and verify_at for the batch, the publish, and CI, and the baseline and the template's file give today's behaviour (in-progress).
- S-0814 The review page and the board show the suites each step needs and ran, and follow a state's job in the activity stream: held by `after`; the root of its chain is S-0806 The workflow file takes verify on a transition, runs on a state, and verify_at for the batch, the publish, and CI, and the baseline and the template's file give today's behaviour (in-progress).
- S-0815 In a project made from the template, review needs only the fast suites and an integration state between review and done deploys where a cluster is configured and runs the integration and smoke suites: held by `after`; the root of its chain is S-0806 The workflow file takes verify on a transition, runs on a state, and verify_at for the batch, the publish, and CI, and the baseline and the template's file give today's behaviour (in-progress).
- S-0822 A rebase and resolve job that passes records stayed in review without reading the story, and a move out of review neither stops the job nor is noticed by it: a draft; finalize it. It waits on TH-1274 Should a move out of review stop a running rebase and resolve job? and TH-1275 S-0822 stays a draft: no touches, forecast, or cost of delay value yet.
- S-0824 A rebase and resolve job runs the whole close-out again each time main moves under it, so one took three rounds and 33 minutes while other stories were accepted: a draft; finalize it. It waits on TH-1287 Plan for S-0824: a resolve job's last check passes over a main that sends back no tier, six tasks in four layers.
- S-0829 A sub-agent's flai kluster status falls through to the guard's generic refusal, against kluster.md's rule that status is a read every agent may run: a draft; finalize it.
- S-0835 flai serve's review pass deploys again a story in review whose deploy passed at its head and only its checks have not, where ADR-0292 §3 runs the checks alone: held by `after`; the root of its chain is S-0806 The workflow file takes verify on a transition, runs on a state, and verify_at for the batch, the publish, and CI, and the baseline and the template's file give today's behaviour (in-progress).
- S-0837 A publish the orchestrator starts is a release_publish call and makes no story, and the baseline publish formula is retired: a draft; finalize it.
- S-0838 flaiover still compares state names in the epic filters and the move route, so an operator's own accepted or cancelled lane is not seen as closed there: a draft; finalize it. It waits on TH-1297 Plan for S-0838: read each lane's kind from the board's flow in the epic filters, the move route, and ProposeChange.
- S-0839 No flai command writes agent.image.args, so a project's layer build arguments can only be written by hand: a draft; finalize it. It waits on TH-1298 S-0839 plan: flai manifest set writes agent.image.layer and agent.image.args.<NAME>.
- S-0840 flai kluster up cannot refresh a stale kubeconfig: kluster up refuses an existing cluster before flai writes it: a draft; finalize it. It waits on TH-1294 S-0840 plan: three tasks in three layers to make flai kluster up write a running cluster's kubeconfig again.
- S-0841 scripts/flai-test.sh's repository check fails on this host while story/S-0001, S-0002, S-0003, and S-0005 are never merged: a draft; finalize it.
- S-0842 workitem's Item.Closed and PlanArchive read the baseline's done and cancelled names, so on a board that renames them flai accept moves a story to the accepted state and then cannot archive it: a draft; finalize it. It waits on TH-1296 Plan for S-0842: Closed and PlanArchive read the workflow's closed kinds.
- S-0843 Callers of the workflow-aware guard, forecast replay, and share hold still pass the baseline, so a project's own workflow file will not reach them: a draft; finalize it. It waits on TH-1295 Plan for S-0843: four tasks in three layers pass the project's workflow to the baseline callers left.
- S-0844 ADR-0301 rules for a project's own workflow file that no code enforces yet: operator-only acceptance of a change to it, the dashboard's refusal to save it on main, board.wip-state, and WIP limits on a state the project adds: a draft; finalize it.
- S-0845 A story in review whose close-out fails after main moved waits for the operator to send it back, while its after chain holds the backlog and the board starves: a draft; finalize it.
- S-0846 flai serve's auto-upgrade restart interrupts the close-out jobs it was running, which start again from scratch: a draft; finalize it.
- S-0847 flai serve opens a new board-starvation thread at each flip of the board, 48 in five days, each a thread to the operator: a draft; finalize it. It waits on TH-1299 Plan for S-0847: a settle window for the starved board's thread, in three layers.

### 2026-11-17T12:34:24Z flai serve
The list changed. No story in ready can start, no backlog story can be promoted, and 1 in-progress lane stands idle.

This host has no reminder channel, so without one nothing reminds you of this thread off the dashboard. Set one in a shell on the host with `flai serve notify set --url <webhook>` or `flai serve notify set --command <command>`.

Each backlog story, in the pull order, and the decision it waits on:

- S-0241 Measure planned stories again with forked task sub-agents once a headless session offers forks: a draft; finalize it. It waits on TH-1165 S-0241 stays a draft: it is blocked until a headless session offers forks.
- S-0813 A tier's queue and when are read as the suites of the batch and the move to review, flai upgrade writes them into a project's file, and the two keys are retired: held by `after`; the root of its chain is S-0806 The workflow file takes verify on a transition, runs on a state, and verify_at for the batch, the publish, and CI, and the baseline and the template's file give today's behaviour (in-progress).
- S-0814 The review page and the board show the suites each step needs and ran, and follow a state's job in the activity stream: held by `after`; the root of its chain is S-0806 The workflow file takes verify on a transition, runs on a state, and verify_at for the batch, the publish, and CI, and the baseline and the template's file give today's behaviour (in-progress).
- S-0815 In a project made from the template, review needs only the fast suites and an integration state between review and done deploys where a cluster is configured and runs the integration and smoke suites: held by `after`; the root of its chain is S-0806 The workflow file takes verify on a transition, runs on a state, and verify_at for the batch, the publish, and CI, and the baseline and the template's file give today's behaviour (in-progress).
- S-0822 A rebase and resolve job that passes records stayed in review without reading the story, and a move out of review neither stops the job nor is noticed by it: a draft; finalize it. It waits on TH-1274 Should a move out of review stop a running rebase and resolve job? and TH-1275 S-0822 stays a draft: no touches, forecast, or cost of delay value yet.
- S-0824 A rebase and resolve job runs the whole close-out again each time main moves under it, so one took three rounds and 33 minutes while other stories were accepted: a draft; finalize it. It waits on TH-1287 Plan for S-0824: a resolve job's last check passes over a main that sends back no tier, six tasks in four layers.
- S-0829 A sub-agent's flai kluster status falls through to the guard's generic refusal, against kluster.md's rule that status is a read every agent may run: a draft; finalize it.
- S-0835 flai serve's review pass deploys again a story in review whose deploy passed at its head and only its checks have not, where ADR-0292 §3 runs the checks alone: held by `after`; the root of its chain is S-0806 The workflow file takes verify on a transition, runs on a state, and verify_at for the batch, the publish, and CI, and the baseline and the template's file give today's behaviour (in-progress).
- S-0837 A publish the orchestrator starts is a release_publish call and makes no story, and the baseline publish formula is retired: a draft; finalize it.
- S-0838 flaiover still compares state names in the epic filters and the move route, so an operator's own accepted or cancelled lane is not seen as closed there: a draft; finalize it. It waits on TH-1297 Plan for S-0838: read each lane's kind from the board's flow in the epic filters, the move route, and ProposeChange.
- S-0839 No flai command writes agent.image.args, so a project's layer build arguments can only be written by hand: a draft; finalize it. It waits on TH-1298 S-0839 plan: flai manifest set writes agent.image.layer and agent.image.args.<NAME>.
- S-0840 flai kluster up cannot refresh a stale kubeconfig: kluster up refuses an existing cluster before flai writes it: a draft; finalize it. It waits on TH-1294 S-0840 plan: three tasks in three layers to make flai kluster up write a running cluster's kubeconfig again.
- S-0841 scripts/flai-test.sh's repository check fails on this host while story/S-0001, S-0002, S-0003, and S-0005 are never merged: a draft; finalize it.
- S-0842 workitem's Item.Closed and PlanArchive read the baseline's done and cancelled names, so on a board that renames them flai accept moves a story to the accepted state and then cannot archive it: a draft; finalize it. It waits on TH-1296 Plan for S-0842: Closed and PlanArchive read the workflow's closed kinds.
- S-0843 Callers of the workflow-aware guard, forecast replay, and share hold still pass the baseline, so a project's own workflow file will not reach them: a draft; finalize it. It waits on TH-1295 Plan for S-0843: four tasks in three layers pass the project's workflow to the baseline callers left.
- S-0844 ADR-0301 rules for a project's own workflow file that no code enforces yet: operator-only acceptance of a change to it, the dashboard's refusal to save it on main, board.wip-state, and WIP limits on a state the project adds: a draft; finalize it.
- S-0845 A story in review whose close-out fails after main moved waits for the operator to send it back, while its after chain holds the backlog and the board starves: a draft; finalize it.
- S-0846 flai serve's auto-upgrade restart interrupts the close-out jobs it was running, which start again from scratch: a draft; finalize it. It waits on TH-1300 Plan for S-0846: drain close-out jobs before the auto-upgrade restart, and ask again what a restart still interrupts.
- S-0847 flai serve opens a new board-starvation thread at each flip of the board, 48 in five days, each a thread to the operator: a draft; finalize it. It waits on TH-1299 Plan for S-0847: a settle window for the starved board's thread, in three layers.

### 2026-11-17T12:35:10Z flai serve
The list changed. No story in ready can start, no backlog story can be promoted, and 1 in-progress lane stands idle.

This host has no reminder channel, so without one nothing reminds you of this thread off the dashboard. Set one in a shell on the host with `flai serve notify set --url <webhook>` or `flai serve notify set --command <command>`.

Each backlog story, in the pull order, and the decision it waits on:

- S-0241 Measure planned stories again with forked task sub-agents once a headless session offers forks: a draft; finalize it. It waits on TH-1165 S-0241 stays a draft: it is blocked until a headless session offers forks.
- S-0813 A tier's queue and when are read as the suites of the batch and the move to review, flai upgrade writes them into a project's file, and the two keys are retired: held by `after`; the root of its chain is S-0806 The workflow file takes verify on a transition, runs on a state, and verify_at for the batch, the publish, and CI, and the baseline and the template's file give today's behaviour (in-progress).
- S-0814 The review page and the board show the suites each step needs and ran, and follow a state's job in the activity stream: held by `after`; the root of its chain is S-0806 The workflow file takes verify on a transition, runs on a state, and verify_at for the batch, the publish, and CI, and the baseline and the template's file give today's behaviour (in-progress).
- S-0815 In a project made from the template, review needs only the fast suites and an integration state between review and done deploys where a cluster is configured and runs the integration and smoke suites: held by `after`; the root of its chain is S-0806 The workflow file takes verify on a transition, runs on a state, and verify_at for the batch, the publish, and CI, and the baseline and the template's file give today's behaviour (in-progress).
- S-0822 A rebase and resolve job that passes records stayed in review without reading the story, and a move out of review neither stops the job nor is noticed by it: a draft; finalize it. It waits on TH-1274 Should a move out of review stop a running rebase and resolve job? and TH-1275 S-0822 stays a draft: no touches, forecast, or cost of delay value yet.
- S-0824 A rebase and resolve job runs the whole close-out again each time main moves under it, so one took three rounds and 33 minutes while other stories were accepted: a draft; finalize it. It waits on TH-1287 Plan for S-0824: a resolve job's last check passes over a main that sends back no tier, six tasks in four layers.
- S-0829 A sub-agent's flai kluster status falls through to the guard's generic refusal, against kluster.md's rule that status is a read every agent may run: a draft; finalize it.
- S-0835 flai serve's review pass deploys again a story in review whose deploy passed at its head and only its checks have not, where ADR-0292 §3 runs the checks alone: held by `after`; the root of its chain is S-0806 The workflow file takes verify on a transition, runs on a state, and verify_at for the batch, the publish, and CI, and the baseline and the template's file give today's behaviour (in-progress).
- S-0837 A publish the orchestrator starts is a release_publish call and makes no story, and the baseline publish formula is retired: a draft; finalize it.
- S-0838 flaiover still compares state names in the epic filters and the move route, so an operator's own accepted or cancelled lane is not seen as closed there: a draft; finalize it. It waits on TH-1297 Plan for S-0838: read each lane's kind from the board's flow in the epic filters, the move route, and ProposeChange.
- S-0839 No flai command writes agent.image.args, so a project's layer build arguments can only be written by hand: a draft; finalize it. It waits on TH-1298 S-0839 plan: flai manifest set writes agent.image.layer and agent.image.args.<NAME>.
- S-0840 flai kluster up cannot refresh a stale kubeconfig: kluster up refuses an existing cluster before flai writes it: a draft; finalize it. It waits on TH-1294 S-0840 plan: three tasks in three layers to make flai kluster up write a running cluster's kubeconfig again.
- S-0841 scripts/flai-test.sh's repository check fails on this host while story/S-0001, S-0002, S-0003, and S-0005 are never merged: a draft; finalize it.
- S-0842 workitem's Item.Closed and PlanArchive read the baseline's done and cancelled names, so on a board that renames them flai accept moves a story to the accepted state and then cannot archive it: a draft; finalize it. It waits on TH-1296 Plan for S-0842: Closed and PlanArchive read the workflow's closed kinds.
- S-0843 Callers of the workflow-aware guard, forecast replay, and share hold still pass the baseline, so a project's own workflow file will not reach them: a draft; finalize it.
- S-0844 ADR-0301 rules for a project's own workflow file that no code enforces yet: operator-only acceptance of a change to it, the dashboard's refusal to save it on main, board.wip-state, and WIP limits on a state the project adds: a draft; finalize it.
- S-0845 A story in review whose close-out fails after main moved waits for the operator to send it back, while its after chain holds the backlog and the board starves: a draft; finalize it.
- S-0846 flai serve's auto-upgrade restart interrupts the close-out jobs it was running, which start again from scratch: a draft; finalize it. It waits on TH-1300 Plan for S-0846: drain close-out jobs before the auto-upgrade restart, and ask again what a restart still interrupts.
- S-0847 flai serve opens a new board-starvation thread at each flip of the board, 48 in five days, each a thread to the operator: a draft; finalize it. It waits on TH-1299 Plan for S-0847: a settle window for the starved board's thread, in three layers.

The list changed. No story in ready can start, no backlog story can be promoted, and 1 in-progress lane stands idle.

This host has no reminder channel, so without one nothing reminds you of this thread off the dashboard. Set one in a shell on the host with `flai serve notify set --url <webhook>` or `flai serve notify set --command <command>`.

Each backlog story, in the pull order, and the decision it waits on:

- S-0241 Measure planned stories again with forked task sub-agents once a headless session offers forks: a draft; finalize it. It waits on TH-1165 S-0241 stays a draft: it is blocked until a headless session offers forks.
- S-0813 A tier's queue and when are read as the suites of the batch and the move to review, flai upgrade writes them into a project's file, and the two keys are retired: held by `after`; the root of its chain is S-0806 The workflow file takes verify on a transition, runs on a state, and verify_at for the batch, the publish, and CI, and the baseline and the template's file give today's behaviour (in-progress).
- S-0814 The review page and the board show the suites each step needs and ran, and follow a state's job in the activity stream: held by `after`; the root of its chain is S-0806 The workflow file takes verify on a transition, runs on a state, and verify_at for the batch, the publish, and CI, and the baseline and the template's file give today's behaviour (in-progress).
- S-0815 In a project made from the template, review needs only the fast suites and an integration state between review and done deploys where a cluster is configured and runs the integration and smoke suites: held by `after`; the root of its chain is S-0806 The workflow file takes verify on a transition, runs on a state, and verify_at for the batch, the publish, and CI, and the baseline and the template's file give today's behaviour (in-progress).
- S-0822 A rebase and resolve job that passes records stayed in review without reading the story, and a move out of review neither stops the job nor is noticed by it: a draft; finalize it. It waits on TH-1274 Should a move out of review stop a running rebase and resolve job? and TH-1275 S-0822 stays a draft: no touches, forecast, or cost of delay value yet.
- S-0824 A rebase and resolve job runs the whole close-out again each time main moves under it, so one took three rounds and 33 minutes while other stories were accepted: a draft; finalize it. It waits on TH-1287 Plan for S-0824: a resolve job's last check passes over a main that sends back no tier, six tasks in four layers.
- S-0829 A sub-agent's flai kluster status falls through to the guard's generic refusal, against kluster.md's rule that status is a read every agent may run: a draft; finalize it.
- S-0835 flai serve's review pass deploys again a story in review whose deploy passed at its head and only its checks have not, where ADR-0292 §3 runs the checks alone: held by `after`; the root of its chain is S-0806 The workflow file takes verify on a transition, runs on a state, and verify_at for the batch, the publish, and CI, and the baseline and the template's file give today's behaviour (in-progress).
- S-0837 A publish the orchestrator starts is a release_publish call and makes no story, and the baseline publish formula is retired: a draft; finalize it.
- S-0838 flaiover still compares state names in the epic filters and the move route, so an operator's own accepted or cancelled lane is not seen as closed there: a draft; finalize it. It waits on TH-1297 Plan for S-0838: read each lane's kind from the board's flow in the epic filters, the move route, and ProposeChange.
- S-0839 No flai command writes agent.image.args, so a project's layer build arguments can only be written by hand: a draft; finalize it. It waits on TH-1298 S-0839 plan: flai manifest set writes agent.image.layer and agent.image.args.<NAME>.
- S-0840 flai kluster up cannot refresh a stale kubeconfig: kluster up refuses an existing cluster before flai writes it: a draft; finalize it. It waits on TH-1294 S-0840 plan: three tasks in three layers to make flai kluster up write a running cluster's kubeconfig again.
- S-0841 scripts/flai-test.sh's repository check fails on this host while story/S-0001, S-0002, S-0003, and S-0005 are never merged: a draft; finalize it.
- S-0842 workitem's Item.Closed and PlanArchive read the baseline's done and cancelled names, so on a board that renames them flai accept moves a story to the accepted state and then cannot archive it: a draft; finalize it.
- S-0843 Callers of the workflow-aware guard, forecast replay, and share hold still pass the baseline, so a project's own workflow file will not reach them: a draft; finalize it.
- S-0844 ADR-0301 rules for a project's own workflow file that no code enforces yet: operator-only acceptance of a change to it, the dashboard's refusal to save it on main, board.wip-state, and WIP limits on a state the project adds: a draft; finalize it.
- S-0845 A story in review whose close-out fails after main moved waits for the operator to send it back, while its after chain holds the backlog and the board starves: a draft; finalize it.
- S-0846 flai serve's auto-upgrade restart interrupts the close-out jobs it was running, which start again from scratch: a draft; finalize it. It waits on TH-1300 Plan for S-0846: drain close-out jobs before the auto-upgrade restart, and ask again what a restart still interrupts.
- S-0847 flai serve opens a new board-starvation thread at each flip of the board, 48 in five days, each a thread to the operator: a draft; finalize it. It waits on TH-1299 Plan for S-0847: a settle window for the starved board's thread, in three layers.

### 2026-11-17T12:35:11Z flai serve
The list changed. No story in ready can start, no backlog story can be promoted, and 1 in-progress lane stands idle.

This host has no reminder channel, so without one nothing reminds you of this thread off the dashboard. Set one in a shell on the host with `flai serve notify set --url <webhook>` or `flai serve notify set --command <command>`.

Each backlog story, in the pull order, and the decision it waits on:

- S-0241 Measure planned stories again with forked task sub-agents once a headless session offers forks: a draft; finalize it. It waits on TH-1165 S-0241 stays a draft: it is blocked until a headless session offers forks.
- S-0813 A tier's queue and when are read as the suites of the batch and the move to review, flai upgrade writes them into a project's file, and the two keys are retired: held by `after`; the root of its chain is S-0806 The workflow file takes verify on a transition, runs on a state, and verify_at for the batch, the publish, and CI, and the baseline and the template's file give today's behaviour (in-progress).
- S-0814 The review page and the board show the suites each step needs and ran, and follow a state's job in the activity stream: held by `after`; the root of its chain is S-0806 The workflow file takes verify on a transition, runs on a state, and verify_at for the batch, the publish, and CI, and the baseline and the template's file give today's behaviour (in-progress).
- S-0815 In a project made from the template, review needs only the fast suites and an integration state between review and done deploys where a cluster is configured and runs the integration and smoke suites: held by `after`; the root of its chain is S-0806 The workflow file takes verify on a transition, runs on a state, and verify_at for the batch, the publish, and CI, and the baseline and the template's file give today's behaviour (in-progress).
- S-0822 A rebase and resolve job that passes records stayed in review without reading the story, and a move out of review neither stops the job nor is noticed by it: a draft; finalize it. It waits on TH-1274 Should a move out of review stop a running rebase and resolve job? and TH-1275 S-0822 stays a draft: no touches, forecast, or cost of delay value yet.
- S-0824 A rebase and resolve job runs the whole close-out again each time main moves under it, so one took three rounds and 33 minutes while other stories were accepted: a draft; finalize it. It waits on TH-1287 Plan for S-0824: a resolve job's last check passes over a main that sends back no tier, six tasks in four layers.
- S-0829 A sub-agent's flai kluster status falls through to the guard's generic refusal, against kluster.md's rule that status is a read every agent may run: a draft; finalize it.
- S-0835 flai serve's review pass deploys again a story in review whose deploy passed at its head and only its checks have not, where ADR-0292 §3 runs the checks alone: held by `after`; the root of its chain is S-0806 The workflow file takes verify on a transition, runs on a state, and verify_at for the batch, the publish, and CI, and the baseline and the template's file give today's behaviour (in-progress).
- S-0837 A publish the orchestrator starts is a release_publish call and makes no story, and the baseline publish formula is retired: a draft; finalize it.
- S-0838 flaiover still compares state names in the epic filters and the move route, so an operator's own accepted or cancelled lane is not seen as closed there: a draft; finalize it. It waits on TH-1297 Plan for S-0838: read each lane's kind from the board's flow in the epic filters, the move route, and ProposeChange.
- S-0839 No flai command writes agent.image.args, so a project's layer build arguments can only be written by hand: a draft; finalize it.
- S-0840 flai kluster up cannot refresh a stale kubeconfig: kluster up refuses an existing cluster before flai writes it: a draft; finalize it. It waits on TH-1294 S-0840 plan: three tasks in three layers to make flai kluster up write a running cluster's kubeconfig again.
- S-0841 scripts/flai-test.sh's repository check fails on this host while story/S-0001, S-0002, S-0003, and S-0005 are never merged: a draft; finalize it.
- S-0842 workitem's Item.Closed and PlanArchive read the baseline's done and cancelled names, so on a board that renames them flai accept moves a story to the accepted state and then cannot archive it: a draft; finalize it.
- S-0843 Callers of the workflow-aware guard, forecast replay, and share hold still pass the baseline, so a project's own workflow file will not reach them: a draft; finalize it.
- S-0844 ADR-0301 rules for a project's own workflow file that no code enforces yet: operator-only acceptance of a change to it, the dashboard's refusal to save it on main, board.wip-state, and WIP limits on a state the project adds: a draft; finalize it.
- S-0845 A story in review whose close-out fails after main moved waits for the operator to send it back, while its after chain holds the backlog and the board starves: a draft; finalize it.
- S-0846 flai serve's auto-upgrade restart interrupts the close-out jobs it was running, which start again from scratch: a draft; finalize it. It waits on TH-1300 Plan for S-0846: drain close-out jobs before the auto-upgrade restart, and ask again what a restart still interrupts.
- S-0847 flai serve opens a new board-starvation thread at each flip of the board, 48 in five days, each a thread to the operator: a draft; finalize it. It waits on TH-1299 Plan for S-0847: a settle window for the starved board's thread, in three layers.

The list changed. No story in ready can start, no backlog story can be promoted, and 1 in-progress lane stands idle.

This host has no reminder channel, so without one nothing reminds you of this thread off the dashboard. Set one in a shell on the host with `flai serve notify set --url <webhook>` or `flai serve notify set --command <command>`.

Each backlog story, in the pull order, and the decision it waits on:

- S-0241 Measure planned stories again with forked task sub-agents once a headless session offers forks: a draft; finalize it. It waits on TH-1165 S-0241 stays a draft: it is blocked until a headless session offers forks.
- S-0813 A tier's queue and when are read as the suites of the batch and the move to review, flai upgrade writes them into a project's file, and the two keys are retired: held by `after`; the root of its chain is S-0806 The workflow file takes verify on a transition, runs on a state, and verify_at for the batch, the publish, and CI, and the baseline and the template's file give today's behaviour (in-progress).
- S-0814 The review page and the board show the suites each step needs and ran, and follow a state's job in the activity stream: held by `after`; the root of its chain is S-0806 The workflow file takes verify on a transition, runs on a state, and verify_at for the batch, the publish, and CI, and the baseline and the template's file give today's behaviour (in-progress).
- S-0815 In a project made from the template, review needs only the fast suites and an integration state between review and done deploys where a cluster is configured and runs the integration and smoke suites: held by `after`; the root of its chain is S-0806 The workflow file takes verify on a transition, runs on a state, and verify_at for the batch, the publish, and CI, and the baseline and the template's file give today's behaviour (in-progress).
- S-0822 A rebase and resolve job that passes records stayed in review without reading the story, and a move out of review neither stops the job nor is noticed by it: a draft; finalize it. It waits on TH-1274 Should a move out of review stop a running rebase and resolve job? and TH-1275 S-0822 stays a draft: no touches, forecast, or cost of delay value yet.
- S-0824 A rebase and resolve job runs the whole close-out again each time main moves under it, so one took three rounds and 33 minutes while other stories were accepted: a draft; finalize it. It waits on TH-1287 Plan for S-0824: a resolve job's last check passes over a main that sends back no tier, six tasks in four layers.
- S-0829 A sub-agent's flai kluster status falls through to the guard's generic refusal, against kluster.md's rule that status is a read every agent may run: a draft; finalize it.
- S-0835 flai serve's review pass deploys again a story in review whose deploy passed at its head and only its checks have not, where ADR-0292 §3 runs the checks alone: held by `after`; the root of its chain is S-0806 The workflow file takes verify on a transition, runs on a state, and verify_at for the batch, the publish, and CI, and the baseline and the template's file give today's behaviour (in-progress).
- S-0837 A publish the orchestrator starts is a release_publish call and makes no story, and the baseline publish formula is retired: a draft; finalize it.
- S-0838 flaiover still compares state names in the epic filters and the move route, so an operator's own accepted or cancelled lane is not seen as closed there: a draft; finalize it. It waits on TH-1297 Plan for S-0838: read each lane's kind from the board's flow in the epic filters, the move route, and ProposeChange.
- S-0839 No flai command writes agent.image.args, so a project's layer build arguments can only be written by hand: a draft; finalize it.
- S-0840 flai kluster up cannot refresh a stale kubeconfig: kluster up refuses an existing cluster before flai writes it: a draft; finalize it. It waits on TH-1294 S-0840 plan: three tasks in three layers to make flai kluster up write a running cluster's kubeconfig again.
- S-0841 scripts/flai-test.sh's repository check fails on this host while story/S-0001, S-0002, S-0003, and S-0005 are never merged: a draft; finalize it.
- S-0842 workitem's Item.Closed and PlanArchive read the baseline's done and cancelled names, so on a board that renames them flai accept moves a story to the accepted state and then cannot archive it: a draft; finalize it.
- S-0843 Callers of the workflow-aware guard, forecast replay, and share hold still pass the baseline, so a project's own workflow file will not reach them: a draft; finalize it.
- S-0844 ADR-0301 rules for a project's own workflow file that no code enforces yet: operator-only acceptance of a change to it, the dashboard's refusal to save it on main, board.wip-state, and WIP limits on a state the project adds: a draft; finalize it.
- S-0845 A story in review whose close-out fails after main moved waits for the operator to send it back, while its after chain holds the backlog and the board starves: a draft; finalize it.
- S-0846 flai serve's auto-upgrade restart interrupts the close-out jobs it was running, which start again from scratch: a draft; finalize it. It waits on TH-1300 Plan for S-0846: drain close-out jobs before the auto-upgrade restart, and ask again what a restart still interrupts.
- S-0847 flai serve opens a new board-starvation thread at each flip of the board, 48 in five days, each a thread to the operator: a draft; finalize it.

### 2026-11-17T12:35:12Z flai serve
The list changed. No story in ready can start, no backlog story can be promoted, and 1 in-progress lane stands idle.

This host has no reminder channel, so without one nothing reminds you of this thread off the dashboard. Set one in a shell on the host with `flai serve notify set --url <webhook>` or `flai serve notify set --command <command>`.

Each backlog story, in the pull order, and the decision it waits on:

- S-0241 Measure planned stories again with forked task sub-agents once a headless session offers forks: a draft; finalize it. It waits on TH-1165 S-0241 stays a draft: it is blocked until a headless session offers forks.
- S-0813 A tier's queue and when are read as the suites of the batch and the move to review, flai upgrade writes them into a project's file, and the two keys are retired: held by `after`; the root of its chain is S-0806 The workflow file takes verify on a transition, runs on a state, and verify_at for the batch, the publish, and CI, and the baseline and the template's file give today's behaviour (in-progress).
- S-0814 The review page and the board show the suites each step needs and ran, and follow a state's job in the activity stream: held by `after`; the root of its chain is S-0806 The workflow file takes verify on a transition, runs on a state, and verify_at for the batch, the publish, and CI, and the baseline and the template's file give today's behaviour (in-progress).
- S-0815 In a project made from the template, review needs only the fast suites and an integration state between review and done deploys where a cluster is configured and runs the integration and smoke suites: held by `after`; the root of its chain is S-0806 The workflow file takes verify on a transition, runs on a state, and verify_at for the batch, the publish, and CI, and the baseline and the template's file give today's behaviour (in-progress).
- S-0822 A rebase and resolve job that passes records stayed in review without reading the story, and a move out of review neither stops the job nor is noticed by it: a draft; finalize it. It waits on TH-1274 Should a move out of review stop a running rebase and resolve job? and TH-1275 S-0822 stays a draft: no touches, forecast, or cost of delay value yet.
- S-0824 A rebase and resolve job runs the whole close-out again each time main moves under it, so one took three rounds and 33 minutes while other stories were accepted: a draft; finalize it. It waits on TH-1287 Plan for S-0824: a resolve job's last check passes over a main that sends back no tier, six tasks in four layers.
- S-0829 A sub-agent's flai kluster status falls through to the guard's generic refusal, against kluster.md's rule that status is a read every agent may run: a draft; finalize it.
- S-0835 flai serve's review pass deploys again a story in review whose deploy passed at its head and only its checks have not, where ADR-0292 §3 runs the checks alone: held by `after`; the root of its chain is S-0806 The workflow file takes verify on a transition, runs on a state, and verify_at for the batch, the publish, and CI, and the baseline and the template's file give today's behaviour (in-progress).
- S-0837 A publish the orchestrator starts is a release_publish call and makes no story, and the baseline publish formula is retired: a draft; finalize it.
- S-0838 flaiover still compares state names in the epic filters and the move route, so an operator's own accepted or cancelled lane is not seen as closed there: a draft; finalize it.
- S-0839 No flai command writes agent.image.args, so a project's layer build arguments can only be written by hand: a draft; finalize it.
- S-0840 flai kluster up cannot refresh a stale kubeconfig: kluster up refuses an existing cluster before flai writes it: a draft; finalize it. It waits on TH-1294 S-0840 plan: three tasks in three layers to make flai kluster up write a running cluster's kubeconfig again.
- S-0841 scripts/flai-test.sh's repository check fails on this host while story/S-0001, S-0002, S-0003, and S-0005 are never merged: a draft; finalize it.
- S-0842 workitem's Item.Closed and PlanArchive read the baseline's done and cancelled names, so on a board that renames them flai accept moves a story to the accepted state and then cannot archive it: a draft; finalize it.
- S-0843 Callers of the workflow-aware guard, forecast replay, and share hold still pass the baseline, so a project's own workflow file will not reach them: a draft; finalize it.
- S-0844 ADR-0301 rules for a project's own workflow file that no code enforces yet: operator-only acceptance of a change to it, the dashboard's refusal to save it on main, board.wip-state, and WIP limits on a state the project adds: a draft; finalize it.
- S-0845 A story in review whose close-out fails after main moved waits for the operator to send it back, while its after chain holds the backlog and the board starves: a draft; finalize it.
- S-0846 flai serve's auto-upgrade restart interrupts the close-out jobs it was running, which start again from scratch: a draft; finalize it. It waits on TH-1300 Plan for S-0846: drain close-out jobs before the auto-upgrade restart, and ask again what a restart still interrupts.
- S-0847 flai serve opens a new board-starvation thread at each flip of the board, 48 in five days, each a thread to the operator: a draft; finalize it.

### 2026-11-17T12:35:28Z flai serve
The list changed. No story in ready can start, no backlog story can be promoted, and 1 in-progress lane stands idle.

This host has no reminder channel, so without one nothing reminds you of this thread off the dashboard. Set one in a shell on the host with `flai serve notify set --url <webhook>` or `flai serve notify set --command <command>`.

Each backlog story, in the pull order, and the decision it waits on:

- S-0241 Measure planned stories again with forked task sub-agents once a headless session offers forks: a draft; finalize it. It waits on TH-1165 S-0241 stays a draft: it is blocked until a headless session offers forks.
- S-0813 A tier's queue and when are read as the suites of the batch and the move to review, flai upgrade writes them into a project's file, and the two keys are retired: held by `after`; the root of its chain is S-0806 The workflow file takes verify on a transition, runs on a state, and verify_at for the batch, the publish, and CI, and the baseline and the template's file give today's behaviour (in-progress).
- S-0814 The review page and the board show the suites each step needs and ran, and follow a state's job in the activity stream: held by `after`; the root of its chain is S-0806 The workflow file takes verify on a transition, runs on a state, and verify_at for the batch, the publish, and CI, and the baseline and the template's file give today's behaviour (in-progress).
- S-0815 In a project made from the template, review needs only the fast suites and an integration state between review and done deploys where a cluster is configured and runs the integration and smoke suites: held by `after`; the root of its chain is S-0806 The workflow file takes verify on a transition, runs on a state, and verify_at for the batch, the publish, and CI, and the baseline and the template's file give today's behaviour (in-progress).
- S-0822 A rebase and resolve job that passes records stayed in review without reading the story, and a move out of review neither stops the job nor is noticed by it: a draft; finalize it. It waits on TH-1274 Should a move out of review stop a running rebase and resolve job? and TH-1275 S-0822 stays a draft: no touches, forecast, or cost of delay value yet.
- S-0824 A rebase and resolve job runs the whole close-out again each time main moves under it, so one took three rounds and 33 minutes while other stories were accepted: a draft; finalize it. It waits on TH-1287 Plan for S-0824: a resolve job's last check passes over a main that sends back no tier, six tasks in four layers.
- S-0829 A sub-agent's flai kluster status falls through to the guard's generic refusal, against kluster.md's rule that status is a read every agent may run: a draft; finalize it.
- S-0835 flai serve's review pass deploys again a story in review whose deploy passed at its head and only its checks have not, where ADR-0292 §3 runs the checks alone: held by `after`; the root of its chain is S-0806 The workflow file takes verify on a transition, runs on a state, and verify_at for the batch, the publish, and CI, and the baseline and the template's file give today's behaviour (in-progress).
- S-0837 A publish the orchestrator starts is a release_publish call and makes no story, and the baseline publish formula is retired: a draft; finalize it.
- S-0838 flaiover still compares state names in the epic filters and the move route, so an operator's own accepted or cancelled lane is not seen as closed there: a draft; finalize it.
- S-0839 No flai command writes agent.image.args, so a project's layer build arguments can only be written by hand: a draft; finalize it.
- S-0840 flai kluster up cannot refresh a stale kubeconfig: kluster up refuses an existing cluster before flai writes it: a draft; finalize it. It waits on TH-1294 S-0840 plan: three tasks in three layers to make flai kluster up write a running cluster's kubeconfig again.
- S-0841 scripts/flai-test.sh's repository check fails on this host while story/S-0001, S-0002, S-0003, and S-0005 are never merged: a draft; finalize it.
- S-0842 workitem's Item.Closed and PlanArchive read the baseline's done and cancelled names, so on a board that renames them flai accept moves a story to the accepted state and then cannot archive it: a draft; finalize it.
- S-0843 Callers of the workflow-aware guard, forecast replay, and share hold still pass the baseline, so a project's own workflow file will not reach them: a draft; finalize it.
- S-0844 ADR-0301 rules for a project's own workflow file that no code enforces yet: operator-only acceptance of a change to it, the dashboard's refusal to save it on main, board.wip-state, and WIP limits on a state the project adds: a draft; finalize it.
- S-0845 A story in review whose close-out fails after main moved waits for the operator to send it back, while its after chain holds the backlog and the board starves: a draft; finalize it. It waits on TH-1301 Plan for S-0845: the orchestrator asks for rebase and resolve when a review story's close-out did not pass after main moved.
- S-0846 flai serve's auto-upgrade restart interrupts the close-out jobs it was running, which start again from scratch: a draft; finalize it. It waits on TH-1300 Plan for S-0846: drain close-out jobs before the auto-upgrade restart, and ask again what a restart still interrupts.
- S-0847 flai serve opens a new board-starvation thread at each flip of the board, 48 in five days, each a thread to the operator: a draft; finalize it.

### 2026-11-17T12:35:32Z flai serve
Resolved: S-0838 flaiover still compares state names in the epic filters and the move route, so an operator's own accepted or cancelled lane is not seen as closed there can be promoted from the backlog, so the board is no longer starved.
