---
id: T-2740
type: task
nature: remediation
title: flai formula run gives the instance story its own story file as its touch, with a test that the instance is not held by a story in progress
status: done
parent: S-0620
owner: alex
created: 2026-11-09T17:05:59Z
updated: 2026-11-09T17:28:53Z
transitions:
  - to: ready
    at: 2026-11-09T17:11:47Z
    by: agent-S-0620
  - to: in-progress
    at: 2026-11-09T17:11:48Z
    by: agent-S-0620
  - to: done
    at: 2026-11-09T17:28:53Z
    by: agent-S-0620
stream: S-0620
tags: [flai, formulas, holds]
touches: [flai/cmd/formula.go, flai/cmd/formula_test.go]
usage:
  source: log
  seconds: 1025
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 71
      output: 10459
      cache_read: 8419054
      cache_write: 153138
      cost: 3.0823
---
# T-2740 flai formula run gives the instance story its own story file as its touch, with a test that the instance is not held by a story in progress

## Work

In `formulaRun` (`flai/cmd/formula.go`), give the instance story its own story file as its one touch once `itemnew.Create` has made it and its path is known. Do this before `formulaReady` moves it under `--ready`, and before the commit, so that the commit holds the touch. The story file is already in `made`, so the commit and the undo cover it. `flai formula run`, `flai import`'s import instance (`flai/cmd/import_commit.go`), and the orchestrator's `publish` all go through `formulaRun`.

In `flai/cmd/formula_test.go`, add a test that reproduces I-0282. Run a formula with `--ready` while another story is in progress. Assert that the instance's touches name its own story file. Assert that `Holds.Of` (`flai/internal/workitem/hold.go`) gives the instance no `no-touches` hold. Extend `TestFormulaRun` to assert the touch.

Waits on no task.

## Done when

- An instance made by `flai formula run` carries its own story file as its touch, with or without `--ready`.
- The new test fails without the change and passes with it.
- `flai test flai/cmd/formula.go flai/cmd/formula_test.go` passes.

## Notes
