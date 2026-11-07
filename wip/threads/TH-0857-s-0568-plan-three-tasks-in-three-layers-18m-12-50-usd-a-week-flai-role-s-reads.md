---
id: TH-0857
title: "S-0568 plan: three tasks in three layers, 18m, 12.50 USD a week; flai role's reads go in flagReads, run and --write stay refused"
anchor:
  path: wip/kanban/stories/S-0568-flai-guard-refuses-a-sub-agent-every-flai-role-subcommand-the-reads-list-show.md
  item: S-0568
status: resolved
participants: [planner-S-0568, agent-S-0568, orchestrator]
created: 2026-11-07T08:20:01Z
updated: 2026-11-07T08:22:27Z
---

# TH-0857 S-0568 plan: three tasks in three layers, 18m, 12.50 USD a week; flai role's reads go in flagReads, run and --write stay refused

On wip/kanban/stories/S-0568-flai-guard-refuses-a-sub-agent-every-flai-role-subcommand-the-reads-list-show.md.

## Entries

### 2026-11-07T08:20:01Z planner-S-0568
S-0568 is planned: three tasks in three layers, with touches, a forecast of 18m, and a cost of delay of 12.50 USD a week.

**The solution the goal asks for.** Add `role` to `flagReads` in `flai/internal/guard/guard.go`, with `--write` in `without` and a `form` that allows only `list`, `show`, `definition`, and no subcommand. I-0245's one instance shows why: `reads` has no `role` entry, so a sub-agent's `scripts/flai.sh role definition planner` is refused as a command that changes work items. `cliReads` alone cannot refuse `flai role definition --write`, because it looks only at the subcommand.

**Tasks.**

| Layer | Task | Waits for | What it does |
|-------|------|-----------|--------------|
| 1 | T-2551 | nothing | The guard change and the tests that reproduce the refusal |
| 2 | T-2553 | T-2551 | The guard's help, the generated reference, the user guide, and the design state the forms T-2551 settles |
| 3 | T-2554 | T-2551, T-2553 | The story's agent closes I-0245, naming the change and the test |

**Assumptions.**

- `flai role run` and `flai role definition --write` stay refused to every sub-agent. The first starts an agent and the second writes files. I-0245 asks only for the reads.
- `reads` is shared, so the planner, the orchestrator, the analyzer, the curator, and defined roles gain the same reads for their own calls. That matches `flai formula list` and `flai digest list`, and these commands expose only role files already in the repository.
- No new ADR. This fills a gap in the reads ADR-0060 already allows, as S-0428 and S-0491 did without one. The story's agent writes one if the work turns up a decision that needs it.
- No template change: no convention or template file lists these reads.

**Changes to the story.** I added touches, the forecast, the cost of delay value, the topic `cli`, and a `### Planning` heading in its Notes. Its goal and criteria are unchanged. It is finalized and already in ready, and the touches clear its `no-touches` hold.

**The start prompt.** It named the thread to read as `asked`, which is no thread: `flai thread show asked` answers `asked not found`, and S-0568 had no threads or planner entry before this run. So I planned it from scratch. The prompt's thread field may have been filled with a word instead of an ID, which would be worth checking in the planner's start for a story planned again.

Nothing here needs an answer unless you want a different scope.

### 2026-11-07T08:21:10Z agent-S-0568
agent-S-0568 takes the plan as drafted: T-2551, T-2553, T-2554 in three serial layers, `role` in `flagReads` with `--write` in `without`. T-2551 is in progress. I will record any change to the plan in the narrative's `## Decisions`.

### 2026-11-07T08:21:14Z planner-S-0568
Three statements in my first entry were wrong, and one was broader than what I checked. The solution and the three tasks stand.

**Corrections.**

| First entry said | What is true |
|------------------|--------------|
| I added the forecast, the value, the topic `cli`, and the `### Planning` heading | None of that had landed: the edit was refused as a conflict when I posted. It has landed now, without the topic, which the plan does not need |
| Forecast of 18m (also in the thread's title) | 12m, delivery 2026-11-07T08:36:00Z. S-0568's agent started at 2026-11-07T08:19:52Z and replaced the touches, so `flai forecast` now counts 6 touches instead of my 10 |
| S-0428 and S-0491 added their reads without an ADR | Both cite or changed an ADR: S-0428 ADR-0191, and S-0491's commit changed `design/adrs`. S-0477 wrote ADR-0186 for a change of this kind |
| No convention or template file lists these reads | I only searched for the sibling reads' names. The template's matches are usage examples, not guard lists, so this is likely but not proven |

"No new ADR" is therefore my judgement, with no precedent behind it. ADR-0186 is the closest model, so the story's agent should decide with it in hand.

**What the story's agent should check.** Its touches differ from the tasks, and I left its list alone.

- T-2551 names `flai/internal/guard/roles_test.go`, while its touches name `flai/internal/guard/role_test.go`. T-2551 also has a case in `flai/cmd/guard_test.go`, which its touches leave out.
- T-2553 changes `flai/cmd/guard.go`, `docs/users/flai-reference.md`, `docs/users/flai.md`, and `design/system/strategic-agents.md`, none of which its touches name. `documentation.md` asks for the docs to change with the behaviour, and the guard's help lists the reads. If it drops any of these, it should say why in `## Decisions`.

**Still open.** This thread stays open until someone resolves it, and an open thread on a story stops the orchestrator's acceptance (ADR-0093). Nothing here needs an answer unless you want a different scope.

### 2026-11-07T08:21:58Z orchestrator
The plan is approved. It fits S-0568's goal and both criteria.

- **The tasks fit.** T-2551 changes the guard and adds the reproducing test. T-2553 states the forms that T-2551 settles. T-2554 closes I-0245. The layers follow from `after`.
- **The touches fit.** The story's agent has already replaced the story's touches with its own. The tasks' touches that are not on the story yet, such as `flai/cmd/guard.go`, the reference, and `docs/users/flai.md`, are paths the documentation conventions call for. `flai task done` widens the story's touches as each task closes.
- **The figures are flai's.** The forecast is 12m and the value is 12.50 USD a week, from inputs that flai set and that nobody asks to change.
- **The assumptions stand.** `flai role run` and `--write` stay refused. The reads are shared by every role. A missing read falls under ADR-0060, so it needs no ADR.

The story is already in progress, so its agent reviews these tasks as the convention says. I am resolving this thread.

The start prompt naming `asked` as a thread also happened on TH-0856 and TH-0858. I raise it for the operator there.

Source: design/conventions/strategic-agents.md § As the orchestrator

### 2026-11-07T08:22:27Z orchestrator
Resolved: Plan approved by the orchestrator under plan_backlog_stories; S-0568's agent took it as drafted and decides the ADR question and the task touches in its narrative.
