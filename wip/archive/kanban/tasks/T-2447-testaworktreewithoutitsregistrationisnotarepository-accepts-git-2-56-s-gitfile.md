---
id: T-2447
type: task
nature: feature
title: TestAWorktreeWithoutItsRegistrationIsNotARepository accepts git 2.56's gitfile message as well as the older not-a-git-repository one
status: done
parent: S-0572
owner: alex
created: 2026-11-02T23:59:56Z
updated: 2026-11-03T00:05:22Z
transitions:
  - to: ready
    at: 2026-11-03T00:00:33Z
    by: system-flow
  - to: in-progress
    at: 2026-11-03T00:00:34Z
    by: system-flow
  - to: done
    at: 2026-11-03T00:05:22Z
    by: system-flow
stream: S-0572
tags: []
touches: [flai/internal/storygit/mounts_test.go]
---
# T-2447 TestAWorktreeWithoutItsRegistrationIsNotARepository accepts git 2.56's gitfile message as well as the older not-a-git-repository one

## Work

In `flai/internal/storygit/mounts_test.go`, the test asserts that `git status` in a worktree whose registration was moved away fails with output containing `not a git repository`. git 2.56.0, on the ubuntu runner, says `gitfile does not point to a valid repository` instead (`setup.c` changed between v2.47.3 and v2.56.0). Make the assertion accept either wording, keeping the check that git refuses the worktree. Waits for nothing: it changes a file no other task touches.

## Done when

- `go test -race -run TestAWorktreeWithoutItsRegistrationIsNotARepository ./internal/storygit/` passes on this host with git 2.47.3
- the assertion names both messages, so it passes with git 2.56's wording too

## Notes
