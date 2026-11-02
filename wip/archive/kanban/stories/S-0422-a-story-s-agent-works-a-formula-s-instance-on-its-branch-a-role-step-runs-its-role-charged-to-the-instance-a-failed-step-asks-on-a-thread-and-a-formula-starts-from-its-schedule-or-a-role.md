---
id: S-0422
type: story
nature: feature
title: A story's agent works a formula's instance on its branch, a role step runs its role charged to the instance, a failed step asks on a thread, and a formula starts from its schedule or a role
status: done
parent: E-0024
owner: alex
created: 2026-10-16T05:40:50Z
updated: 2026-11-02T04:07:12Z
transitions:
  - to: ready
    at: 2026-10-28T13:11:31Z
    by: alex
  - to: backlog
    at: 2026-10-28T22:20:45Z
    by: alex
  - to: ready
    at: 2026-11-01T23:43:54Z
    by: alex
  - to: in-progress
    at: 2026-11-02T03:02:51Z
    by: agent-S-0422
  - to: review
    at: 2026-11-02T03:55:08Z
    by: agent-S-0422
  - to: done
    at: 2026-11-02T04:07:12Z
    by: orchestrator
tags: [cli, template]
topics: [cli, orchestration, conventions, template]
touches: [design/conventions/strategic-agents.md, template/root/design/conventions/strategic-agents.md, design/conventions/work-management.md, template/root/design/conventions/work-management.md, design/system/strategic-agents.md, design/system/flai-cli.md, docs/operators/index.md, docs/operators/settings.md, docs/users/flai.md, docs/users/flai-reference.md, template/CHANGELOG.md, flai/cmd/accept.go, flai/cmd/accept_formula_test.go, flai/cmd/formula.go, flai/cmd/formula_test.go, flai/cmd/role.go, flai/cmd/role_test.go, flai/cmd/serve_actions.go, flai/cmd/serve_actions_test.go, flai/internal/check/check.go, flai/internal/check/formulas_test.go, flai/internal/formulas/formulas.go, flai/internal/formulas/formulas_test.go, flai/internal/guard/roles.go, flai/internal/guard/roles_test.go, flai/internal/harness/harness.go, flai/internal/harness/harness_test.go, flai/internal/hostapi/writes.go, flai/internal/hostapi/writes_test.go, flai/internal/mcpserver/role.go, flai/internal/mcpserver/role_test.go, flai/internal/preview/orchestrator.go, flai/internal/preview/accept_test.go, flai/internal/serve/activity.go, flai/internal/serve/activity_test.go, flai/internal/serve/agents.go, flai/internal/serve/formula_schedule.go, flai/internal/serve/formula_schedule_test.go, flai/internal/serve/role.go, flai/internal/serve/role_test.go, flai/internal/serve/serve.go, flai/internal/roles/roles.go, design/adrs/0239-a-formula-s-instance-is-worked-by-its-story-s-agent-step-by-step-a-role-step.md, design/adrs/README.md, design/workflows/formulas/README.md, template/root/design/workflows/formulas/README.md, design/issues/I-0250-the-orchestrator-s-conventions-sit-within-about-1-5kb-of-its-80kb-pack-budget.md, design/issues/summary.md, design/issues/I-0246-testaheldwaitfailsonatruncatedstory-and-testservestartsadefinedroleonaneventando.md, design/issues/I-0248-testtheoperatorsstopisonestoppedagentevent-s-waiting-for-an-answer-case-fails.md]
after: [S-0421]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 3157
  priced_by: harness
  turns:
    - day: 2026-11-02
      ceremony: 5
      test_runs: 1
      hand_edits: 4
      work: 90
  models:
    - model: claude-opus-5-5
      input: 914
      output: 261603
      cache_read: 66459126
      cache_write: 1399779
      cost: 26.6804
  strategic:
    - kind: orchestrator
      seconds: 845
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 275
          output: 4225
          cache_read: 64007755
          cache_write: 357964
          cost: 14.2243
        - model: claude-sonnet-5-5
          input: 6
          output: 33
          cache_read: 52548
          cache_write: 48289
          cost: 0.0703
cost_of_delay:
  value: 64.07
  by: planner-S-0422
  at: 2026-10-21T12:57:41Z
forecast:
  duration: 1h25m
  delivery: 2026-11-02T05:10:00Z
  basis: "Its own forecast of 1h25m; 4th in the pull order with an in-progress limit of 5, behind S-0421, S-0553, S-0423, S-0424 and S-0554."
  by: flai
  at: 2026-11-02T02:09:16Z
finalized:
  by: alex
  at: 2026-10-28T13:11:21Z
---
# S-0422 A story's agent works a formula's instance on its branch, a role step runs its role charged to the instance, a failed step asks on a thread, and a formula starts from its schedule or a role

## Goal

