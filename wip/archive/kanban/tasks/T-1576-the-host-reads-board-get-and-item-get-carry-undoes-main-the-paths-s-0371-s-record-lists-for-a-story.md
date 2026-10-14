---
id: T-1576
type: task
nature: feature
title: The host reads board.get and item.get carry undoes_main, the paths S-0371's record lists for a story
status: done
parent: S-0388
owner: alex
created: 2026-10-14T21:08:41Z
updated: 2026-10-14T23:15:23Z
transitions:
  - to: ready
    at: 2026-10-14T23:10:18Z
    by: agent-S-0388
  - to: in-progress
    at: 2026-10-14T23:10:18Z
    by: agent-S-0388
  - to: done
    at: 2026-10-14T23:15:23Z
    by: agent-S-0388
stream: S-0388
tags: [cli]
touches: [flai/internal/workitem/boardview.go, flai/internal/hostapi/hostapi.go, flai/internal/hostapi/hostapi_test.go, design/system/flai-cli.md]
usage:
  source: log
  seconds: 305
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 59
      output: 17986
      cache_read: 2611247
      cache_write: 96577
      cost: 1.4441
---
# T-1576 The host reads board.get and item.get carry undoes_main, the paths S-0371's record lists for a story

## Work

- Add `UndoesMain []string` with the JSON name `undoes_main,omitempty` to `workitem.BoardCard` in `flai/internal/workitem/boardview.go`.
- In `flai/internal/hostapi/hostapi.go`, have `board.get` fill it on each open story's card from the per-story record that S-0371's T-1477 writes under `.flai-cache/`. Read the record with the reader T-1477 adds in `flai/internal/storygit/undone.go`, so the path and format stay S-0371's.
- Fill the cards in `hostapi.go` after `workitem.NewBoardView` returns. Leave its signature alone, so `flai/cmd/board.go`, `flai/internal/inbox/inbox.go`, and `flai/internal/serve/agents.go` keep their calls.
- Have `item.get` give the same paths on `ItemWithChildren` as `undoes_main`, for a story only.
- Add the field to the `board.get` and `item.get` rows of the host reads in `design/system/flai-cli.md`.
- It waits for no task of this story. It needs S-0371's record, which the story's `after` on S-0371 gives it.

## Done when

- [ ] Tests in `flai/internal/hostapi/hostapi_test.go` show a story with a record carrying its paths on its board card and in `item.get`, and a story with none leaving `undoes_main` out
- [ ] `design/system/flai-cli.md` names `undoes_main` on both reads
- [ ] `flai test flai/internal/hostapi flai/internal/workitem design/system/flai-cli.md` passes

## Notes
