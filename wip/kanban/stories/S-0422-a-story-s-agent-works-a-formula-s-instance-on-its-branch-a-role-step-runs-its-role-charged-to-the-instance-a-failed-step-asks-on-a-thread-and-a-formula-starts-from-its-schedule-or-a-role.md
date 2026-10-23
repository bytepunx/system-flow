---
id: S-0422
type: story
nature: feature
title: A story's agent works a formula's instance on its branch, a role step runs its role charged to the instance, a failed step asks on a thread, and a formula starts from its schedule or a role
status: backlog
parent: E-0024
owner: alex
created: 2026-10-16T05:40:50Z
updated: 2026-10-23T20:20:24Z
transitions: []
tags: [cli, template]
topics: [cli, orchestration, conventions, template]
touches: [design/conventions/strategic-agents.md, template/root/design/conventions/strategic-agents.md, design/conventions/work-management.md, template/root/design/conventions/work-management.md, design/system/strategic-agents.md, design/system/flai-cli.md, docs/operators/index.md, docs/operators/settings.md, docs/users/flai.md, docs/users/flai-reference.md, template/CHANGELOG.md, flai/cmd/accept.go, flai/cmd/accept_formula_test.go, flai/cmd/formula.go, flai/cmd/formula_test.go, flai/cmd/role.go, flai/cmd/role_test.go, flai/cmd/serve_actions.go, flai/cmd/serve_actions_test.go, flai/internal/check/check.go, flai/internal/check/formulas_test.go, flai/internal/formulas/formulas.go, flai/internal/formulas/formulas_test.go, flai/internal/guard/roles.go, flai/internal/guard/roles_test.go, flai/internal/harness/harness.go, flai/internal/harness/harness_test.go, flai/internal/hostapi/writes.go, flai/internal/hostapi/writes_test.go, flai/internal/mcpserver/role.go, flai/internal/mcpserver/role_test.go, flai/internal/preview/orchestrator.go, flai/internal/preview/accept_test.go, flai/internal/serve/activity.go, flai/internal/serve/activity_test.go, flai/internal/serve/agents.go, flai/internal/serve/formula_schedule.go, flai/internal/serve/formula_schedule_test.go, flai/internal/serve/role.go, flai/internal/serve/role_test.go, flai/internal/serve/serve.go]
after: [S-0421]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: sum
  seconds: 0
  models: []
  strategic:
    - kind: orchestrator
      seconds: 612
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 194
          output: 3439
          cache_read: 55675456
          cache_write: 38642
          cost: 12.302
draft: true
cost_of_delay:
  value: 64.07
  by: planner-S-0422
  at: 2026-10-21T12:57:41Z
forecast:
  duration: 1h25m
  delivery: 2026-10-24T10:53:00Z
  basis: "Its own forecast of 1h25m; 17th in the pull order with an in-progress limit of 5, behind S-0430, S-0241, S-0406, S-0407, S-0408, S-0409, S-0410, S-0411, S-0412, S-0414, S-0415, S-0416, S-0417, S-0418, S-0419, S-0420 and S-0421."
  by: flai
  at: 2026-10-23T20:20:24Z
---
# S-0422 A story's agent works a formula's instance on its branch, a role step runs its role charged to the instance, a failed step asks on a thread, and a formula starts from its schedule or a role

## Goal

A formula's instance is worked as any story is: by a story's agent, on its branch and worktree, as TH-0471 settled. The agent's prompt and the convention tell it how. It runs each step as its task's `## Work` says, layer by layer: a flai command or a script under `scripts/` from the worktree, a role through `role_run` with the step's prompt and inputs. It closes each task when its step passes and moves the story to review when all are done, where it is accepted as any story is, even when its branch changes nothing. A step that fails blocks its task with its exit and output, asks on a thread, and runs nothing after it; the answer starts the agent again on that step. A formula starts from `flai formula run`, from its file's schedule while the `formula` host action is on, or from a role permitted to run it.

