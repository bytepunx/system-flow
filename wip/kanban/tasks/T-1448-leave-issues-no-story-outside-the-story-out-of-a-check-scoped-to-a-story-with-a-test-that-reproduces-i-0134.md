---
id: T-1448
type: task
nature: improvement
title: Leave issues.no-story outside the story out of a check scoped to a story, with a test that reproduces I-0134
status: done
parent: S-0365
owner: alex
created: 2026-10-12T14:56:34Z
updated: 2026-10-12T15:04:28Z
transitions:
  - to: ready
    at: 2026-10-12T14:58:36Z
    by: agent-S-0365
  - to: in-progress
    at: 2026-10-12T14:59:31Z
    by: agent-S-0365
  - to: done
    at: 2026-10-12T15:04:28Z
    by: agent-S-0365
stream: S-0365
tags: [flai, check]
touches: [flai/internal/check/scope.go, flai/internal/check/scope_test.go, flai/cmd/check.go, flai/cmd/check_test.go, docs/users/flai-reference.md]
after: [T-1447]
usage:
  source: log
  seconds: 297
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 50
      output: 349
      cache_read: 1924321
      cache_write: 96763
      cost: 0.9078
---
# T-1448 Leave issues.no-story outside the story out of a check scoped to a story, with a test that reproduces I-0134

## Work

In `ScopeToStory` (`flai/internal/check/scope.go`), add `issues.no-story` to the rules dropped with `res.drop` when `!inside(f.Path)`, beside `item.archive` and `board.wip-limit`, as T-1447's ADR decides. If the ADR decides a note instead, keep the finding outside and leave it unrecorded, as `wip.overlap` is. Extend the function's doc comment with the ADR and I-0134. The finding comes from `issuesNoStory` in `flai/internal/check/check.go`; leave its production alone.

In `flai/internal/check/scope_test.go`, reproduce I-0134. Set up a repository with an open issue first reported longer ago than `issues.story_after` that no open story links, and a story that does not change it. Checked with `--story` on that story, it keeps no `issues.no-story` and counts none as outside. Checked with the story's changed paths naming the issue file, it keeps the finding. Checked without `--story`, it still warns.

In `flai/cmd/check.go`, update `recordOutside`'s comment and the `--record-issues` text of the command's long help to say an `issues.no-story` outside the story never reaches an issue. Add a case to `flai/cmd/check_test.go` that runs `flai check --story S-nnnn --record-issues` over such an issue and finds no issue opened or bumped. Regenerate `docs/users/flai-reference.md` with `make flai-reference`.

It waits for T-1447 so that the comments and the help cite the ADR's number.

## Done when

- `flai check --story S-nnnn --record-issues` over an old open issue with no story, outside the story, prints no `issues.no-story` and records no issue, shown by the new tests.
- `flai check` without `--story` still warns `issues.no-story` on it.
- `flai test` on the changed paths passes.

## Notes
