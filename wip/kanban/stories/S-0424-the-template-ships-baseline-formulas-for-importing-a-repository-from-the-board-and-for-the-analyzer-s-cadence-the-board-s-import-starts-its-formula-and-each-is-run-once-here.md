---
id: S-0424
type: story
nature: improvement
title: The template ships baseline formulas for importing a repository from the board and for the analyzer's cadence, the board's import starts its formula, and each is run once here
status: backlog
parent: E-0024
owner: alex
created: 2026-10-16T05:41:19Z
updated: 2026-10-17T07:43:16Z
transitions: []
tags: [cli, template, dashboard]
topics: [cli, orchestration, template, conventions]
touches: [design/formulas/analysis-cadence.md, design/formulas/import.md, design/system/continuous-improvement.md, design/system/project-manifest.md, design/system/strategic-agents.md, docs/operators/index.md, docs/operators/settings.md, docs/users/flaiover.md, flai/internal/hostapi/import.go, flai/internal/hostapi/import_test.go, flai/internal/manifest/manifest.go, flai/internal/manifest/manifest_test.go, flai/internal/serve/analysis_schedule.go, flai/internal/serve/analysis_schedule_test.go, template/CHANGELOG.md, template/root/design/formulas/analysis-cadence.md, template/root/design/formulas/import.md, template/template.yaml, design/system/flai-cli.md, design/system/flaiover-dashboard.md]
after: [S-0422]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
draft: true
forecast:
  duration: 50m
  delivery: 2026-10-18T00:54:00Z
  basis: "Its own forecast of 50m; 32nd in the pull order with an in-progress limit of 5, behind S-0395, S-0396, S-0397, S-0393, S-0394, S-0241, S-0398, S-0399, S-0400, S-0401, S-0402, S-0403, S-0404, S-0405, S-0406, S-0407, S-0408, S-0409, S-0410, S-0411, S-0412, S-0414, S-0415, S-0416, S-0417, S-0418, S-0419, S-0420, S-0421, S-0422 and S-0423."
  by: flai
  at: 2026-10-17T07:43:16Z
---
# S-0424 The template ships baseline formulas for importing a repository from the board and for the analyzer's cadence, the board's import starts its formula, and each is run once here

## Goal

Two more processes that live in prose and Go become formulas. `import` brings a repository under system-flow as the board's Import does today, its steps shown as a story; `analysis-cadence` runs the analyzer and then lists what it filed for the operator's choice, as `continuous-improvement.md` § Cadence and `strategic-agents.md` § The analysis schedule describe. The board's Import starts the `import` formula, and `analysis.schedule` may start the cadence formula. Each is run once in this repository before the epic is accepted.

## Acceptance criteria

- [ ] `design/formulas/import.md`, in the template and here, takes the input `path` and has the steps the board's Import runs today (`flai import` with the path, the check, and the first commit), and the host API's import starts `flai formula run import --input path=<path>`, so its progress shows as a story; the import tests pass.
- [ ] `design/formulas/analysis-cadence.md`, in the template and here, takes the input `focus` and has a role step that runs the analyzer with it and a step that lists the open issues its report names, for the operator to choose stories from at acceptance; the analysis schedule may name the formula in place of the analyzer.
- [ ] `flai check --strict` reports nothing on either formula, here and in a project made from the template.
- [ ] Each formula is run once in this repository: an import of a scratch repository under a folder named for import, and one analyzer cadence, each instance worked to review and accepted, with its Notes naming what it produced.
- [ ] `design/system/strategic-agents.md` § The analysis schedule, `design/system/continuous-improvement.md` § Cadence, `docs/operators/index.md`, `docs/users/flaiover.md` where it describes Import, and `template/CHANGELOG.md` say that both are formulas.

## Tasks
- T-1785 design/formulas/import.md takes a path and runs the board's import steps, and the host API's import starts the formula so its progress shows as a story
- T-1786 design/formulas/analysis-cadence.md runs the analyzer with a focus and lists the open issues its report names, and analysis.schedule may name the formula
- T-1787 The import formula is run once from the board on a scratch repository under a folder named for import, its instance worked to review and accepted
- T-1788 The analysis-cadence formula is run once in this repository, its instance worked to review and accepted, with the report and the issues it names recorded
- T-1789 The design, the manifest's design, the operator and user docs, and the changelog say that import and the analyzer's cadence are formulas

## Notes

- Drafted by the planner for E-0024. Waits for the story that works a formula's instance; runs beside the publish formula story, with which it shares only the shared docs.
- The analyzer's cadence run spends a model session; the operator chooses its focus when finalizing this story, or the run takes none, which looks for all three.

### Planning

- **Touches.** The story declared none. Every touch comes from its five tasks.
  - **From layout:** the two formula files in both copies, the template's version, the host API's `import.go` and its test, `analysis_schedule.go` and its test, and the manifest for `analysis.formula`.
  - **From design:** `strategic-agents.md`, `continuous-improvement.md`, `project-manifest.md`, the operator docs and `docs/operators/settings.md`, `docs/users/flaiover.md`, and the changelog.
  - **Co-change, added:** `design/system/flai-cli.md` (45%), and `design/system/flaiover-dashboard.md` (37%), which describes the board's Import.
  - **Co-change, not added:** `docs/users/flai.md` and `flai-reference.md`. No command or flag changes.
  - The two run tasks touch nothing.
  - No folder touch.
- **Forecast.** Raised from flai's 36m to 50m for two live runs, one an analyzer session. The forecast's basis gives the figures, and its delivery is flai's play-out moved by any raise times the cycle factor.
- **Cost of delay.** Not set. Neither the story nor E-0024 has inputs, so `flai cod` refuses. The epic's input is asked on TH-0471, question 1, and the value follows its answer.
