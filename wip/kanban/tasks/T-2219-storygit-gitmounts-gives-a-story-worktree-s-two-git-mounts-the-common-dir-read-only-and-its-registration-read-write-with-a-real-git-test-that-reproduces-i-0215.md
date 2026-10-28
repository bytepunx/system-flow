---
id: T-2219
type: task
nature: remediation
title: storygit.GitMounts gives a story worktree's two git mounts, the common dir read-only and its registration read-write, with a real-git test that reproduces I-0215
status: backlog
parent: S-0517
owner: alex
created: 2026-10-28T13:20:19Z
updated: 2026-10-28T13:20:19Z
transitions: []
stream: S-0517
tags: [cli, docker]
touches: [flai/internal/storygit/mounts.go, flai/internal/storygit/mounts_test.go]
after: [T-2218]
---
# T-2219 storygit.GitMounts gives a story worktree's two git mounts, the common dir read-only and its registration read-write, with a real-git test that reproduces I-0215

## Work

Add `flai/internal/storygit/mounts.go` with `GitMounts(worktree string) ([]Mount, error)`, where `Mount` is a host path and whether it is read-only, as the ADR of T-2218 has it. Read `<worktree>/.git`: a file `gitdir: <path>` names the registration, absolute or relative to the worktree (relative under `worktrees.relative_paths`, ADR-0022); resolve it, read its `commondir`, resolve that against the registration, and return the common dir read-only and the registration read-write, both as cleaned absolute paths. A `.git` that is a folder, a normal checkout, returns no mounts, since the worktree carries its repository; a missing `.git`, a `gitdir` that does not exist, or a registration without `commondir` is an error naming the path, so a start fails with the reason rather than a container without git. Keep it beside `open.go`, which makes the worktree the helper describes, and leave `workitem.mainCheckout` (`flai/internal/workitem/repo.go`) as it is: it answers where the project is, not what to mount.

Test it with the package's real git, as `open_test.go` does: a repository with a worktree made by `git worktree add` gives the two paths, and a relative-path worktree (skipped with a reason where git is older than 2.48, as ADR-0022's tests are) gives the same absolute ones. Reproduce I-0215 in one test: with the common dir made read-only below the registration (chmod, skipped as root, since root ignores file modes) and the registration left writable, `git status` and `git log` in the worktree succeed and `git commit` is refused by git; and, as the issue's shape, that a worktree whose `gitdir` is absent fails `git status` with `not a git repository`, which is what the mounts prevent.

## Done when

- `storygit.GitMounts` returns the common dir read-only and the story's registration read-write for an absolute-path and a relative-path worktree, nothing for a normal checkout, and an error naming the path for a worktree whose registration is absent
- A real-git test shows `git status` works and `git commit` is refused with the common dir read-only and the registration writable, and `flai test flai/internal/storygit` passes

## Notes

Layer 2, with T-2220 (the script): both wait for the ADR, and neither touches the other's paths. T-2221, the adapter, calls this helper.
