---
id: S-0562
type: story
nature: improvement
title: The publish formula is run once in this repository, its instance worked to review and accepted, and the release or the evaluation that stopped it recorded
status: backlog
parent: E-0024
owner: alex
created: 2026-11-02T04:14:44Z
updated: 2026-11-07T04:45:40Z
transitions: []
tags: [template]
topics: [orchestration]
touches: [wip/kanban/stories/S-0562-the-publish-formula-is-run-once-in-this-repository-its-instance-worked-to.md, design/issues, design/issues/summary.md]
after: [S-0423]
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
      seconds: 57
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 12
          output: 73
          cache_read: 761554
          cache_write: 10270
          cost: 0.1863
draft: true
cost_of_delay:
  inputs:
    penalty_per_week: 25
    by: orchestrator
    at: 2026-11-07T03:49:49Z
forecast:
  duration: 45m
  delivery: 2026-11-07T10:45:00Z
  basis: "Its own forecast of 45m; 14th in the pull order with an in-progress limit of 5, behind S-0565, S-0567, S-0569, S-0570, S-0406, S-0407, S-0408, S-0409, S-0410, S-0411, S-0412, S-0241, S-0543, S-0556, S-0559, S-0560 and S-0561."
  by: flai
  at: 2026-11-07T04:45:40Z
---
# S-0562 The publish formula is run once in this repository, its instance worked to review and accepted, and the release or the evaluation that stopped it recorded

## Goal

E-0024's `publish` formula, shipped by S-0423, is run once in this repository before the epic is accepted. This shows that a publish is a story worked, measured, and archived: the formula's `in: main` steps, its `--require-met` evaluation, and the `push` gate in `flai release --pending` all run for real.

## Acceptance criteria

- [ ] The formula is run once in this repository: its instance story is worked to review and accepted, and its Notes name the release it published, or, when nothing was pending, the evaluation that stopped it.
- [ ] This story's Notes record the instance's ID, each step's exit, and the versions and tags published, or the evaluation that stopped the instance.
- [ ] What the run finds wrong is fixed in a task of this story, or recorded as an issue.

## Tasks
- T-2504 Check the publish run's preconditions and start one instance of the publish formula from the main checkout
- T-2505 Follow the publish instance to review and acceptance, and record each step's exit and the release or the evaluation that stopped it
- T-2506 File each fault the publish run found as an issue, or as a fix task of S-0562, and say in the Notes when it found none

## Notes

- Moved out of S-0423 on TH-0803 (T-1783 and S-0423's criterion 4). The run needs `flai release --evaluate --require-met` and the `push` gate, which exist only once S-0423 is accepted, and the instance's acceptance is the operator's. This follows D-0422, where a run that needs the story's own change on the host moves to a follow-up story.
- `after: S-0423`: this story reads S-0423's output, the `publish` formula file and the flai changes its steps run.
- How to run it: with the `formula` and `push` host actions on, run `flai formula run publish --ready` in the main checkout, or let the orchestrator start it under `publish` when `release_evaluate` is met. `flai serve` starts the instance's agent, which runs the `in: main` steps from the main checkout. If those steps run the installed `flai` rather than `scripts/flai.sh`, the run waits until S-0423's flai is published too.
- Criterion 1 needs the operator's `push` host action on. When no accepted story is unreleased, the run shows the evaluation stopping the instance cleanly.

### Planning

Planned by planner-S-0562 on 2026-11-07. The plan's thread on S-0562 names the tasks, their layers, the assumptions, and the questions for the operator.

What this story's own branch changes. The instance is another story, made by `flai formula run` in the main checkout and worked by its own agent. Its steps tag and push from the main checkout. So this story's branch holds only the issues and fix tasks the run leads to. Its Notes, where criterion 2's record goes, are written with `flai edit` in the main checkout.

Tasks, in three layers, one task each:

1. T-2504 checks the preconditions and starts one instance (criterion 1).
2. T-2505 follows the instance to acceptance and records each step's exit and the release or the evaluation (criteria 1 and 2). It waits for T-2504.
3. T-2506 files each fault the run found as an issue or as a fix task of this story (criterion 3). It waits for T-2505.

Touches. The story declared none. `flai touches suggest S-0562` needs a seed, so it ran on this story's file, `design/issues/summary.md`, and `design/workflows/formulas/publish.md` (713 of 3012 commits). Its top paths are older issue files, `design/system/flai-cli.md`, and `docs/users/flai.md`, at 6% or less. Each changes with flai's commands, which this story does not change, so none is added.

- `wip/kanban/stories/S-0562-the-publish-formula-is-run-once-in-this-repository-its-instance-worked-to.md` (layout): criterion 2's record, written by T-2504 and T-2505.
- `design/issues`, folder kept (layout): criterion 3 makes each fault an issue with `flai issue new`, which names each file from its title. No task can name those files yet, so the folder stays on the story alone. `flai task done` widens the touches with each issue committed.
- `design/issues/summary.md` (layout, co-change): `flai issue new` regenerates it. It is in the manifest's shared paths.

A fix T-2506 writes as a task names its own files, such as `design/workflows/formulas/publish.md` and `template/root/design/workflows/formulas/publish.md`. `flai task done` widens the touches with them.

`after`. S-0423, declared: this story runs S-0423's `publish` formula file, `flai release --evaluate --require-met`, and the `push` gate in `flai release --pending`. S-0423 is accepted, and from flai 1.80.2 the installed flai holds its change. The host ran 1.83.0 on 2026-11-07.

TH-0811 bears on criterion 1 and is still open. This project's release policy is `judgement`, because `orchestration.release` has `value: 1000` and no `policy` key. Under judgement, `flai release --evaluate --require-met` always stops, so every instance stops at step 2 even with accepted stories pending. Criterion 1 allows a stopped evaluation only when nothing was pending. T-2504 therefore runs nothing until TH-0811 is answered.

Forecast. `flai forecast S-0562` gave 4m: median 80 s per unit of size over 13 done improvement stories on claude-opus-5-5, size 3 (3 criteria, 0 touches), 17th in the pull order. I raised it to 45m of agent time. The size counts no waiting, and the work is a precondition check, one instance started, its run watched to review and acceptance, its exits recorded, and its faults filed. The instance's agent time belongs to the instance, not to this story. The delivery is flai's 2026-11-07T09:02Z moved by the added time, 2026-11-07T09:45Z. It holds only once TH-0811 is answered and the operator finalizes this draft.

Cost of delay. `flai cod S-0562` cannot give a figure. Neither this story nor E-0024 has cost of delay inputs, so none is written. The plan's thread asks the operator for them.
