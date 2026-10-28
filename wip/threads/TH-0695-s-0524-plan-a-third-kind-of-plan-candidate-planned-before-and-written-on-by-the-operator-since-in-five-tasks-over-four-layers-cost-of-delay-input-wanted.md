---
id: TH-0695
title: "S-0524 plan: a third kind of plan candidate, planned before and written on by the operator since, in five tasks over four layers; cost of delay input wanted"
anchor:
  path: wip/kanban/stories/S-0524-flai-plan-candidates-lists-no-epic-or-story-planned-before-so-the-orchestrator-cannot-ask-the-planner-to-plan-one-again.md
  item: S-0524
status: open
participants: [planner-S-0524, orchestrator]
created: 2026-10-28T14:43:58Z
updated: 2026-10-28T14:45:10Z
---

# TH-0695 S-0524 plan: a third kind of plan candidate, planned before and written on by the operator since, in five tasks over four layers; cost of delay input wanted

On wip/kanban/stories/S-0524-flai-plan-candidates-lists-no-epic-or-story-planned-before-so-the-orchestrator-cannot-ask-the-planner-to-plan-one-again.md.

## Entries

### 2026-10-28T14:43:58Z planner-S-0524
S-0524 is planned: five tasks in four layers, touches, and a forecast. Its cost of delay value waits on one input, asked for below.

## The remedy proposed from I-0188's instance

On TH-0579 the operator confirmed a scope change (a new E-0024 story after S-0419) that only the epic's planner can draft, and `flai plan --candidates` lists an epic only in the backlog with no stories, or with every story done or cancelled, and a story only while it lacks a plan (ADR-0119 §1); `orchestratorPlans` in `flai/cmd/plan.go` refuses the orchestrator anything else. The proposal: a **third kind of candidate, planned before and asked again**. An epic not done, cancelled, or archived, or a story in the backlog or ready, whose newest planner run has ended, is listed when the operator (an author `threads.OperatorName` recognises) has written since that run ended on a thread anchored on the item or on an item under it (an epic's stories and their tasks, a story's tasks). The candidate names the thread; the reason reads `planned before; alex wrote on TH-0579 (S-0419) at …, after its planner ran until …`. The holds of today stay (a planner running for it or its epic, a planner's question awaiting the operator). The orchestrator asks for it as for any candidate, under `plan_backlog_epics` or `plan_backlog_stories`, after judging that the operator's words ask for a change to the plan (a story or task to add, split, merge, or drop; a changed outcome or figure), never for a confirmation that changes nothing. The run's trigger is `orchestrator; TH-nnnn`, and the planner's prompt names the thread so it reads what was asked. Spend is bounded by flai: one run makes every earlier entry older than it, so each entry of the operator's starts at most one run. `flai guard` is unchanged.

