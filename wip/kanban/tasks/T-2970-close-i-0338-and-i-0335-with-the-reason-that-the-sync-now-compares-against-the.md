---
id: T-2970
type: task
nature: remediation
title: Close I-0338 and I-0335 with the reason that the sync now compares against the commit its rebase took
status: backlog
parent: S-0712
owner: alex
created: 2026-11-10T10:28:27Z
updated: 2026-11-11T09:05:13Z
transitions: []
stream: S-0712
tags: [issues]
touches: [design/issues/I-0338-flai-task-done-s-sync-reports-undoes-main-for-paths-the-story-never-changed.md, design/issues/I-0335-flai-stream-sync-reports-a-path-as-undone-by-the-branch-when-main-gains-a.md, design/issues/summary.md]
after: [T-2965, T-2967]
---
# T-2970 Close I-0338 and I-0335 with the reason that the sync now compares against the commit its rebase took

## Work

Close both issues once the fix, its tests, and its docs are in: `flai issue close I-0338 --reason "<why>"` and `flai issue close I-0335 --reason "<why>"`. Each reason names the fix: `storygit.Rebase` takes main's head once, as a commit, and the undone check and the changelog settle use that commit, so a commit main gains during the sync is not reported (T-2963). It names the test that reproduces the cause (T-2965). `flai issue close` writes each issue file and regenerates `design/issues/summary.md`.

It waits for T-2965 and T-2967: the reasons cite the tests, and the issues are closed only once the fix is tested and described.

## Done when

- I-0338's and I-0335's status is closed, each with a reason naming the fix and the test.
- `design/issues/summary.md` lists neither.
- `flai check --strict` is clean for the story.

## Notes
