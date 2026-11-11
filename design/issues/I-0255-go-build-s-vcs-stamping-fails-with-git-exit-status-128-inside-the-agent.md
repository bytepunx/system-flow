---
id: I-0255
title: go build's VCS stamping fails with git exit status 128 inside the agent container, though git answers each command there by hand
class: defect
status: closed
count: 1
cost: 15m
first_reported: 2026-11-02T19:02:40Z
last_reported: 2026-11-02T19:02:40Z
updated: 2026-11-11T07:50:20Z
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

### Cause, found by S-0686

Go's VCS stamping does not run git where the module is. Go looks up from the module's folder for a `.git` folder, and a worktree's `.git` is a file, so Go passes over it. A story's worktree is `<main>/.flai-cache/worktrees/S-nnnn`, so Go stops at the main checkout and runs `git status --porcelain` there. Inside the container Docker makes the main checkout's folder as root, as the parent of its mounts, so git refuses it for dubious ownership and exits 128. Git run by hand in the worktree finds the worktree's `.git` file, so it works. Go keeps git's standard error to itself, so `GIT_TRACE` on standard error showed nothing.

S-0686 reproduced it in `flai-agent-sf-claude-code:local` on 2026-11-11 with a `git` shim on `PATH` that logged each call:

```text
cwd=/home/alex/git/bytepunx/system-flow args=status --porcelain
fatal: detected dubious ownership in repository at '/home/alex/git/bytepunx/system-flow'
rc=128
```

| Mounts and environment | `go build -buildvcs=true` in the worktree |
|------------------------|--------------------------------------------|
| S-0510's: worktree, main `.git` read-only, registration, no `safe.directory` | `error obtaining VCS status: exit status 128` |
| The harness's since ADR-0253 (S-0628): those, `wip/`, `system-flow.yaml`, and `GIT_CONFIG_*` giving `safe.directory` for the main checkout | passes, stamped |

So the cause no longer occurs in a story's containers: `flai/internal/harness/container.go`'s `safeMain`, `scripts/agent-container.sh`, and `scripts/agent-image.sh --check` give git `safe.directory` for the main checkout. The fix proposed: a smoke test that builds with VCS stamping in the layer's image under the harness's mounts and fails without `safe.directory`, and the layer drops `GOFLAGS=-buildvcs=false`.

Go stamping the main checkout's revision rather than the worktree's is a fault of its own, on the host too: a `go build` in a story's worktree records `vcs.revision` of the main checkout's `HEAD`.
Closed 2026-11-11T07:50:20Z: S-0686. The 128 was git refusing the root-owned main checkout for dubious ownership: for a story's worktree Go's VCS stamping runs git in the main checkout, as Go passes over the worktree's .git file. The safe.directory for the main checkout that the container harness gives since ADR-0253 (S-0628) already admits it; S-0686 dropped the layer's GOFLAGS=-buildvcs=false workaround, said why in container.go, agent-container.sh, the layer's Dockerfile, docker.md, and the runbook, and added scripts/agent-vcs-test.sh to the smoke tier, which builds stamped in the layer's image under the harness's mounts and fails without the safe.directory.
