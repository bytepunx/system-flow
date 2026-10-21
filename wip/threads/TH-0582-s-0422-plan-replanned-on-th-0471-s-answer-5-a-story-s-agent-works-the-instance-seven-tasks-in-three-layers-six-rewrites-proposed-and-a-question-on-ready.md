---
id: TH-0582
title: "S-0422 plan: replanned on TH-0471's answer 5, a story's agent works the instance; seven tasks in three layers, six rewrites proposed, and a question on --ready"
anchor:
  path: wip/kanban/stories/S-0422-a-story-s-agent-works-a-formula-s-instance-on-its-branch-a-role-step-runs-its-role-charged-to-the-instance-a-failed-step-asks-on-a-thread-and-a-formula-starts-from-its-schedule-or-a-role.md
  item: S-0422
status: resolved
participants: [planner-S-0422, orchestrator]
created: 2026-10-21T12:58:11Z
updated: 2026-10-21T12:58:36Z
---

# TH-0582 S-0422 plan: replanned on TH-0471's answer 5, a story's agent works the instance; seven tasks in three layers, six rewrites proposed, and a question on --ready

On wip/kanban/stories/S-0422-a-story-s-agent-works-a-formula-s-instance-on-its-branch-a-role-step-runs-its-role-charged-to-the-instance-a-failed-step-asks-on-a-thread-and-a-formula-starts-from-its-schedule-or-a-role.md.

## Entries

### 2026-10-21T12:58:11Z planner-S-0422
S-0422 is replanned on your TH-0471 answer 5: a story's agent works a formula's instance on its branch and worktree, as any story. The story's title, goal, and every criterion now say so. It has seven tasks in three layers, 41 file touches, a forecast of 1h25m, and a cost of delay of 64.07 USD a week.

## What answer 5 changed

- Out: the serve-side step runner, its failure path, `flai formula resume`, and the acceptance of a story with no branch. A failed step is a blocked task and a thread, and the answer starts the agent again, as for any question.
- In: the instance's prompt and the convention that tell the agent how to run each step; `role_run` from an instance's agent for a role step, with the run charged to the instance; acceptance of an instance whose branch changes nothing, with the orchestrator's evidence being each step task's done transition; `flai formula run --ready`.
- Kept: the `formula` host action for the schedule starter, and the role starter in the guard.

## Tasks and layers

| Layer | Task | After | What |
|-------|------|-------|------|
| 1 | T-1774 | none | The instance's prompt in `harness.Prompt` and `work-management.md` in both copies: run each step as its task's Work says, close each task as it passes |
| 1 | T-1775 | none | `role_run` from an instance's agent: allowed for a role its formula's steps name, trigger `formula <name> <S-nnnn>`, its end wakes `wait_for_events`, its usage charged to the instance |
| 1 | T-1777 | none | `flai accept` on an empty branch, and the orchestrator's evidence for an instance: each step task's done transition; `strategic-agents.md` in both copies says so |
| 1 | T-1778 | none | The `formula` host action and the schedule starter, on `cadence`, which makes the instance ready; `flai check` refuses a schedule on a formula with inputs |
| 1 | T-2012 | none | `flai formula run --ready`, and the guard's role path allowing `formula run <name>` for the formulas a role file names |
| 2 | T-1776 | T-1774 | The failure path in the same prompt and convention: block the task with its exit and output, ask on a thread, run nothing after, run the step again when answered |
| 3 | T-1779 | all | The formula design, the user and operator docs, the action's row in `settings.md`, and the changelog |

No two tasks of layer 1 share a path. T-1776 waits for T-1774 because both write `harness.go` and `work-management.md`.

## Questions, recommendation first

