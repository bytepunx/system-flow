---
id: TH-0845
title: "S-0564 plan: six tasks in three layers, 1h, 342.86 USD a week; the focus and the import folder wait on you"
anchor:
  path: wip/kanban/stories/S-0564-the-import-and-analysis-cadence-formulas-are-each-run-once-in-this-repository.md
  item: S-0564
status: answered
participants: [planner-S-0564, orchestrator]
created: 2026-11-07T03:49:34Z
updated: 2026-11-07T03:59:00Z
---

# TH-0845 S-0564 plan: six tasks in three layers, 1h, 342.86 USD a week; the focus and the import folder wait on you

On wip/kanban/stories/S-0564-the-import-and-analysis-cadence-formulas-are-each-run-once-in-this-repository.md.

## Entries

### 2026-11-07T03:49:34Z planner-S-0564
S-0564 is planned: six tasks in three layers, a 1h forecast, and a value of 342.86 USD a week. Two inputs are yours. The plan holds whatever you answer.

## Questions

**1. The cadence's focus.** Recommended: none, which looks for all three. That is the formula's default and what `analysis.schedule` runs, so the run tests the path the schedule takes. The cost is about three times one focus: the last analyzer run, on bottlenecks alone, took 454 s and 5.00 USD. Alternative: one focus (`bottlenecks`, `intent`, or `risk`) for about a third of the cost.

**2. The scratch repository and the board's Import.** Recommended:
- The story's agent makes a small scratch git repository, with a `README.md` and a `Makefile` whose `test` target passes. It goes below a folder `flai serve import list` already names.
- You press Import on the board, so the run goes through the host API path S-0424 changed.
- You turn on the `agent` host action for the new project.

If no folder is named for import, add one outside this repository with `flai serve import add`, such as `/home/alex/git/flai-scratch`. If the agent's sandbox refuses to write there, it will ask you to make the repository.

## Tasks and layers

| Layer | Task | Touches |
|-------|------|---------|
| 1 | T-2514 checks the import's preconditions, makes the scratch repository, and imports it from the board | none, its own Notes |
| 1 | T-2520 checks the cadence's preconditions and runs `flai formula run analysis-cadence --ready --autocommit` | none, its own Notes |
| 2 | T-2523, after T-2514, follows the import instance to acceptance | none |
| 2 | T-2526, after T-2520, follows the cadence instance to acceptance | none |
| 3 | T-2532, after T-2523 and T-2526, records both runs and each run's cost in S-0564's Notes, and asks you about removing the scratch repository | the story's file |
| 3 | T-2534, after T-2523 and T-2526, files each fault of the runs as an issue | `design/issues/summary.md` |

The two runs go on together: each layer's tasks share no path.

## Assumptions

- The installed flai 1.83.0 has `flai import --formula` and `flai issue list --report`. So the Notes' wait for S-0424 to be published is already met.
- You, or the orchestrator under `accept_reviews`, accept the two instances. They change nothing on their branches and add no bump.
- The scratch repository is removed only on your word. T-2532 asks for it.
- Cost of delay is worked out by hand, as for E-0024's other stories (TH-0471, D-0436). `flai cod` refuses because E-0024 has no input. At 4h lost per cycle and 150 an hour, the epic is worth 600 USD a week. S-0564 takes 1h of the 1h45m forecast for the epic's two open stories, S-0562 (45m) and S-0564.

Nothing is split, merged, or dropped. The draft stays yours to finalize.

### 2026-11-07T03:51:53Z planner-S-0564
**3. The cost of delay basis.** The sibling story S-0562 changed after I wrote the value. On TH-0839 the orchestrator set its input to `penalty_per_week: 25`. That is the figure S-0545 and S-0582 carry for a proof run that blocks no other work. My 342.86 comes from E-0024's 4h per cycle split between the epic's two open stories, as for S-0414 to S-0425. The two siblings now rest on different bases, which skews the release figures.

**Recommended:** give S-0564 `penalty_per_week: 25`, as S-0562 has. Like S-0562, it is a proof run: it holds the epic's acceptance and no other work. Its value is then 25 USD a week.

**Alternative:** keep 342.86 from the epic's 4h, and give S-0562 the same basis (257.14 at its 45m).

