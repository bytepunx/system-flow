---
id: I-0255
title: go build's VCS stamping fails with git exit status 128 inside the agent container, though git answers each command there by hand
class: defect
status: open
count: 1
cost: 15m
first_reported: 2026-11-02T19:02:40Z
last_reported: 2026-11-02T19:02:40Z
updated: 2026-11-10T08:37:12Z
---

# I-0255 go build's VCS stamping fails with git exit status 128 inside the agent container, though git answers each command there by hand

## Description
go build's VCS stamping fails with git exit status 128 inside the agent container, though git answers each command there by hand

## Instances

### 2026-11-02T19:02:40Z
Story: S-0510.
S-0510, running `scripts/flai.sh test` inside flai-agent-sf-claude-code:local on its worktree. The mounts were the worktree read-write, the main checkout's .git read-only, and the story's registration read-write. Two tests failed: TestSubprojectEditAddsSetsAndRemovesCommittingNothing and TestSubprojectEditWorksInTheStorysWorktree (flai/internal/mcpserver/subproject_test.go:136 and :232). Both failed on `go build: error obtaining VCS status: exit status 128`. In the same container, as the same uid, `git status --porcelain`, `git -c log.showsignature=false log -1 --format=%H:%ct`, and `git rev-parse --show-toplevel` all exit 0, and GIT_TRACE shows no git call from go. Worked around in the layer with GOFLAGS=-buildvcs=false, since a test's build needs no stamp. The cause is not found.

## Remediation

Story S-0686 remediates this issue, created from it at 2026-11-10T08:37:12Z.
