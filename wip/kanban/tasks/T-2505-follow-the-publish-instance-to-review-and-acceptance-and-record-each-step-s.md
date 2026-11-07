---
id: T-2505
type: task
nature: improvement
title: Follow the publish instance to review and acceptance, and record each step's exit and the release or the evaluation that stopped it
status: backlog
parent: S-0562
owner: alex
created: 2026-11-07T03:45:19Z
updated: 2026-11-07T03:45:19Z
transitions: []
stream: S-0562
tags: [template]
touches: [wip/kanban/stories/S-0562-the-publish-formula-is-run-once-in-this-repository-its-instance-worked-to.md]
after: [T-2504]
---
# T-2505 Follow the publish instance to review and acceptance, and record each step's exit and the release or the evaluation that stopped it

## Work

Watch the instance T-2504 started. Its own agent works its three steps (`fetch`, `evaluate`, `publish`) from the main checkout and moves it to review. Its acceptance is the operator's, or the orchestrator's with `orchestration.permissions.accept_reviews` on. Hold `wait_for_events` for the instance's changes of kind `agent` and `moved`, and do not work the instance's tasks yourself.

When the instance is accepted, write in S-0562's Notes with `flai edit`:

- the instance's ID, and when it reached review and was accepted;
- each step's exit, as the instance's task Notes give it;
- the versions and tags step 3 printed, or, when step 2 stopped the run, the evaluation's reason and its count and value.

Then check that the instance's Notes name the release it published or the evaluation that stopped it, as criterion 1 asks. When a step fails or the instance is sent back, record what happened for T-2506 rather than fixing it here.

This task waits for T-2504, because there is no instance to follow until T-2504 starts one.

## Done when

- The instance is accepted, and its Notes name the release or the evaluation that stopped it.
- S-0562's Notes hold the instance's ID, each step's exit, and the versions and tags published, or the evaluation that stopped it.

## Notes
