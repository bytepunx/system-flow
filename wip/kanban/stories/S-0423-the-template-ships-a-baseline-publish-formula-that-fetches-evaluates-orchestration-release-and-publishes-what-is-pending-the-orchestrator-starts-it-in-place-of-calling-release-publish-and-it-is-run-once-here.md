---
id: S-0423
type: story
nature: improvement
title: The template ships a baseline publish formula that fetches, evaluates orchestration.release, and publishes what is pending, the orchestrator starts it in place of calling release_publish, and it is run once here
status: backlog
parent: E-0024
owner: alex
created: 2026-10-16T05:41:11Z
updated: 2026-10-25T21:20:31Z
transitions: []
tags: [cli, template]
topics: [cli, orchestration, template, conventions]
touches: [design/adrs, design/adrs/README.md, design/workflows/formulas/publish.md, template/root/design/workflows/formulas/publish.md, template/template.yaml, template/CHANGELOG.md, design/workflows/roles/orchestrator.md, template/root/design/workflows/roles/orchestrator.md, design/conventions/strategic-agents.md, template/root/design/conventions/strategic-agents.md, design/conventions/work-management.md, template/root/design/conventions/work-management.md, design/system/strategic-agents.md, design/system/flai-cli.md, docs/operators/index.md, docs/users/flai.md, docs/users/flai-reference.md, flai/cmd/release.go, flai/cmd/release_evaluate_test.go, flai/cmd/release_push_action_test.go, flai/cmd/orchestrate_permissions_test.go, flai/cmd/formula.go, flai/cmd/formula_test.go, flai/internal/formulas/formulas.go, flai/internal/formulas/formulas_test.go, flai/internal/check/check.go, flai/internal/check/formulas_test.go, flai/internal/guard/roles.go, flai/internal/guard/roles_test.go, flai/internal/mcpserver/orchestrate.go, flai/internal/mcpserver/orchestrate_test.go]
after: [S-0422]
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
    - kind: planner
      seconds: 1746
      estimated: true
      models:
        - model: claude-fable-5-1
          input: 1090
          output: 50181
          cache_read: 10121204
          cache_write: 346401
          cost: 11.9783
        - model: claude-haiku-5-5
          input: 84
          output: 28177
          cache_read: 2955958
          cache_write: 109701
          cost: 0.123
    - kind: orchestrator
      seconds: 785
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 192
          output: 3232
          cache_read: 62996281
          cache_write: 40349
          cost: 13.9187
draft: true
cost_of_delay:
  value: 44.67
  by: planner-S-0423
  at: 2026-10-21T13:13:04Z
forecast:
  duration: 1h
  delivery: 2026-10-26T13:36:00Z
  basis: "Its own forecast of 1h; 18th in the pull order with an in-progress limit of 5, behind S-0241, S-0406, S-0407, S-0408, S-0409, S-0410, S-0411, S-0412, S-0414, S-0415, S-0416, S-0417, S-0418, S-0419, S-0420, S-0421 and S-0422."
  by: flai
  at: 2026-10-25T21:20:31Z
---
# S-0423 The template ships a baseline publish formula that fetches, evaluates orchestration.release, and publishes what is pending, the orchestrator starts it in place of calling release_publish, and it is run once here

## Goal

Publishing a release stops being prose. A baseline formula, `publish`, spells out what `work-management.md` and `strategic-agents.md` § Publishing a release say in words: fetch, evaluate the release policy, publish what is accepted and not yet released, and record the versions and tags. The orchestrator, under its `publish` permission, starts the formula when `release_evaluate` is met, in place of calling `release_publish` itself, so each publish is a story worked, measured, and archived. It is run once in this repository before the epic is accepted.

## Acceptance criteria

- [ ] `design/workflows/formulas/publish.md`, in the template and here, has the steps `git fetch`, `flai release --evaluate` (which passes only when the policy is met), `flai release --pending`, and a last step that records the versions and tags on the instance, each after the one before, and `flai check --strict` reports nothing on it.
- [ ] A formula step of `flai release --pending` is refused, as `release_publish` is, while the `push` host action is off, and a publish refused for a moved remote stops the instance with flai's `conflict:` message on its task and a thread on the story.
- [ ] Under `publish`, the orchestrator's role file and the convention tell it to run `flai formula run publish` when `release_evaluate` is met, the guard lets it run that formula and no other way to the remote, and `release_publish` stays for a project with no `publish` formula; the orchestrator tests pass.
- [ ] The formula is run once in this repository: its instance story is worked to review and accepted, and its Notes name the release it published, or, when nothing was pending, the evaluation that stopped it.
- [ ] `design/system/strategic-agents.md` § Publishing a release, the convention `strategic-agents.md` and `work-management.md` where they describe publishing, `docs/operators/index.md`, and `template/CHANGELOG.md` say that publishing is the `publish` formula.

## Tasks
- T-1780 design/workflows/formulas/publish.md, in the template and here, fetches, evaluates the release policy, publishes what is pending, and records the versions and tags
- T-1781 A formula step that publishes is refused while the push host action is off, and a publish refused for a moved remote stops the instance with flai's conflict message
- T-1782 Under publish, the orchestrator's role file tells it to run flai formula run publish when release_evaluate is met, and the guard lets it run that formula
- T-1783 The publish formula is run once in this repository, its instance worked to review and accepted, and the release or the evaluation that stopped it recorded
- T-1784 The design, the conventions that describe publishing, the operator docs, and the changelog say that publishing is the publish formula
- T-2013 A formula step may say in: main, which flai formula run writes into its task's Work as the main checkout's path and flai check validates
- T-2014 release_publish refuses, naming flai formula run publish, while the project defines a publish formula, and stays for a project with none
- T-2015 An ADR refining ADR-0094 records that under publish the orchestrator starts the publish formula in place of release_publish, which stays for a project with none

