---
id: T-0145
type: task
nature: improvement
title: "ADR-0021 and living design: ready without tasks, tasks written in-progress, required at review"
status: done
parent: S-0049
owner: alex
created: 2026-09-18T18:08:29Z
updated: 2026-09-18T18:11:03Z
transitions:
  - to: ready
    at: 2026-09-18T18:10:06Z
    by: alex
  - to: in-progress
    at: 2026-09-18T18:10:06Z
    by: alex
  - to: done
    at: 2026-09-18T18:11:03Z
    by: alex
stream: S-0049
tags: []
touches: [design/adrs, design/system]
---

# T-0145 ADR-0021 and living design: ready without tasks, tasks written in-progress, required at review

## Work
Write `design/adrs/0021-story-ready-without-tasks.md` from `0000-template.md` and add it to the index in `design/adrs/README.md`. Decision: a story is ready with a goal, acceptance criteria, and a parent epic that is not cancelled; the agent that pulls it moves it to in-progress, then writes the tasks; a story cannot reach review without at least one task. Update `design/system/workflow.md` (pull policy steps, definition of ready, transitions table), `design/system/work-hierarchy.md` (the ready constraint becomes a review constraint), and the "Defer detail" principle in `design/system/overview.md`. Link the ADR from each.

## Done when
- ADR-0021 exists, is `accepted`, and is in the ADR index
- `workflow.md`, `work-hierarchy.md`, and `overview.md` state the new rule and nowhere say ready needs a task
- `flai check --strict` is clean

## Notes
