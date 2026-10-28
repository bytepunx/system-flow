---
id: S-0424
type: story
nature: improvement
title: The template ships baseline formulas for importing a repository from the board and for the analyzer's cadence, the board's import starts its formula, and each is run once here
status: backlog
parent: E-0024
owner: alex
created: 2026-10-16T05:41:19Z
updated: 2026-10-28T12:40:13Z
transitions: []
tags: [cli, template, dashboard]
topics: [cli, orchestration, template, conventions]
touches: [design/workflows/formulas/analysis-cadence.md, design/workflows/formulas/import.md, template/root/design/workflows/formulas/analysis-cadence.md, template/root/design/workflows/formulas/import.md, template/template.yaml, template/CHANGELOG.md, flai/internal/hostapi/import.go, flai/internal/hostapi/import_test.go, flai/cmd/import.go, flai/cmd/import_commit.go, flai/cmd/import_test.go, flai/cmd/import_commit_test.go, flai/cmd/issue.go, flai/cmd/issue_test.go, flai/cmd/formula.go, flai/cmd/formula_test.go, flai/internal/formulas/formulas.go, flai/internal/formulas/formulas_test.go, flai/internal/check/check.go, flai/internal/check/formulas_test.go, flai/internal/serve/analysis_schedule.go, flai/internal/serve/analysis_schedule_test.go, flai/internal/manifest/manifest.go, flai/internal/manifest/manifest_test.go, design/system/strategic-agents.md, design/system/continuous-improvement.md, design/system/project-manifest.md, design/system/flai-cli.md, design/system/flaiover-dashboard.md, docs/operators/index.md, docs/operators/settings.md, docs/users/flaiover.md, docs/users/flai.md, docs/users/flai-reference.md]
after: [S-0422, S-0423]
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
      seconds: 734
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 100
          output: 1498
          cache_read: 24183926
          cache_write: 39724
          cost: 5.3487
draft: true
cost_of_delay:
  value: 50.85
  by: planner-S-0424
  at: 2026-10-21T16:11:12Z
forecast:
  duration: 1h10m
  delivery: 2026-10-29T06:21:00Z
  basis: "Its own forecast of 1h10m; 19th in the pull order with an in-progress limit of 5, behind S-0503, S-0241, S-0406, S-0407, S-0408, S-0409, S-0410, S-0411, S-0412, S-0414, S-0415, S-0416, S-0417, S-0418, S-0419, S-0420, S-0421, S-0422 and S-0423."
  by: flai
  at: 2026-10-28T12:40:13Z
---
# S-0424 The template ships baseline formulas for importing a repository from the board and for the analyzer's cadence, the board's import starts its formula, and each is run once here

## Goal

Two more processes that live in prose and Go become formulas. `import` brings a repository under system-flow as the board's Import does today, its steps shown as a story; `analysis-cadence` runs the analyzer and then lists what it filed for the operator's choice, as `continuous-improvement.md` § Cadence and `strategic-agents.md` § The analysis schedule describe. The board's Import starts the `import` formula, and `analysis.schedule` may start the cadence formula. Each is run once in this repository before the epic is accepted.

## Acceptance criteria

- [ ] `design/workflows/formulas/import.md`, in the template and here, takes the input `path` and has the steps the board's Import runs today (`flai import` with the path, the check, and the first commit), and the host API's import starts `flai formula run import --input path=<path>`, so its progress shows as a story; the import tests pass.
- [ ] `design/workflows/formulas/analysis-cadence.md`, in the template and here, takes the input `focus` and has a role step that runs the analyzer with it and a step that lists the open issues its report names, for the operator to choose stories from at acceptance; the analysis schedule may name the formula in place of the analyzer.
- [ ] `flai check --strict` reports nothing on either formula, here and in a project made from the template.
- [ ] Each formula is run once in this repository: an import of a scratch repository under a folder named for import, and one analyzer cadence, each instance worked to review and accepted, with its Notes naming what it produced.
- [ ] `design/system/strategic-agents.md` § The analysis schedule, `design/system/continuous-improvement.md` § Cadence, `docs/operators/index.md`, `docs/users/flaiover.md` where it describes Import, and `template/CHANGELOG.md` say that both are formulas.

## Tasks
- T-1785 design/workflows/formulas/import.md takes a path and runs the board's import steps, and the host API's import starts the formula so its progress shows as a story
- T-1786 design/workflows/formulas/analysis-cadence.md runs the analyzer with a focus and lists the open issues its report names, and analysis.schedule may name the formula
- T-1787 The import formula is run once from the board on a scratch repository under a folder named for import, its instance worked to review and accepted
- T-1788 The analysis-cadence formula is run once in this repository, its instance worked to review and accepted, with the report and the issues it names recorded
- T-1789 The design, the manifest's design, the operator and user docs, and the changelog say that import and the analyzer's cadence are formulas
- T-2016 A formula input may carry a default, which flai formula run takes for a missing input and flai check passes over when the formula is scheduled
- T-2017 flai import --formula commits the template's files before any test and starts the import formula in the new project, and flai import test runs the tests that --commit runs

## Notes

- Drafted by the planner for E-0024. Waits for the story that works a formula's instance; runs beside the publish formula story, with which it shares only the shared docs.
- The analyzer's cadence run spends a model session; the operator chooses its focus when finalizing this story, or the run takes none, which looks for all three.