## Notes

- Drafted by the planner for E-0024. Waits for the story that works a formula's instance.
- Replanned on 2026-10-21 on TH-0471's answers 2 and 5: the formula lives under `design/workflows/formulas/`, and its instance is worked by a story's agent on its branch and worktree, as any story (S-0422). Criterion 1's path changed with answer 2.
- `flai release --pending` acts with the host's credentials on the main branch, never from a story's worktree: it tags and pushes the branch that is checked out. How a step says it runs in the main checkout is asked on TH-0583, question 1; the plan assumes `in: main` on the step (T-2013).
- Criterion 4 needs the operator's `push` host action on and an accepted, unreleased story to publish; with none pending, the run shows the evaluation stopping the instance cleanly.

### Planning

- **Plan thread.** TH-0583, 2026-10-21: the replan on TH-0471's answers 2 and 5, with four questions: where a main-checkout step runs (`in: main`, T-2013), the `push` gate (`flai release --pending` refusing an agent `flai serve` started, T-1781), the last step (three steps and the formula's body telling the agent to record the versions and tags), and the rewrite of T-1780 to T-1783's words, which predate the answers. The orchestrator approved the plan and the rewrites that apply TH-0471 alone, which are made: T-1780's and T-1782's paths and T-1783's instance worked by a story's agent. Questions 1 to 3 went to the operator as a recommendation to take the three answers; T-1781's words, the `in: main` and three-step parts of T-1780, and criterion 1 change once it is confirmed. The orchestrator asked for an ADR refining ADR-0094, which is T-2015.
- **Layers.** T-2013 and T-2015 first; T-1780 after T-2013, since the formula file carries `in: main`, and T-2014 after T-2015, which records the decision it enforces; T-1781 and T-1782 after T-1780, T-1781 because both change `flai/cmd/release.go` and T-1782 because its test runs the formula; T-1783 after T-1781, T-1782, and T-2014, so that the run exercises the gate, the rule, and the refusal; T-1784 after T-2013, T-1782, and T-2015, so that the docs can say what `in: main` and the orchestrator's rule do and link the ADR. No two tasks of a layer share a path.
- **Touches.** 31: 30 files and one folder. The story declared none of its own; every touch comes from its eight tasks.
  - **Declared, kept:** 14 of the 20 earlier touches: the conventions in both copies, `design/system/strategic-agents.md` and `flai-cli.md`, the operator and user docs with `flai-reference.md`, the changelog, `template.yaml`, `flai/cmd/release.go` with `release_evaluate_test.go`, and `orchestrate_permissions_test.go`.
  - **Declared, moved:** `design/formulas/publish.md` and `design/roles/orchestrator.md`, in both copies, now under `design/workflows/` (answer 2), on T-1780 and T-1782.
  - **Declared, dropped:** `flai/internal/serve/formula.go` and `formula_publish_test.go`: the serve-side step runner answer 5 took out of S-0422.
  - **From layout, added:** `flai/internal/formulas/formulas.go`, `flai/internal/check/check.go`, and `flai/cmd/formula.go`, with their tests, for a step's `in` field (T-2013) and flai's hold on the orchestrator's `formula run` (T-1782); `flai/cmd/release_push_action_test.go`, new, for the `push` gate (T-1781); `flai/internal/guard/roles.go` and its test, S-0416's role path, in case the `publish` permission on `formula run publish` needs code beside the role file (T-1782); `flai/internal/mcpserver/orchestrate.go` and `orchestrate_test.go`, where `release_publish` and `mayPublish` live (T-2014); `design/adrs/README.md` for the ADR's row (T-2015).
  - **From design:** the role file in both copies, since S-0419 reads the orchestrator's prompt from it; the conventions and the docs criterion 5 names.
  - **Co-change, kept:** `docs/users/flai.md` and `design/system/flai-cli.md`, for `--require-met` and the refusal of `--pending`.
  - **Co-change, not added:** `docs/operators/settings.md` (19%: no new setting; `in` is a formula field and `push` has its row), `design/system/flaiover-dashboard.md` and `docs/users/flaiover.md` (no page changes), `design/system/workflow.md` (publishing is described in `strategic-agents.md` and the conventions), `flai/internal/harness/harness.go` (S-0419 moves the orchestrator's prompt into its role file before this story).
  - The run task, T-1783, touches nothing: the instance is a story of its own, and its agent's commits are its own.
  - **Folder touch kept:** `design/adrs`, on T-2015, because the ADR's file name is known only once `flai adr new` makes it. The folder is in `claims.shared`, so it holds no story.
- **Forecast.** flai's 38m was for size 25; at size 36 (5 criteria, 31 touches) the same 89 s a unit gives 53m, raised to 1h because T-1783 waits on an instance another agent works and on its acceptance, which no touch counts, as the earlier plan raised it by 14m for the same reason. The delivery is flai's play-out after S-0422, 18th in the pull order, with the duration times the cycle factor 1.15: 2026-10-22T04:59Z.
- **Cost of delay.** 44.67 USD a week, worked out by hand from TH-0471's answer 1, as the planners of S-0414 to S-0422 did, because `flai cod` refuses while E-0024 has no inputs (TH-0567, question 3, confirmed and not yet written on the epic). 4h lost per cycle at `planning.hour_rate` 150 over the default 168h cycle gives E-0024 600 USD a week; S-0423's share is its 1h of the 13h26m forecast for E-0024's twelve open stories: 60 / 806 × 600. Once the input is on E-0024, `flai cod S-0423` gives the same figure while the forecasts stand.
