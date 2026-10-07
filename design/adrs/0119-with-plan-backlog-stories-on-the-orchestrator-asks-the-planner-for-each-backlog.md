---
id: ADR-0119
title: "With plan_backlog_stories on, the orchestrator asks the planner for each backlog story that lacks a plan, answers and resolves the threads a story's planner opens, and sets a backlog story's cost of delay"
status: accepted
date: 2026-10-07
supersedes: []
superseded_by: []
refines: [ADR-0087]
---

# ADR-0119 With plan_backlog_stories on, the orchestrator asks the planner for each backlog story that lacks a plan, answers and resolves the threads a story's planner opens, and sets a backlog story's cost of delay

## Context

[ADR-0087](0087-flai-serve-runs-one-orchestrator-per-project-behind-the-orchestrate-host-action.md) gives the orchestrator permissions, each off by default. `plan_backlog_epics` lets it ask the planner for an epic `flai plan --candidates` lists, and for nothing else. The MCP tool `plan` refuses it a story. In the backlog, stories without touches, a forecast, a cost of delay value, or tasks then wait for the operator to ask the planner for each. Its threads TH-0252 to TH-0313 recommended exactly that, one story at a time. When the planner asks for a cost of delay input, the answer is the operator's too: under `answer_threads` set to `autonomous`, a cost of delay input is the operator's judgement, so the orchestrator only recommends. S-0328 asks for a permission that hands the orchestrator this loop.

## Decision

`orchestration.permissions.plan_backlog_stories`, a boolean, off by default, lets the orchestrator plan the backlog's stories and settle what their planner asks.

1. **Candidates.** `flai plan --candidates` lists stories after the epics: each story in the backlog, not archived, that lacks a plan. A story lacks a plan when it has no touches, no forecast duration, no cost of delay value, or no task that is not cancelled, and the reason names what it lacks. Left out, with why: a story a planner runs for now, or whose epic a planner runs for now; one whose newest planner run ended asking on a thread whose last word is still the planner's; and one whose newest planner run ended and which has not changed since, so that a planner that could not finish is not started again until something changes.
2. **Planning.** With the permission on, the orchestrator starts the planner on a story the candidates list, one at a time, with the MCP tool `plan` or `flai plan`. flai refuses it any other story (`orchestratorPlans`), and the guard refuses it every story while the permission is off.
3. **A story planner's threads.** On a thread a story's planner opened, its opener `planner-S-nnnn`, the orchestrator may reply as an answer, with or without a source, or as a recommendation, and resolve the thread. That covers approving the plan the planner posts, answering its questions, and choosing a cost of delay input, whether by accepting the figure the planner recommends or one of its alternatives. On every other thread `answer_threads` still holds. Confirming a recommendation stays the operator's alone ([ADR-0090](0090-a-thread-entry-is-marked-a-recommendation-in-its-heading-and-cites-its-source.md)).
4. **Cost of delay.** The orchestrator may set a backlog story's cost of delay, its inputs and its value, with `item_edit`'s `cost_of_delay` or `flai edit`'s cost of delay flags, and change nothing else of the story in that call. flai holds this as well as the guard, so a call the guard does not see is held all the same.

## Consequences

- An operator who turns the permission on hands the orchestrator a money decision: the cost of delay inputs of backlog stories. Its risk on the settings page says so.
- The planner runs, and spends, on every unplanned backlog story, one at a time; drafts are planned as well as finalized stories.
- Each answer, resolution, and cost of delay set is logged with `activity_log`, so the operator can review what the orchestrator chose.
- Threads that a story's planner opens no longer wait for the operator while the permission is on.

## Alternatives considered

- Widening `plan_backlog_epics` to stories: an operator who allowed epics would silently allow stories, and the cost of delay decision with them.
- Leaving the cost of delay inputs to the operator under this permission: the planner cannot value a story without them, so the loop would still stop on the operator for every story without an epic.
- Letting the orchestrator answer every planner thread, epics' included: the permission is about stories; an epic planner's threads stay under `answer_threads`.