## Acceptance criteria

- [ ] The prompt of a story with `formula:` (`harness.Prompt`) says it is an instance of that formula: work its tasks as its steps, layer by layer, running each as its `## Work` names it, a flai command or a script under `scripts/` from the worktree, a role step through the MCP tool `role_run`; close each task with `task_done` when its step exits 0, with the exit in the task's Notes; move the story to review when all are done. `work-management.md`, here and in the template, says the same.
- [ ] A role step: an instance's agent starts the step's role with `role_run`, or `flai role run`, with the step's prompt and inputs and the instance as the run's item. flai lets that agent start it, though `flai serve` started it, for a role a step of its formula names and no other, and still refuses every other agent `flai serve` started. The run records the trigger `formula <name> <S-nnnn>`, its end wakes the agent's `wait_for_events`, and its usage is charged to the instance under `usage.strategic` as a planner's is to the item it planned.
- [ ] A step that fails, a command or script that exits non-zero or a role run that ends `failed` or `asked`: the agent blocks the task with its exit and the last lines of its output in the task's Notes, opens a thread on the story naming the step, runs nothing after it, and ends its session as for any question. Started again on the answer, its prompt tells it to run that step again and go on.
- [ ] `flai accept` accepts an instance whose branch changes nothing, and the orchestrator's acceptance under `accept_reviews` takes an instance as any story, with the done transition of each step's task as the evidence for its criterion in place of changed files, as the convention `strategic-agents.md` says.
- [ ] A formula starts from `flai formula run`, which with `--ready` moves the instance to ready; from its file's `schedule` while the `formula` host action is on, which makes the instance ready with the trigger `schedule <spec>` and a journal entry under `formula`, and which `flai check` refuses a formula with inputs; and from a role whose role file permits `formula run <name>`, which the guard allows for that formula alone.
- [ ] Harness tests cover an instance's prompt through a pass, a failure, and the answer; a serve test starts a fixture formula from its schedule; and the formula section of `design/system/strategic-agents.md`, `docs/users/flai.md` with `flai-reference.md`, `docs/operators/index.md` with `settings.md`, and `template/CHANGELOG.md` say how an instance is worked, how a step fails and goes on, and how a formula starts.

## Tasks
- T-1774 The prompt of a story with formula tells its agent to work the instance step by step on its branch, closing each task as its step passes, and work-management.md says the same
- T-1775 An instance's agent starts a role step's role through role_run with the step's prompt and inputs, allowed for that formula's roles alone, woken when it ends, its usage charged to the instance
- T-1776 A failed step: the prompt and the convention have the agent block the task with its exit and output, ask on a thread, run nothing after, and run the step again when answered
- T-1777 flai accept accepts an instance whose branch changes nothing, and the orchestrator's acceptance takes each step task's done transition as its criterion's evidence
- T-1778 The formula host action and the schedule starter: a formula's schedule makes its instance ready while the action is on, with the trigger schedule spec, refused on a formula with inputs
- T-1779 The formula design and the user and operator docs say how a story's agent works an instance, how a step fails and goes on, and how a formula starts
- T-2012 A formula starts from a role whose file permits formula run with its name, which the guard allows for that formula alone, and flai formula run --ready makes its instance ready

## Notes

- Drafted by the planner for E-0024 on TH-0471's recommended answer to question 5, that `flai serve` works an instance itself. Replanned on 2026-10-21 on the operator's answer, the alternative: a story's agent works the instance on its branch and worktree, as any story. The title, the goal, and every criterion changed with it.
- Waits for the formula file story.
- A step runs where the agent works: in the instance's worktree. A step that must act in the main checkout with the host's credentials, such as `flai release --pending`, is the publish formula story's to settle (S-0423).

### Planning

