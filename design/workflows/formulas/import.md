---
title: "Check the import of {{path}} and run its tests"
updated: 2026-11-02
status: active
name: import
description: "Check the repository flai import brought under system-flow and run its tests, now that the import is committed."
nature: improvement
inputs:
  - name: path
    description: The root of the imported repository, as flai import --formula gives it.
steps:
  - id: check
    title: flai check --strict reports nothing on the import of {{path}}
    run: {flai: check --strict}
  - id: test
    title: The tests of the import of {{path}} pass
    run: {flai: import test}
    after: [check]
---

# Import

`flai import --formula` starts this formula in the project it has just imported, as `flai formula run import --input path=<root> --ready --autocommit`, and the board's Import runs `flai import --formula` ([ADR-0035](../../adrs/0035-repositories-under-folders-the-operator-names-can-be-imported-from-the-board.md)). The import is committed before its tests run (TH-0588), so the board shows the tests' progress as this story's, and both steps run from the story's worktree.

- Record in the story's Notes what step 1 reported, and each test step 2 ran with how it went: `flai import test` names the tests it found and where it found them.
- A failing step blocks its task with the output and asks the operator on a thread: a finding of `flai check --strict`, or a failed test, on which `flai import test` exits 5. Fixing the repository is a story of its own.
- This story's branch changes nothing. Its acceptance adds no bump, and the next publish releases nothing for it.

<!-- system-flow:end-of-baseline -->

## Project additions
