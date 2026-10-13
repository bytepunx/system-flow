---
id: T-1514
type: task
nature: feature
title: flai doc show and doc save take --story and edit the file in that story's worktree, committing on its branch
status: done
parent: S-0378
owner: alex
created: 2026-10-13T16:56:25Z
updated: 2026-10-13T17:49:24Z
transitions:
  - to: ready
    at: 2026-10-13T17:30:36Z
    by: agent-S-0378
  - to: in-progress
    at: 2026-10-13T17:30:37Z
    by: agent-S-0378
  - to: done
    at: 2026-10-13T17:49:24Z
    by: agent-S-0378
stream: S-0378
tags: [cli]
touches: [flai/cmd/doc.go, flai/cmd/doc_test.go, flai/internal/docedit/docedit.go, flai/internal/docedit/docedit_test.go, flai/internal/hostapi/writes.go, flai/internal/hostapi/writes_test.go, flai/cmd/doc_lint_test.go, docs/users/flai-reference.md, docs/operators/settings.md]
usage:
  source: log
  seconds: 1127
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 111
      output: 499
      cache_read: 5294201
      cache_write: 163041
      cost: 2.4979
---
# T-1514 flai doc show and doc save take --story and edit the file in that story's worktree, committing on its branch

## Work

- Add `Story` to `docedit.SaveOptions` and a `--story S-nnnn` flag to `flai doc show` and `flai doc save` in `flai/cmd/doc.go`. With it, the path resolves in the story's worktree (`.flai-cache/worktrees/S-nnnn`), the hash and the check run there, and the save commits that one path on `story/S-nnnn`.
- Refuse `--story` for a story with no worktree, and for a path under the wip folder: wip lives on main (ADR-0019). `ModeOf` resolves wip against the main checkout today, so a worktree's wip file would get mode `full`.
- Make `docedit.Save` run `check.Run` against the worktree's root when a story is given, so the before-and-after comparison is the branch's.
- Give the host methods `doc.show` and `doc.save` in `flai/internal/hostapi/writes.go` an optional `story`, and run them with `dir` set to the worktree, as `testDir` does for `test.run`.
- Waits for nothing. It is the first layer, beside the lint task, with which it shares no file.

## Done when

- `flai doc save design/x.md --story S-nnnn --hash <h>` writes and commits the file on the story branch, leaves the main checkout unchanged, and refuses a stale hash with exit 3.
- A wip path or a story without a worktree is refused with the reason.
- `flai/cmd/doc_test.go`, `flai/internal/docedit/docedit_test.go`, and `flai/internal/hostapi/writes_test.go` cover the story save, its conflict, and the refusals, and pass.

## Notes

- Criterion 3 of S-0378: the ADR-0140 refusal names files on the story branch after the rebase and its fix, so they are corrected there.
