---
id: T-3406
type: task
nature: remediation
title: The preflight checks main again with the manifest's flai.command when the running flai is behind, and records the cause on a hold it cannot clear that way
status: done
parent: S-0797
owner: alex
created: 2026-11-16T01:44:46Z
updated: 2026-11-16T03:33:56Z
transitions:
  - to: ready
    at: 2026-11-16T03:14:36Z
    by: agent-S-0797
  - to: in-progress
    at: 2026-11-16T03:14:36Z
    by: agent-S-0797
  - to: done
    at: 2026-11-16T03:33:56Z
    by: agent-S-0797
stream: S-0797
tags: [flai, merge-queue]
touches: [flai/internal/mergequeue/preflight.go, flai/internal/mergequeue/preflight_test.go, flai/internal/mergequeue/state.go, flai/internal/mergequeue/state_test.go, flai/internal/mergequeue/batch_test.go, flai/internal/mergequeue/land_test.go, flai/internal/check/output.go, flai/internal/check/output_test.go, flai/internal/taskdone/taskdone.go, flai/internal/taskdone/taskdone_test.go]
after: [T-3404]
usage:
  source: log
  seconds: 1160
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 54
      output: 18165
      cache_read: 4152586
      cache_write: 152818
      cost: 2.2219
---
# T-3406 The preflight checks main again with the manifest's flai.command when the running flai is behind, and records the cause on a hold it cannot clear that way

## Work

Change `preflightCheck` in `flai/internal/mergequeue/preflight.go` as ADR-0309's rules 2 and 3 say (T-3403 writes the ADR).

- When `check.Run` gives error findings, ask T-3404's function whether the running flai is behind main's flai sub-project.
- When it is behind and the manifest names `flai.command`, run `<command> check --json` in a scratch worktree at main's head. Reuse `confirmLint`'s checkout and removal under `PreflightWorktreeName`, carrying the main checkout's uncommitted files as that does, so no file of the main checkout changes. Keep the error findings it reports, by `checkFindings`' rule, in place of the running flai's. Read its JSON as `taskdone`'s `checkResult` does. Move that reader into a package both can import only if it stays small; otherwise read the JSON here.
- When it is behind and there is no command, or the command exits without a result, keep the findings. Add a `Cause` to `PreflightHold` in `flai/internal/mergequeue/state.go`, `flai_behind`, with the running version, its commit, and main's commit count beyond it, and set it. `preflightLook` must keep the cause across looks and drop it when the hold clears.
- When it is not behind, or cannot tell, hold as today.

It waits for T-3404, whose function it calls.

## Done when

- A test in `flai/internal/mergequeue/preflight_test.go` reproduces I-0382. A role file in a form the running check reports `roles.file` on, a running flai stubbed as behind, and a stub `flai.command` that reports clean together give no hold.
- With the stub command giving no result, the hold carries the `flai_behind` cause.
- With the running flai not behind, the hold is as before, with no cause.
- `flai/internal/mergequeue/state_test.go` round-trips the cause through `state.json`.
- The main checkout's files are unchanged after each test.
- `flai test flai/internal/mergequeue` passes.

## Notes
