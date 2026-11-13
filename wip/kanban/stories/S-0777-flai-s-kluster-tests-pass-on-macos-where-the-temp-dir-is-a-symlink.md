---
id: S-0777
type: story
nature: remediation
title: flai's kluster tests pass on macOS, where the temp dir is a symlink
status: backlog
owner: alex
created: 2026-11-13T19:58:46Z
updated: 2026-11-13T20:02:35Z
transitions: []
tags: [flai]
touches: [flai/cmd/kluster_test.go, flai/internal/kluster/deploy_test.go, flai/internal/kluster/kluster_test.go, flai/internal/kluster/agentconfig_test.go]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: sum
  seconds: 0
  models: []
  strategic:
    - kind: orchestrator
      seconds: 55
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 14
          output: 78
          cache_read: 436814
          cache_write: 67308
          cost: 0.1477
cost_of_delay:
  inputs:
    time_lost_per_cycle: 30m
    by: alex
    at: 2026-11-13T20:02:24Z
finalized:
  by: alex
  at: 2026-11-13T20:02:35Z
---
# S-0777 flai's kluster tests pass on macOS, where the temp dir is a symlink

## Goal

The `test (macos-latest)` job of the `flai` workflow fails its "Behavior tests" step on every push to `main`: run 43495183748 (commit `9b63bacf3`, 2026-11-13) and the run before it (flai 1.94.1) fail the same eight tests, and the job was already red at flai 1.90.0. While it is red, the macOS job tells nobody anything, and a real macOS defect would go unseen.

The eight tests are:

- `flai/cmd`: `TestKlusterDeploy`, `TestKlusterTeardown`, `TestKlusterPromoteAndDeployments`
- `flai/internal/kluster`: `TestAPassingDeployWritesTheNamespacesKubeconfig`, `TestDeployBuildsLoadsAndDeploysIntoTheStorysNamespace`, `TestADeployWithoutABuildLoadsNothing`, `TestATeardownRunsTheTeardownCommandThenDeletesTheNamespace`, `TestPromoteDeploysMainsHeadIntoTheDefaultNamespace`

Each fails on a path's spelling alone. On macOS the temp dir `/var/folders/...` is a symlink to `/private/var/folders/...`. The fake `build`, `deploy`, `teardown`, and `git` scripts log their working directory with `$(pwd)`, which prints the resolved path there, and the tests build the lines they expect from `t.TempDir()`, which is not resolved:

```text
 got "build /private/var/folders/.../005 --all"
want "build /var/folders/.../005 --all"
```

The tests are wrong, not flai: nothing in the log shows flai misbehaving on macOS. This is the defect S-0180 fixed elsewhere (I-0045, closed), come back in tests written since.

## Acceptance criteria
- [ ] The eight tests above pass on macOS: both sides of each path comparison are resolved, or neither is
- [ ] No other test under `flai/` compares a `$(pwd)` a fake script logged with an unresolved temp path; those found are fixed the same way
- [ ] The `test (macos-latest)` job of the `flai` workflow passes on the story's branch, or, where the branch cannot run it, the narrative says how the fix was verified without a macOS host
- [ ] The tests still pass on Linux, with `scripts/flai-test.sh`

## Tasks
- T-3323 The kluster package's fixtures resolve their temp folders, so pwd and the expected lines agree on macOS
- T-3324 flai/cmd's kluster fixtures resolve the project root, so the fake scripts' pwd lines match on macOS

## Notes

- Where `pwd` is logged: `flai/internal/kluster/deploy_test.go` (~41, ~790), `flai/internal/kluster/kluster_test.go` (~22), `flai/cmd/kluster_test.go` (~24, ~252-272).
- The smallest fix is to resolve the temp root once where each fixture makes it, with `filepath.EvalSymlinks(t.TempDir())`, as `flai/internal/agentimage/command_test.go` and `flai/internal/serve/plan_test.go` do already. Expected lines mix paths flai passes as arguments (unresolved) with paths `pwd` prints (resolved), so resolving the root makes both agree.
- Not confirmed: why `pwd` resolves for these steps and not for others. Go sets `PWD` for a child only when `cmd.Env` is nil, and the kluster steps pass an environment of their own, which would explain it. The fix does not depend on it.
- The failing job's log: `gh api repos/bytepunx/system-flow/actions/jobs/132126294410/logs`.
- No change to flai's behaviour, so no design or user documentation changes.
