---
id: I-0045
title: "Three flai tests fail on a macOS host, on main too: the temp dir's /private/var path and a git identity the host has"
class: defect
status: closed
count: 4
cost: 4m
first_reported: 2026-09-26T03:11:49Z
last_reported: 2026-10-01T07:46:11Z
updated: 2026-10-01T08:26:05Z
---

# I-0045 Three flai tests fail on a macOS host, on main too: the temp dir's /private/var path and a git identity the host has

## Description
Three flai tests fail on a macOS host, on main too: the temp dir's /private/var path and a git identity the host has

## Instances

### 2026-09-26T03:11:49Z
S-0118, 2026-09-26: make flai-test on this macOS host fails TestRunChecksSubstitutesStoryAndRootAndRunsInTheWorktree (got /var/folders/..., want /private/var/folders/...), TestAcceptRefusesBeforeChangingAnythingWithoutIdentity (the acceptance went ahead instead of refusing for a missing identity), and install-test.sh's self-upgrade path check (mktemp -d gives /var, flai reports /private/var). All three fail the same way on a clean checkout of main, so the story's own change is not the cause. The remaining tiers were run one by one to cover the story's change.

### 2026-09-26T03:21:40Z
S-0119, 2026-09-26: go test -short ./... fails the same two Go tests on this macOS host (accept identity, checks /private/var); the story's change is in flai/cmd/dashboard.go and does not touch them.

### 2026-09-26T03:33:18Z
S-0117, 2026-09-26: make test fails TestRunChecksSubstitutesStoryAndRootAndRunsInTheWorktree (/var vs /private/var) and make smoke stops at install-test's self-upgrade path check, the same way. The story changes flaiover/src and design/ only; flaiover's tests, the template render, flai check --strict, and the markdown lint were run and pass.

### 2026-10-01T07:46:11Z
S-0173, 2026-10-01: make flai-test in the story worktree failed the same three tests (internal/serve TestRunChecksSubstitutesStoryAndRootAndRunsInTheWorktree, TestTheOperatorStopsAStorysAgent; cmd TestAcceptRefusesBeforeChangingAnythingWithoutIdentity), and each fails on main unchanged. The worktree's bin/ also had no golangci-lint v2, so ~/go/bin's v1 refused the config until main's bin/golangci-lint was copied in.

## Remediation
Closed 2026-10-01T08:26:05Z: S-0180: checks_test reads pwd -P and install-test.sh resolves its mktemp path, so both sides of each path comparison are resolved; the identity test isolates git with internal/gittest and sets user.useConfigOnly, so git no longer works out an identity from the hostname; the stop test waits until the stub has set its TERM trap, which macOS's slower-starting bash had not when the signal came; flai.yml now runs the tests on macos-latest too.
