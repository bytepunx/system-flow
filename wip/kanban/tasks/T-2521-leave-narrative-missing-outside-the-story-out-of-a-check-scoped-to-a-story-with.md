---
id: T-2521
type: task
nature: improvement
title: Leave narrative.missing outside the story out of a check scoped to a story, with a test that reproduces I-0235
status: done
parent: S-0555
owner: alex
created: 2026-11-07T03:47:27Z
updated: 2026-11-07T04:13:45Z
transitions:
  - to: ready
    at: 2026-11-07T03:59:34Z
    by: system-flow
  - to: in-progress
    at: 2026-11-07T03:59:36Z
    by: system-flow
  - to: done
    at: 2026-11-07T04:13:45Z
    by: agent-S-0555
stream: S-0555
tags: [flai, check]
touches: [flai/internal/check/scope.go, flai/internal/check/scope_test.go, flai/cmd/check.go, flai/cmd/check_test.go, docs/users/flai-reference.md]
after: [T-2518]
usage:
  source: log
  seconds: 849
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 40
      output: 308
      cache_read: 4936502
      cache_write: 303340
      cost: 1.8728
---
# T-2521 Leave narrative.missing outside the story out of a check scoped to a story, with a test that reproduces I-0235

## Work

Build what T-2518's ADR decides. If it decides to leave the finding out, do the following.

In `ScopeToStories` (`flai/internal/check/scope.go`), add `narrative.missing` to the rules dropped with `res.drop` when `!inside(f.Path)`, beside `item.archive`, `board.wip-limit`, `issues.no-story`, and `threads.archived`. The finding's path is the other story's item file, so `inside` already tells the two apart. Extend the function's doc comment with the ADR and I-0235. The finding comes from `narratives()` in `flai/internal/check/check.go`. Leave it as it is.

In `flai/internal/check/scope_test.go`, reproduce I-0235 beside `TestScopedCheckLeavesOutAThreadsArchivedOutsideTheStory`. Set up two stories in progress, one with a narrative and one without, as S-0549 was. Check with `--story` on the one with a narrative. The result keeps no `narrative.missing` and counts none as outside. A check scoped to the story with no narrative keeps its own finding. A check without `--story` still reports it.

In `flai/cmd/check.go`, update the comment of `recordOutside` and the `--record-issues` text of the command's long help. Both should say that a `narrative.missing` outside the story never reaches an issue. Add `TestCheckRecordIssuesLeavesOutANarrativeMissingOutsideTheStory` to `flai/cmd/check_test.go`. It runs `flai check --story S-nnnn --record-issues` over such a story and finds no issue opened or bumped. Regenerate `docs/users/flai-reference.md` with `make flai-reference`.

It waits for T-2518 so that the comments and the help cite the ADR's number.

## Done when

- `flai check --story S-nnnn --record-issues`, run while another story is in progress with no narrative, prints no `narrative.missing` and records no issue, as the new tests show.
- A check scoped to the story that lacks a narrative, and `flai check` without `--story`, still report it.
- `flai test` on the changed paths passes.

## Notes
