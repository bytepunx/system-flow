---
id: T-3593
type: task
nature: remediation
title: Confirm S-0781's fix keeps I-0370's cause from occurring, and record the evidence in the story's Notes
status: done
parent: S-0841
owner: alex
created: 2026-11-17T12:35:16Z
updated: 2026-11-17T12:55:05Z
transitions:
  - to: ready
    at: 2026-11-17T12:54:35Z
    by: agent-S-0841
  - to: in-progress
    at: 2026-11-17T12:54:38Z
    by: agent-S-0841
  - to: done
    at: 2026-11-17T12:55:05Z
    by: agent-S-0841
stream: S-0841
tags: [issues, testing]
usage:
  source: log
  seconds: 27
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 6
      output: 136
      cache_read: 524043
      cache_write: 3523
      cost: 0.2018
---
# T-3593 Confirm S-0781's fix keeps I-0370's cause from occurring, and record the evidence in the story's Notes

## Work

I-0370's cause is already fixed. S-0777's discarded test run set TMPDIR inside its worktree. Its temp projects had no `git init` of their own, so git walked up into the worktree, and `flai story start` in the tests made `story/S-0001`, `S-0002`, `S-0003`, and `S-0005` in the main repository. S-0781 fixed that leak (I-0371, D-0617): the TestMain of `flai/cmd` and of `flai/internal/mcpserver` sets `GIT_CEILING_DIRECTORIES` to the temporary folder. `TestStoryStartInATemporaryFolderInsideARepositoryMakesNoBranch` in `flai/cmd/story_start_test.go` reproduces the leak. The four branches were deleted on the operator's word on TH-1190 (D-0625).

This task proves that, and changes no code:

1. In the main checkout, `git branch --list 'story/S-000*'` lists nothing.
2. The main checkout's unscoped `scripts/flai.sh check --strict` reports no `story.unaccepted`. It may fail on findings outside I-0370, such as the `threads.archived` on TH-1232. Note those and leave them alone.
3. From the worktree's `flai/` folder, run `go test ./cmd -run 'TestStoryStartInATemporaryFolderInsideARepositoryMakesNoBranch' -count=1` without `-short`, since the test skips under `-short`. It must pass.

The planner's explorer found no other test package that writes refs in a temp folder it did not `git init` itself. If you find one, write the test and the `GIT_CEILING_DIRECTORIES` TestMain that D-0617 asks for, and add its files to this task's touches.

Write each command's result in the story's `## Notes` under a `### Verification` heading, then tick criterion 1 with `flai criteria tick S-0841 1`.

## Done when

- `git branch --list 'story/S-000*'` lists nothing in the main checkout.
- The main checkout's unscoped `flai check --strict` reports no `story.unaccepted`.
- `TestStoryStartInATemporaryFolderInsideARepositoryMakesNoBranch` passes without `-short`.
- The story's Notes record the three results under `### Verification`, and criterion 1 is ticked.

## Notes
