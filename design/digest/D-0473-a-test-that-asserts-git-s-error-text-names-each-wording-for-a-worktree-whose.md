---
id: D-0473
title: "A test that asserts git's error text names each wording: for a worktree whose registration is gone, git 2.56 says gitfile does not point to a valid repository where 2.47 says not a git repository"
kind: decision
sources: [S-0572]
topics: [testing, ci, git, go, code]
paths: [flai/internal/storygit/mounts_test.go]
date: 2026-11-03
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0473 A test that asserts git's error text names each wording: for a worktree whose registration is gone, git 2.56 says gitfile does not point to a valid repository where 2.47 says not a git repository

`TestAWorktreeWithoutItsRegistrationIsNotARepository` in `flai/internal/storygit/mounts_test.go` (S-0517) expected git to say `not a git repository` for a worktree whose registration is gone. git 2.56.0 on the `flai` workflow's ubuntu runner says `gitfile does not point to a valid repository`, while git 2.47.3 on this host says the old text, so the test passed when it was written and failed in CI and in GoReleaser's `go test ./...` hook. The wording changed in git's `setup.c` between those versions (S-0572's story, Goal and Notes).

S-0572's fix, from its Decisions (`wip/archive/agents/S-0572.md`):

```text
T-2447: the test accepts either of git's wordings rather than checking only the exit status, so it still shows git refuses the worktree as not a repository, which the test is about.
```

The host's git is older than CI's, so a test that matches git's message passes here and can still fail there. Match every wording the supported versions print. Do not fall back to the exit status alone.

Source: S-0572.
