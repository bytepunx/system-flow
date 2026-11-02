---
id: T-1786
type: task
nature: improvement
title: design/workflows/formulas/analysis-cadence.md runs the analyzer with a focus and lists the open issues its report names, and analysis.schedule may name the formula
status: done
parent: S-0424
owner: alex
created: 2026-10-16T05:48:59Z
updated: 2026-11-02T06:19:54Z
transitions:
  - to: ready
    at: 2026-11-02T05:52:06Z
    by: agent-S-0424
  - to: in-progress
    at: 2026-11-02T05:52:07Z
    by: agent-S-0424
  - to: done
    at: 2026-11-02T06:19:54Z
    by: agent-S-0424
stream: S-0424
tags: [cli, template]
touches: [design/workflows/formulas/analysis-cadence.md, template/root/design/workflows/formulas/analysis-cadence.md, flai/internal/serve/analysis_schedule.go, flai/internal/serve/analysis_schedule_test.go, flai/internal/manifest/manifest.go, flai/internal/manifest/manifest_test.go, template/template.yaml, flai/cmd/issue.go, flai/cmd/issue_test.go, docs/users/flai-reference.md, docs/operators/settings.md, flai/internal/check/check.go, flai/internal/check/formulas_test.go]
after: [T-2016]
usage:
  source: log
  seconds: 1667
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 215
      output: 66455
      cache_read: 13296357
      cache_write: 413956
      cost: 6.3156
---
# T-1786 design/workflows/formulas/analysis-cadence.md runs the analyzer with a focus and lists the open issues its report names, and analysis.schedule may name the formula

## Work

Write the `analysis-cadence` formula, `design/workflows/formulas/analysis-cadence.md` here and `template/root/design/workflows/formulas/analysis-cadence.md` in the template, with the input `focus`, whose default is empty, all three focuses (TH-0588, question 2). Its steps:

1. A role step that runs the analyzer with the focus: `role: analyzer` with a prompt naming `{{focus}}`.
2. After it, `flai issue list --report <report> --json`, `in: main`, listing the open issues the report names. The body tells the instance's agent to put the report's path in the step's command, and to record the report and the issues in the instance's Notes for the operator's choice at acceptance.

Add `flai issue list --report <path>` (`flai/cmd/issue.go`), filtering to the issues an instance names that report (`issues.NamingReport`), and regenerate `docs/users/flai-reference.md` with `make flai-reference`.

Add `analysis.formula` to the manifest (`flai/internal/manifest/manifest.go`): a formula's name, checked as `manifest.analysis`. When it is set, the analysis scheduler (`flai/internal/serve/analysis_schedule.go`) starts that formula when `analysis.schedule` comes round, by running `flai formula run <name> --ready --autocommit --json` in the main checkout through `os.Executable` as `serve/formula_schedule.go` does (D-0459), in place of starting the analyzer, behind the same `analyze` host action. `flai check` reports `manifest.analysis` on a formula name no file defines, or one with an input that carries no default (`flai/internal/check/check.go`). Bump the template's version in `template/template.yaml`.

Second layer: waits for T-2016, since the file declares `focus` with a default and may be scheduled. Runs beside T-1785, since they share no path.

## Done when

`analysis_schedule_test.go`, `manifest_test.go`, `issue_test.go`, and `check/formulas_test.go` pass, their existing cases unchanged, with cases for the formula started in place of the analyzer and for `--report`. `flai check --strict` reports nothing on the formula here, and the template test renders it.

## Notes

Rewritten on 2026-10-21 by the planner of S-0424 for TH-0471's answer 2, and on 2026-11-02 by agent-S-0424 for TH-0588's answer to question 2.
