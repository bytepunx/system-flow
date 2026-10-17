---
id: S-0422
type: story
nature: feature
title: flai serve works a formula's instance step by step, commands and scripts in the main checkout and roles through their runs, and a formula starts from a host action, a schedule, or a role
status: backlog
parent: E-0024
owner: alex
created: 2026-10-16T05:40:50Z
updated: 2026-10-17T03:57:03Z
transitions: []
tags: [cli]
topics: [cli, orchestration]
touches: [design/system/strategic-agents.md, docs/operators/index.md, docs/operators/settings.md, docs/users/flai.md, docs/users/flai-reference.md, flai/cmd/accept_formula_test.go, flai/cmd/accept.go, flai/cmd/formula.go, flai/cmd/formula_test.go, flai/internal/formulas/formulas.go, flai/internal/formulas/formulas_test.go, flai/internal/guard/roles.go, flai/internal/guard/roles_test.go, flai/internal/hostapi/writes.go, flai/internal/hostapi/writes_test.go, flai/internal/preview/accept.go, flai/internal/preview/orchestrator.go, flai/internal/serve/formula_fail.go, flai/internal/serve/formula_fail_test.go, flai/internal/serve/formula.go, flai/internal/serve/formula_role.go, flai/internal/serve/formula_role_test.go, flai/internal/serve/formula_schedule.go, flai/internal/serve/formula_schedule_test.go, flai/internal/serve/formula_test.go, flai/internal/serve/serve.go, flai/internal/serve/usage.go, flai/internal/serve/usage_test.go, flai/cmd/serve_actions.go, flai/cmd/serve_actions_test.go, design/system/flai-cli.md]
after: [S-0421]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
draft: true
forecast:
  duration: 1h15m
  delivery: 2026-10-17T21:41:00Z
  basis: "Its own forecast of 1h15m; 33rd in the pull order with an in-progress limit of 5, behind S-0444, S-0443, S-0241, S-0392, S-0393, S-0394, S-0395, S-0396, S-0397, S-0398, S-0399, S-0400, S-0401, S-0402, S-0403, S-0404, S-0405, S-0406, S-0407, S-0408, S-0409, S-0410, S-0411, S-0412, S-0414, S-0415, S-0416, S-0417, S-0418, S-0419, S-0420 and S-0421."
  by: flai
  at: 2026-10-17T03:57:03Z
---
# S-0422 flai serve works a formula's instance step by step, commands and scripts in the main checkout and roles through their runs, and a formula starts from a host action, a schedule, or a role

## Goal

A formula's instance runs itself. `flai serve` works the story `flai formula run` made, layer by layer: a command or script step runs in the main checkout and its task is done when it exits 0; a role step starts the role with the step's prompt and inputs and is done when the run ends `worked`. When every task is done the story goes to review, where it is accepted and archived as any story is, with nothing to merge. A step that fails stops the instance with its output on the task and a thread on the story. A formula starts from its host action, its schedule, or a role permitted to run it.

## Acceptance criteria

- [ ] `flai serve` works an instance story behind the `formula` host action: it moves the story to in-progress, runs each layer's steps in `after` order, a command or a script under `scripts/` in the main checkout with the formula's inputs in its environment, a role through `role.run` with the step's prompt and inputs; it moves each task to done when its step passes, and the story to review when all are done.
- [ ] A step that fails, a command or script that exits non-zero or a role run that ends `failed` or `asked`, blocks its task with the exit and the last lines of its output in the task's Notes, opens a thread on the story naming the step, and runs nothing after it; `flai formula resume <S-nnnn>` runs the step again once the operator has answered.
- [ ] An instance has no branch or worktree, and `flai accept` accepts it by moving it to done and archiving it with nothing to merge; the orchestrator's acceptance under `accept_reviews` takes it as any story, its evidence each step's exit.
- [ ] Its usage is measured as any story's: each role step's run is charged to the story, and each command and script step's seconds are recorded on its task.
- [ ] A formula starts from `flai formula run`, from its file's schedule while the `formula` host action is on, and from a role whose role file permits `formula run <name>`, which the guard allows only for that formula.
- [ ] A serve test works a fixture formula of two command steps, a script step, and a role step with a fake harness, through a pass and through a failure resumed; `design/system/strategic-agents.md` or the formula design, `docs/users/flai.md` with `flai-reference.md`, and `docs/operators/index.md` say how an instance is worked.

## Tasks
- T-1774 serve works a formula instance's command and script steps layer by layer in the main checkout, behind the formula host action, moving each task and the story as its steps pass
- T-1775 A formula's role step starts its role with the step's prompt and inputs through the defined-role path and is done when the run ends worked, its cost charged to the instance
- T-1776 A failed formula step blocks its task with its exit and output, opens a thread on the instance, and flai formula resume runs it again once answered
- T-1777 flai accept accepts a formula instance with no branch by moving it to done and archiving it, and the orchestrator's acceptance takes its steps' exits as evidence
- T-1778 A formula starts on its file's schedule while the formula action is on, and from a role whose file permits formula run with its name, which the guard allows for that formula alone
- T-1779 The formula design and the user and operator docs say how flai serve works an instance, how a step fails and resumes, and how a formula starts

## Notes

- Drafted by the planner for E-0024 on TH-0471's recommended answer to question 5, that `flai serve` works an instance itself. If the operator chooses a story's agent instead, this story becomes the prompt and convention changes that tell the agent how to work an instance, and criteria 1 to 4 change with it.
- Waits for the formula file story.

### Planning

- **Touches.** The story declared none. Every touch comes from its six tasks.
  - **From layout:** five new serve files (`formula.go`, `formula_role.go`, `formula_fail.go`, `formula_schedule.go`, each with a test), `serve.go`, and `usage.go`; `flai/cmd/formula.go` for `resume`; `flai/cmd/accept.go` with a new `accept_formula_test.go`; `preview/accept.go` and `orchestrator.go`; `guard/roles.go`; `formulas.go` for the schedule check; and the host API's `writes.go`.
  - **From design:** `strategic-agents.md` (or the formula design S-0421 makes), the user and operator docs, and `docs/operators/settings.md` for the `formula` action.
  - **Co-change, added:** `flai/cmd/serve_actions.go` and its test (8%), which list the host actions, and `design/system/flai-cli.md` (48%).
  - **Co-change, not added:** `flaiover/src/lib/server/agent.ts` and the dashboard's design. An instance shows on the board as any story does.
  - No folder touch.
- **Forecast.** Raised from flai's 59m to 1h15m. Accepting a story with no branch changes `flai accept`, which seventeen test files hold. The forecast's basis gives the figures, and its delivery is flai's play-out moved by any raise times the cycle factor.
- **Cost of delay.** Not set. Neither the story nor E-0024 has inputs, so `flai cod` refuses. The epic's input is asked on TH-0471, question 1, and the value follows its answer.
