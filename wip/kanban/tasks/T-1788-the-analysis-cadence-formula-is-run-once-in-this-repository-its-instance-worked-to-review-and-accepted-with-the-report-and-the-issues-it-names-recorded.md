---
id: T-1788
type: task
nature: improvement
title: The analysis-cadence formula is run once in this repository, its instance worked to review and accepted, with the report and the issues it names recorded
status: backlog
parent: S-0424
owner: alex
created: 2026-10-16T05:49:09Z
updated: 2026-10-21T16:12:44Z
transitions: []
stream: S-0424
tags: [cli]
after: [T-1786]
---
# T-1788 The analysis-cadence formula is run once in this repository, its instance worked to review and accepted, with the report and the issues it names recorded

## Work

With the `formula` and `analyze` host actions on, run `flai formula run analysis-cadence --input focus=<focus> --ready` in the main checkout, with the focus the operator chose when finalizing the story, or none. The instance is worked by a story's agent on its branch and worktree as any story is (TH-0471, answer 5): it starts the analyzer step's role through `role_run` and waits for the run, then runs the listing step. Watch the instance reach review, and have it accepted. Record in this story's Notes:

- the instance's ID;
- the report it wrote;
- the issues it names;
- the run's cost.

Waits for the cadence formula. Runs beside the import run.

## Done when

The instance is accepted and archived, and this story's Notes record the run.

## Notes

The analyzer writes its report and issues in the main checkout itself; this task touches no file of the story.

Rewritten on 2026-10-21 by the planner of S-0424 for TH-0471's answer 5, a story's agent working the instance in place of `flai serve`, as approved on TH-0588.
