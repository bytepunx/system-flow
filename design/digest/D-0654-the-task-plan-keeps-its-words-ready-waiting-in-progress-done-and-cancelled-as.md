---
id: D-0654
title: "The task plan keeps its words ready, waiting, in-progress, done, and cancelled as its own JSON vocabulary, not state names: each is chosen by the kind of the task's state, and renaming them would change flai's JSON"
kind: decision
sources: [S-0739]
topics: [flai, go, code, planning]
paths: [flai/internal/workitem/taskplan.go]
date: 2026-11-17
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0654 The task plan keeps its words ready, waiting, in-progress, done, and cancelled as its own JSON vocabulary, not state names: each is chosen by the kind of the task's state, and renaming them would change flai's JSON

S-0739 moved the task plan (S-0176) off state names and onto a state's kind (ADR-0301). Its criterion 1 asks for no comparison with a state's name, and its narrative's Decisions (`wip/archive/agents/S-0739.md`) record why `taskplan.go`'s constants stay:

> taskplan.go keeps its plan vocabulary, PlanReady, PlanWaiting, PlanInProgress, PlanDone, PlanCancelled and TaskSummary's fields: they are the task plan's own JSON words (S-0176), not workflow states, and each is now chosen by the kind of the state the task is in. Renaming them would change flai's JSON and its callers; criterion 1 is read as no comparison with a state's name.

In `flai/internal/workitem/taskplan.go`, `ready` is a task in a queue state that waits for no open task, `waiting` one in a queue state whose `after:` names an open task, `in-progress` one in an active or a waiting state, `done` one in the accepted state, and `cancelled` one in the cancelled state. A later E-0034 story, or the dashboard's task plan (S-0758), reads these values as the plan's words. It does not rename them after a project's states, and it does not take a project's state name for a plan word.

Source: S-0739's narrative, `wip/archive/agents/S-0739.md`, § Decisions.
