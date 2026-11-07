---
id: T-2535
type: task
nature: remediation
title: flai stream sync and flai accept put a rebased branch's changelog lines back under main's top Unreleased before checking what the rebase undoes
status: done
parent: S-0570
owner: alex
created: 2026-11-07T03:48:21Z
updated: 2026-11-07T04:49:20Z
transitions:
  - to: ready
    at: 2026-11-07T04:16:34Z
    by: system-flow
  - to: in-progress
    at: 2026-11-07T04:35:25Z
    by: system-flow
  - to: done
    at: 2026-11-07T04:49:19Z
    by: agent-S-0570
stream: S-0570
tags: [cli, git]
touches: [flai/internal/storygit/sync.go, flai/internal/storygit/undone.go, flai/internal/storygit/sync_test.go, flai/internal/issues/generated.go, flai/cmd/stream_sync.go, flai/cmd/stream_sync_test.go, flai/cmd/stream.go]
after: [T-2531]
usage:
  source: log
  seconds: 834
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 71
      output: 22091
      cache_read: 9742353
      cache_write: 147169
      cost: 3.5038
---
# T-2535 flai stream sync and flai accept put a rebased branch's changelog lines back under main's top Unreleased before checking what the rebase undoes

## Work

Waits for T-2531: it calls the changelog merge T-2531 writes.

`storygit.Rebase` in `flai/internal/storygit/sync.go` is the rebase both `flai stream sync` and `flai accept` (`flai/cmd/branch.go` `syncStoryBranch`) run. Before the `git rebase`, it reads the old base (`oldBase`) and the head. After a clean rebase, and before `checkUndone`, add the step T-2524's ADR names:

1. For each `CHANGELOG.md` the branch changed between the old base and the old head, read the lines it added with `git diff`.
2. Write the file again from the base branch's version with those lines, through a hook on `storygit.Files` that `issues.FilesIn` in `flai/internal/issues/generated.go` fills with T-2531's function.
3. Commit the change on the branch as one commit of its own, `docs: [S-nnnn] keep CHANGELOG.md's lines under ## Unreleased`, only when the file changed.
4. Report the paths in a new `SyncResult` field. `flai/cmd/stream_sync.go` prints it and gives it in `--json`.

In `flai/internal/storygit/undone.go`, leave the rewritten changelogs out of `checkUndone`, as the generated and merged files are left out.

Tests with real git:

- **`flai/internal/storygit/sync_test.go`.** A branch whose line sits under `## Unreleased`, rebased over a main commit that renames the section to a version and adds an empty `## Unreleased` on top, as `flai release` does. The rebased file must equal main's with the branch's line under the top `## Unreleased`, and the sync must report nothing undone.
- **`flai/cmd/stream_sync_test.go`.** Next to `unionStories`, the same with I-0251's second instance: two releases in a row, and a version heading main holds twice. Check the CLI output names the rewritten changelog.

## Done when

- Both tests fail on the code before the change and pass after it.
- `flai test` passes on the paths this task changes.

## Notes

- Drafted by the planner. `fold` regenerates only when there are folds (`sync.go`), so the changelog step is a call of its own, not one of `Files.Generated`.