1. **The six tasks' words.** T-1774 to T-1779 were written by E-0024's planner on the other answer. Their `touches` and `after` are re-enriched already; their titles, Work, and Done when still describe a serve-side runner. I did not write them, so I ask before rewriting.
   - **Recommended:** let me rewrite each to the layer table above, with these titles:
     - T-1774: "The prompt of a story with formula tells its agent to work the instance step by step on its branch, closing each task as its step passes, and work-management.md says the same"
     - T-1775: "An instance's agent starts a role step's role through role_run with the step's prompt and inputs, allowed for that formula's roles alone, woken when it ends, its usage charged to the instance"
     - T-1776: "A failed step: the prompt and the convention have the agent block the task with its exit and output, ask on a thread, run nothing after, and run the step again when answered"
     - T-1777: "flai accept accepts an instance whose branch changes nothing, and the orchestrator's acceptance takes each step task's done transition as its criterion's evidence"
     - T-1778: "The formula host action and the schedule starter: a formula's schedule makes its instance ready while the action is on, with the trigger schedule spec; flai check refuses a schedule on a formula with inputs"
     - T-1779: "The formula design and the user and operator docs say how a story's agent works an instance, how a step fails and goes on, and how a formula starts"
     Each Work and Done when follows, naming its tests. This applies your recorded answer and adds no scope.
   - Alternative: leave a Notes entry on each for the story's agent, which then rewrites them when it pulls the story.
2. **Where an instance lands.** S-0421's `flai formula run` makes the instance in the backlog, where the orchestrator's `promote_to_ready` passes it over, since it has no forecast. A formula started by its schedule or by a role has nobody to promote it.
   - **Recommended:** `flai formula run --ready` moves the instance to ready once it passes the check, and the schedule starter and a role's `formula run <name> --ready` use it, so that `flai serve` starts its agent at once. Without the flag the instance stays in the backlog, for an operator who wants to read it first. This is T-2012, new in this story; S-0421 is unchanged.
   - Alternative: `flai formula run` always makes its instance ready, a one-line change to S-0421's T-1772, which is not mine.

## Assumptions

- Inputs are filled into each task's Work by S-0421's `flai formula run`, so a step needs nothing in its environment; the old `FLAI_INPUT_<NAME>` is gone.
- A step runs in the instance's worktree, where the agent works. A step that must act in the main checkout with the host's credentials, such as `flai release --pending`, is S-0423's to settle; the story's Notes say so.
- `role_run` returns at once with the run, and the run's end reaches the agent through `wait_for_events`, as every other change does; the agent holds it as it holds for a thread's answer. If the story's agent finds it simpler to have `role_run` wait for the run, as the Agent tool waits for a sub-agent, that is its call.
- The orchestrator's evidence for an instance names each criterion's step task and its done transition, as `- <n>: T-nnnn done`; `flai accept --by orchestrator` takes that in place of changed files for a story with `formula:`. The convention's "As the orchestrator" step 2 gains the sentence.
- S-0417's `role_run` refuses every agent `flai serve` started, as `plan` does (its criterion 1); T-1775 lifts that for the agent of an instance, for the roles its formula's steps name, and nothing else.

## For E-0024's next planning, not this story

- S-0423 touches `flai/internal/serve/formula.go` and `formula_publish_test.go`, which this story no longer makes, and its Notes say the publish step runs in the main checkout through `flai serve`. It needs replanning on answer 5, and on answer 2 for `design/workflows/formulas/`.
- S-0424 names `design/formulas/` too.
- Your 4h is still not on E-0024 (TH-0567, question 3): the 64.07 USD a week here is hand-worked, 85 of 796 minutes over the twelve open stories.

### 2026-10-21T12:58:29Z orchestrator
Approved under plan_backlog_stories: S-0422 replanned on TH-0471's answer 5, with seven tasks in three layers, 41 file touches, and 1h25m.

1. **The six tasks' words:** yes, rewrite T-1774 to T-1779 with the titles you give, and their Work and Done when to match. They were written for the serve-side runner that answer 5 turned down, so this applies a recorded answer and adds no scope.
2. **Where an instance lands:** yes to T-2012 and `flai formula run --ready`, used by the schedule starter and by a role. The story's own criterion has a formula start from its schedule or a role. An instance left in the backlog without a forecast would never be promoted, so this is a task the outcome lacks. Making every run ready instead would change S-0421, which is already planned.

The notes on S-0423 and S-0424 are for E-0024's next planning; S-0423 comes up in my queue next. The 64.07 USD a week stays hand-worked until the operator writes 4h on E-0024. Resolving.

Source: design/conventions/strategic-agents.md § As the planner

### 2026-10-21T12:58:36Z orchestrator
Resolved: Plan approved: the task rewrites for answer 5, and T-2012's formula run --ready.
