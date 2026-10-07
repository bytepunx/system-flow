---
id: T-0946
type: task
nature: feature
title: The dashboard design and user guide describe the parallelism, hold-time, and touches-drift charts
status: done
parent: S-0214
owner: alex
created: 2026-10-05T05:45:44Z
updated: 2026-10-07T08:43:09Z
transitions:
  - to: ready
    at: 2026-10-07T08:41:20Z
    by: agent-S-0214
  - to: in-progress
    at: 2026-10-07T08:41:21Z
    by: agent-S-0214
  - to: done
    at: 2026-10-07T08:43:09Z
    by: agent-S-0214
stream: S-0214
tags: [dashboard, docs]
touches: [design/system/flaiover-dashboard.md, docs/users/flaiover.md]
after: [T-0943]
usage:
  source: log
  seconds: 108
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 25
      output: 10745
      cache_read: 1242661
      cache_write: 46596
      cost: 0.736
---
# T-0946 The dashboard design and user guide describe the parallelism, hold-time, and touches-drift charts

## Work

Describe the three charts under Planning in `design/system/flaiover-dashboard.md` § Views, and add rows for them to the table in `docs/users/flaiover.md` § Charts. Say what each shows, the `flai stats` values it reads (`metrics.md` § Claims and touches), and how to read it: held stories beside the in-progress limit, hold hours by reason, and drift as touches the planner or agent should have declared or dropped. Waits for T-0943, so that the docs describe the page as built.

## Done when

- [ ] Both documents describe the three charts and link the metrics they read, and their `updated` dates are today
- [ ] `flai check --strict` and the markdown lint pass

## Notes
