---
id: T-2744
type: task
nature: improvement
title: A formula step's stop lists the exits that stop its instance cleanly, flai check validates it, and run and show write it
status: backlog
parent: S-0621
owner: alex
created: 2026-11-09T17:06:37Z
updated: 2026-11-09T17:06:37Z
transitions: []
stream: S-0621
tags: [flai]
touches: [flai/internal/formulas/formulas.go, flai/internal/formulas/formulas_test.go, flai/internal/check/formulas_test.go, flai/cmd/formula.go, flai/cmd/formula_test.go]
after: [T-2743]
---
# T-2744 A formula step's stop lists the exits that stop its instance cleanly, flai check validates it, and run and show write it

## Work

Add the step field T-2743's ADR settles, proposed `stop`, to the formula format.

- In `flai/internal/formulas/formulas.go`, add it to `Step` (about line 128) and to `Step.UnmarshalYAML`'s shape message (about line 336). Today the loader reads past an unknown key, so a formula that names `stop` parses on an older flai and ignores it.
- Check it in the step checks (about lines 712 to 738, beside the `in` check): each exit is an integer from 1 to 255, none twice, and the step's `run` is a `flai` or `script` step. Each problem is a `formulas.file` error on the step's line.
- In `flai/cmd/formula.go`, `formulaTaskSections` (about line 565) writes the stop into the task's `## Done when`: the step exits 0, or exits with one of its stop exits, which stops the run cleanly through the ADR's command. `printFormula` (about line 131) prints the stop beside `in:`.
- Tests: `flai/internal/formulas/formulas_test.go` (parse, and each validation problem), `flai/internal/check/formulas_test.go` (a bad stop is an error on its line), `flai/cmd/formula_test.go` (`TestFormulaRun`'s Done-when text with and without a stop, and `show`).

This task waits for T-2743, whose ADR names the field and its rules.

## Done when

- A formula step with a valid stop parses, `flai check` passes it, and `flai formula run` writes its stop exits into the step task's `## Done when`.
- Each invalid stop is a `formulas.file` error on its line, each with a test.
- `flai test` on the changed paths passes.

## Notes
