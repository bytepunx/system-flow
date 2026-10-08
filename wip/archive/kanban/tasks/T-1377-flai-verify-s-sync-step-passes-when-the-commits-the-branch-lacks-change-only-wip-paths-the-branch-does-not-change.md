---
id: T-1377
type: task
nature: improvement
title: flai verify's sync step passes when the commits the branch lacks change only wip paths the branch does not change
status: done
parent: S-0347
owner: alex
created: 2026-10-08T08:44:02Z
updated: 2026-10-08T10:36:09Z
transitions:
  - to: ready
    at: 2026-10-08T10:31:18Z
    by: agent-S-0347
  - to: in-progress
    at: 2026-10-08T10:31:18Z
    by: agent-S-0347
  - to: done
    at: 2026-10-08T10:35:55Z
    by: agent-S-0347
stream: S-0347
tags: [cli]
touches: [flai/internal/verify/story.go, flai/internal/verify/story_test.go, flai/internal/verify/paths.go, flai/internal/verify/paths_test.go]
after: [T-1375]
usage:
  source: log
  seconds: 277
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 34
      output: 12065
      cache_read: 1971748
      cache_write: 63487
      cost: 1.0512
---
# T-1377 flai verify's sync step passes when the commits the branch lacks change only wip paths the branch does not change

## Work

`storyRun.sync` in `flai/internal/verify/story.go` counts `git rev-list --count HEAD..<base>` and fails on any commit the branch lacks. Build the rule T-1375's ADR decides:

- In `flai/internal/verify/paths.go`, beside `ChangedPaths`, add a helper that lists the commits of the base the branch lacks (`git rev-list HEAD..<base>`) with the paths each changes (`git diff-tree --no-commit-id --name-only -r`, both sides of a rename), and splits them into those that change only paths under the manifest's `wip` folder (resolved through the project, never hard-coded) and the rest.
- `sync` passes when there is no rest and none of those `wip` paths is among the branch's changed paths. It then adds a note to the report naming how many commits it passed over and their short hashes. Otherwise it fails as now, naming the commits outside `wip` first.
- Reproduce I-0119 in `story_test.go` on a temporary repository: a story branch is verified, then main gains a commit that changes only `wip/kanban/...`; `sync` passes with the note. A main commit that changes a file outside `wip/`, and one that changes a `wip/` file the branch also changes, each fail it.
- Unit-test the helper in `paths_test.go`, including a merge commit and a rename.

Waits for T-1375, whose ADR settles the rule. Runs beside T-1379: they share no path.

## Done when

- The new tests in `story_test.go` and `paths_test.go` pass, and the I-0119 reproduction fails on the code before the change.
- `flai verify` on a story whose branch lacks only `wip`-only commits passes its `sync` step and prints the note; on any other commit it still stops there.
- `flai test flai/internal/verify` passes.

## Notes

Drafted by the planner for S-0347. S-0341, in progress, changes `story.go` and `story_test.go` too; sync first and build on its record.
