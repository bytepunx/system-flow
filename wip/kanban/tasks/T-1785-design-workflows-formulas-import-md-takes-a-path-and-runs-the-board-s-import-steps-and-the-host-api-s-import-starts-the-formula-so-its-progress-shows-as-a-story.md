---
id: T-1785
type: task
nature: improvement
title: design/workflows/formulas/import.md takes a path and runs the board's import steps, and the host API's import starts the formula so its progress shows as a story
status: backlog
parent: S-0424
owner: alex
created: 2026-10-16T05:48:54Z
updated: 2026-10-21T16:12:23Z
transitions: []
stream: S-0424
tags: [cli, template, dashboard]
touches: [design/workflows/formulas/import.md, template/root/design/workflows/formulas/import.md, flai/internal/hostapi/import.go, flai/internal/hostapi/import_test.go]
after: [T-2017]
---
# T-1785 design/workflows/formulas/import.md takes a path and runs the board's import steps, and the host API's import starts the formula so its progress shows as a story

## Work

Write the `import` formula, `design/workflows/formulas/import.md` here and `template/root/design/workflows/formulas/import.md` in the template, with the input `path`. Its steps are the ones the host API's import runs today, as `flai/internal/hostapi/import.go`, `flai/cmd/import.go`, and `import_commit.go` do them: `flai import` on the path, the check, and the first commit. Make the host API's import run `flai formula run import --input path=<path>` and return the instance's ID beside what it returns today, keeping its refusals. The board shows the instance as a story.

Second layer: waits for T-2017, which gives `flai import` the flag the host API's import switches to. Runs beside the cadence formula task, since they share no path.

## Done when

`import_test.go` passes, its existing cases unchanged, with a case for the instance. `flai check --strict` reports nothing on the formula.

## Notes

Rewritten on 2026-10-21 by the planner of S-0424 for TH-0471's answer 2, the folder, as approved on TH-0588. Which steps the formula has and how the instance is made in a repository that is no project yet (TH-0588, question 1: `flai import --formula`, with the check and `flai import test` as the steps after the import's commit) are the operator's to answer; the Work above keeps the earlier words until then.
