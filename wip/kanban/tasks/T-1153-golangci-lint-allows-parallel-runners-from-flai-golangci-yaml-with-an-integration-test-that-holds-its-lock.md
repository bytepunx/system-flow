---
id: T-1153
type: task
nature: remediation
title: golangci-lint allows parallel runners from flai/.golangci.yaml, with an integration test that holds its lock
status: backlog
parent: S-0308
owner: alex
created: 2026-10-07T02:19:48Z
updated: 2026-10-07T02:19:48Z
transitions: []
stream: S-0308
tags: [flai, lint, test]
touches: [flai/.golangci.yaml, flai/tests/integration/golangci_lock_test.go]
---
# T-1153 golangci-lint allows parallel runners from flai/.golangci.yaml, with an integration test that holds its lock

## Work

- Write the test first, in `flai/tests/integration/golangci_lock_test.go`, for unix hosts only. It skips under `-short` and when no golangci-lint v2 is on `PATH`, as `scripts/flai-test.sh` finds it.
- The test points `TMPDIR` at a folder of its own, so it never takes the host's real lock. It holds an exclusive `flock` on `golangci-lint.lock` there, as a run on another worktree would. Then it runs `golangci-lint run --config <repo>/flai/.golangci.yaml` on a tiny clean package it writes.
- Watch it fail with exit 3 ("parallel golangci-lint is running") before the fix.
- Add `run: allow-parallel-runners: true` to `flai/.golangci.yaml`, with a comment naming I-0101. The test then passes, and `make flai-test` lints as before.
- Waits for nothing: this is the first layer.

## Done when

- The test fails without the config key and passes with it.
- `flai test flai/.golangci.yaml flai/tests/integration/golangci_lock_test.go` passes.

## Notes

golangci-lint v2 locks `golangci-lint.lock` in `os.TempDir()` unless `run.allow-parallel-runners` is set, and exits 3 after 5 seconds without it (`pkg/commands/run.go`, `acquireFileLock`).
