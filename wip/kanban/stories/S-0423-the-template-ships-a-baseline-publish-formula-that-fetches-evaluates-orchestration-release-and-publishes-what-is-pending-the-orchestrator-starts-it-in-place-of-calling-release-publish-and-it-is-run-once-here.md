---
id: S-0423
type: story
nature: improvement
title: The template ships a baseline publish formula that fetches, evaluates orchestration.release, and publishes what is pending, the orchestrator starts it in place of calling release_publish, and it is run once here
status: backlog
parent: E-0024
owner: alex
created: 2026-10-16T05:41:11Z
updated: 2026-10-16T06:07:38Z
transitions: []
tags: [cli, template]
topics: [cli, orchestration, template, conventions]
touches: [design/conventions/strategic-agents.md, design/conventions/work-management.md, design/formulas/publish.md, design/roles/orchestrator.md, design/system/strategic-agents.md, docs/operators/index.md, docs/users/flai-reference.md, flai/cmd/orchestrate_permissions_test.go, flai/cmd/release_evaluate_test.go, flai/cmd/release.go, flai/internal/serve/formula.go, flai/internal/serve/formula_publish_test.go, template/CHANGELOG.md, template/root/design/conventions/strategic-agents.md, template/root/design/conventions/work-management.md, template/root/design/formulas/publish.md, template/root/design/roles/orchestrator.md, template/template.yaml, design/system/flai-cli.md, docs/users/flai.md]
after: [S-0422]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
draft: true
forecast:
  duration: 50m
  delivery: 2026-10-17T01:57:00Z
  basis: "Its own forecast of 50m; 32nd in the pull order with an in-progress limit of 5, behind S-0426, S-0427, S-0241, S-0392, S-0393, S-0394, S-0395, S-0396, S-0397, S-0398, S-0399, S-0400, S-0401, S-0402, S-0403, S-0404, S-0405, S-0406, S-0407, S-0408, S-0409, S-0410, S-0411, S-0412, S-0414, S-0415, S-0416, S-0417, S-0418, S-0419, S-0420, S-0421 and S-0422."
  by: flai
  at: 2026-10-16T06:07:38Z
---
# S-0423 The template ships a baseline publish formula that fetches, evaluates orchestration.release, and publishes what is pending, the orchestrator starts it in place of calling release_publish, and it is run once here

## Goal

Publishing a release stops being prose. A baseline formula, `publish`, spells out what `work-management.md` and `strategic-agents.md` § Publishing a release say in words: fetch, evaluate the release policy, publish what is accepted and not yet released, and record the versions and tags. The orchestrator, under its `publish` permission, starts the formula when `release_evaluate` is met, in place of calling `release_publish` itself, so each publish is a story worked, measured, and archived. It is run once in this repository before the epic is accepted.

## Acceptance criteria

- [ ] `design/formulas/publish.md`, in the template and here, has the steps `git fetch`, `flai release --evaluate` (which passes only when the policy is met), `flai release --pending`, and a last step that records the versions and tags on the instance, each after the one before, and `flai check --strict` reports nothing on it.
- [ ] A formula step of `flai release --pending` is refused, as `release_publish` is, while the `push` host action is off, and a publish refused for a moved remote stops the instance with flai's `conflict:` message on its task and a thread on the story.
- [ ] Under `publish`, the orchestrator's role file and the convention tell it to run `flai formula run publish` when `release_evaluate` is met, the guard lets it run that formula and no other way to the remote, and `release_publish` stays for a project with no `publish` formula; the orchestrator tests pass.
- [ ] The formula is run once in this repository: its instance story is worked to review and accepted, and its Notes name the release it published, or, when nothing was pending, the evaluation that stopped it.
- [ ] `design/system/strategic-agents.md` § Publishing a release, the convention `strategic-agents.md` and `work-management.md` where they describe publishing, `docs/operators/index.md`, and `template/CHANGELOG.md` say that publishing is the `publish` formula.

## Tasks
- T-1780 design/formulas/publish.md, in the template and here, fetches, evaluates the release policy, publishes what is pending, and records the versions and tags
- T-1781 A formula step that publishes is refused while the push host action is off, and a publish refused for a moved remote stops the instance with flai's conflict message
- T-1782 Under publish, the orchestrator's role file tells it to run flai formula run publish when release_evaluate is met, and the guard lets it run that formula
- T-1783 The publish formula is run once in this repository, its instance worked to review and accepted, and the release or the evaluation that stopped it recorded
- T-1784 The design, the conventions that describe publishing, the operator docs, and the changelog say that publishing is the publish formula

## Notes

- Drafted by the planner for E-0024. Waits for the story that works a formula's instance.
- `flai release --pending` acts with the host's credentials: the formula runs it in the main checkout through `flai serve`, as `publish.run` does, never from a story's worktree.
- Criterion 4 needs the operator's `push` host action on and an accepted, unreleased story to publish; with none pending, the run shows the evaluation stopping the instance cleanly.

### Planning

- **Touches.** The story declared none. Every touch comes from its five tasks.
  - **From layout:** `design/formulas/publish.md` in both copies and the template's version; `flai/cmd/release.go` and `release_evaluate_test.go` for `--require-met`, which `evaluateRelease` lacks today (it exits 0 whether or not the policy is met); serve's `formula.go` with a new `formula_publish_test.go`; the orchestrator's role file in both copies; and `orchestrate_permissions_test.go`.
  - **From design:** `strategic-agents.md`, the conventions `strategic-agents.md` and `work-management.md` in both copies, the operator docs, the changelog, and `flai-reference.md`.
  - **Co-change, added:** `docs/users/flai.md` and `design/system/flai-cli.md` (48% each), for the new flag.
  - **Co-change, not added:** `docs/operators/settings.md` (no new key) and the dashboard's design.
  - The run task touches nothing. The instance's items are written by `flai serve` in the main checkout.
  - No folder touch.
- **Forecast.** Raised from flai's 36m to 50m for the live run, whose agent waits on the instance and its acceptance. The forecast's basis gives the figures, and its delivery is flai's play-out moved by any raise times the cycle factor.
- **Cost of delay.** Not set. Neither the story nor E-0024 has inputs, so `flai cod` refuses. The epic's input is asked on TH-0471, question 1, and the value follows its answer.
