---
id: T-3440
type: task
nature: feature
title: The user guide says how a project's CI runs its suites, and the CLI and workflow design describe flai test --at
status: done
parent: S-0812
owner: alex
created: 2026-11-16T23:47:40Z
updated: 2026-11-18T09:04:10Z
transitions:
  - to: ready
    at: 2026-11-18T08:57:03Z
    by: agent-S-0812
  - to: in-progress
    at: 2026-11-18T08:57:03Z
    by: agent-S-0812
  - to: done
    at: 2026-11-18T09:04:10Z
    by: agent-S-0812
stream: S-0812
tags: [flai]
touches: [docs/users/flai.md, design/system/flai-cli.md, design/system/workflow.md, docs/operators/runbooks/change-the-board.md]
after: [T-3437]
usage:
  source: log
  seconds: 427
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 36
      output: 10221
      cache_read: 2513437
      cache_write: 89636
      cost: 1.2654
---
# T-3440 The user guide says how a project's CI runs its suites, and the CLI and workflow design describe flai test --at

## Work

ADR-0312 §6, criterion 4. In `docs/users/flai.md`, under "Run the tests for what changed", describe `flai test --at <point>`: what each point is, that the suites come from the workflow file's `verify_at:`, that `all` runs every tier but the manual ones, and how `--suite` narrows a point. Add how a project's CI runs its suites: the template's CI file installs flai and calls `flai test --at ci`. A project chooses what CI runs by listing suites under `verify_at.ci` in its workflow file, and one job can run its share with `--suite`.

In `design/system/flai-cli.md`, add `--at` to the `flai test` command. In `design/system/workflow.md`, say that `verify_at.ci` is read by `flai test --at ci`.

It waits for T-3437, because it describes the command as that task builds it.

## Done when

- `docs/users/flai.md` documents `--at` and has a passage on how a project's CI runs its suites.
- `design/system/flai-cli.md` and `design/system/workflow.md` describe `--at` and `verify_at.ci` as built.
- The markdown lint passes on the three files.

## Notes
