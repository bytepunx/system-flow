---
id: T-3020
type: task
nature: improvement
title: Time a warm flai check against main, record it in server-performance.md and the docs, and close I-0312
status: done
parent: S-0677
owner: alex
created: 2026-11-10T10:32:07Z
updated: 2026-11-10T20:10:27Z
transitions:
  - to: ready
    at: 2026-11-10T20:06:47Z
    by: system-flow
  - to: in-progress
    at: 2026-11-10T20:06:47Z
    by: system-flow
  - to: done
    at: 2026-11-10T20:10:27Z
    by: agent-S-0677
stream: S-0677
tags: [flai, performance, docs]
touches: [design/system/server-performance.md, design/system/flai-cli.md, docs/operators/runbooks/backup.md, design/issues/I-0312-a-whole-flai-check-takes-3-3-to-4-9-s-on-a-loaded-host-after-s-0601-not-under-2.md, design/issues/summary.md]
after: [T-3010, T-3016]
usage:
  source: log
  seconds: 220
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 61
      output: 20336
      cache_read: 7867393
      cache_write: 250127
      cost: 3.4604
---
# T-3020 Time a warm flai check against main, record it in server-performance.md and the docs, and close I-0312

## Work

Measure what the other three tasks changed, record it, and close the issue. Waits for T-3010 and T-3016, and through T-3016 for T-3012: it times their combined change and describes their kept files.

- Build `main` and the story's branch. Time a warm `flai check --strict` of each as a fresh process, in turn, five times each, after a first run of each that writes the kept files, as `design/system/server-performance.md` § "Cause 11" did. Note the host's load average. Both must give the same output and exit status. Time the steps in process too: the list, the first step alone, the digest rule, the conflict-marker rule, and `check.Run` in all.
- Add a section after "Cause 11" to `design/system/server-performance.md` with the fix, both tables, and what is left. `Repo.List`'s decode is S-0676's to fix (I-0310), so name it as what remains, if it does.
- Add the two new kept files to `docs/operators/runbooks/backup.md`'s table of `.flai-cache` files, and to `design/system/flai-cli.md` § "Internal structure", where the kept items and lint findings are described.
- Close I-0312 with `flai issue close I-0312 --reason "<what fixed it, with the before and after times>"`, which rewrites `design/issues/summary.md`.

## Done when

- `design/system/server-performance.md` records the measured times before and after, and the steps' times.
- `backup.md` and `flai-cli.md` name the new kept files.
- I-0312 is closed with a reason, `flai check --strict` is clean, and the markdown lint passes on the changed files.

## Notes

Drafted by the planner for S-0677.
