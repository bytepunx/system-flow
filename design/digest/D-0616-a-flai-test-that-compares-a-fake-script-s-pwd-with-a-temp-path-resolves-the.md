---
id: D-0616
title: "A flai test that compares a fake script's $(pwd) with a temp path resolves the fixture's temp root with filepath.EvalSymlinks, since macOS's temp dir is a symlink that $(pwd) prints resolved, and a TMPDIR that is a symlink reproduces it on Linux"
kind: decision
sources: [S-0777]
topics: [flai, cli, go, code, testing, ci]
paths: [flai/cmd/kluster_test.go, flai/internal/kluster/deploy_test.go, flai/internal/kluster/kluster_test.go]
date: 2026-11-13
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0616 A flai test that compares a fake script's $(pwd) with a temp path resolves the fixture's temp root with filepath.EvalSymlinks, since macOS's temp dir is a symlink that $(pwd) prints resolved, and a TMPDIR that is a symlink reproduces it on Linux

S-0777 fixed eight kluster tests that failed on the flai workflow's macOS job. Each compared a fake script's `$(pwd)` with a path under `t.TempDir()`. On macOS the temp dir is under `/var`, a symlink to `/private/var`: `t.TempDir` names the link and `$(pwd)` prints the resolved folder. The comparison failed on the path's spelling alone. D-0077 recorded the same cause for three earlier tests (I-0045).

The decision, from S-0777's Decisions (`wip/archive/agents/S-0777.md`):

```text
Fix in the tests only: resolve each fixture's temp root with filepath.EvalSymlinks, as agentimage/command_test.go and serve/plan_test.go do; flai's behaviour does not change, so no design or user docs change.
```

How it was verified without a macOS host, since the story's branch is not pushed before acceptance and the `test (macos-latest)` job cannot run on it:

```text
On Linux a TMPDIR that is a symlink to a real folder reproduces macOS's /var to /private/var link ... Before the fix: go test -run on the eight tests with TMPDIR a symlink failed all eight ... After the fix: the eight pass, and the whole short suite from flai/ (go test -short -count=1 -skip '^TestGateway' ./...) passes with 0 failures with TMPDIR a symlink.
```

A later test that compares a path a fake or a child process prints with a temp path resolves the temp root first. A story that fixes a macOS-only path failure reproduces it here with a symlinked TMPDIR, kept outside any git repository (see the entry on where that TMPDIR lives), and words its criteria to that run as D-0475 says.

Sources: S-0777 (wip/archive/agents/S-0777.md, Decisions).
