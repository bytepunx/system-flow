---
id: T-3108
type: task
nature: remediation
title: flai check reports a work item file that does not parse as one finding on its path and checks the rest
status: done
parent: S-0716
owner: alex
created: 2026-11-11T08:56:45Z
updated: 2026-11-11T15:08:45Z
transitions:
  - to: ready
    at: 2026-11-11T14:55:07Z
    by: agent-S-0716
  - to: in-progress
    at: 2026-11-11T14:55:08Z
    by: agent-S-0716
  - to: done
    at: 2026-11-11T15:08:45Z
    by: agent-S-0716
stream: S-0716
tags: [flai, check]
touches: [flai/internal/check/check.go, flai/internal/check/check_items.go, flai/internal/check/scope.go, flai/internal/check/items_test.go, flai/internal/check/scope_test.go, flai/internal/verify/story_test.go]
after: [T-3107]
usage:
  source: log
  seconds: 817
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 96
      output: 25259
      cache_read: 4757737
      cache_write: 191266
      cost: 2.603
---
# T-3108 flai check reports a work item file that does not parse as one finding on its path and checks the rest

## Work

`check.Run` in `flai/internal/check/check.go` returns `repo.List`'s error, so the close-out stopped with "flai check could not run: …T-3079-…md: no front matter". Once T-3107 lands, it goes on. Report each file the repository could not read as one error finding on that file's path, in `check_items.go`, and run every other rule over the items that did load. The file never reaches `oneItem`, so this is a pass of its own, as the threads rule's per-file pass in `check_records.go` is. Choose the rule's name: `item.front-matter` or a new `item.parse`.

In `scope.go`, a check scoped to a story keeps the finding only when the story's branch changes that file. Otherwise it is a finding outside the story, a note the close-out records in an issue as for any other. With no front matter the file names no parent, so the scope goes by its path alone. Test both cases in `scope_test.go`, and in `verify/story_test.go` show that a close-out no longer stops at check while another story's task file holds only `PLACEHOLDER`.

This task waits for T-3107, whose `Repo.List` stops failing and records the unreadable files this rule reports.

## Done when

- `flai check` on a project with a `PLACEHOLDER` task file reports one error finding on that file and the findings of every other item.
- A check scoped to another story lists it as outside the story. A check scoped to a story whose branch changes the file keeps it.
- The close-out's check step runs to its verdict with such a file present.
- `flai test` on the changed files passes.

## Notes
