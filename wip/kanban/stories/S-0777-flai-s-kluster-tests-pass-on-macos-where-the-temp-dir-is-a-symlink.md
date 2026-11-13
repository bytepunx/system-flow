---
id: S-0777
type: story
nature: remediation
title: flai's kluster tests pass on macOS, where the temp dir is a symlink
status: in-progress
owner: alex
created: 2026-11-13T19:58:46Z
updated: 2026-11-13T20:17:55Z
transitions:
  - to: ready
    at: 2026-11-13T20:07:08Z
    by: orchestrator
  - to: in-progress
    at: 2026-11-13T20:07:22Z
    by: agent-S-0777
tags: [flai]
touches: [flai/cmd/kluster_test.go, flai/internal/kluster/deploy_test.go, flai/internal/kluster/kluster_test.go]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 637
  estimated: true
  priced_by: estimate
  turns:
    - day: 2026-11-13
      ceremony: 1
      work: 10
  models:
    - model: claude-opus-5-5
      input: 56
      output: 282
      cache_read: 1866118
      cache_write: 267119
      cost: 0.8005
  strategic:
    - kind: planner
      seconds: 305
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 52
          output: 25775
          cache_read: 4056908
          cache_write: 201392
          cost: 2.9382
    - kind: orchestrator
      seconds: 250
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 76
          output: 584
          cache_read: 2485729
          cache_write: 350835
          cost: 0.8317
cost_of_delay:
  inputs:
    time_lost_per_cycle: 30m
    by: alex
    at: 2026-11-13T20:02:24Z
  value: 75
  by: planner-S-0777
  at: 2026-11-13T20:04:52Z
forecast:
  duration: 19m
  delivery: 2026-11-14T02:53:00Z
  basis: "flai forecast: median 137 s per unit of size over 65 done remediation stories on claude-opus-5-5, times size 8, 21st in the pull order; kept, as the change is two test fixtures and a verification run."
  by: planner-S-0777
  at: 2026-11-13T20:04:52Z
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
- T-3325 No other flai test compares a fake script's pwd with an unresolved temp path, checked under a symlinked TMPDIR and on Linux

## Notes

- Where `pwd` is logged: `flai/internal/kluster/deploy_test.go` (~41, ~790), `flai/internal/kluster/kluster_test.go` (~22), `flai/cmd/kluster_test.go` (~24, ~252-272).
- The smallest fix is to resolve the temp root once where each fixture makes it, with `filepath.EvalSymlinks(t.TempDir())`, as `flai/internal/agentimage/command_test.go` and `flai/internal/serve/plan_test.go` do already. Expected lines mix paths flai passes as arguments (unresolved) with paths `pwd` prints (resolved), so resolving the root makes both agree.
- Not confirmed: why `pwd` resolves for these steps and not for others. Go sets `PWD` for a child only when `cmd.Env` is nil, and the kluster steps pass an environment of their own, which would explain it. The fix does not depend on it.
- The failing job's log: `gh api repos/bytepunx/system-flow/actions/jobs/132126294410/logs`.
- No change to flai's behaviour, so no design or user documentation changes.

### Planning

Touches, all files, no folder:

| Touch | From | Why |
|-------|------|-----|
| `flai/cmd/kluster_test.go` | declared | `klusterProject` and the fakes that log `$(pwd)`; T-3324 |
| `flai/internal/kluster/kluster_test.go` | declared | `setup` makes the root the fakes run in; T-3323 |
| `flai/internal/kluster/deploy_test.go` | declared | `deploySetup` makes the worktree; T-3323 |
| `flai/internal/kluster/agentconfig_test.go` | declared | holds `TestAPassingDeployWritesTheNamespacesKubeconfig`, which may pass from the fixture alone; T-3323 |

`flai touches suggest` listed co-changed production files: `flai/internal/kluster/deploy.go`, `flai/cmd/kluster.go`, `flai/internal/serve/kluster.go`, and the docs. None is added: the story changes no behaviour of flai, so they are not touched. The planner's survey of every test under `flai/` that logs `pwd` found no other file to fix. The rest resolve their temp roots already, compare only a base name (`flai/internal/serve/mergequeue_test.go`), or rely on `PWD` that flai sets (`TestKlusterExecRunsInTheWorktreeWithTheNamespacesVariables`). T-3325 confirms it with the suite under a symlinked `TMPDIR` and adds a file if one turns up.

Forecast: 19m, delivered 2026-11-14T02:53:00Z, as `flai forecast` gave it: median 137 s per unit of size over 65 done remediation stories on claude-opus-5-5, times size 8 (4 criteria, 4 touches), 21st in the pull order. Kept: the change is two fixtures in test files and a verification run, the size the median covers.

Cost of delay: 75 USD a week, as `flai cod` gave it from the operator's input, `time_lost_per_cycle` 30m (TH-1175): 30m lost per 168h cycle at 150 USD an hour, one cycle a week. Kept: the story names no revenue or penalty to add.

After: none. The story reads nothing another open story writes.

Tasks, in two layers:

1. T-3323 (kluster package fixtures) and T-3324 (`flai/cmd` fixtures), together: no `after` between them and no path in common.
2. T-3325, after both: the whole suite under a symlinked `TMPDIR`, `scripts/flai-test.sh`, and the narrative's account of how the macOS criterion was verified.
