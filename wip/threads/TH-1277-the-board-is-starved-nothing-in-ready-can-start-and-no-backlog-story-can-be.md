---
id: TH-1277
title: "The board is starved: nothing in ready can start and no backlog story can be promoted"
anchor:
  path: wip/kanban/board.md
status: resolved
participants: [flai serve]
created: 2026-11-17T08:29:12Z
updated: 2026-11-17T08:30:10Z
---

# TH-1277 The board is starved: nothing in ready can start and no backlog story can be promoted

On wip/kanban/board.md.

## Entries

### 2026-11-17T08:29:12Z flai serve
No story in ready can start, no backlog story can be promoted, and 1 in-progress lane stands idle. flai serve adds an entry here when the list below changes, and resolves this thread once a story can start.

This host has no reminder channel, so without one nothing reminds you of this thread off the dashboard. Set one in a shell on the host with `flai serve notify set --url <webhook>` or `flai serve notify set --command <command>`.

Each backlog story, in the pull order, and the decision it waits on:

- S-0241 Measure planned stories again with forked task sub-agents once a headless session offers forks: a draft; finalize it. It waits on TH-1165 S-0241 stays a draft: it is blocked until a headless session offers forks.
- S-0813 A tier's queue and when are read as the suites of the batch and the move to review, flai upgrade writes them into a project's file, and the two keys are retired: held by `after`; the root of its chain is S-0806 The workflow file takes verify on a transition, runs on a state, and verify_at for the batch, the publish, and CI, and the baseline and the template's file give today's behaviour (in-progress).
- S-0814 The review page and the board show the suites each step needs and ran, and follow a state's job in the activity stream: held by `after`; the root of its chain is S-0806 The workflow file takes verify on a transition, runs on a state, and verify_at for the batch, the publish, and CI, and the baseline and the template's file give today's behaviour (in-progress).
- S-0815 In a project made from the template, review needs only the fast suites and an integration state between review and done deploys where a cluster is configured and runs the integration and smoke suites: held by `after`; the roots of its chain are S-0806 The workflow file takes verify on a transition, runs on a state, and verify_at for the batch, the publish, and CI, and the baseline and the template's file give today's behaviour (in-progress), S-0752 flai stats reports each timer a workflow defines beside the baseline's figures, for every item and over the window, past items included (in-progress) and S-0758 The dashboard's item page, inbox, review, activity, deployments, task plan, and cards read a state's kind from flai instead of comparing its name (review).
- S-0822 A rebase and resolve job that passes records stayed in review without reading the story, and a move out of review neither stops the job nor is noticed by it: a draft; finalize it. It waits on TH-1274 Should a move out of review stop a running rebase and resolve job? and TH-1275 S-0822 stays a draft: no touches, forecast, or cost of delay value yet.
- S-0823 A story the operator moves from review to in-progress by hand gets no agent and no thread, so it waits in in-progress with nothing to bring it back: a draft; finalize it.
- S-0824 A rebase and resolve job runs the whole close-out again each time main moves under it, so one took three rounds and 33 minutes while other stories were accepted: a draft; finalize it.
- S-0827 TestRepositoryManifestLintsChangedFlaioverFiles skips under -short, so a new tier it must list first fails in the close-out's integration tier: a draft; finalize it. It waits on TH-1276 Plan for S-0827: run TestRepositoryManifestLintsChangedFlaioverFiles under -short, two tasks in two layers.
- S-0828 board.get's undone.read runs one git diff per recorded path of each story's undone record, one after another, on every call: a draft; finalize it.
- S-0829 A sub-agent's flai kluster status falls through to the guard's generic refusal, against kluster.md's rule that status is a read every agent may run: a draft; finalize it.
- S-0830 The dashboard's kluster setup form does not show kluster.agent_namespaces: a draft; finalize it.
- S-0831 The kluster send-back's texts to the agent do not say it can ask for its own deploy before review: a draft; finalize it.
- S-0832 flai guard reads a subcommand's value flag before flai's command as the command, so a rule can miss the command flai runs: a draft; finalize it.
- S-0833 flai release --pending fails when another process holds the main checkout's git index lock, and the publish formula's step then stops on exit 1: a draft; finalize it.
- S-0834 flai kluster exec and log read .flai-cache/kluster, which a story's agent in a container does not mount: a draft; finalize it.
- S-0835 flai serve's review pass deploys again a story in review whose deploy passed at its head and only its checks have not, where ADR-0292 §3 runs the checks alone: a draft; finalize it.
- S-0836 flai's markdown lint reports MD038 on a code span holding escaped backticks that markdownlint passes, so TestRepositoryLintsClean fails on a curator's digest entry: a draft; finalize it.
- S-0837 A publish the orchestrator starts is a release_publish call and makes no story, and the baseline publish formula is retired: a draft; finalize it.