### Planning

- **Plan thread.** TH-0588, 2026-10-21, approved by the orchestrator under `plan_backlog_stories`: the replan on TH-0471's answers 2 and 5, seven tasks in three layers, the new `after` on S-0423, and the 1h10m forecast. Two questions went to the operator as a recommendation to take the planner's answers, and the thread stays open for them: where the import formula's instance lives and what its steps are (`flai import --formula` committing the import and starting the formula in the new project, with `flai check --strict` and `flai import test` as the steps, T-2017), and the cadence's shape (a role step with `focus` defaulting to all three, then `flai issue list --report` `in: main`, and `analysis.formula` in the manifest, T-2016 and T-1786). Criterion 1's words change once the first is answered; T-2017 is dropped if it is answered otherwise.
- **Replanned on 2026-10-21** by the planner of S-0424 on TH-0471's answers. Answer 2 moves the formulas from `design/formulas/` to `design/workflows/formulas/`: criteria 1 and 2, every touch, and T-1785's and T-1786's titles and Work now say so. Answer 5 has a story's agent work an instance on its branch and worktree: T-1787's and T-1788's Work now say so, as approved on TH-0588, and each rewritten task carries a Notes line naming what still waits on the operator. The first Notes bullet predates the answers: the story now also waits for S-0423, which adds the step field `in: main` (T-2013) the cadence's listing step needs and changes `flai/internal/formulas`, `flai check`, and `flai formula run`, which T-2016 changes too.
- **Layers.** T-2016 and T-2017 first, sharing no path; then T-1785 after T-2017, since the host API switches to `flai import --formula`, and T-1786 after T-2016, since the cadence file declares `focus` with a default; then T-1787 after T-1785, T-1788 after T-1786, and T-1789 after both formula tasks. The two run tasks touch nothing, so the last layer shares no path.
- **Touches.** 34, all files. The story declared none of its own at first; every touch comes from its seven tasks.
  - **Declared, kept:** 16 of the 20 earlier touches.
  - **Declared, moved:** the two formula files in both copies, now under `design/workflows/formulas/` (answer 2), on T-1785 and T-1786.
  - **From layout, added:** `flai/cmd/import.go`, `import_commit.go`, and their tests, where `--commit` renders, tests, and commits and where `--formula` and `flai import test` go (T-2017); `flai/internal/formulas/formulas.go`, `flai/cmd/formula.go`, and `flai/internal/check/check.go` with their tests, for an input's `default` (T-2016); `flai/cmd/issue.go` and its test, for `flai issue list --report`, since the listing step needs the issues that name the report and `issues.NamingReport` has no command (T-1786).
  - **From design:** `strategic-agents.md`, `continuous-improvement.md`, `project-manifest.md`, the operator docs and `docs/operators/settings.md` for `analysis.formula`, `docs/users/flaiover.md` for the board's Import, and the changelog (T-1789).
  - **Co-change, kept:** `design/system/flai-cli.md` and `design/system/flaiover-dashboard.md`, which describe the import flow and the board's Import.
  - **Co-change, added:** `docs/users/flai.md` (48% of the seed commits) and `docs/users/flai-reference.md` (22%): `flai import --formula`, `flai import test`, `flai issue list --report`, and `analysis.formula` are command and setting changes, which the documentation convention has documented in the same story.
  - **Co-change, not added:** `design/adrs/README.md` (no ADR here: the formula file is S-0421's and `in: main` S-0423's), `design/system/workflow.md` and `work-management.md` (an instance is worked as any story, S-0422's), `flai/internal/harness/harness.go` (the instance's prompt is S-0422's), `flai/internal/hostapi/writes.go` (`import.run` lives in `hostapi/import.go`), and `flai/cmd/serve_actions.go` (no new host action).
  - **Checked, not added:** `flai/internal/serve/formula_schedule.go`, S-0422's schedule starter, which T-1786 calls for `analysis.formula`; `flai task done` widens the touches if it changes. `flai/internal/importer`: `DetectTests` is reused as it is.
  - No folder touch.
- **Forecast.** flai's 58m at size 39 (5 criteria, 34 touches), 89 s a unit over 86 done improvement stories on this model in the large band, raised to 1h10m for the two live runs the touches do not count: the analyzer session the cadence instance's agent waits on, and the import of a scratch repository, started from the board, worked by another agent, and accepted, as the earlier plan raised 36m to 50m for the same runs and S-0423 raised 53m to 1h for one. The delivery is flai's play-out after S-0423, 19th in the pull order: S-0423's delivery of 2026-10-22T07:44Z plus 1h10m times the cycle factor 1.15, 2026-10-22T09:05Z.
- **Cost of delay.** 50.85 USD a week, worked out by hand from TH-0471's answer 1, as the planners of S-0414 to S-0423 did, because `flai cod` refuses while E-0024 has no inputs (TH-0567, question 3, confirmed; TH-0587 asks the operator to write it on the epic). 4h lost per cycle at `planning.hour_rate` 150 over the default 168h cycle gives E-0024 600 USD a week; S-0424's share is its 1h10m of the 13h46m forecast for E-0024's twelve open stories: 70 / 826 × 600. Once the input is on E-0024, `flai cod S-0424` gives the same figure while the forecasts stand.
