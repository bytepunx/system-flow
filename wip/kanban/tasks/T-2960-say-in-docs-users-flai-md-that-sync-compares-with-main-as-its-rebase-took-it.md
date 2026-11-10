---
id: T-2960
type: task
nature: remediation
title: Say in docs/users/flai.md that sync compares with main as its rebase took it, and close I-0335
status: backlog
parent: S-0709
owner: alex
created: 2026-11-10T10:27:35Z
updated: 2026-11-10T10:27:35Z
transitions: []
stream: S-0709
tags: [cli]
touches: [docs/users/flai.md, design/issues/I-0335-flai-stream-sync-reports-a-path-as-undone-by-the-branch-when-main-gains-a.md, design/issues/summary.md]
after: [T-2959]
---
# T-2960 Say in docs/users/flai.md that sync compares with main as its rebase took it, and close I-0335

## Work

In `docs/users/flai.md`, under Story branches, add to the "Main's change undone" paragraph that sync compares the rebased branch with main's commit as the rebase went onto it, so a commit that reaches main while the sync runs is not reported as undone; the next sync rebases onto it. Then run `flai issue close I-0335 --reason` naming T-2959's fix and its test, which updates `design/issues/I-0335-flai-stream-sync-reports-a-path-as-undone-by-the-branch-when-main-gains-a.md` and `design/issues/summary.md`.

Waits for T-2959: the paragraph describes its behaviour, and the reason names its fix and test, which must exist first.

## Done when

- The paragraph says what sync compares with, and the markdown lint passes on `docs/users/flai.md`.
- I-0335 is closed, its reason naming the fix and the test, and `design/issues/summary.md` no longer lists it open.
- `flai check --strict` is clean.

## Notes
