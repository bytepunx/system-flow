---
id: T-3325
type: task
nature: remediation
title: No other flai test compares a fake script's pwd with an unresolved temp path, checked under a symlinked TMPDIR and on Linux
status: in-progress
parent: S-0777
owner: alex
created: 2026-11-13T20:02:41Z
updated: 2026-11-13T20:24:00Z
transitions:
  - to: ready
    at: 2026-11-13T20:24:00Z
    by: agent-S-0777
  - to: in-progress
    at: 2026-11-13T20:24:00Z
    by: agent-S-0777
stream: S-0777
tags: [flai]
after: [T-3323, T-3324]
---
# T-3325 No other flai test compares a fake script's pwd with an unresolved temp path, checked under a symlinked TMPDIR and on Linux

## Work

- Run the whole Go suite from `flai/` with `TMPDIR` a symlink to a real folder, `TMPDIR=<symlink> go test -short ./...`, which reproduces macOS's `/var` to `/private/var` link on Linux. Any test that fails there on a path's spelling alone is fixed the same way, by resolving its fixture's temp root with `filepath.EvalSymlinks`, and its file is added to this task's touches with `flai touches`.
- The planner's survey found no other case: the other tests that log `pwd` resolve their temp roots already (`flai/internal/agentimage`, `flai/internal/serve`, `flai/internal/mcpserver`, `flai/internal/taskdone`), `TestRunCommandAnswersStandardErrorAndTheExit` in `flai/internal/serve/mergequeue_test.go` compares only the base name, and `TestKlusterExecRunsInTheWorktreeWithTheNamespacesVariables` passes because `flai kluster exec` sets `PWD`. The run above confirms or refutes it.
- Run `scripts/flai-test.sh` on Linux.
- The story's branch is not pushed until it is accepted and published, so the `test (macos-latest)` job cannot run on it. Say in the narrative's `## Decisions` how the fix was verified without a macOS host: the symlinked `TMPDIR` runs, before and after, and their results.

Waits for T-3323 and T-3324: the suite under a symlinked `TMPDIR` judges what is left only once both fixtures are fixed.

## Done when

- `go test -short ./...` from `flai/` passes with `TMPDIR` a symlink.
- `scripts/flai-test.sh` passes on Linux.
- The narrative's `## Decisions` names how the macOS criterion was verified and what each run gave.

## Notes
