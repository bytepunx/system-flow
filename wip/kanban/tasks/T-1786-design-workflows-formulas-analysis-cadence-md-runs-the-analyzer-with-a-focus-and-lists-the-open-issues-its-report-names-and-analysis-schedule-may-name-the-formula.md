---
id: T-1786
type: task
nature: improvement
title: design/workflows/formulas/analysis-cadence.md runs the analyzer with a focus and lists the open issues its report names, and analysis.schedule may name the formula
status: backlog
parent: S-0424
owner: alex
created: 2026-10-16T05:48:59Z
updated: 2026-10-21T16:12:31Z
transitions: []
stream: S-0424
tags: [cli, template]
touches: [design/workflows/formulas/analysis-cadence.md, template/root/design/workflows/formulas/analysis-cadence.md, flai/internal/serve/analysis_schedule.go, flai/internal/serve/analysis_schedule_test.go, flai/internal/manifest/manifest.go, flai/internal/manifest/manifest_test.go, template/template.yaml, flai/cmd/issue.go, flai/cmd/issue_test.go, docs/users/flai-reference.md]
after: [T-2016]
---
# T-1786 design/workflows/formulas/analysis-cadence.md runs the analyzer with a focus and lists the open issues its report names, and analysis.schedule may name the formula

## Work

Write the `analysis-cadence` formula, `design/workflows/formulas/analysis-cadence.md` here and `template/root/design/workflows/formulas/analysis-cadence.md` in the template, with the input `focus`, empty for all three. It has two steps:

1. A role step that runs the analyzer with the focus.
2. After it, `flai issue list --json` filtered to the issues that name the report, written into the instance's Notes for the operator's choice at acceptance.

Let `analysis.formula: analysis-cadence` in the manifest make the analysis schedule start the formula in place of the analyzer. Check the key as `manifest.analysis`. Bump the template's version.

Second layer: waits for T-2016, since the file declares `focus` with a default so that the schedule can start the formula with no input. Runs beside the import formula task, since they share no path.

## Done when

`analysis_schedule_test.go` and `manifest_test.go` pass, their existing cases unchanged. `flai check --strict` reports nothing on the formula.

## Notes

Rewritten on 2026-10-21 by the planner of S-0424 for TH-0471's answer 2, the folder, as approved on TH-0588. The listing step's `in: main`, a `--report` filter on `flai issue list`, and the input's default in the file (TH-0588, question 2) are the operator's to answer; the Work above keeps the earlier words until then.
