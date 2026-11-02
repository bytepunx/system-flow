---
id: T-1785
type: task
nature: improvement
title: design/workflows/formulas/import.md takes a path and runs the check and the tests after the import's commit, and the host API's import runs flai import --formula so its progress shows as a story
status: in-progress
parent: S-0424
owner: alex
created: 2026-10-16T05:48:54Z
updated: 2026-11-02T05:52:04Z
transitions:
  - to: ready
    at: 2026-11-02T05:52:01Z
    by: agent-S-0424
  - to: in-progress
    at: 2026-11-02T05:52:04Z
    by: agent-S-0424
stream: S-0424
tags: [cli, template, dashboard]
touches: [design/workflows/formulas/import.md, template/root/design/workflows/formulas/import.md, flai/internal/hostapi/import.go, flai/internal/hostapi/import_test.go]
after: [T-2017]
---
# T-1785 design/workflows/formulas/import.md takes a path and runs the check and the tests after the import's commit, and the host API's import runs flai import --formula so its progress shows as a story

## Work

Write the `import` formula, `design/workflows/formulas/import.md` here and `template/root/design/workflows/formulas/import.md` in the template, with the input `path`, the imported repository's root, used in the formula's title. Its steps are the ones the import runs after its commit (TH-0588, question 1): `flai: check --strict`, then `flai: import test`, each from the instance's worktree. The body tells the instance's agent what a failing test means: block the task, quote the output, and ask on a thread, as any instance does.

Make the host API's `import.run` (`flai/internal/hostapi/import.go`) run `flai import . --yes --formula` with the same `--var` and `--trailer` arguments in place of `--commit`, describe the outcome with the instance's ID (`describeImport`), and keep its refusals; the `NotCommitted` exit no longer comes from `--formula`, so drop or keep its mapping by what `flai import --formula` exits. The board shows the instance as a story of the imported project.

Second layer: waits for T-2017, which gives `flai import` the `--formula` flag. Runs beside T-1786, since they share no path.

## Done when

`hostapi/import_test.go` passes, its existing cases changed only where `--commit` became `--formula`, with a case for the instance in the outcome. `flai check --strict` reports nothing on the formula here, and the template test renders it.

## Notes

Rewritten on 2026-10-21 by the planner of S-0424 for TH-0471's answer 2, and on 2026-11-02 by agent-S-0424 for TH-0588's answer to question 1.