### 2026-11-17T08:29:35Z flai serve
The list changed. No story in ready can start, no backlog story can be promoted, and 1 in-progress lane stands idle.

This host has no reminder channel, so without one nothing reminds you of this thread off the dashboard. Set one in a shell on the host with `flai serve notify set --url <webhook>` or `flai serve notify set --command <command>`.

Each backlog story, in the pull order, and the decision it waits on:

- S-0241 Measure planned stories again with forked task sub-agents once a headless session offers forks: a draft; finalize it. It waits on TH-1165 S-0241 stays a draft: it is blocked until a headless session offers forks.
- S-0813 A tier's queue and when are read as the suites of the batch and the move to review, flai upgrade writes them into a project's file, and the two keys are retired: held by `after`; the root of its chain is S-0806 The workflow file takes verify on a transition, runs on a state, and verify_at for the batch, the publish, and CI, and the baseline and the template's file give today's behaviour (in-progress).
- S-0814 The review page and the board show the suites each step needs and ran, and follow a state's job in the activity stream: held by `after`; the root of its chain is S-0806 The workflow file takes verify on a transition, runs on a state, and verify_at for the batch, the publish, and CI, and the baseline and the template's file give today's behaviour (in-progress).
- S-0815 In a project made from the template, review needs only the fast suites and an integration state between review and done deploys where a cluster is configured and runs the integration and smoke suites: held by `after`; the roots of its chain are S-0806 The workflow file takes verify on a transition, runs on a state, and verify_at for the batch, the publish, and CI, and the baseline and the template's file give today's behaviour (in-progress), S-0752 flai stats reports each timer a workflow defines beside the baseline's figures, for every item and over the window, past items included (in-progress) and S-0758 The dashboard's item page, inbox, review, activity, deployments, task plan, and cards read a state's kind from flai instead of comparing its name (review).
- S-0822 A rebase and resolve job that passes records stayed in review without reading the story, and a move out of review neither stops the job nor is noticed by it: a draft; finalize it. It waits on TH-1274 Should a move out of review stop a running rebase and resolve job? and TH-1275 S-0822 stays a draft: no touches, forecast, or cost of delay value yet.
- S-0823 A story the operator moves from review to in-progress by hand gets no agent and no thread, so it waits in in-progress with nothing to bring it back: a draft; finalize it.
- S-0824 A rebase and resolve job runs the whole close-out again each time main moves under it, so one took three rounds and 33 minutes while other stories were accepted: a draft; finalize it.
- S-0827 TestRepositoryManifestLintsChangedFlaioverFiles skips under -short, so a new tier it must list first fails in the close-out's integration tier: a draft; finalize it. It waits on TH-1276 Plan for S-0827: run TestRepositoryManifestLintsChangedFlaioverFiles under -short, two tasks in two layers.
- S-0828 board.get's undone.read runs one git diff per recorded path of each story's undone record, one after another, on every call: a draft; finalize it. It waits on TH-1278 Plan for S-0828: batch the undone record's git diffs.
- S-0829 A sub-agent's flai kluster status falls through to the guard's generic refusal, against kluster.md's rule that status is a read every agent may run: a draft; finalize it.
- S-0830 The dashboard's kluster setup form does not show kluster.agent_namespaces: a draft; finalize it.
- S-0831 The kluster send-back's texts to the agent do not say it can ask for its own deploy before review: a draft; finalize it.
- S-0832 flai guard reads a subcommand's value flag before flai's command as the command, so a rule can miss the command flai runs: a draft; finalize it.
- S-0833 flai release --pending fails when another process holds the main checkout's git index lock, and the publish formula's step then stops on exit 1: a draft; finalize it.
- S-0834 flai kluster exec and log read .flai-cache/kluster, which a story's agent in a container does not mount: a draft; finalize it.
- S-0835 flai serve's review pass deploys again a story in review whose deploy passed at its head and only its checks have not, where ADR-0292 §3 runs the checks alone: a draft; finalize it.
- S-0836 flai's markdown lint reports MD038 on a code span holding escaped backticks that markdownlint passes, so TestRepositoryLintsClean fails on a curator's digest entry: a draft; finalize it.
- S-0837 A publish the orchestrator starts is a release_publish call and makes no story, and the baseline publish formula is retired: a draft; finalize it.

