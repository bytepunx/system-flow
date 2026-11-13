---
id: T-3270
type: task
nature: improvement
title: Leave item.blocked-closed outside the story out of a check scoped to a story, with a test that reproduces I-0354
status: done
parent: S-0761
owner: alex
created: 2026-11-13T07:29:45Z
updated: 2026-11-13T08:58:02Z
transitions:
  - to: ready
    at: 2026-11-13T08:24:48Z
    by: agent-S-0761
  - to: in-progress
    at: 2026-11-13T08:24:48Z
    by: agent-S-0761
  - to: done
    at: 2026-11-13T08:58:02Z
    by: agent-S-0761
stream: S-0761
tags: [flai, check]
touches: [flai/internal/check/scope.go, flai/internal/check/scope_test.go, flai/cmd/check.go, flai/cmd/check_test.go, docs/users/flai-reference.md]
after: [T-3269]
usage:
  source: log
  seconds: 1994
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 19
      output: 3901
      cache_read: 1627561
      cache_write: 68873
      cost: 0.8394
---
# T-3270 Leave item.blocked-closed outside the story out of a check scoped to a story, with a test that reproduces I-0354

## Work

Build what T-3269's ADR decides. If it decides to leave the finding out, do the following.

In `ScopeToStories` (`flai/internal/check/scope.go`), add `item.blocked-closed` to the rules dropped with `res.drop` when `!inside(f.Path)`, beside `item.archive`, `board.wip-limit`, `issues.no-story`, `threads.archived`, and `narrative.missing`. The finding's path is the other item's file, archived or not (`check_items.go`), so `inside` already tells the story's own item and tasks apart from others. Extend the function's doc comment with the ADR and I-0354.

In `flai/internal/check/scope_test.go`, reproduce I-0354 beside `TestScopedCheckLeavesOutANarrativeMissingOutsideTheStory`. Archive a done story with an open blocked interval, as S-0678 was, and check with `--story` on another story. The result keeps no `item.blocked-closed` and counts none as outside. A done task of the scoped story with an open interval keeps its own finding. A check without `--story` still reports both.

In `flai/cmd/check.go`, update the comment of `recordOutside` and the `--record-issues` text of the command's long help to say an `item.blocked-closed` outside the story never reaches an issue. Add `TestCheckRecordIssuesLeavesOutAnItemBlockedClosedOutsideTheStory` to `flai/cmd/check_test.go`: `flai check --story S-nnnn --record-issues` over such a repository opens or bumps no issue. Regenerate `docs/users/flai-reference.md` with `make flai-reference`.

It waits for T-3269 so that the comments and the help cite the ADR's number, and so that it builds what the ADR decides.

## Done when

- `flai check --story S-nnnn --record-issues`, run while another closed item has an open blocked interval, prints no `item.blocked-closed` and records no issue, as the new tests show.
- One on the scoped story's own item or tasks, and `flai check` without `--story`, still report it.
- `flai test` on the changed paths passes.

## Notes
