---
id: T-3617
type: task
nature: remediation
title: Close I-0407 saying what fixed it
status: backlog
parent: S-0848
owner: alex
created: 2026-11-27T23:56:47Z
updated: 2026-11-27T23:56:47Z
transitions: []
stream: S-0848
tags: [issues]
touches: [design/issues/I-0407-flai-stream-sync-panics-while-continuing-a-rebase-whose-stop-the-agent-resolved.md, design/issues/summary.md]
after: [T-3615, T-3616]
---
# T-3617 Close I-0407 saying what fixed it

## Work

Run `flai issue close I-0407 --reason "<what fixed it>"`. The reason names the input that panicked the section merge and its fix from T-3615's notes, or says that none was found. It also names the guard from T-3616, which leaves a stop waiting with its paths when a resolver panics.

It waits for T-3615 and T-3616, because the reason states what both found and changed.

## Done when

- I-0407's status is closed, with a reason that names both fixes.
- `design/issues/summary.md` no longer lists I-0407 as open.
- `flai check --strict` is clean on the two paths.

## Notes
