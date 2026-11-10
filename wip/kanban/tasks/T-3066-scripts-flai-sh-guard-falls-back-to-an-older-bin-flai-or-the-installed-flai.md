---
id: T-3066
type: task
nature: improvement
title: scripts/flai.sh guard falls back to an older bin/flai or the installed flai when go build fails, and exits 127 with neither, with a test that makes the build fail
status: backlog
parent: S-0666
owner: alex
created: 2026-11-10T10:35:46Z
updated: 2026-11-10T10:35:46Z
transitions: []
stream: S-0666
tags: [guard, scripts]
touches: [scripts/flai.sh, flai/cmd/guard_test.go]
after: [T-3061]
---
# T-3066 scripts/flai.sh guard falls back to an older bin/flai or the installed flai when go build fails, and exits 127 with neither, with a test that makes the build fail

## Work

In `scripts/flai.sh`, when `go` is present and `scripts/flai-build.sh` fails, and the first argument is `guard`, do not exit with the build's status. Run `bin/flai` when it is there, older than the sources, or else the installed `flai` from `PATH`, each with a one-line warning on standard error naming the failed build, and exit 127 when neither is there, so the hooks in `.claude/settings.json` refuse the call. Every other command keeps exiting with the build's status. Keep the script POSIX `sh` under `set -eu`, and keep the hooks' text as it is: `checkGuardHooks` in `flai/cmd/guard_test.go` pins it.

Reproduce I-0275 in `TestTheGuardHooksRefuseACallWhenTheGuardCannotStart`, in `flai/cmd/guard_test.go`, which runs both hook commands with a copied `scripts/flai.sh`, `env.sh`, and `flai-build.sh` and no `go` on `PATH`. Add cases with a fake `go` that exits 1 and a stale `bin/flai`:

- `bin/flai` there: the hook gives the verdict of the `bin/flai` guard, 0 or 2.
- no `bin/flai`, a fake installed `flai`: the hook gives its verdict.
- neither: the hook refuses with `flai guard did not run (exit 127)`.

Check that the first of these fails on the current `scripts/flai.sh`, since it lets the call through. Write every file the test runs before `t.Parallel`, as the test's comment asks.

It waits for T-3061, whose ADR records the decision this code follows.

## Done when

- With a failing `go build`, `scripts/flai.sh guard` runs the older `bin/flai` or the installed flai, or exits 127 when neither is there, and the hook refuses in that last case.
- Any other command still exits non-zero on a failed build.
- The new cases in `TestTheGuardHooksRefuseACallWhenTheGuardCannotStart` pass, and the first failed before the change.
- `flai test scripts/flai.sh flai/cmd/guard_test.go` passes.

## Notes