### 2026-11-17T08:29:41Z flai serve
The list changed. No story in ready can start, no backlog story can be promoted, and 1 in-progress lane stands idle.

This host has no reminder channel, so without one nothing reminds you of this thread off the dashboard. Set one in a shell on the host with `flai serve notify set --url <webhook>` or `flai serve notify set --command <command>`.

Each backlog story, in the pull order, and the decision it waits on:

- S-0241 Measure planned stories again with forked task sub-agents once a headless session offers forks: a draft; finalize it. It waits on TH-1165 S-0241 stays a draft: it is blocked until a headless session offers forks.
- S-0813 A tier's queue and when are read as the suites of the batch and the move to review, flai upgrade writes them into a project's file, and the two keys are retired: held by `after`; the root of its chain is S-0806 The workflow file takes verify on a transition, runs on a state, and verify_at for the batch, the publish, and CI, and the baseline and the template's file give today's behaviour (in-progress).
- S-0814 The review page and the board show the suites each step needs and ran, and follow a state's job in the activity stream: held by `after`; the root of its chain is S-0806 The workflow file takes verify on a transition, runs on a state, and verify_at for the batch, the publish, and CI, and the baseline and the template's file give today's behaviour (in-progress).
- S-0815 In a project made from the template, review needs only the fast suites and an integration state between review and done deploys where a cluster is configured and runs the integration and smoke suites: held by `after`; the roots of its chain are S-0806 The workflow file takes verify on a transition, runs on a state, and verify_at for the batch, the publish, and CI, and the baseline and the template's file give today's behaviour (in-progress), S-0752 flai stats reports each timer a workflow defines beside the baseline's figures, for every item and over the window, past items included (in-progress) and S-0758 The dashboard's item page, inbox, review, activity, deployments, task plan, and cards read a state's kind from flai instead of comparing its name (review).
- S-0822 A rebase and resolve job that passes records stayed in review without reading the story, and a move out of review neither stops the job nor is noticed by it: a draft; finalize it. It waits on TH-1274 Should a move out of review stop a running rebase and resolve job? and TH-1275 S-0822 stays a draft: no touches, forecast, or cost of delay value yet.
- S-0823 A story the operator moves from review to in-progress by hand gets no agent and no thread, so it waits in in-progress with nothing to bring it back: a draft; finalize it.
- S-0824 A rebase and resolve job runs the whole close-out again each time main moves under it, so one took three rounds and 33 minutes while other stories were accepted: a draft; finalize it.
- S-0827 TestRepositoryManifestLintsChangedFlaioverFiles skips under -short, so a new tier it must list first fails in the close-out's integration tier: a draft; finalize it. It waits on TH-1276 Plan for S-0827: run TestRepositoryManifestLintsChangedFlaioverFiles under -short, two tasks in two layers.
- S-0828 board.get's undone.read runs one git diff per recorded path of each story's undone record, one after another, on every call: a draft; finalize it. It waits on TH-1278 Plan for S-0828: batch the undone record's git diffs.
- S-0829 A sub-agent's flai kluster status falls through to the guard's generic refusal, against kluster.md's rule that status is a read every agent may run: a draft; finalize it.
- S-0830 The dashboard's kluster setup form does not show kluster.agent_namespaces: a draft; finalize it.
- S-0831 The kluster send-back's texts to the agent do not say it can ask for its own deploy before review: a draft; finalize it.
- S-0832 flai guard reads a subcommand's value flag before flai's command as the command, so a rule can miss the command flai runs: a draft; finalize it.
- S-0833 flai release --pending fails when another process holds the main checkout's git index lock, and the publish formula's step then stops on exit 1: a draft; finalize it.
- S-0834 flai kluster exec and log read .flai-cache/kluster, which a story's agent in a container does not mount: a draft; finalize it.
- S-0835 flai serve's review pass deploys again a story in review whose deploy passed at its head and only its checks have not, where ADR-0292 §3 runs the checks alone: a draft; finalize it.
- S-0836 flai's markdown lint reports MD038 on a code span holding escaped backticks that markdownlint passes, so TestRepositoryLintsClean fails on a curator's digest entry: a draft; finalize it. It waits on TH-1279 Plan for S-0836: fix mdlint's MD038 false positive on D-0592, then close I-0363.
- S-0837 A publish the orchestrator starts is a release_publish call and makes no story, and the baseline publish formula is retired: a draft; finalize it.

### 2026-11-17T08:29:55Z flai serve
The list changed. No story in ready can start, no backlog story can be promoted, and 1 in-progress lane stands idle.

