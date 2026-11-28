---
id: T-3596
type: task
nature: remediation
title: flai doc save, and so the dashboard's doc.save, refuses the workflow file on main and takes it only on a story's branch
status: done
parent: S-0844
owner: alex
created: 2026-11-17T12:36:39Z
updated: 2026-11-28T00:21:58Z
transitions:
  - to: ready
    at: 2026-11-27T23:38:19Z
    by: agent-S-0844
  - to: in-progress
    at: 2026-11-27T23:38:19Z
    by: agent-S-0844
  - to: done
    at: 2026-11-28T00:21:58Z
    by: agent-S-0844
stream: S-0844
tags: [flai, dashboard]
touches: [flai/internal/docedit/docedit.go, flai/internal/docedit/docedit_test.go, flai/internal/hostapi/writes_test.go]
usage:
  source: log
  seconds: 2594
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 99
      output: 474
      cache_read: 4474375
      cache_write: 202957
      cost: 1.801
---
# T-3596 flai doc save, and so the dashboard's doc.save, refuses the workflow file on main and takes it only on a story's branch

## Work

This task waits for no other task. It shares no file with the other tasks of the first layer.

ADR-0301 §7 says the dashboard's document and file saves refuse `design/workflows/board.md` on main and take it only on a story's branch (ADR-0145). No code does this today (I-0380, instance 2).

- The dashboard's `doc.save` runs `flai doc save`, which calls `docedit.Save` in `flai/internal/docedit/docedit.go`. Put the refusal there, so the CLI and the dashboard refuse alike. When `opt.Story` is empty and the path is the workflow file, return a `RefusedError`. Find the file's path with `workflow.Rel` of the design folder. The error must say that the file changes only on a story's branch, name ADR-0301 §7, and say to name a story.
- `file.save` already refuses every `.md` file (`FileReach` in `flai/internal/docedit/file.go`). It needs no change: confirm that in a test or in the task's notes.
- If `Show` can tell that the editor is on main, show the file read-only there with the same reason. If it cannot, leave `ModeOf` as it is and say so in the task's notes.
- Add tests to `flai/internal/docedit/docedit_test.go`: a save on main is refused, and a save that names a story commits on its branch. Add a refusal row to the `doc.save` refusal table in `flai/internal/hostapi/writes_test.go`.

## Done when

- `flai doc save design/workflows/board.md` in the main checkout with no `--story` is refused, naming ADR-0301 §7, and writes nothing.
- The same save with `--story S-nnnn` commits on that story's branch.
- The dashboard's `doc.save` with no story is refused with that message.
- `flai test` on the changed paths passes.

## Notes
