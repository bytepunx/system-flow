---
id: S-0421
type: story
nature: feature
title: A formula is a markdown file under design/formulas that flai check validates, and flai formula run makes a story from it with one task per step and the steps' after between them
status: backlog
parent: E-0024
owner: alex
created: 2026-10-16T05:40:40Z
updated: 2026-10-21T09:04:46Z
transitions: []
tags: [cli, template]
topics: [cli, orchestration, template, conventions]
touches: [design/adrs, design/adrs/README.md, design/conventions/strategic-agents.md, design/formulas/README.md, design/system/repository-layout.md, design/system/strategic-agents.md, design/system/work-hierarchy.md, docs/operators/index.md, docs/users/flai.md, docs/users/flai-reference.md, flai/cmd/formula.go, flai/cmd/formula_test.go, flai/cmd/root.go, flai/internal/check/check.go, flai/internal/check/formulas_test.go, flai/internal/formulas/formulas.go, flai/internal/formulas/formulas_test.go, flai/internal/workitem/fields_test.go, flai/internal/workitem/front-matter-fields.txt, flai/internal/workitem/item.go, template/CHANGELOG.md, template/root/design/conventions/strategic-agents.md, template/root/design/formulas/README.md, template/template.yaml, design/system/flai-cli.md, design/system/project-manifest.md]
after: [S-0420]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
draft: true
forecast:
  duration: 51m
  delivery: 2026-10-21T20:33:00Z
  basis: "Its own forecast of 51m; 19th in the pull order with an in-progress limit of 5, behind S-0459, S-0461, S-0462, S-0460, S-0241, S-0406, S-0407, S-0408, S-0409, S-0410, S-0411, S-0412, S-0414, S-0415, S-0416, S-0417, S-0418, S-0419 and S-0420."
  by: flai
  at: 2026-10-21T09:04:46Z
---
# S-0421 A formula is a markdown file under design/formulas that flai check validates, and flai formula run makes a story from it with one task per step and the steps' after between them

## Goal

A repeatable process is a file, not prose. A formula under `design/formulas/` is a named sequence of steps with `after` between them, each a flai command, a script under `scripts/`, or a role run with a prompt and inputs, and the inputs the formula takes. `flai formula run <name>` makes a story from it, its steps as tasks in the same layers, so that it is planned, measured, and archived as any story is. `flai check` refuses a formula that could not run.

## Acceptance criteria

- [ ] An ADR records the formula file: where it lives (`design/formulas/<name>.md`, or the folder TH-0471 settles), its front matter (name, description, nature, inputs, and steps with `id`, `title`, `run` as a flai command, a script, or a role with its prompt, and `after`), its body as the formula's description, the baseline marker, and the instance's nature and its `formula` field as TH-0471 settles them.
- [ ] `flai/internal/formulas` reads a project's formulas, and `flai formula list [--json]` and `flai formula show <name> [--json]` print each with its steps and layers.
- [ ] `flai formula run <name> [--input key=value]... [--epic E-nnnn] [--json]` makes a story with the formula's nature and `formula: <name>`, one task per step with `## Work` naming what the step runs and `## Done when` its exit, the steps' `after` as the tasks' `after`, through `itemnew` so that `flai check --strict` passes on it, and prints the story's ID. A missing input, or one the formula does not declare, is refused before anything is written.
- [ ] `flai check` reports `formulas.file` for a formula that does not parse or misses a field, a step whose `run` is none of the three kinds, a script not under `scripts/`, a role no file defines, an `after` naming no step, a cycle, and an input a step uses that the formula does not declare.
- [ ] `formula` is listed in `front-matter-fields.txt` and in `design/system/work-hierarchy.md`'s schema, and an older flai reads past it.
- [ ] `design/system/repository-layout.md`, a new section of `design/system/strategic-agents.md` or a design document of its own, `docs/users/flai.md` with `flai-reference.md`, `docs/operators/index.md`, the convention `strategic-agents.md`, and `template/CHANGELOG.md` with the template's `design/formulas/README.md` say what a formula holds and how to run one.

## Tasks
- T-1768 An ADR records the formula file: where it lives, its steps and inputs, its body, its marker, and the instance's nature and formula field
- T-1769 flai/internal/formulas reads and validates a project's formulas into steps, inputs, and layers
- T-1770 A story carries formula in its front matter, listed in front-matter-fields.txt and the work hierarchy's schema, and an older flai reads past it
- T-1771 flai check reports formulas.file on each formula that could not run
- T-1772 flai formula list, show, and run: run makes a story with the formula's nature and one task per step, through itemnew, and prints its ID
- T-1773 The template and this repository have design/formulas with its README, and the design, the convention, and the docs say what a formula holds and how to run one

## Notes

- Drafted by the planner for E-0024. Follows the epic's order: after the dashboard's role page, though it reads only the role file story's loader for the role steps it validates.
- This story makes the instance; the next one works it.

### Planning

- **Touches.** The story declared none. Every touch comes from its six tasks.
  - **From layout:** a new `flai/internal/formulas` package; `flai/cmd/formula.go` and `root.go`; the check with a new `formulas_test.go`; `workitem/item.go`, `front-matter-fields.txt`, and `fields_test.go`; the formulas folder's READMEs; and the template's version.
  - **From design:** the ADR's index, `work-hierarchy.md`, `repository-layout.md`, `strategic-agents.md`, the convention in both copies, the user and operator docs, and the changelog.
  - **Co-change, added:** `design/system/flai-cli.md` (44%) and `design/system/project-manifest.md` (6%), the latter for the folder the manifest may name.
  - **Co-change, not added:** the dashboard's documents and the harness. A formula changes neither here.
  - **Folder touch kept:** `design/adrs`, for the ADR's new file; it is shared.
- **Forecast.** Kept at flai's 51m. The forecast's basis gives the figures, and its delivery is flai's play-out moved by any raise times the cycle factor.
- **Cost of delay.** Not set. Neither the story nor E-0024 has inputs, so `flai cod` refuses. The epic's input is asked on TH-0471, question 1, and the value follows its answer.
