---
id: T-3594
type: task
nature: remediation
title: Close I-0370 with a reason naming S-0781's fix and the branches deleted on TH-1190
status: in-progress
parent: S-0841
owner: alex
created: 2026-11-17T12:35:26Z
updated: 2026-11-17T12:55:25Z
transitions:
  - to: ready
    at: 2026-11-17T12:55:25Z
    by: agent-S-0841
  - to: in-progress
    at: 2026-11-17T12:55:25Z
    by: agent-S-0841
stream: S-0841
tags: [issues]
touches: [design/issues/I-0370-scripts-flai-test-sh-s-repository-check-fails-on-this-host-while-story-s-0001-s.md, design/issues/summary.md]
after: [T-3593]
---
# T-3594 Close I-0370 with a reason naming S-0781's fix and the branches deleted on TH-1190

## Work

This task waits for T-3593, because the reason it writes cites T-3593's verification. Without it, the issue would be closed on a claim nobody checked.

In the story's worktree, run `flai issue close I-0370 --reason "<reason>"`. The reason must say three things:

- The cause was test temp projects with no `git init` of their own. Under a TMPDIR inside a worktree, `flai story start` in those tests made `story/S-0001`, `S-0002`, `S-0003`, and `S-0005` in the main repository.
- S-0781 (I-0371) fixed it: the TestMain of `flai/cmd` and of `flai/internal/mcpserver` sets `GIT_CEILING_DIRECTORIES`, and `TestStoryStartInATemporaryFolderInsideARepositoryMakesNoBranch` reproduces the leak.
- The four branches were deleted on TH-1190, and S-0841 confirmed that none is left and that the unscoped check reports no `story.unaccepted`.

flai rewrites `design/issues/summary.md` when it closes the issue. Commit both files on the story branch. Then tick criterion 2 with `flai criteria tick S-0841 2`.

## Done when

- I-0370's front matter reads `status: closed`, and its reason names S-0781's fix and T-3593's verification.
- `design/issues/summary.md` no longer lists I-0370.
- Both files are committed on `story/S-0841`, and criterion 2 is ticked.

## Notes
