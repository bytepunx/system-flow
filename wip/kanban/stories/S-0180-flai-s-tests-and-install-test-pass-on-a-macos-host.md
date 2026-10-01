---
id: S-0180
type: story
nature: remediation
title: flai's tests and install test pass on a macOS host
status: ready
owner: alex
created: 2026-10-01T08:00:32Z
updated: 2026-10-01T08:12:26Z
transitions:
  - to: ready
    at: 2026-10-01T08:12:26Z
    by: alex
tags: [flai]
touches: [flai/internal/serve/checks_test.go, flai/internal/serve/stop_test.go, flai/cmd/accept_resume_test.go, flai/internal/preview, scripts/install-test.sh]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
---
# S-0180 flai's tests and install test pass on a macOS host

## Goal

Three flai tests and the install test fail on a macOS host, on `main` too (I-0045, four occurrences, last 2026-10-01), so every story an agent works on a Mac either ignores red tests or chases them:

- `flai/internal/serve/checks_test.go` (~68-72) compares `pwd` output, the logical `/var/...` path because exec sets `PWD`, against `filepath.EvalSymlinks`, which gives `/private/var/...`.
- `scripts/install-test.sh` (~47-50) greps for a `mktemp -d` path without resolving it, with the same mismatch.
- `flai/cmd/accept_resume_test.go` (~90) expects acceptance to be refused for want of a git identity, but on macOS git works one out from the hostname, so the check in `flai/internal/preview/accept.go` (~70) finds one and the acceptance goes ahead.
- `TestTheOperatorStopsAStorysAgent` (`flai/internal/serve/stop_test.go` ~25) fails for a cause not yet known: process-group or exit-143 handling on darwin are suspects.

## Acceptance criteria
- [ ] Path comparisons in tests and scripts resolve both sides (`pwd -P`, `EvalSymlinks` on both, or `realpath` in shell)
- [ ] The identity test controls git's identity itself (an isolated `GIT_CONFIG_GLOBAL` and no hostname fallback, for example `user.useConfigOnly`) instead of relying on the host having none; a shared git-isolation test helper is used by the tests that need one
- [ ] The stop test's darwin failure is diagnosed and fixed, with the cause recorded in the story's narrative
- [ ] `scripts/flai-test.sh` and `scripts/install-test.sh` pass on a macOS host and in CI
- [ ] I-0045 is closed with what fixed it

## Tasks

## Notes

A macOS runner in CI (the flai workflow) would keep this from coming back; add one if it is cheap.
