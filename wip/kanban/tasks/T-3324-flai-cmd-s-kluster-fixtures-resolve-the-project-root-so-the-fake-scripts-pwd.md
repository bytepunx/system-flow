---
id: T-3324
type: task
nature: remediation
title: flai/cmd's kluster fixtures resolve the project root, so the fake scripts' pwd lines match on macOS
status: done
parent: S-0777
owner: alex
created: 2026-11-13T20:02:30Z
updated: 2026-11-13T20:18:00Z
transitions:
  - to: ready
    at: 2026-11-13T20:08:28Z
    by: agent-S-0777
  - to: in-progress
    at: 2026-11-13T20:08:28Z
    by: agent-S-0777
  - to: done
    at: 2026-11-13T20:18:00Z
    by: agent-S-0777
stream: S-0777
tags: [flai]
touches: [flai/cmd/kluster_test.go]
usage:
  source: log
  seconds: 567
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 24
      output: 140
      cache_read: 713476
      cache_write: 42918
      cost: 0.2839
---
# T-3324 flai/cmd's kluster fixtures resolve the project root, so the fake scripts' pwd lines match on macOS

## Work

- In `flai/cmd/kluster_test.go`, `klusterProject` takes `root` from `tempProject(t)`; resolve it with `filepath.EvalSymlinks` before it is returned, so `deployProject` and every test built on them get a resolved root. Leave `tempProject` alone: it is shared by every test in `flai/cmd`.
- The fakes `fakeGit`, `klusterBuildScript`, `klusterDeployScript`, and `klusterTeardownScript` log `$(pwd)`; check that each line the three tests expect is built from the resolved root or the worktree under it.
- `TestKlusterExecRunsInTheWorktreeWithTheNamespacesVariables` compares `pwd` with the worktree too. It passes on macOS today, because `flai kluster exec` sets `PWD`. Keep it passing; `execProject` needs no change unless the resolved root breaks it.
- Reproduce macOS on Linux as the kluster package's task does: `TMPDIR=<symlink to a real folder> go test ./cmd/ -run 'Kluster'` from `flai/`.

Waits for nothing: it changes only `flai/cmd/kluster_test.go`, which the kluster package's task does not touch.

## Done when

- `TestKlusterDeploy`, `TestKlusterTeardown`, and `TestKlusterPromoteAndDeployments` pass with `TMPDIR` a symlink.
- `go test ./cmd/ -run Kluster` passes with the default `TMPDIR`.
- No file but `flai/cmd/kluster_test.go` changes.

## Notes
