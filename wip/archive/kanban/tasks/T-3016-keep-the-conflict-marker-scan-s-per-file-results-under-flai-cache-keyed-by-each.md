---
id: T-3016
type: task
nature: improvement
title: Keep the conflict-marker scan's per-file results under .flai-cache, keyed by each file's size and modification time
status: done
parent: S-0677
owner: alex
created: 2026-11-10T10:31:44Z
updated: 2026-11-10T20:06:22Z
transitions:
  - to: ready
    at: 2026-11-10T19:52:05Z
    by: system-flow
  - to: in-progress
    at: 2026-11-10T19:52:05Z
    by: system-flow
  - to: done
    at: 2026-11-10T20:06:22Z
    by: agent-S-0677
stream: S-0677
tags: [flai, performance, check]
touches: [flai/internal/conflictmark/kept.go, flai/internal/conflictmark/kept_test.go, flai/internal/check/check_records.go, flai/internal/check/check_test.go, flai/cmd/promote_test.go]
after: [T-3012]
usage:
  source: log
  seconds: 857
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 53
      output: 17811
      cache_read: 6890328
      cache_write: 219063
      cost: 3.0307
---
# T-3016 Keep the conflict-marker scan's per-file results under .flai-cache, keyed by each file's size and modification time

## Work

Once the digest rule keeps its results, `conflictMarkers` (`flai/internal/check/check_records.go`) is the slowest rule left run at once, 0.45 to 1.10 s on a loaded host (I-0312). It reads every markdown file under the design, docs, and wip folders and at the root.

Keep each file's conflict-marker lines between processes, in `flai/internal/conflictmark/kept.go`, by the rules of the digest task's kept file and S-0600's lint findings:

- keyed per file by its path, size, and modification time, and as a whole by the scanned folder and a format version;
- only a read made at least two seconds after the file last changed is kept;
- written again only when it changed, and only into a `.flai-cache` that exists;
- a miss is never an error.

A file removed since the kept run drops out. A checked write's run keeps scanning only the files it wrote.

Waits for T-3012: both change `flai/internal/check/check_records.go`, and this follows the kept file T-3012 shapes.

## Done when

- A second whole check reads no unchanged markdown file for conflict markers. A test in `flai/internal/conflictmark/kept_test.go` shows it.
- A marker added to a kept file is still reported, in a test in `flai/internal/check/check_test.go`.
- `flai test` passes on the changed paths.

## Notes

Drafted by the planner for S-0677.
