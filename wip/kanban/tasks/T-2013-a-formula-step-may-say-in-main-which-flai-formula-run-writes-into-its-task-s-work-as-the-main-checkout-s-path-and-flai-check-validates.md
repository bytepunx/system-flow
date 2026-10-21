---
id: T-2013
type: task
nature: improvement
title: "A formula step may say in: main, which flai formula run writes into its task's Work as the main checkout's path and flai check validates"
status: backlog
parent: S-0423
owner: alex
created: 2026-10-21T13:10:10Z
updated: 2026-10-21T13:10:10Z
transitions: []
stream: S-0423
tags: [cli]
touches: [flai/internal/formulas/formulas.go, flai/internal/formulas/formulas_test.go, flai/internal/check/check.go, flai/internal/check/formulas_test.go, flai/cmd/formula.go, flai/cmd/formula_test.go]
---
# T-2013 A formula step may say in: main, which flai formula run writes into its task's Work as the main checkout's path and flai check validates

## Work

A formula's step runs where the instance's agent works: its worktree (S-0422). `flai release --pending` must run in the main checkout, on the main branch with the host's credentials: `publishPending` in `flai/cmd/release.go` tags and pushes the branch that is checked out, so from a worktree it would tag and push the story branch.

Give a step an optional `in` field in `flai/internal/formulas`: `worktree`, the default, or `main`. `flai check` reports `formulas.file` on any other value. `flai formula run` writes into the `## Work` of a `main` step's task that the step runs from the project's main checkout, with its path, and leaves a `worktree` step's as it is. `flai formula show` prints a step's `in`.

Waits for nothing: first layer. Runs beside the `release_publish` task, since they share no path.

## Done when

`go test ./internal/formulas/... ./internal/check/... ./cmd/...` passes, with a case for each value and for the refusal of another. `flai formula run` on a fixture with a `main` step names the main checkout in that task's Work.

## Notes

Asked on the plan thread as question 1. The alternative is `flai release --pending` acting on the main checkout wherever it is run, which needs no field.