- **Plan thread.** TH-0582, 2026-10-21, approved and resolved by the orchestrator under `plan_backlog_stories`: the replan on TH-0471's answer 5, the rewrite of T-1774 to T-1779, and T-2012 with `flai formula run --ready`. The six tasks were then rewritten as approved; each carries a Notes line saying so. Layers: T-1774, T-1775, T-1777, T-1778, and T-2012 first; T-1776 after T-1774, since both write `harness.go` and `work-management.md`; T-1779 last.
- **TH-0471's answer 5** takes the serve-side step runner, the failure path, the no-branch acceptance, and `flai formula resume` out of the story. An instance's agent runs the steps; a failed step is a blocked task and a thread, and the answer starts the agent again, as for any question. What stays: the role step's start, now from the agent through `role_run`; the acceptance of an instance, now one whose branch may change nothing; the starters, with T-1778's role starter split into T-2012.
- **Touches.** 41, all files. The story declared none of its own; every touch comes from its seven tasks.
  - **Declared, kept:** 17 of the 31 earlier touches: the accept files, `cmd/formula.go` and its test, `formulas.go` and its test, `guard/roles.go` and its test, `hostapi/writes.go` and its test, `preview/orchestrator.go`, `serve/formula_schedule.go` and its test, `serve/serve.go`, `serve_actions.go` and its test, the design and the docs.
  - **Declared, dropped:** `serve/formula.go`, `formula_role.go`, `formula_fail.go`, and `serve/usage.go` with their tests, and `preview/accept.go`: the serve-side runner, its role and failure paths, and the no-branch acceptance that answer 5 took out.
  - **From layout, added:** `harness/harness.go` and its test, where `Prompt` varies by the story's fields, and `serve/agents.go`, where the story's request is built (T-1774, T-1776); `cmd/role.go`, `mcpserver/role.go`, and `serve/role.go` with their tests, S-0417's start path, whose served-agent refusal T-1775 lifts for an instance; `serve/activity.go` and its test, whose `chargePlanned` charges a run to its item (T-1775); `preview/accept_test.go` (T-1777); `check/check.go` and `check/formulas_test.go`, for a schedule on a formula with inputs (T-1778).
  - **From design:** the conventions `work-management.md` and `strategic-agents.md` in both copies, since the story's agent reads the first and the orchestrator the second; `template/CHANGELOG.md` with the template's conventions.
  - **Co-change, kept:** `design/system/flai-cli.md` (shared, holds nothing) and `docs/operators/settings.md`, which `make flai-reference` regenerates (T-2012) and which gains the action's row (T-1779).
  - **Co-change, not added:** `flaiover-dashboard.md`, `docs/users/flaiover.md`, `flaiover/src/lib/server/agent.ts`, `workflow.md`, `work-hierarchy.md`: an instance shows on the board and its pages as any story does, and the `formula` field is S-0421's. `mcpserver/folder.go` and `server.go`: `role_run` registers in S-0417; waking `wait_for_events` on a run's end is the agent's to place, and `task_done` widens the touches to it.
  - No folder touch.
- **Forecast.** flai's 1h25m is kept: 107 s a unit over 85 done feature stories on this model in the large band, times size 47 (6 criteria, 41 touches). The two convention tasks are prose, which overstates the size a little; lifting the served-agent refusal for one role and one agent (T-1775) and the orchestrator's evidence for an empty diff (T-1777) each need a careful test, which offsets it. The delivery is flai's play-out, 17th in the pull order after S-0421: 2026-10-22T03:22Z.
- **Cost of delay.** 64.07 USD a week, worked out by hand from TH-0471's answer 1, as the planners of S-0414 to S-0421 did, because `flai cod` refuses while E-0024 has no inputs (TH-0567, question 3, confirmed but not yet written on the epic). 4h lost per cycle at `planning.hour_rate` 150 over the default 168h cycle gives E-0024 600 USD a week; S-0422's share is its 1h25m of the 13h16m forecast for E-0024's twelve open stories: 85 / 796 × 600. Once the input is on E-0024, `flai cod S-0422` gives the same figure while the forecasts stand.
