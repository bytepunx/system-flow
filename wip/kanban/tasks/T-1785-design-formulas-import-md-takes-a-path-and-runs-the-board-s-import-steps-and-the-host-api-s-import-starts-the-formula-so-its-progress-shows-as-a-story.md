---
id: T-1785
type: task
nature: improvement
title: design/formulas/import.md takes a path and runs the board's import steps, and the host API's import starts the formula so its progress shows as a story
status: backlog
parent: S-0424
owner: alex
created: 2026-10-16T05:48:54Z
updated: 2026-10-16T05:48:54Z
transitions: []
stream: S-0424
tags: [cli, template, dashboard]
touches: [design/formulas/import.md, template/root/design/formulas/import.md, flai/internal/hostapi/import.go, flai/internal/hostapi/import_test.go]
---
# T-1785 design/formulas/import.md takes a path and runs the board's import steps, and the host API's import starts the formula so its progress shows as a story

## Work

Write the `import` formula with the input `path`. Its steps are the ones the host API's import runs today, as `flai/internal/hostapi/import.go`, `flai/cmd/import.go`, and `import_commit.go` do them: `flai import` on the path, the check, and the first commit. Make the host API's import run `flai formula run import --input path=<path>` and return the instance's ID beside what it returns today, keeping its refusals. The board shows the instance as a story.

First layer: runs beside the cadence formula task, since they share no path.

## Done when

`import_test.go` passes, its existing cases unchanged, with a case for the instance. `flai check --strict` reports nothing on the formula.

## Notes
