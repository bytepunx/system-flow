---
id: T-2970
type: task
nature: remediation
title: Close I-0338 with the reason that the sync now compares against the commit its rebase took
status: backlog
parent: S-0712
owner: alex
created: 2026-11-10T10:28:27Z
updated: 2026-11-10T10:28:27Z
transitions: []
stream: S-0712
tags: [issues]
touches: [design/issues/I-0338-flai-task-done-s-sync-reports-undoes-main-for-paths-the-story-never-changed.md, design/issues/summary.md]
after: [T-2965, T-2967]
---
# T-2970 Close I-0338 with the reason that the sync now compares against the commit its rebase took

## Work

Close the issue once the fix, its test, and its docs are in: `flai issue close I-0338 --reason "<why>"`. The reason names the fix: `storygit.Rebase` takes main's head once, as a commit, and the undone check compares against that commit, so a commit main gains during the sync is not reported (T-2963). It names the test that reproduces the cause (T-2965). `flai issue close` writes the issue file and regenerates `design/issues/summary.md`.

It waits for T-2965 and T-2967: the reason cites the test, and the issue is closed only once the fix is tested and described.

## Done when

- I-0338's status is closed, with a reason naming the fix and the test.
- `design/issues/summary.md` no longer lists I-0338.
- `flai check --strict` is clean for the story.

## Notes
