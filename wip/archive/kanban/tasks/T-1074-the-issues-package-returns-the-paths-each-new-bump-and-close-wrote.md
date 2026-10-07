---
id: T-1074
type: task
nature: improvement
title: The issues package returns the paths each new, bump, and close wrote
status: done
parent: S-0275
owner: alex
created: 2026-10-06T22:52:24Z
updated: 2026-10-07T08:20:15Z
transitions:
  - to: ready
    at: 2026-10-07T08:15:12Z
    by: agent-S-0275
  - to: in-progress
    at: 2026-10-07T08:15:12Z
    by: agent-S-0275
  - to: done
    at: 2026-10-07T08:20:15Z
    by: agent-S-0275
stream: S-0275
tags: [cli]
touches: [flai/internal/issues/issues.go, flai/internal/issues/issues_test.go, flai/internal/issues/record.go, flai/internal/issues/record_test.go]
usage:
  source: log
  seconds: 303
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 23
      output: 9371
      cache_read: 1444319
      cache_write: 39442
      cost: 0.705
---
# T-1074 The issues package returns the paths each new, bump, and close wrote

## Work

`adr.New` already returns the paths it wrote (`Result.Changed`). `issues.New`, `BumpWith`, `Close`, and `NewOrBump` do not, so `--commit` cannot know what to commit. Make each of them report the files it wrote: the issue file and `design/issues/summary.md` where the caller regenerates it. This is layer 1 and waits for no task: it shares no path with T-1072.

- Return the paths, or expose them on the result, without changing what is written.
- Keep the existing callers (`flai/cmd/issue.go`, `flai/internal/mcpserver/issues.go`) compiling. Change their call sites only as far as the new signature needs.
- Numbering stays as S-0252 left it, from the whole repository.

## Done when

- Tests in `issues_test.go` and `record_test.go` assert the paths returned by a new issue, a bump, a close, and a new-or-bump that bumps.
- `scripts/flai-test.sh` passes.

## Notes

Written by the planner. If the signature change reaches a caller outside these files, widen the task's touches to it.
