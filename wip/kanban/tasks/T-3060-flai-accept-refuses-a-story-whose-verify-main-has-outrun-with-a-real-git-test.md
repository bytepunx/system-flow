---
id: T-3060
type: task
nature: remediation
title: flai accept refuses a story whose verify main has outrun, with a real-git test that reproduces I-0229
status: backlog
parent: S-0683
owner: alex
created: 2026-11-10T10:35:32Z
updated: 2026-11-10T10:35:32Z
transitions: []
stream: S-0683
tags: [flai, acceptance, tests]
touches: [flai/cmd/accept.go, flai/cmd/accept_stale_verify_test.go]
after: [T-3058]
---
# T-3060 flai accept refuses a story whose verify main has outrun, with a real-git test that reproduces I-0229

## Work

It waits for T-3058, which computes the blocker `verify_stale` that this task makes `flai accept` refuse on.

Make `flai accept` and `flai move <story> done` refuse on the blocker before anything is merged, the way they refuse on the preview's other blockers. Print its message, and add it to `--dry-run` and `--json`. Change `flai/cmd/accept.go` only where the refusal or its output needs it.

In the new `flai/cmd/accept_stale_verify_test.go`, a real-git test reproduces I-0229 (it skips under `-short`, as the other real-git tests do):

- Package `a` holds a string.
- Story X changes the string and is accepted.
- Story Y, verified before X landed, adds a test in package `b`, which imports `a`, that asserts the old string.
- Accepting Y is refused with `verify_stale`, naming the go-test tier and `a`'s file.
- After Y's verify runs again at its rebased head, the acceptance proceeds or fails on its own test, never landing red unseen.
- A second case: when main brought in only `wip/` paths, Y is accepted.

## Done when

- The new test fails on main as it is, before T-3058, and passes after.
- `flai accept --dry-run --json` shows `verify_stale` with its tiers and paths.
- `flai test flai/cmd/accept.go flai/cmd/accept_stale_verify_test.go` passes.

## Notes
