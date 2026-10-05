---
id: T-0946
type: task
nature: feature
title: The dashboard design and user guide describe the parallelism, hold-time, and touches-drift charts
status: backlog
parent: S-0214
owner: alex
created: 2026-10-05T05:45:44Z
updated: 2026-10-05T05:45:44Z
transitions: []
stream: S-0214
tags: [dashboard, docs]
touches: [design/system/flaiover-dashboard.md, docs/users/flaiover.md]
after: [T-0943]
---
# T-0946 The dashboard design and user guide describe the parallelism, hold-time, and touches-drift charts

## Work

Describe the three charts under Planning in `design/system/flaiover-dashboard.md` § Views, and add rows for them to the table in `docs/users/flaiover.md` § Charts. Say what each shows, the `flai stats` values it reads (`metrics.md` § Claims and touches), and how to read it: held stories beside the in-progress limit, hold hours by reason, and drift as touches the planner or agent should have declared or dropped. Waits for T-0943, so that the docs describe the page as built.

## Done when

- [ ] Both documents describe the three charts and link the metrics they read, and their `updated` dates are today
- [ ] `flai check --strict` and the markdown lint pass

## Notes
