---
id: ADR-0222
title: "With plan_backlog_epics or plan_backlog_stories on, the orchestrator may ask the planner again for an epic or a story planned before when the operator has written on a thread on it, or on an item under it, since its newest planner run ended, and the run records the thread"
status: accepted
date: 2026-10-30
supersedes: []
superseded_by: []
refines: [ADR-0087, ADR-0119]
topics: [planning, orchestration]
---

# ADR-0222 With plan_backlog_epics or plan_backlog_stories on, the orchestrator may ask the planner again for an epic or a story planned before when the operator has written on a thread on it, or on an item under it, since its newest planner run ended, and the run records the thread

## Context

`flai plan --candidates` lists two kinds of epic and one kind of story ([ADR-0119](0119-with-plan-backlog-stories-on-the-orchestrator-asks-the-planner-for-each-backlog.md) §1):

- an epic in the backlog with no stories;
- an epic not done or cancelled whose stories are all done or cancelled, with one done;
- a story in the backlog that lacks a plan.

The orchestrator asks for the planner on what the list holds alone (`orchestratorPlans` in `flai/cmd/plan.go`, [ADR-0087](0087-flai-serve-runs-one-orchestrator-per-project-behind-the-orchestrate-host-action.md)). So it never asks again for an epic that has stories, or for a story planned once.

I-0188 records the cost. On TH-0579 the operator confirmed a new E-0024 story after S-0419, which only E-0024's planner can draft. The operator asked there that the orchestrator be allowed to start the planner on items already planned. The story waited for the operator to run `flai plan E-0024` by hand.

The question: when may the orchestrator ask the planner again for an item planned before, bounded so that it does not spend on planner runs without the operator's word?

## Decision

`flai plan --candidates` gains a third kind of candidate, *planned before, asked again*.

1. **Candidates.** An item is listed when all of these hold:
   - It is an epic not done, cancelled, or archived, or a story in the backlog or ready.
   - Neither of the other kinds lists it.
   - Its newest planner run on this host has ended.
   - An entry was written on a thread anchored on the item, or on an item under it, after that run ended. For an epic, the items under it are its stories and their tasks; for a story, its tasks.
   - That entry's author is one of the operator's names, as `threads.OperatorName` reads them: the manifest's owner, the story's owner, or `designer`.

   The candidate carries `thread`, the thread with the newest such entry. Its reason names the author, the thread, the item the thread is anchored on, the entry's time, and the run's end, such as `planned before; alex wrote on TH-0579 (S-0419) at 2026-10-21T12:40:26Z, after its planner ran until 2026-10-21T12:16:34Z`. The holds of ADR-0119 §1 stay: an item is left out while a planner runs for it or for its epic, or while its planner's question awaits the operator.
2. **Asking.** The orchestrator asks for such a candidate with the MCP tool `plan` or `flai plan`, as for the other kinds: an epic under `plan_backlog_epics`, a story under `plan_backlog_stories`. It reads the thread first and asks only when the operator's words ask for a change to the plan. Such a change is a story or task to add, split, merge, or drop, a changed outcome, or a changed figure. It never asks for a confirmation that changes nothing. `flai guard` is unchanged: it holds the tool to the permission by the ID's type and leaves the list to flai.
3. **The run.** The run's trigger is `orchestrator; TH-nnnn`, in `agents.json` and in the `- Trigger:` line of its entry in `wip/agents/planner.md`. A run of the other kinds keeps `orchestrator`. The planner's start prompt names the thread and says that the item was planned before, and that the operator wrote on that thread since. The planner reads the thread and plans what the operator asked for, as the convention says.
4. **The bound.** A run that starts makes every earlier entry older than its end once it ends. So each of the operator's entries can start at most one planner run, and only through the orchestrator's judgement.

## Consequences

- The operator's scope change on a thread reaches the planner without the operator running `flai plan` by hand. I-0188 closes.
- The orchestrator's spend on planner runs stays tied to the operator's written word: one run at most for each operator entry.
- A candidate the orchestrator judged to ask for no change stays listed until the item is planned again, or until a later entry or run changes it. The orchestrator judges it again each run that looks at the list. That costs the orchestrator's reading, not a planner run.
- The rule reads the threads under the item for each item with an ended planner run, each time the list is made.
- A story whose own words changed after its plan, with no entry from the operator, is not a candidate of this kind. The replanner plans it again on the edit behind the `plan` host action ([ADR-0084](0084-flai-serve-plans-again-on-its-own-behind-the-plan-host-action-on-an-edit-when.md)).
- The first planning of an epic with stories, or of a story that already has a plan, stays the operator's. This kind needs a planner run before it.

## Alternatives considered

- **Let the orchestrator ask for any open epic or story.** Lost: nothing bounds its spend, and the first planning of an item planned by hand stays the operator's.
- **Trigger on an edit of an epic's own words.** Lost: it misses I-0188's instance, whose words were on a story's thread, and the replanner already plans a story again on an edit.
- **Start the planner on every operator entry, without the orchestrator.** Lost: it spends on every confirmation and every question, with no judgement of whether the words ask for a change.