This host has no reminder channel, so without one nothing reminds you of this thread off the dashboard. Set one in a shell on the host with `flai serve notify set --url <webhook>` or `flai serve notify set --command <command>`.

Each backlog story, in the pull order, and the decision it waits on:

- S-0241 Measure planned stories again with forked task sub-agents once a headless session offers forks: a draft; finalize it. It waits on TH-1165 S-0241 stays a draft: it is blocked until a headless session offers forks.
- S-0813 A tier's queue and when are read as the suites of the batch and the move to review, flai upgrade writes them into a project's file, and the two keys are retired: held by `after`; the root of its chain is S-0806 The workflow file takes verify on a transition, runs on a state, and verify_at for the batch, the publish, and CI, and the baseline and the template's file give today's behaviour (in-progress).
- S-0814 The review page and the board show the suites each step needs and ran, and follow a state's job in the activity stream: held by `after`; the root of its chain is S-0806 The workflow file takes verify on a transition, runs on a state, and verify_at for the batch, the publish, and CI, and the baseline and the template's file give today's behaviour (in-progress).
- S-0815 In a project made from the template, review needs only the fast suites and an integration state between review and done deploys where a cluster is configured and runs the integration and smoke suites: held by `after`; the roots of its chain are S-0806 The workflow file takes verify on a transition, runs on a state, and verify_at for the batch, the publish, and CI, and the baseline and the template's file give today's behaviour (in-progress), S-0752 flai stats reports each timer a workflow defines beside the baseline's figures, for every item and over the window, past items included (in-progress) and S-0758 The dashboard's item page, inbox, review, activity, deployments, task plan, and cards read a state's kind from flai instead of comparing its name (review).
- S-0822 A rebase and resolve job that passes records stayed in review without reading the story, and a move out of review neither stops the job nor is noticed by it: a draft; finalize it. It waits on TH-1274 Should a move out of review stop a running rebase and resolve job? and TH-1275 S-0822 stays a draft: no touches, forecast, or cost of delay value yet.
- S-0823 A story the operator moves from review to in-progress by hand gets no agent and no thread, so it waits in in-progress with nothing to bring it back: a draft; finalize it.
- S-0824 A rebase and resolve job runs the whole close-out again each time main moves under it, so one took three rounds and 33 minutes while other stories were accepted: a draft; finalize it.
- S-0827 TestRepositoryManifestLintsChangedFlaioverFiles skips under -short, so a new tier it must list first fails in the close-out's integration tier: a draft; finalize it. It waits on TH-1276 Plan for S-0827: run TestRepositoryManifestLintsChangedFlaioverFiles under -short, two tasks in two layers.
- S-0828 board.get's undone.read runs one git diff per recorded path of each story's undone record, one after another, on every call: a draft; finalize it. It waits on TH-1278 Plan for S-0828: batch the undone record's git diffs.
- S-0829 A sub-agent's flai kluster status falls through to the guard's generic refusal, against kluster.md's rule that status is a read every agent may run: a draft; finalize it.
- S-0830 The dashboard's kluster setup form does not show kluster.agent_namespaces: a draft; finalize it. It waits on TH-1280 Plan for S-0830: show and save kluster.agent_namespaces in the setup form.
- S-0831 The kluster send-back's texts to the agent do not say it can ask for its own deploy before review: a draft; finalize it.
- S-0832 flai guard reads a subcommand's value flag before flai's command as the command, so a rule can miss the command flai runs: a draft; finalize it.
- S-0833 flai release --pending fails when another process holds the main checkout's git index lock, and the publish formula's step then stops on exit 1: a draft; finalize it.
- S-0834 flai kluster exec and log read .flai-cache/kluster, which a story's agent in a container does not mount: a draft; finalize it.
- S-0835 flai serve's review pass deploys again a story in review whose deploy passed at its head and only its checks have not, where ADR-0292 §3 runs the checks alone: a draft; finalize it.
- S-0836 flai's markdown lint reports MD038 on a code span holding escaped backticks that markdownlint passes, so TestRepositoryLintsClean fails on a curator's digest entry: a draft; finalize it. It waits on TH-1279 Plan for S-0836: fix mdlint's MD038 false positive on D-0592, then close I-0363.
- S-0837 A publish the orchestrator starts is a release_publish call and makes no story, and the baseline publish formula is retired: a draft; finalize it.

### 2026-11-17T08:30:10Z flai serve
Resolved: S-0827 TestRepositoryManifestLintsChangedFlaioverFiles skips under -short, so a new tier it must list first fails in the close-out's integration tier in ready can start, so the board is no longer starved.
