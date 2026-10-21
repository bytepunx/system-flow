---
id: T-2016
type: task
nature: improvement
title: A formula input may carry a default, which flai formula run takes for a missing input and flai check passes over when the formula is scheduled
status: backlog
parent: S-0424
owner: alex
created: 2026-10-21T16:09:10Z
updated: 2026-10-21T16:09:10Z
transitions: []
stream: S-0424
tags: [cli]
touches: [flai/internal/formulas/formulas.go, flai/internal/formulas/formulas_test.go, flai/cmd/formula.go, flai/cmd/formula_test.go, flai/internal/check/check.go, flai/internal/check/formulas_test.go]
---
# T-2016 A formula input may carry a default, which flai formula run takes for a missing input and flai check passes over when the formula is scheduled

## Work

An input a formula declares in its front matter may carry `default`, read in `flai/internal/formulas/formulas.go`. `flai formula run` (`flai/cmd/formula.go`) takes the default for an input not given with `--input`, and still refuses a missing input that has none, before anything is written. The `formulas.file` finding that refuses a `schedule` on a formula with inputs (S-0422) passes over an input that carries a default, so a scheduled formula may take an input such as the cadence formula's `focus`, which the analysis schedule starts with none.

First layer: runs beside the `flai import --formula` task; they share no path. The cadence formula task waits for this, since its file declares `focus` with an empty default.

## Done when

`formulas_test.go`, `formula_test.go`, and `check/formulas_test.go` pass, their existing cases unchanged, with cases for a defaulted input taken when it is not given, a missing input without a default refused, and a scheduled formula whose inputs all carry defaults accepted. `flai check --strict` reports nothing.

## Notes

Drafted by the planner of S-0424 on 2026-10-21, so that the `analysis-cadence` formula can be started from `analysis.schedule` with no input.
