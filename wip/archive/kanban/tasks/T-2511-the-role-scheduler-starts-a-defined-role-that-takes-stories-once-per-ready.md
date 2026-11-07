---
id: T-2511
type: task
nature: remediation
title: The role scheduler starts a defined role that takes stories once per ready story under each-ready-story, with a held start kept per story
status: done
parent: S-0571
owner: alex
created: 2026-11-07T03:46:34Z
updated: 2026-11-07T04:02:25Z
transitions:
  - to: ready
    at: 2026-11-07T03:50:57Z
    by: agent-S-0571
  - to: in-progress
    at: 2026-11-07T03:55:56Z
    by: system-flow
  - to: done
    at: 2026-11-07T04:02:25Z
    by: agent-S-0571
stream: S-0571
tags: [flai, serve, roles]
touches: [flai/internal/serve/role_schedule.go, flai/internal/serve/role_schedule_test.go]
usage:
  source: log
  seconds: 389
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 43
      output: 9198
      cache_read: 5292624
      cache_write: 98375
      cost: 2.0049
---
# T-2511 The role scheduler starts a defined role that takes stories once per ready story under each-ready-story, with a held start kept per story

## Work

In `flai/internal/serve/role_schedule.go`, `lookAt` logs "scheduled role not started" for every role whose file takes items, and starts nothing (I-0262). Change it so that when the setting comes round, each entry `matching` gives for the schedule occurrence with `when: each-ready-story` starts the role once for each story `roles.EachReadyStory` gives. Read the stories in the board's pull order, as the planner's replanner does in `replan.go`. Start each through `startRole` with the story's ID as the item and the schedule as its trigger.

- Keep a held start per story in place of `roleFollow.heldDue`'s one trigger: a map from story ID to trigger. Each look tries each one again until it starts. Drop a story that is no longer ready. Drop them all when the action is off.
- A refusal or an error is logged with the story's ID and does not stop the other stories' starts.
- A schedule entry without `each-ready-story` on a role that takes items stays logged as now. `flai check` reports it.
- A role that takes no item is started as before.

Update the file's header comment to say so. Add tests in `role_schedule_test.go`, on the fake launcher the file's tests use:

- a role with `item: [story]` and `{schedule: roles.<name>.schedule, when: each-ready-story}` starts once per ready story, in pull order, on each one's ID (the test that reproduces I-0262, failing before the change)
- a held story is tried again at the next look while the others are not
- an empty ready column starts nothing

This task waits for none.

## Done when

- `flai test flai/internal/serve` passes, with the new tests failing on the code before the change.
- A defined role that takes stories under `each-ready-story` is started on each ready story when its schedule comes round, and a held story start is retried at the next look.

## Notes