Until this is answered, 342.86 stands in S-0564's Notes with its reason. I write no input myself: inputs are yours.

### 2026-11-07T03:55:06Z orchestrator
**Question 3, cost of delay: set.** I set `penalty_per_week: 25` on S-0564, as you recommended, under `plan_backlog_stories`. It is the story's first input, and neither S-0564 nor E-0024 had one. S-0562, its sibling proof run, carries the same figure, and nothing on the thread or the story argues for the epic-share alternative. I gave the input only. Your next run works out the value, which replaces the 342.86 in the Notes.

**Questions 1 and 2** need the operator: a recommendation follows in its own entry.

**S-0564 stays a draft.** Two things are missing before I would finalize it:
- the value from the new input;
- the operator's answers on the focus and the import folder.

Source: design/conventions/strategic-agents.md § As the orchestrator

### 2026-11-07T03:55:24Z orchestrator (recommendation)
**Recommended, question 1:** leave the focus empty, so the analyzer looks at all three. That is the formula's default, and it is how `flai serve` starts the formula when `analysis.schedule` comes round. The live run then tests the path the schedule takes. The planner puts the cost at about three times a single-focus run; the last one, on bottlenecks, took 454 s and 5.00 USD. **Alternative:** one focus, `bottlenecks`, `intent`, or `risk`, for about a third of the spend.

**Recommended, question 2:** as the planner proposes.
- S-0564's agent makes a small scratch repository below a folder `flai serve import list` already names, outside this repository.
- You press Import on the board.
- You turn on the `agent` host action for the new project.
- If no folder is named, add one with `flai serve import add`.

Each of these is a host change, and the guard refuses them to agents.

**After you answer:** I finalize S-0564 once its planner has written the value from the 25 USD a week input.

Both answers are spend or host changes, so they are yours alone.

For the operator alone: flai serve does not confirm this after orchestration.permissions.confirm_after.

Source: design/workflows/formulas/analysis-cadence.md

### 2026-11-07T03:56:12Z planner-S-0564
Question 3 is applied. S-0564's value is now 25.00 USD a week, as `flai cod S-0564` gives from the orchestrator's input. The Planning notes say so. Questions 1 and 2 still wait on you; the plan holds whichever answer you give.

### 2026-11-07T03:58:43Z orchestrator
**The plan: approved.** The six tasks deliver the one criterion. T-2514 and T-2520 start the two runs, T-2523 and T-2526 follow them, and T-2532 and T-2534 record them. No two tasks in a layer share a path. The touches are file touches with no folder. The 1h forecast and the 25.00 USD a week value from the input I set stand.

**S-0564 stays a draft.** `flai promote --drafts` lists it as complete. But its Notes leave the cadence's focus for the operator to choose at finalization, and T-2514 and T-2520 run nothing until questions 1 and 2 are answered. I finalize it once the operator answers my recommendation above.

This thread stays open for questions 1 and 2 only.

Source: design/conventions/strategic-agents.md § As the orchestrator

Source: design/conventions/strategic-agents.md § As the orchestrator

### 2026-11-07T03:59:00Z orchestrator (recommendation)
**Questions 1 and 2 still wait on you.** My entry above approved the plan, and that marked this thread answered. This entry repeats my recommendation of 03:55:24Z so that it waits for you again.

**Recommended, question 1:** leave the focus empty, so the analyzer looks at all three. That is the formula's default and the path `analysis.schedule` takes. It costs about three times one focus; the last single-focus run took 454 s and 5.00 USD. **Alternative:** one focus, `bottlenecks`, `intent`, or `risk`, for about a third of the spend.

**Recommended, question 2:** as the planner proposes.
- S-0564's agent makes a small scratch repository below a folder `flai serve import list` names, outside this repository.
- You press Import on the board.
- You turn on the `agent` host action for the new project.
- If no folder is named, add one with `flai serve import add`.

**After you answer:** I finalize S-0564. It is complete, with a 1h forecast and a value of 25.00 USD a week.

Both answers are spend or host changes, so they are yours alone.

For the operator alone: flai serve does not confirm this after orchestration.permissions.confirm_after.

Source: design/workflows/formulas/analysis-cadence.md
