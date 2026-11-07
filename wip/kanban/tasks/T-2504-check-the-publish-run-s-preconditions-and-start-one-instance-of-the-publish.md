---
id: T-2504
type: task
nature: improvement
title: Check the publish run's preconditions and start one instance of the publish formula from the main checkout
status: done
parent: S-0562
owner: alex
created: 2026-11-07T03:45:06Z
updated: 2026-11-07T08:04:22Z
transitions:
  - to: ready
    at: 2026-11-07T08:03:06Z
    by: agent-S-0562
  - to: in-progress
    at: 2026-11-07T08:03:08Z
    by: agent-S-0562
  - to: done
    at: 2026-11-07T08:04:22Z
    by: agent-S-0562
stream: S-0562
tags: [template]
touches: [wip/kanban/stories/S-0562-the-publish-formula-is-run-once-in-this-repository-its-instance-worked-to.md]
usage:
  source: log
  seconds: 73
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 16
      output: 61
      cache_read: 722644
      cache_write: 10846
      cost: 0.2625
---
# T-2504 Check the publish run's preconditions and start one instance of the publish formula from the main checkout

## Work

Before anything runs, check these on the host. Ask the operator on a thread on S-0562 for any that does not hold, and run nothing until it is answered:

- The installed `flai` holds S-0423's change: `flai release --evaluate --require-met` exists, and `flai release --pending` has the `push` gate. Any flai from 1.80.2 does.
- The `formula` and `push` host actions are on. `flai release --pending` refuses an agent `flai serve` started while `push` is off.
- TH-0811 is answered. Under the `judgement` release policy, the formula's step 2, `flai release --evaluate --require-met`, always stops. Run under judgement, the instance only shows the evaluation stopping, even with accepted stories pending.

Then start one instance in the main checkout with `flai formula run publish --ready --autocommit`, or let the orchestrator start it under its `publish` permission once `release_evaluate` is met (ADR-0240). Record the instance's story ID and how it was started in S-0562's Notes with `flai edit`. Run no second instance.

This task waits for no other task: it is the plan's first layer.

## Done when

- Each precondition above is checked, and the result is written in S-0562's Notes.
- One instance of the `publish` formula exists, in ready or later, and S-0562's Notes name its ID.

## Notes
