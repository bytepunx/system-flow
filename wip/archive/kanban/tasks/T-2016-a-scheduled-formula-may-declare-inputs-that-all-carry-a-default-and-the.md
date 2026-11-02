---
id: T-2016
type: task
nature: improvement
title: A scheduled formula may declare inputs that all carry a default, and the template test passes over a formula's input placeholders
status: done
parent: S-0424
owner: alex
created: 2026-10-21T16:09:10Z
updated: 2026-11-02T05:48:20Z
transitions:
  - to: ready
    at: 2026-11-02T05:31:52Z
    by: agent-S-0424
  - to: in-progress
    at: 2026-11-02T05:31:55Z
    by: agent-S-0424
  - to: done
    at: 2026-11-02T05:48:20Z
    by: agent-S-0424
stream: S-0424
tags: [cli]
touches: [flai/internal/formulas/formulas.go, flai/internal/formulas/formulas_test.go, flai/cmd/formula.go, flai/cmd/formula_test.go, flai/internal/check/check.go, flai/internal/check/formulas_test.go, flai/cmd/new_test.go]
usage:
  source: log
  seconds: 985
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 23
      output: 6748
      cache_read: 1717441
      cache_write: 50869
      cost: 0.7909
---
# T-2016 A scheduled formula may declare inputs that all carry a default, and the template test passes over a formula's input placeholders

## Work

An input's `default` is read already (S-0421): `formulas.Input.Default`, `Required`, and `flai formula run` taking a default for an input not given. What is left:

- The `formulas.file` finding on a `schedule` (`checker.schedule` in `flai/internal/formulas/formulas.go`) passes over an input that carries a default, and still names each input that has none. A scheduled formula may then take an input such as the cadence formula's `focus`, which the analysis schedule starts with none.
- `TestRenderPrototypeTemplate` in `flai/cmd/new_test.go` passes over the `{{name}}` input placeholders of a file under `design/workflows/formulas/` (ADR-0238 §4), and still refuses `{{` anywhere else in a rendered project (I-0249). The import and cadence formulas the template ships use an input.
- Check that `flai formula run` (`flai/cmd/formula.go`) and `flai formula show` say the default; add a case only where one is missing.

First layer: runs beside T-2017; they share no path. T-1786 waits for this, since its file declares `focus` with an empty default and may be scheduled.

## Done when

`formulas_test.go`, `formula_test.go`, `check/formulas_test.go`, and `new_test.go` pass, their existing cases unchanged, with cases for a scheduled formula whose inputs all carry defaults accepted, one with an input without a default refused, and a rendered formula file with `{{name}}` passed over. `flai check --strict` reports nothing.

## Notes

Drafted by the planner of S-0424 on 2026-10-21, so that the `analysis-cadence` formula can be started from `analysis.schedule` with no input. Rewritten by agent-S-0424 on 2026-11-02: S-0421 already reads the default, and I-0249's test fix joins this task.
