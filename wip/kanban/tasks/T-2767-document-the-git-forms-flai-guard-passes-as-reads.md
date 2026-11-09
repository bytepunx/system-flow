---
id: T-2767
type: task
nature: remediation
title: Document the git forms flai guard passes as reads
status: done
parent: S-0623
owner: alex
created: 2026-11-09T21:56:05Z
updated: 2026-11-09T22:23:12Z
transitions:
  - to: ready
    at: 2026-11-09T22:21:17Z
    by: system-flow
  - to: in-progress
    at: 2026-11-09T22:21:19Z
    by: system-flow
  - to: done
    at: 2026-11-09T22:23:12Z
    by: agent-S-0623
stream: S-0623
tags: [flai, guard, docs]
touches: [flai/cmd/guard.go, docs/users/flai-reference.md, design/system/flai-cli.md]
after: [T-2766]
usage:
  source: log
  seconds: 113
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 7
      output: 1107
      cache_read: 809931
      cache_write: 24298
      cost: 0.3433
---
# T-2767 Document the git forms flai guard passes as reads

## Work

Add the read forms T-2766 settles (`git tag` that lists, `git worktree list`, `git branch` that lists) to the git reads the help of `flai guard` in `flai/cmd/guard.go` names, regenerate `docs/users/flai-reference.md` with `make flai-reference`, and add them to the `flai guard` row of `design/system/flai-cli.md`, which lists the same reads. Waits for T-2766, because it documents the forms that task's code decides.

## Done when

- The help, the generated reference, and `design/system/flai-cli.md` name the three read forms in the same words.
- `make flai-reference` leaves `docs/users/flai-reference.md` unchanged after the commit.
- `flai test` on the three paths passes.

## Notes
