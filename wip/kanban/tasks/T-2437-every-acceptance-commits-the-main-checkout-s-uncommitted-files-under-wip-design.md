---
id: T-2437
type: task
nature: feature
title: Every acceptance commits the main checkout's uncommitted files under wip, design, and the shared paths, and leaves the rest as they are
status: done
parent: S-0566
owner: alex
created: 2026-11-02T18:35:38Z
updated: 2026-11-02T19:52:15Z
transitions:
  - to: ready
    at: 2026-11-02T19:03:07Z
    by: system-flow
  - to: in-progress
    at: 2026-11-02T19:03:08Z
    by: system-flow
  - to: done
    at: 2026-11-02T19:52:15Z
    by: agent-S-0566
stream: S-0566
tags: []
touches: [flai/internal/storygit/storygit.go, flai/internal/storygit/storygit_test.go, flai/internal/preview/accept.go, flai/internal/preview/accept_test.go, flai/cmd/accept.go, flai/cmd/accept_batch_dirty.go, flai/cmd/accept_batch_dirty_test.go, flai/cmd/accept_markdown.go, flai/cmd/accept_uncommitted_test.go, flai/cmd/accept_orchestrator_dirty_test.go, flai/cmd/accept_test.go, flai/internal/mergequeue/eligible.go, flai/internal/mergequeue/eligible_test.go, flai/internal/mergequeue/land.go, flai/internal/mergequeue/land_test.go, flai/internal/mergequeue/preflight.go, flai/internal/mergequeue/preflight_test.go, flai/internal/hostapi/writes.go, docs/users/flai-reference.md, flai/internal/harness/harness.go, flai/internal/harness/harness_test.go]
after: [T-2436]
usage:
  source: log
  seconds: 2934
  estimated: true
  priced_by: estimate
  models:
    - model: claude-fable-5-1
      input: 3274
      output: 1255
      cache_read: 26034200
      cache_write: 343968
      cost: 39.9973
---
# T-2437 Every acceptance commits the main checkout's uncommitted files under wip, design, and the shared paths, and leaves the rest as they are

## Work

ADR-0248 (`design/adrs/0248-*.md`), the uncommitted-files part. Today `storygit.DirtyOutsideWip` lists the main checkout's uncommitted paths outside the wip folder; the acceptance by hand refuses them unless `--yes` (`acceptItem` in `flai/cmd/accept.go`, "working tree has uncommitted changes outside wip"), and `--yes` stages everything with `git add -A`; the orchestrator's acceptance and the queue's landing leave them out of the commit (`leftUncommitted` and `addArgs` in `flai/cmd/accept_batch_dirty.go`) and are refused or held back when the story's change would overwrite one (`dirtyMainBlocker` in `flai/internal/preview/accept.go`, `mainDirt.holds` in `flai/internal/mergequeue/eligible.go`, `overwrittenBy` in `flai/internal/mergequeue/land.go`); the queue's preflight lints the files outside wip (`flai/internal/mergequeue/preflight.go`).

Build the shared set: in `storygit`, replace `DirtyOutsideWip` with a reader of the main checkout's uncommitted and untracked files (as `DirtyFiles` reads them) split into those under the manifest's `wip` folder, its `design` folder, or under or at a path in `claims.shared` (the shared set, which every acceptance commits) and the rest, with a predicate every caller can use. Then:

- `preview.Acceptance` gains `SharedUncommitted []string` (`shared_uncommitted`): the shared set; `Uncommitted` (`uncommitted`) becomes the rest, left as they are. The `dirty_main` blocker applies to every acceptance but `--batch`, not the orchestrator's alone, and reads only the rest: a shared file never blocks, since the acceptance commits it first. Reword its message and the dry-run lines in `printAccept` for the operator ("the acceptance leaves them as they are; it is refused only while the story's change would overwrite one … `--yes` includes them").
- In `acceptItem`, drop the refusal of uncommitted files outside wip. `leftUncommitted` answers the rest for every acceptance, except the acceptance by hand with `--yes`, which includes them (carry the choice on the preview, say `IncludesUncommitted`), and `addArgs` excludes them from the acceptance commit as it does today for the queue; the lint fix (`refuseUnfixedUncommitted`, `fixUncommitted` in `flai/cmd/accept_markdown.go`) runs over what will be committed, the shared set included, and never over the rest.
- Before the merge (step 0 in `acceptItem`, for a story with a branch, under local acceptance and `--batch` alike), when the story's change `main...story/S-nnnn` would overwrite a file of the shared set (`storygit.ChangedAmong` over the shared set), run the lint fix over the shared set and commit it on the main branch as `chore: [S-nnnn] commit the main checkout's wip and design before accepting`, through `storygit.RunPastIndexLock` as the acceptance commit is, log an `acceptStep` for it, and go on to the rebase and merge; a commit that fails refuses the acceptance with git's output and what to do. When nothing would be overwritten, the shared set goes with the acceptance commit as `wip/` did. `answerRelease` (the release-due event) is read before that commit.
- The merge queue's eligibility (`mainDirt`), the landing's `overwrittenBy`, and the preflight's `left` count only the rest: a dirty file in the shared set holds no story back and is lint-fixed by the preflight as a file the landing commits.
- `flai/internal/hostapi/writes.go`: `accept.run` still takes `include_uncommitted` from older dashboards and passes `--yes`; update its comment to say what `--yes` includes now.
- Update the `flai accept` help text in `accept.go` (the paragraphs on `--yes`, "It leaves the main checkout's uncommitted files outside wip …", and `dirty_main`) to ADR-0248's rule.

Tests (standard library `testing`, in the files named in the touches): the split of dirty files by the shared set, with a `claims.shared` path and an untracked file under `design/`; an acceptance by hand with a dirty `design/issues/summary.md` the branch also changes makes the pre-commit with its subject and then the acceptance commit, and leaves a dirty code file outside the set uncommitted; an acceptance by hand with a dirty code file the branch changes is refused as `dirty_main` naming the file; `--yes` includes it; the dry run's JSON has `shared_uncommitted` and `uncommitted`; the queue's eligibility and landing do not hold a story back for a dirty shared file. Keep `flai/cmd/accept_orchestrator_dirty_test.go` and `accept_batch_dirty_test.go` passing, adjusting them to the rule. Run `flai test` on the paths you changed.

This task waits for the send-back task of S-0566, which changes `flai/cmd/accept.go` and `flai/internal/preview/accept.go` too; it is layer 2.

## Done when

- An acceptance by hand or from the dashboard with uncommitted files under `wip/`, `design/`, or `claims.shared` in the main checkout commits them with the acceptance, or before the merge when the story's change would overwrite one, asking nothing; one with an uncommitted file elsewhere leaves it, and is refused only when the story's change would overwrite it, naming it as `dirty_main`.
- The orchestrator's acceptance and the queue's landing, eligibility, and preflight treat a dirty shared file as clean and commit it the same way.
- `flai accept --dry-run --json` prints `shared_uncommitted` and `uncommitted`.
- The help text and the tests named above are updated and `flai test` passes on the changed paths.

## Notes

Acceptance criterion 3 of S-0566.