A formula's instance is worked as any story is: by a story's agent, on its branch and worktree, as TH-0471 settled. The agent's prompt and the convention tell it how. It runs each step as its task's `## Work` says, layer by layer: a flai command or a script under `scripts/` from the worktree, a role through `role_run` with the step's prompt and inputs. It closes each task when its step passes and moves the story to review when all are done, where it is accepted as any story is, even when its branch changes nothing. A step that fails blocks its task with its exit and output, asks on a thread, and runs nothing after it; the answer starts the agent again on that step. A formula starts from `flai formula run`, from its file's schedule while the `formula` host action is on, or from a role permitted to run it.

## Acceptance criteria

- [x] The prompt of a story with `formula:` (`harness.Prompt`) says it is an instance of that formula: work its tasks as its steps, layer by layer, running each as its `## Work` names it, a flai command or a script under `scripts/` from the worktree, a role step through the MCP tool `role_run`; close each task with `task_done` when its step exits 0, with the exit in the task's Notes; move the story to review when all are done. `work-management.md`, here and in the template, says the same.
- [x] A role step: an instance's agent starts the step's role with `role_run`, or `flai role run`, with the step's prompt and inputs and the instance as the run's item. flai lets that agent start it, though `flai serve` started it, for a role a step of its formula names and no other, and still refuses every other agent `flai serve` started. The run records the trigger `formula <name> <S-nnnn>`, and its end wakes the agent's `wait_for_events`.
- [x] A step that fails, a command or script that exits non-zero or a role run that ends `failed` or `asked`: the agent blocks the task with its exit and the last lines of its output in the task's Notes, opens a thread on the story naming the step, runs nothing after it, and ends its session as for any question. Started again on the answer, its prompt tells it to run that step again and go on.
- [x] `flai accept` accepts an instance whose branch changes nothing, and the orchestrator's acceptance under `accept_reviews` takes an instance as any story, with the done transition of each step's task as the evidence for its criterion in place of changed files, as the convention `strategic-agents.md` says.
- [x] A formula starts from `flai formula run`, which with `--ready` moves the instance to ready; from its file's `schedule` while the `formula` host action is on, which makes the instance ready with the trigger `schedule <spec>` and a journal entry under `formula`, and which `flai check` refuses a formula with inputs; and from a role whose role file permits `formula run <name>`, which the guard allows for that formula alone.
- [x] Harness tests cover an instance's prompt through a pass, a failure, and the answer; a serve test starts a fixture formula from its schedule; and the formula section of `design/system/strategic-agents.md`, `docs/users/flai.md` with `flai-reference.md`, `docs/operators/index.md` with `settings.md`, and `template/CHANGELOG.md` say how an instance is worked, how a step fails and goes on, and how a formula starts.

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

- **After S-0421.** S-0422 reads the formula file format, `flai formula run`, and the instance story's `formula` field, which S-0421 writes, to work a formula's instance step by step (TH-0787, D-0454).
- **Usage charge moved to S-0425.** Criterion 2 no longer charges a role step's usage to the instance: `usage.strategic` takes a defined role's name only through S-0425's ADR on `metrics.md`, whose criterion 2 charges a defined role's run to its item, the instance here. The operator confirmed it on TH-0800 on 2026-11-02 and adds the clause to S-0425's criteria.

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

### Accepted by the orchestrator

- Verified: e8660210cab367829c682f689176e7b4cd0722aa
- Batch: 20261102-035950
- Enqueued verified: 2f1ca8fa7df5cf22be9ed0e1ca5425e6540e7e21
- At: 2026-11-02T04:07:12Z

Verdict: pass. flai verify S-0422 passed every step at 2f1ca8fa7df5cf22be9ed0e1ca5425e6540e7e21, and the verifier matched all six criteria to the diff.
- 1: flai/internal/harness/harness.go, flai/internal/serve/agents.go, design/conventions/work-management.md, template/root/design/conventions/work-management.md
- 2: flai/cmd/role.go, flai/internal/serve/role.go, flai/internal/mcpserver/role.go, flai/cmd/role_test.go, flai/internal/serve/role_test.go, flai/internal/mcpserver/role_test.go
- 3: flai/internal/harness/harness.go, flai/internal/harness/harness_test.go, design/conventions/work-management.md, template/root/design/conventions/work-management.md, template/CHANGELOG.md
- 4: flai/cmd/accept.go, flai/cmd/accept_formula_test.go, flai/internal/preview/orchestrator.go, flai/internal/preview/accept_test.go, design/conventions/strategic-agents.md, template/root/design/conventions/strategic-agents.md
- 5: flai/cmd/formula.go, flai/cmd/formula_test.go, flai/internal/serve/formula_schedule.go, flai/internal/serve/serve.go, flai/internal/hostapi/writes.go, flai/cmd/serve_actions_test.go, flai/internal/roles/roles.go, flai/internal/formulas/formulas.go, flai/internal/check/formulas_test.go, flai/internal/guard/roles.go, flai/internal/guard/roles_test.go
- 6: flai/internal/harness/harness_test.go, flai/internal/serve/formula_schedule_test.go, design/system/strategic-agents.md, docs/users/flai.md, docs/users/flai-reference.md, docs/operators/index.md, docs/operators/settings.md, template/CHANGELOG.md
