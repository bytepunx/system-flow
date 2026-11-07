---
id: S-0562
type: story
nature: improvement
title: The publish formula is run once in this repository, its instance worked to review and accepted, and the release or the evaluation that stopped it recorded
status: in-progress
parent: E-0024
owner: alex
created: 2026-11-02T04:14:44Z
updated: 2026-11-07T09:40:04Z
transitions:
  - to: ready
    at: 2026-11-07T08:02:09Z
    by: alex
  - to: in-progress
    at: 2026-11-07T08:02:15Z
    by: agent-S-0562
blocked:
  - from: 2026-11-07T09:40:04Z
    reason: awaiting the operator on TH-0867
tags: [template]
topics: [orchestration]
touches: [wip/kanban/stories/S-0562-the-publish-formula-is-run-once-in-this-repository-its-instance-worked-to.md, design/issues, design/issues/summary.md, design/issues/I-0283-a-publish-instance-with-nothing-to-release-stops-at-its-evaluate-step-down-the.md]
after: [S-0423]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 5869
  priced_by: harness
  turns:
    - day: 2026-11-07
      ceremony: 4
      hand_edits: 2
      work: 114
  models:
    - model: claude-opus-5-5
      input: 248
      output: 35998
      cache_read: 15763471
      cache_write: 185202
      cost: 5.3553
  strategic:
    - kind: orchestrator
      seconds: 167
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 49
          output: 402
          cache_read: 1972514
          cache_write: 164020
          cost: 0.5267
cost_of_delay:
  inputs:
    penalty_per_week: 25
    by: orchestrator
    at: 2026-11-07T03:49:49Z
forecast:
  duration: 45m
  delivery: 2026-11-07T13:26:00Z
  basis: "Its own forecast of 45m; 12th in the pull order with an in-progress limit of 5, behind S-0406, S-0407, S-0408, S-0409, S-0410, S-0411, S-0412, S-0241, S-0543, S-0560 and S-0561."
  by: flai
  at: 2026-11-07T07:55:49Z
finalized:
  by: alex
  at: 2026-11-07T08:01:33Z
---
# S-0562 The publish formula is run once in this repository, its instance worked to review and accepted, and the release or the evaluation that stopped it recorded

## Goal

E-0024's `publish` formula, shipped by S-0423, is run once in this repository before the epic is accepted. This shows that a publish is a story worked, measured, and archived: the formula's `in: main` steps, its `--require-met` evaluation, and the `push` gate in `flai release --pending` all run for real.

## Acceptance criteria

- [ ] The formula is run once in this repository: its instance story is worked to review and accepted, and its Notes name the release it published, or, when nothing was pending, the evaluation that stopped it.
- [x] This story's Notes record the instance's ID, each step's exit, and the versions and tags published, or the evaluation that stopped the instance.
- [x] What the run finds wrong is fixed in a task of this story, or recorded as an issue.

## Tasks
- T-2504 Check the publish run's preconditions and start one instance of the publish formula from the main checkout
- T-2505 Follow the publish instance to review and acceptance, and record each step's exit and the release or the evaluation that stopped it
- T-2506 File each fault the publish run found as an issue, or as a fix task of S-0562, and say in the Notes when it found none

## Notes

- Moved out of S-0423 on TH-0803 (T-1783 and S-0423's criterion 4). The run needs `flai release --evaluate --require-met` and the `push` gate, which exist only once S-0423 is accepted, and the instance's acceptance is the operator's. This follows D-0422, where a run that needs the story's own change on the host moves to a follow-up story.
- `after: S-0423`: this story reads S-0423's output, the `publish` formula file and the flai changes its steps run.
- How to run it: with the `formula` and `push` host actions on, run `flai formula run publish --ready` in the main checkout, or let the orchestrator start it under `publish` when `release_evaluate` is met. `flai serve` starts the instance's agent, which runs the `in: main` steps from the main checkout. If those steps run the installed `flai` rather than `scripts/flai.sh`, the run waits until S-0423's flai is published too.
- Criterion 1 needs the operator's `push` host action on. When no accepted story is unreleased, the run shows the evaluation stopping the instance cleanly.

### Run

Preconditions, checked by T-2504 on 2026-11-07 at 08:03Z on the host:

| Precondition | Result |
|--------------|--------|
| Installed flai holds S-0423's change | Holds. `flai version` gives 1.84.0. `flai release --evaluate --require-met` exists, and `flai release --help` documents the `push` gate for agents `flai serve` started. |
| `push` host action on | On for this project (`flai serve actions`). |
| `formula` host action on | Off. It governs only the dashboard's scheduled runs, so a run from the shell does not need it. |
| TH-0811 answered | Not answered. This does not block the run: `release_evaluate` gave policy judgement, count 0, value 0 USD/week, with nothing accepted waiting for a release, because flai 1.84.1 had just been published. With nothing pending, step 2 stops under judgement or threshold alike, and criterion 1 accepts that stop as written. |

The live evaluation before the run, `flai release --evaluate --require-met` in the main checkout, exited 2: "release policy judgement not met: … with no accepted story waiting for a release".

Instance: S-0585, "Publish what is accepted and not yet released". It was started by hand by agent-S-0562 at 08:04Z with `flai formula run publish --ready --autocommit` in the main checkout. That made it ready, with tasks T-2546 (fetch), T-2548 (evaluate), and T-2549 (publish), and committed it as ff4e47115. No second instance was run.

What the instance's steps exited with, as S-0585's task Notes give them:

| Step | Task | Exit | What happened |
|------|------|------|---------------|
| 1 `fetch` | T-2546 | 0 | `scripts/fetch.sh` ran `git fetch --tags` in the main checkout. The remote had nothing new. Done at 08:06:07Z. |
| 2 `evaluate` | T-2548 | 2 | `flai release --evaluate --require-met` in the main checkout: "release policy judgement: not met … with no accepted story waiting for a release", count 0, value 0 USD/week. Blocked at 08:06:26Z as a failed step. On TH-0854 the operator confirmed agent-S-0562's recommendation, and the task was closed with the stop in its Notes at 08:15:59Z. |
| 3 `publish` | T-2549 | not run | Cancelled at 08:16:06Z: "step 2 stopped the run: nothing accepted is waiting for a release". Nothing was tagged or pushed, so the `push` gate was not reached. |

No release was published, and no versions or tags. The evaluation that stopped the instance is step 2's, above. S-0585's Notes name it.

S-0585 reached review at 08:16:55Z. The orchestrator could not accept it, because its criterion 3 (step 3) is unticked. It asked the operator on TH-0855 at 08:18Z.

Faults the run found, for criterion 3:

- I-0283, recorded on this story's branch and bumped once. A publish instance with nothing to release stops at step 2 down the formula's failure path, so it blocks and needs the operator's answer every time. Its unticked criterion for the skipped step then blocks acceptance. The fix is in flai, as a clean-stop exit for a step, so it is an issue and not a task of this story.
- I-0282, recorded by S-0564 on its branch. A formula instance made with no touches is held (no-touches) while any story is in progress, including the story that started it and waits on it. S-0585 was held this way until agent-S-0562 gave it its own story file as touches at 08:04Z.

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
