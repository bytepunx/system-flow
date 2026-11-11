---
id: T-3093
type: task
nature: remediation
title: The smoke tier's flai-build-stamp-test.sh reproduces I-0348 in a scratch main checkout with a nested story worktree
status: done
parent: S-0724
owner: alex
created: 2026-11-11T08:54:31Z
updated: 2026-11-11T14:40:33Z
transitions:
  - to: ready
    at: 2026-11-11T14:38:23Z
    by: agent-S-0724
  - to: in-progress
    at: 2026-11-11T14:38:23Z
    by: agent-S-0724
  - to: done
    at: 2026-11-11T14:40:33Z
    by: agent-S-0724
stream: S-0724
tags: [go, build, worktree, test]
touches: [scripts/flai-build-stamp-test.sh, scripts/smoke.sh, scripts/README.md]
usage:
  source: log
  seconds: 130
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 17
      output: 3925
      cache_read: 1200500
      cache_write: 37882
      cost: 0.6015
---
# T-3093 The smoke tier's flai-build-stamp-test.sh reproduces I-0348 in a scratch main checkout with a nested story worktree

## Work

Write `scripts/flai-build-stamp-test.sh`, a test that reproduces I-0348 on the host, offline, with no Docker.

- In a temporary folder, make a scratch main checkout with `git init` and one commit. Give it a copy of `scripts/env.sh` and `scripts/flai-build.sh`, and a small `flai/` module whose module path is `github.com/bytepunx/system-flow/flai`. The module has an `internal/buildinfo` package with `Commit` and `Date` variables that fall back to `debug.ReadBuildInfo`, as the real one does, and a `main` that prints `Commit`.
- Add a story worktree at `.flai-cache/worktrees/S-0001` in it with `git worktree add`, on a branch with a commit of its own. That is where flai puts a story's worktree.
- Run the worktree's `scripts/flai-build.sh` and its `bin/flai`. The commit printed must be the worktree's HEAD, `git -C <worktree> rev-parse HEAD`, and not the main checkout's.
- Skip, saying why, when `go` or `git` is not on `PATH`.

Have `scripts/smoke.sh` run it beside `agent-vcs-test.sh`, and add its row to `scripts/README.md`.

It runs first: it fails against `scripts/flai-build.sh` as it is now, which is the reproduction the first criterion asks for. The next task's fix is checked by it.

## Done when

- `scripts/flai-build-stamp-test.sh` fails against the current `scripts/flai-build.sh`, and its message names both commits: the main checkout's that was printed and the worktree's that was expected.
- `scripts/smoke.sh` runs it, and `scripts/README.md` lists it.
- It cleans up its temporary folder, and it passes `shellcheck` in the shell lint tier.

## Notes
