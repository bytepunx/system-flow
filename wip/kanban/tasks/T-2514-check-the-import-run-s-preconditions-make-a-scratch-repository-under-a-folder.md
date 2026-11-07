---
id: T-2514
type: task
nature: improvement
title: Check the import run's preconditions, make a scratch repository under a folder named for import, and import it from the board's Import
status: done
parent: S-0564
owner: alex
created: 2026-11-07T03:46:58Z
updated: 2026-11-07T08:27:24Z
transitions:
  - to: ready
    at: 2026-11-07T08:03:12Z
    by: system-flow
  - to: in-progress
    at: 2026-11-07T08:03:13Z
    by: system-flow
  - to: done
    at: 2026-11-07T08:27:24Z
    by: agent-S-0564
stream: S-0564
tags: [template]
usage:
  source: log
  seconds: 1233
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 104
      output: 11217
      cache_read: 5482957
      cache_write: 63017
      cost: 1.8255
---
# T-2514 Check the import run's preconditions, make a scratch repository under a folder named for import, and import it from the board's Import

## Work

Check these on the host before anything runs. If one does not hold, ask the operator on a thread on S-0564, and run nothing until it is answered:

- `flai serve import list` names a folder for import outside this repository. Adding one with `flai serve import add` is the operator's.
- The installed `flai` holds S-0424's change: `flai import --help` names `--formula`, and `flai issue list --help` names `--report`. flai 1.83.0 has both.
- TH-0845, the plan's thread, is answered on question 2: who makes the scratch repository, where, and who presses the board's Import.

Make the scratch repository in a new folder below that import folder. Run `git init` and commit a `README.md` and a `Makefile` whose `test` target passes, so that `flai import test` finds a test and runs it. It has no `system-flow.yaml`. If your sandbox refuses the write outside the worktree, ask on the thread for the operator to make it.

Import it from the board's Import, which runs `flai import --formula` through the host API. The import is committed first. Then the `import` formula starts in the new project as `flai formula run import --input path=<root> --ready --autocommit` (ADR-0035, S-0424). Check that the `agent` host action is on for the new project, so `flai serve` starts the instance's agent. Turning it on is the operator's.

Write the scratch repository's path and the instance's ID in this task's Notes with `flai edit`. Run no second import.

This task waits for no other task. It is in the plan's first layer, beside the cadence's start, with which it shares no path.

## Done when

- Each precondition above is checked, and the result is written in this task's Notes.
- The scratch repository is a system-flow project, and its import is committed.
- One instance of the `import` formula exists in the new project, in ready or later, and this task's Notes name its ID and the repository's path.

## Notes
