---
id: T-2940
type: task
nature: improvement
title: flai task done narrows the closing task's touches, and the story's, to the files the branch changes
status: done
parent: S-0674
owner: alex
created: 2026-11-10T08:50:09Z
updated: 2026-11-11T02:00:03Z
transitions:
  - to: ready
    at: 2026-11-11T01:43:21Z
    by: system-flow
  - to: in-progress
    at: 2026-11-11T01:44:27Z
    by: system-flow
  - to: done
    at: 2026-11-11T02:00:03Z
    by: agent-S-0674
stream: S-0674
tags: [flai, cli]
touches: [flai/internal/taskdone/taskdone.go, flai/internal/taskdone/taskdone_test.go, flai/internal/itemedit/widen.go, flai/internal/itemedit/widen_test.go, flai/cmd/task_done.go, docs/users/flai.md, docs/users/flai-reference.md, design/system/flai-cli.md, design/system/workflow.md, design/conventions/work-management.md, design/conventions/git.md, template/root/design/conventions/work-management.md, template/root/design/conventions/git.md, template/CHANGELOG.md]
after: [T-2936]
usage:
  source: log
  seconds: 936
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 105
      output: 29675
      cache_read: 14935250
      cache_write: 345705
      cost: 5.7976
---
# T-2940 flai task done narrows the closing task's touches, and the story's, to the files the branch changes

## Work

Build decision 2 of T-2936's ADR. It waits for T-2936 because it follows that ADR. It shares no path with the drift task, so the two run together.

- In `flai/internal/itemedit/widen.go`, add a narrowing beside `Widen` and `WidenStory`. It removes from the task's touches each file touch that `branchChanges` does not list. It removes that file from the story's touches when no other open task of the story names it. It never removes a folder touch, a sub-project name or tag, or a path under `wip/`. It stamps the edit and records the notice as `Widen` does.
- In `flai/internal/taskdone/taskdone.go`, `touches()` calls it after widening. Run it even when `r.res.Commit` is nil: a task that committed nothing still claims what it predicted. Report what was dropped in `Touches`, for the text and `--json` answers. Update the comment that says the step never removes a touch.
- Tests in `widen_test.go` and `taskdone_test.go`. A task touching `a.go` and `b.go`, whose branch changes only `a.go`, ends with `a.go` alone. The story drops `b.go` unless another open task names it. A folder touch is kept.
- Docs: `docs/users/flai.md` § Closing a task, step 6. The `flai task done` row of `design/system/flai-cli.md`. The task-done rule of `design/conventions/work-management.md` that says it widens the touches, and the same line in `template/root/design/conventions/work-management.md`. A line under `## Unreleased` in `template/CHANGELOG.md`.

## Done when

- The tests pass. A closed task keeps no file touch the branch does not change.
- `flai test` on the changed paths passes.
- The two docs, the convention, and its template copy say that task done widens and narrows.

## Notes
