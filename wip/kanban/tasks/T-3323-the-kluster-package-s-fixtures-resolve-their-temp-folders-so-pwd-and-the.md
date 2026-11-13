---
id: T-3323
type: task
nature: remediation
title: The kluster package's fixtures resolve their temp folders, so pwd and the expected lines agree on macOS
status: done
parent: S-0777
owner: alex
created: 2026-11-13T20:02:26Z
updated: 2026-11-13T20:17:55Z
transitions:
  - to: ready
    at: 2026-11-13T20:08:27Z
    by: agent-S-0777
  - to: in-progress
    at: 2026-11-13T20:08:27Z
    by: agent-S-0777
  - to: done
    at: 2026-11-13T20:17:55Z
    by: agent-S-0777
stream: S-0777
tags: [flai]
touches: [flai/internal/kluster/kluster_test.go, flai/internal/kluster/deploy_test.go]
usage:
  source: log
  seconds: 568
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 16
      output: 81
      cache_read: 568158
      cache_write: 40638
      cost: 0.2285
---
# T-3323 The kluster package's fixtures resolve their temp folders, so pwd and the expected lines agree on macOS

## Work

- In `flai/internal/kluster/kluster_test.go`, `setup` makes `root` with `t.TempDir()`; resolve it once with `filepath.EvalSymlinks`, as `flai/internal/agentimage/command_test.go` does, failing the test on an error.
- In `flai/internal/kluster/deploy_test.go`, `deploySetup` makes `worktree` and `env` with `t.TempDir()`; resolve them the same way. Check each other `root := t.TempDir()` in the file (about lines 636, 662, 690, 707, 1125, 1262) and resolve only those whose path a fake script's `$(pwd)` line is compared with.
- `flai/internal/kluster/agentconfig_test.go` holds `TestAPassingDeployWritesTheNamespacesKubeconfig`, which compares `"deploy " + f.worktree`. It should pass from the fixture alone; change it only if it builds a path of its own.
- A small helper in `kluster_test.go`, such as `resolvedTempDir(t)`, keeps the two fixtures alike.
- Reproduce macOS on Linux: run the package with `TMPDIR` set to a symlink to a real folder, for example `ln -s "$(mktemp -d)" /tmp/flai-symlinked-tmp`, then `TMPDIR=/tmp/flai-symlinked-tmp go test ./internal/kluster/` from `flai/`. The five tests below fail before the change and pass after it.

Waits for nothing: it changes only the kluster package's tests, apart from the `flai/cmd` task.

## Done when

- `TestAPassingDeployWritesTheNamespacesKubeconfig`, `TestDeployBuildsLoadsAndDeploysIntoTheStorysNamespace`, `TestADeployWithoutABuildLoadsNothing`, `TestATeardownRunsTheTeardownCommandThenDeletesTheNamespace`, and `TestPromoteDeploysMainsHeadIntoTheDefaultNamespace` pass with `TMPDIR` a symlink.
- `go test ./internal/kluster/` passes with the default `TMPDIR`.
- No file outside `flai/internal/kluster/*_test.go` changes.

## Notes
