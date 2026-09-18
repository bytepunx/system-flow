---
id: ADR-0021
title: A story is ready without tasks; the pulling agent writes them and review requires them
status: accepted
date: 2026-09-18
supersedes: []
superseded_by: []
---

# ADR-0021 A story is ready without tasks; the pulling agent writes them and review requires them

## Context

The definition of ready in [workflow.md](../system/workflow.md) required a story to have at least one task, each with `## Work` and `## Done when`, and `flai move` and `flai check` enforced it. The designer writes the goal and the acceptance criteria; breaking a story into tasks is implementation planning, and the agent that will do the work is better placed to do it, with the code open, at the moment the work starts. The rule made the designer write tasks, or ask an agent to, before a story could enter `ready`, which put a planning step in front of the queue and produced task detail that was rewritten when the work began.

The rule also had a purpose: a story with no tasks leaves no record of how the work was divided and nothing for the narrative's task transitions to hang on. That purpose does not need the tasks to exist before the work starts, only before it is presented as finished.

[ADR-0004](0004-workflow-states-and-transitions.md) fixes the states and the transition log. It does not fix the definition of ready, so it stands unchanged.

## Decision

A story is ready when its `## Goal` is filled, its `## Acceptance criteria` has at least one checkbox, and its parent epic is not cancelled. Tasks are not required to enter `ready` or `in-progress`. The agent that pulls the story moves it to `in-progress`, opens the narrative, reads the story, and then writes the tasks, each with `## Work` and `## Done when`. A story cannot enter `review` without at least one task; `flai move` refuses the transition and `flai check` raises `story.tasks` for `review` and `done` stories.

When the agent cannot write the tasks because the goal or criteria do not say enough, it does not invent scope. It blocks the story with the reason, opens a thread on the story saying what is missing, and moves to other work.

## Consequences

- The designer moves a story to `ready` from the board or the CLI as soon as the goal and criteria are written.
- Task breakdown happens inside `in-progress`, so it counts towards cycle time. The definition of cycle time in [metrics.md](../system/metrics.md) does not change; what falls inside it does, by the minutes a breakdown takes.
- A `ready` or `in-progress` story with no tasks is valid and passes `flai check --strict`. A story that stays `in-progress` with no tasks for long is visible on the board by its age, not by a check finding.
- The pull step in the conventions gains a task-writing step, and the agent needs a stated way out when it cannot write them.
- Anyone may still write tasks earlier; nothing refuses a `ready` story that already has them.

## Alternatives considered

- Keep the rule and have an agent refine each story to `ready` on request: the step the designer asked to remove.
- Write the tasks while the story is `ready`, before moving it to `in-progress`: keeps every in-progress story with tasks, but books the breakdown as queue time and needs the agent to edit a story it has not yet pulled.
- Drop the task requirement entirely: a story could be delivered with no record of how the work was divided, and task transitions would stop being a signal of progress.
- Require tasks at `in-progress` plus a grace period in `flai check`: a time-based rule in a checker that otherwise reads only state, for no gain over enforcing at `review`.
