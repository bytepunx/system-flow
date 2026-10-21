---
id: T-2017
type: task
nature: improvement
title: flai import --formula commits the template's files before any test and starts the import formula in the new project, and flai import test runs the tests that --commit runs
status: backlog
parent: S-0424
owner: alex
created: 2026-10-21T16:09:19Z
updated: 2026-10-21T16:09:19Z
transitions: []
stream: S-0424
tags: [cli]
touches: [flai/cmd/import.go, flai/cmd/import_commit.go, flai/cmd/import_test.go, flai/cmd/import_commit_test.go, docs/users/flai-reference.md]
---
# T-2017 flai import --formula commits the template's files before any test and starts the import formula in the new project, and flai import test runs the tests that --commit runs

## Work

- `flai import --formula` (`flai/cmd/import.go`, `import_commit.go`): after the template is rendered, the manifest written, and the check run, it commits the import's paths as `--commit` does, with the `--trailer` lines, but runs no test first; has the host serve the project as today; then runs the `import` formula in the new project, as `flai formula run import --input path=<root> --ready` does, and answers with the instance's ID beside what `--json` gives today (`formula: {story: S-nnnn}`). It refuses, before anything changes, a repository that is not git or has uncommitted changes, as `--commit` does, and `--formula` with `--commit`. `--commit` stays as it is.
- `flai import test [path]` runs the tests `--commit` runs, the host's checks when any are set (`flai serve checks`), else what the repository has, prints each outcome, and exits non-zero when one fails, so that a step of the `import` formula can run them in the instance's worktree.

First layer: runs beside the input-default task; they share no path. The host API's switch from `--commit` to `--formula` is the import formula task's, which waits for this.

## Done when

`import_test.go` and `import_commit_test.go` pass, their existing cases unchanged, with cases for `--formula` committing the import before any test and making the instance ready, for `--formula` with `--commit` refused, and for `flai import test` passing and failing. `docs/users/flai-reference.md` is regenerated with `make flai-reference`. `flai check --strict` reports nothing.

## Notes

Drafted by the planner of S-0424 on 2026-10-21 on the plan thread's recommended answer to where the import formula's instance lives: in the imported repository, made after the import's commit, since no project exists before `flai import` runs and a story's agent works an instance on a branch from main (TH-0471, answer 5). Dropped if the operator answers otherwise.
