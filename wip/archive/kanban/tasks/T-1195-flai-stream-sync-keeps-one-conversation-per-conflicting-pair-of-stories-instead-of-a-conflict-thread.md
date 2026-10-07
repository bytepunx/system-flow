---
id: T-1195
type: task
nature: improvement
title: flai stream sync keeps one conversation per conflicting pair of stories instead of a conflict thread
status: done
parent: S-0332
owner: alex
created: 2026-10-07T20:15:03Z
updated: 2026-10-07T23:05:39Z
transitions:
  - to: ready
    at: 2026-10-07T22:50:25Z
    by: agent-S-0332
  - to: in-progress
    at: 2026-10-07T22:50:26Z
    by: agent-S-0332
  - to: done
    at: 2026-10-07T23:05:39Z
    by: agent-S-0332
stream: S-0332
tags: [flai]
touches: [flai/internal/storygit/sync.go, flai/internal/storygit/sync_test.go, flai/cmd/stream_sync.go, flai/cmd/stream_sync_test.go, flai/cmd/stream.go]
after: [T-1194]
usage:
  source: log
  seconds: 913
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 100
      output: 39616
      cache_read: 6062679
      cache_write: 192215
      cost: 3.158
---
# T-1195 flai stream sync keeps one conversation per conflicting pair of stories instead of a conflict thread

## Work

Rework `reportConflicts` in `flai/internal/storygit/sync.go`. It waits for T-1194, whose ADR settles the behaviour. It shares no path with T-1196 or T-1197 and runs beside them.

- A new conflict opens a conversation between the two stories, `about` the conflicting paths, with the text `ConflictText` writes today; a change of paths adds an entry; a clean merge or a story no longer open closes it. The author is flai's conflict author.
- A conflict thread already open is still resolved as today, and no new one is opened.
- `res.Branches[i]` names the conversation instead of the thread; `flai stream sync` prints it, and `--json` returns it.
- Update the help of `flai stream sync` in `flai/cmd/stream.go`.

## Done when

- Tests cover a new conflict, changed paths, a clean merge closing it, a story gone, and an old conflict thread resolved without a new one.
- `flai test` passes on the changed paths.
- `flai check --strict` passes.

## Notes
