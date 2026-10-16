---
id: T-1786
type: task
nature: improvement
title: design/formulas/analysis-cadence.md runs the analyzer with a focus and lists the open issues its report names, and analysis.schedule may name the formula
status: backlog
parent: S-0424
owner: alex
created: 2026-10-16T05:48:59Z
updated: 2026-10-16T05:48:59Z
transitions: []
stream: S-0424
tags: [cli, template]
touches: [design/formulas/analysis-cadence.md, template/root/design/formulas/analysis-cadence.md, flai/internal/serve/analysis_schedule.go, flai/internal/serve/analysis_schedule_test.go, flai/internal/manifest/manifest.go, flai/internal/manifest/manifest_test.go, template/template.yaml]
---
# T-1786 design/formulas/analysis-cadence.md runs the analyzer with a focus and lists the open issues its report names, and analysis.schedule may name the formula

## Work

Write the `analysis-cadence` formula with the input `focus`, empty for all three. It has two steps:

1. A role step that runs the analyzer with the focus.
2. After it, `flai issue list --json` filtered to the issues that name the report, written into the instance's Notes for the operator's choice at acceptance.

Let `analysis.formula: analysis-cadence` in the manifest make the analysis schedule start the formula in place of the analyzer. Check the key as `manifest.analysis`. Bump the template's version.

First layer.

## Done when

`analysis_schedule_test.go` and `manifest_test.go` pass, their existing cases unchanged. `flai check --strict` reports nothing on the formula.

## Notes
