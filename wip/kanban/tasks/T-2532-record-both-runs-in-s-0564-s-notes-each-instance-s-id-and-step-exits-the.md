---
id: T-2532
type: task
nature: improvement
title: "Record both runs in S-0564's Notes: each instance's ID and step exits, the scratch repository, the analyzer's report, the issues it names, and each run's cost"
status: backlog
parent: S-0564
owner: alex
created: 2026-11-07T03:48:03Z
updated: 2026-11-07T03:48:03Z
transitions: []
stream: S-0564
tags: [template]
touches: [wip/kanban/stories/S-0564-the-import-and-analysis-cadence-formulas-are-each-run-once-in-this-repository.md]
after: [T-2523, T-2526]
---
# T-2532 Record both runs in S-0564's Notes: each instance's ID and step exits, the scratch repository, the analyzer's report, the issues it names, and each run's cost

## Work

Write the record of both runs in S-0564's Notes with `flai edit`, under a heading of its own below `### Planning`:

- **The import run.** The scratch repository's path, the new project's key, the instance's ID, and each step's exit with the tests `flai import test` ran. Take them from T-2514's Notes and the instance's Notes.
- **The cadence run.** The instance's ID, its focus, the report the analyzer wrote under `design/analysis`, and the issues the report names, each with its ID, title, and class. Take them from T-2520's Notes and the instance's Notes.
- **Each run's cost.** The usage each instance carries after acceptance, from `flai show <instance>`. The analyzer's cost is charged to the cadence instance (ADR-0239, ADR-0245).

Then ask the operator on a thread on S-0564 whether to remove the scratch repository and stop serving its project. Remove nothing without the operator's word. Write the answer in the Notes.

This task waits for T-2523 and T-2526: both instances must be accepted first. It runs beside the task that files the runs' faults, with which it shares no path.

## Done when

- S-0564's Notes name both instances, their step exits, the scratch repository, the report, the issues it names, and each run's cost.
- The operator's word on the scratch repository is asked for, and the answer is in S-0564's Notes.

## Notes
