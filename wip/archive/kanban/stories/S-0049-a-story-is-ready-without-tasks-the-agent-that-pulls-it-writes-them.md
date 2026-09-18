---
id: S-0049
type: story
nature: improvement
title: "A story is ready without tasks; the agent that pulls it writes them"
status: done
parent: E-0001
owner: alex
created: 2026-09-18T18:01:55Z
updated: 2026-09-18T18:43:03Z
transitions:
  - to: ready
    at: 2026-09-18T18:09:10Z
    by: alex
  - to: in-progress
    at: 2026-09-18T18:09:15Z
    by: alex
  - to: review
    at: 2026-09-18T18:20:59Z
    by: alex
  - to: done
    at: 2026-09-18T18:43:03Z
    by: alex
tags: [cli, conventions]
touches: [flai/internal/workitem, flai/internal/check, template, design/conventions]
---

# S-0049 A story is ready without tasks; the agent that pulls it writes them

## Goal
The designer marks a story ready once its goal and acceptance criteria are written, without having to write tasks. The agent that pulls the story reads it, breaks it into tasks, and starts; when it cannot, it says why and asks instead of guessing.

## Acceptance criteria
- [x] `flai move <story> ready` succeeds for a story with a filled `## Goal`, acceptance criteria with at least one checkbox, and a parent epic that is not cancelled, with no tasks; the same move from the dashboard board succeeds
- [x] `flai move <story> in-progress` succeeds for a ready story with no tasks, from the CLI, the MCP `item_move` tool, and the dashboard board; nothing on the path into `ready` or `in-progress` refuses or warns about missing tasks
- [x] The task requirement moves later instead of disappearing: a story cannot reach `review` with no tasks, `flai move` refuses it with a message naming the rule, and `flai check` raises `story.tasks` for `review` and `done` stories only, so a `ready` or `in-progress` story with no tasks passes `flai check --strict`
- [x] The conventions tell the pulling agent what to do, in this order: move the story to `in-progress` and open the narrative; read the goal, criteria, and notes; create the tasks with `flai task new`, each with `## Work` and `## Done when`; then work them in order. Stated in `work-management.md` and `session-start.md` where pulling is described
- [x] The conventions say what "unable" means and what happens then: the agent does not invent scope; it records the blocker with `flai block --reason`, opens a thread on the story with `thread_open` stating what is missing, and moves to other work
- [x] "Do not write tasks for backlog stories" is reconciled with the new rule so the two read as one policy: tasks are written by whoever starts the story, at the time it is started
- [x] `design/system/workflow.md` (definition of ready, pull policy, transitions table) and `design/system/work-hierarchy.md` (the ready constraint) are updated first, then the baseline in `template/root/design/conventions/`, then this repository's copy, in this story; an ADR records the change to the definition of ready
- [x] `docs/users/flai.md` and the flaiover guide no longer say ready needs a task; the refused-move message on the board matches the new rule
- [x] Tests: `rules.go` and `check.go` cases for ready without tasks, review without tasks, and the unchanged criteria rule; `flai check --strict` clean on this repository and on a project rendered from the template

## Tasks
- T-0145 ADR-0021 and living design: ready without tasks, tasks written in-progress, required at review
- T-0146 flai move: allow ready and in-progress without tasks, refuse review without tasks
- T-0147 flai check: story.tasks applies to review and done only
- T-0148 Conventions: the pulling agent writes the tasks; template baseline first, then this repository
- T-0149 User docs and dashboard wording for the new ready rule
- T-0150 Verify all tiers and strict check on this repository and a rendered template

## Notes
Raised by the operator on 2026-09-18 when `flai move S-0048 ready` was refused with "a story needs at least one task before it is ready". The operator's preference: the agent reads the story and creates the tasks unless it is unable.

Where the rule lives today: `flai/internal/workitem/rules.go` (the `Ready` case in `Move`), `flai/internal/check/check.go` (`story.tasks`, every status except backlog and cancelled), `design/system/workflow.md` (definition of ready: "Tasks exist, each with `## Work` and `## Done when`"), `design/system/work-hierarchy.md` ("A story cannot be `ready` without at least one task"), `work-management.md` in the baseline and the template, and `docs/users/flai.md`. ADR-0004 fixes the states and the transition log, not the definition of ready, so nothing is superseded; a new ADR is enough.

Decided by the operator on 2026-09-18: the agent writes the tasks after moving the story to `in-progress`, and a story with no tasks must not be blocked from `ready` or from `in-progress`. Breaking the story down is work on it, so it counts towards cycle time. Today only the `Ready` case in `rules.go` refuses a move for missing tasks; nothing refuses `in-progress`, but `check.go` raises `story.tasks` for every status except backlog and cancelled, which would fail `flai check --strict` on a freshly pulled story, so that is the second place to change. Cycle time is defined in `design/system/metrics.md`; the definition does not change, only what falls inside it, so no metrics ADR is needed. The ADR for this story records both halves: ready without tasks, and tasks written in-progress with the requirement enforced at `review`.

This edits a baseline convention, which an agent does not do on its own; it is done here because the operator asked for it and it ships as a template release.

Until this lands, S-0048 and any other story still needs a task before `flai move` lets it into ready.

Verification, 2026-09-18: the CLI path was walked end to end in a project rendered from the template (ready, in-progress, `flai check --strict` clean at both, refused at review, accepted at review with a task). The MCP path is covered by `TestStoryWithoutTasksMovesUntilReview` in `flai/internal/mcpserver`. The dashboard path is covered by a new case in `flaiover/src/lib/server/flai.test.ts`, which drives the same `flai move` call the board's move endpoint makes; a card was not dragged in a browser. The running MCP server in this session is the installed flai 1.2.1 and keeps the old rule until the release from this story is installed.
