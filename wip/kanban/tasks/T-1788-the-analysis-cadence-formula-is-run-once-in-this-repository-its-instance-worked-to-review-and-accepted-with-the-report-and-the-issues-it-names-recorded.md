---
id: T-1788
type: task
nature: improvement
title: The analysis-cadence formula is run once in this repository, its instance worked to review and accepted, with the report and the issues it names recorded
status: backlog
parent: S-0424
owner: alex
created: 2026-10-16T05:49:09Z
updated: 2026-10-16T05:49:09Z
transitions: []
stream: S-0424
tags: [cli]
after: [T-1786]
---
# T-1788 The analysis-cadence formula is run once in this repository, its instance worked to review and accepted, with the report and the issues it names recorded

## Work

With the `formula` and `analyze` host actions on, run `flai formula run analysis-cadence --input focus=<focus>` in the main checkout, with the focus the operator chose when finalizing the story, or none. Watch the instance's analyzer step and its listing step reach review, and have it accepted. Record in this story's Notes:

- the instance's ID;
- the report it wrote;
- the issues it names;
- the run's cost.

Waits for the cadence formula. Runs beside the import run.

## Done when

The instance is accepted and archived, and this story's Notes record the run.

## Notes

The analyzer writes its report and issues in the main checkout itself; this task touches no file of the story.
