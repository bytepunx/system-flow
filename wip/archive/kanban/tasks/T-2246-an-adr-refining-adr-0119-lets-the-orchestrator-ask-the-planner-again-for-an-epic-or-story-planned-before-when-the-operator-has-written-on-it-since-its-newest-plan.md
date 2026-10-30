---
id: T-2246
type: task
nature: improvement
title: An ADR refining ADR-0119 lets the orchestrator ask the planner again for an epic or story planned before when the operator has written on it since its newest plan
status: done
parent: S-0524
owner: alex
created: 2026-10-28T14:41:09Z
updated: 2026-10-30T13:09:40Z
transitions:
  - to: ready
    at: 2026-10-30T13:08:48Z
    by: agent-S-0524
  - to: in-progress
    at: 2026-10-30T13:08:51Z
    by: agent-S-0524
  - to: done
    at: 2026-10-30T13:09:40Z
    by: agent-S-0524
stream: S-0524
tags: [flai, orchestrator, planner, adr]
touches: [design/adrs/README.md, design/adrs/0222-with-plan-backlog-epics-or-plan-backlog-stories-on-the-orchestrator-may-ask-the.md]
usage:
  source: log
  seconds: 49
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 5
      output: 1823
      cache_read: 381248
      cache_write: 16458
      cost: 0.2087
---
# T-2246 An ADR refining ADR-0119 lets the orchestrator ask the planner again for an epic or story planned before when the operator has written on it since its newest plan

## Work

First layer: waits for nothing. Write the ADR with `flai adr new` (or the MCP tool `adr_new`), refining ADR-0119 and ADR-0087, that records the remedy S-0524's plan thread proposes for I-0188:

- Today `flai plan --candidates` lists an epic only in the backlog with no stories, or open with every story done or cancelled, and a story only while it lacks a plan (ADR-0119 §1). The orchestrator may ask the planner for what the list holds alone (`orchestratorPlans` in `flai/cmd/plan.go`), so an epic with stories, or a story once planned, is never planned again at its asking, and a scope change the operator decided on a thread (I-0188: TH-0579, a new E-0024 story after S-0419) waits for the operator to run `flai plan` by hand.
- The decision: a third kind of candidate, *planned before, asked again*. An epic that is not done, cancelled, or archived, or a story in the backlog or ready, whose newest planner run has ended, is listed when the operator has written since that run ended on a thread anchored on the item or on an item under it (an epic's stories and their tasks, a story's tasks). The operator's entries are those `threads.OperatorName` recognises: the manifest's owner, the story's owner, or `designer`. The reason names the thread and the entry's time, such as `planned before; alex wrote on TH-0579 (S-0419) at 2026-10-21T12:40:26Z, after its planner ran until 2026-10-21T12:16:34Z`. The holds of ADR-0119 stay: a planner running for it or its epic, or its planner's question still awaiting the operator.
- The orchestrator asks for such a candidate with the MCP tool `plan` or `flai plan` as for any other, under `plan_backlog_epics` for an epic and `plan_backlog_stories` for a story, and judges first whether the operator's words ask for a change to the plan: a story to add, split, merge, or drop, a changed outcome, or a changed figure. A confirmation that changes nothing is not planned again. flai bounds the spend: one planner run makes every earlier entry older than the run, so each entry of the operator's can start at most one run.
- The run records its trigger as `orchestrator; TH-nnnn`, in `serve/agents.json` and the `- Trigger:` line of `wip/agents/planner.md`, and the planner's prompt names the thread, so the planner reads what the operator asked for.
- Alternatives considered, recorded in the ADR: letting the orchestrator ask for any open epic or story (no bound on spend, and the first planning of an item stays the operator's as Planning again says); treating an edit of an epic's own words as the trigger (misses the instance, where the words were on a story's thread, and the replanner already plans a story again on an edit); starting the planner on its own on every operator entry (money spent without the orchestrator's judgement).

`flai adr new` writes the file and its row in `design/adrs/README.md`; the ADR's file is the story's `design/adrs` folder touch, which `flai task done` widens with the committed path. Give the ADR `topics: [planning, orchestration]` with `flai adr topics`.

## Done when

- A new ADR under `design/adrs/` with status `accepted`, `refines: [ADR-0119, ADR-0087]`, and the decision, context, consequences, and alternatives above, indexed in `design/adrs/README.md`.
- `flai check --strict` passes on the ADR and the index.

## Notes

The ADR's number is allotted by `flai adr new` when it runs, so this task names `design/adrs/README.md` alone and the story keeps the `design/adrs` folder touch for the file.
