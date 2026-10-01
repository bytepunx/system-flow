---
id: S-0180
type: story
nature: remediation
title: flai's tests and install test pass on a macOS host
status: review
owner: alex
created: 2026-10-01T08:00:32Z
updated: 2026-10-01T08:30:12Z
transitions:
  - to: ready
    at: 2026-10-01T08:12:26Z
    by: alex
  - to: in-progress
    at: 2026-10-01T08:19:04Z
    by: agent-S-0180
  - to: review
    at: 2026-10-01T08:30:12Z
    by: agent-S-0180
tags: [flai]
touches: [flai/internal/serve/checks_test.go, flai/internal/serve/stop_test.go, flai/cmd/accept_resume_test.go, flai/cmd/doc_test.go, flai/internal/gittest, flai/internal/preview, scripts/install-test.sh, ".github/workflows/flai.yml", design/tech/ci.md, design/issues]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 690
  models:
    - model: claude-opus-5-5
      input: 166
      output: 39007
      cache_read: 8272356
      cache_write: 140238
      cost: 3.5572
---
# S-0180 flai's tests and install test pass on a macOS host

## Goal

Three flai tests and the install test fail on a macOS host, on `main` too (I-0045, four occurrences, last 2026-10-01), so every story an agent works on a Mac either ignores red tests or chases them:

- `flai/internal/serve/checks_test.go` (~68-72) compares `pwd` output, the logical `/var/...` path because exec sets `PWD`, against `filepath.EvalSymlinks`, which gives `/private/var/...`.
- `scripts/install-test.sh` (~47-50) greps for a `mktemp -d` path without resolving it, with the same mismatch.
- `flai/cmd/accept_resume_test.go` (~90) expects acceptance to be refused for want of a git identity, but on macOS git works one out from the hostname, so the check in `flai/internal/preview/accept.go` (~70) finds one and the acceptance goes ahead.
- `TestTheOperatorStopsAStorysAgent` (`flai/internal/serve/stop_test.go` ~25) fails for a cause not yet known: process-group or exit-143 handling on darwin are suspects.

## Acceptance criteria
- [x] Path comparisons in tests and scripts resolve both sides (`pwd -P`, `EvalSymlinks` on both, or `realpath` in shell)
- [x] The identity test controls git's identity itself (an isolated `GIT_CONFIG_GLOBAL` and no hostname fallback, for example `user.useConfigOnly`) instead of relying on the host having none; a shared git-isolation test helper is used by the tests that need one
- [x] The stop test's darwin failure is diagnosed and fixed, with the cause recorded in the story's narrative
- [x] `scripts/flai-test.sh` and `scripts/install-test.sh` pass on a macOS host and in CI
- [x] I-0045 is closed with what fixed it

## Tasks
- T-0628 Path comparisons in tests and scripts resolve both sides
- T-0629 A shared git-isolation helper; the identity test controls git's identity
- T-0630 Diagnose and fix the stop test's darwin failure
- T-0631 Both suites pass on macOS and in CI, a macOS runner in CI, I-0045 closed

## Notes

A macOS runner in CI (the flai workflow) would keep this from coming back; add one if it is cheap.

The repository is public, so a standard macOS runner costs nothing: the flai workflow's test job runs on `macos-latest` too (276de17).

Verification so far (agent-S-0180, 2026-10-01): on this Mac (darwin/arm64), `scripts/flai-test.sh` passes gofmt, vet, golangci-lint, behavior, integration, and the template smoke. Its repository check stops on six warnings about the main checkout's board (four done epics not archived, TH-0026 and TH-0032 on archived stories), not on anything this story touches. The markdown lint and `scripts/install-test.sh` pass when run on their own. On Linux, `go test -race -count=1 ./...` passes in `golang:1.26` in Docker. The workflow runs only on a pull request or a push to main, so nothing from this branch has run in GitHub Actions. The designer chose (TH-0043, option b) to count this Mac and the Docker Linux run for "in CI", and to let CI confirm on the push to main after acceptance.