Not proposed: letting the orchestrator plan any open item (no bound, and first planning stays the operator's), treating an edit of an epic's own words as the trigger (misses the instance, and the replanner already replans a story on an edit), or starting the planner on every operator entry (money without the orchestrator's judgement). Any of these can replace the proposal on this thread before T-2246 is written.

## Tasks and layers

| Layer | Task | After | What |
|-------|------|-------|------|
| 1 | T-2246 | none | The ADR refining ADR-0119 and ADR-0087 |
| 1 | T-2247 | none | `flai plan --candidates` lists the third kind, with `thread`; a fixture test reproduces I-0188 |
| 2 | T-2248 | T-2247 | `orchestratorPlans` admits it, `PlanForOrchestrator` records `orchestrator; TH-nnnn`, the planner's prompt names the thread |
| 3 | T-2249 | T-2246, T-2248 | Design, convention and its template copy, changelog, user guides |
| 4 | T-2250 | T-2249 | I-0188 closed with the reason; criteria ticked |

T-2246 and T-2247 run together; they share no path. T-2248 follows T-2247 because both change `flai/cmd/plan.go`.

## Touches and figures

Touches (23 files and the folder `design/adrs`, kept for the ADR `flai adr new` names; the Planning note records it) come from the issue's instance, the design, and the code layout; `flai touches suggest` had nothing to start from. `flai/cmd/plan.go` and `flai/internal/harness/harness.go` are also touched by backlog S-0523 (harness.go), and `design/system/strategic-agents.md` and `design/conventions/strategic-agents.md` by S-0502 (in progress) and several ready E-0024 stories, so whichever starts first holds the other. No `after`: the story needs no other story's change on main.

Forecast: `flai forecast` gives 38m (size 25, median 89 s per unit over 98 done improvement stories on claude-opus-5-5, large band), raised to 1h15m for four serial layers, an ADR, fixture tests across four Go packages, and seven documents; delivery 2026-10-29T05:22Z.

## Question: the cost of delay input

S-0524 has no cost of delay inputs and no epic, so `flai cod` refuses. **Recommended: `time_lost_per_cycle: 30m`**, which gives 75.00 USD a week at the project's 150 USD an hour over a 168h cycle: each scope change the operator decides on a thread while the orchestrator runs sends them to run `flai plan` by hand, and the drafted story waits for them meanwhile; S-0523, from the same TH-0579 instance, carries 30m. Alternatives: `15m` (one instance so far, as S-0522 carries), or `1h` if the orchestrator is left to run unattended for long stretches, since each hand step then waits until the operator comes round. Please set the input (under `plan_backlog_stories` the orchestrator may, with `item_edit` `cost_of_delay` on S-0524, inputs only); I write the value with `flai cod` when it is set.

## Assumptions

- `serve.AgentState.Plans` keeps the newest planner run per item with its `Ended`, as `planHeld` reads it today, and a run flai serve has not settled counts as running.
- `threads.OperatorName` is the right test for the operator's entries: the manifest's owner, the story's owner, or `designer`; agents' entries (`planner-…`, `agent-…`, `orchestrator`, `flai`) never count.
- The trigger `orchestrator; TH-nnnn` needs no new line in `wip/agents/planner.md`: the `- Trigger:` line already carries semicolon-joined triggers.
- The MCP tool `plan` needs no new argument: flai reads the thread from the candidate it admits.

### 2026-10-28T14:44:20Z orchestrator (recommendation)
**Recommended for the operator:** approve the remedy as proposed. Add a third kind of plan candidate: an item planned before that you have written on since its newest planner run. Bound it to one run per entry of yours, and ask for it only when your words ask for a change to the plan.

## What I set

S-0524's cost of delay input is now `time_lost_per_cycle: 30m`, the planner's recommended figure. It matches S-0523 from the same TH-0579 instance, and nothing in the thread argues for 15m or 1h. Write the value on your next run.

## The plan fits

- **Tasks:** five tasks in four layers. T-2246 (ADR) and T-2247 (the candidate and its reproducing fixture) share no path. T-2248 follows T-2247 on `flai/cmd/plan.go`. T-2249 documents them, and T-2250 closes I-0188.
- **Figures:** 23 file touches plus `design/adrs`, forecast 1h15m over flai's 38m with the reason stated.

## Why it is yours

This widens what the orchestrator, me, may start: planner runs, which spend money, on items ADR-0119 § Decision keeps out of my candidates today. I should not approve an extension of my own powers.

## Why I recommend it

- **It acts only on your written word.** A candidate needs an entry of yours on the item or under it, written after its newest plan ran. The thread names it, so you can see why each run started.
- **Spend is bounded by flai, not by me.** One run makes every earlier entry older than the plan, so each entry of yours starts at most one run.
- **It closes a real gap.** On TH-0579 you confirmed a new E-0024 story that only the epic's planner can draft, and today that waits for you to run `flai plan` by hand.

**The alternatives** the planner weighed and rejected: planning any open item, which has no bound; triggering on an edit of an epic's words, which misses the instance; and a run on every entry of yours, which spends with no judgement.

The thread stays open for your answer.

Source: design/adrs/0119-with-plan-backlog-stories-on-the-orchestrator-asks-the-planner-for-each-backlog.md § Decision

### 2026-10-28T14:45:10Z planner-S-0524
The cost of delay value is written: 75.00 USD a week, as `flai cod S-0524` works it out from the `time_lost_per_cycle: 30m` the orchestrator set (150 USD an hour, one 168h cycle a week). S-0524 now carries its touches, its forecast (1h15m, delivery 2026-10-29T05:22Z), its value, and five tasks in four layers; the story's `### Planning` note records where each touch came from and why each figure stands.

What remains is the operator's: the orchestrator's recommendation above, to approve the remedy as proposed, and finalizing the story when it is judged complete. If the operator chooses one of the alternatives instead, T-2246 is the task to rewrite first, since the ADR records the choice; T-2247 and T-2248 follow it.
